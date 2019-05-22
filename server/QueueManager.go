package server

import (
	"errors"
	"fmt"
	"strconv"

	"../marshaller"
	"../message"
)

// QueueManager ffifui
type QueueManager struct {
	queues map[string]*Topic
	users  map[string]*User
}

// NewQueueManager Creates new Queue Manager.
func NewQueueManager() *QueueManager {
	return &QueueManager{make(map[string]*Topic), make(map[string]*User)}
}

func (manager *QueueManager) Manage(host string, port int) error {
	listener := NewListener(host, strconv.Itoa(port))
	marshaller := &marshaller.JSONMarshaller{}
	for {
		handler := listener.Accept()
		go manager.handleRequests(marshaller, handler)
	}
}

func (manager *QueueManager) handleRequests(marshaller marshaller.Marshaller, handler *RequestHandler) {
	var userIdentifier string
	for {
		request, err := receiveRequest(marshaller, handler)

		if err != nil {
			if connectionSuccessful := handleConnectionIssues(handler, err); !connectionSuccessful {
				manager.disconnectUser(userIdentifier)
				break
			}
		}

		switch request.GetOperation() {
		case message.Register:
			userIdentifier, err = manager.registerUser(request.GetRequestor(), handler)
		case message.Initialize:
			err = manager.initializeQueue(request.GetQueueName())
		case message.Publish:
			err = manager.putPublishing(request.GetQueueName(), request.GetBody())
		case message.Subscribe:
			err = manager.insertSubscriber(request.GetQueueName(), request.GetRequestor())
		case message.CheckExistence:
			err = manager.checkExistence(request.GetQueueName())
		}
		returnResponse(marshaller, handler, message.NewResponse(err))

	}
}

func (manager *QueueManager) disconnectUser(user string) {
	manager.users[user].setConnected(false)
	manager.users[user].setHandler(nil)
}

func handleConnectionIssues(handler *RequestHandler, err error) bool {
	return false // TODO: lidar com desconexões.
}

func receiveRequest(marshaller marshaller.Marshaller, handler *RequestHandler) (message.Request, error) {
	serializedData, err := handler.Receive()
	if err != nil {
		return message.Request{}, err
	}
	return marshaller.UnmarshalRequest(serializedData), nil
}

func returnResponse(marshaller marshaller.Marshaller, handler *RequestHandler, response *message.Response) {
	serializedData := marshaller.MarshalResponse(*response)
	handler.Send(serializedData)
}

/*
 * REQUESTS
 */

func (manager *QueueManager) registerUser(identifier string, handler *RequestHandler) (string, error) {
	// An user is identified by "identifier:host".
	absoluteIdentifier := identifier + ":" + handler.GetHost()

	// If it's a known user, then we just need to update its connection status and handler.
	if user, exists := manager.users[absoluteIdentifier]; exists {
		user.setConnected(true)
		user.setHandler(handler)
	}

	// Otherwise, we create a new one.
	manager.users[absoluteIdentifier] = newUser(absoluteIdentifier, handler)
	fmt.Println("User", absoluteIdentifier, "registered.")
	return absoluteIdentifier, nil
}

func (manager *QueueManager) initializeQueue(queueName string) error {
	if _, exists := manager.queues[queueName]; exists {
		return errors.New("Topic " + queueName + " already exists")
	}
	manager.queues[queueName] = NewTopic(queueName)
	fmt.Println("Queue", queueName, "created successfully.")
	return nil
}

func (manager *QueueManager) putPublishing(queueName string, item interface{}) error {
	topic := manager.queues[queueName]
	if err := topic.push(item); err != nil {
		return errors.New("Publishing unsuccessful, details:" + err.Error())
	}
	return nil
}

func (manager *QueueManager) insertSubscriber(queueName string, subscriber string) error {
	topic := manager.queues[queueName]
	if topic.isSubscriber(subscriber) {
		return errors.New("You are already a subscriber")
	}
	topic.subscribe(subscriber)
	fmt.Println("Subscription successful.")
	return nil
}

func (manager *QueueManager) checkExistence(queueName string) error {
	if _, exists := manager.queues[queueName]; exists {
		return nil
	}
	return errors.New("Topic named " + queueName + " has not been initialized yet")
}
