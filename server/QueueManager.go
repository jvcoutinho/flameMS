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
		// Sempre que um cliente se conecta, dois request handlers são criados.
		requestHandler := listener.Accept() 
		queueHandler := listener.Accept()
		go manager.handleRequests(marshaller, requestHandler, queueHandler)
	}
}

func (manager *QueueManager) handleRequests(marshaller marshaller.Marshaller, requestHandler *RequestHandler, queueHandler *RequestHandler) {
	var userIdentifier string
	for {
		// Espera um request.
		request, err := receiveRequest(marshaller, requestHandler)

		// Checa se houve problemas na conexão.
		if err != nil {
			if connectionSuccessful := handleConnectionIssues(requestHandler, err); !connectionSuccessful {
				manager.disconnectUser(userIdentifier)
				break
			}
		}

		// 	Demultiplexa o request de acordo com o atributo Operation.
		switch request.GetOperation() {
		case message.Register:
			userIdentifier, err = manager.registerUser(request.GetRequestor(), queueHandler)
		case message.Initialize:
			err = manager.initializeQueue(request.GetQueueName())
		case message.Publish:
			err = manager.putPublishing(request.GetQueueName(), request.GetBody(), request.GetPriority())
		case message.Subscribe:
			err = manager.insertSubscriber(request.GetQueueName(), request.GetRequestor()+":"+requestHandler.GetHost())
		case message.CheckExistence:
			err = manager.checkQueueExistence(request.GetQueueName())
		case message.Stream:
			err = manager.prepareStreaming(request.GetQueueName(), request.GetBody().(float64))
		}
		returnResponse(marshaller, requestHandler, message.NewResponse(err))

	}
}

func (manager *QueueManager) disconnectUser(user string) {
	manager.users[user].setConnected(false)
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
		return absoluteIdentifier, nil
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

func (manager *QueueManager) putPublishing(queueName string, item interface{}, priority int) error {
	topic := manager.queues[queueName]
	if err := topic.push(item, priority); err != nil {
		return errors.New("Publishing unsuccessful, details:" + err.Error())
	}
	return nil
}

func (manager *QueueManager) insertSubscriber(queueName string, subscriber string) error {
	topic := manager.queues[queueName]
	user := manager.users[subscriber]
	if topic.isSubscriber(user) {
		return errors.New("You are already a subscriber of the topic '" + queueName + "'")
	}
	topic.subscribe(user)
	go user.handleSubscriptions()
	fmt.Println("Subscription successful.")
	return nil
}

func (manager *QueueManager) checkQueueExistence(queueName string) error {
	if _, exists := manager.queues[queueName]; exists {
		return nil
	}
	return errors.New("Topic named " + queueName + " has not been initialized yet")
}

func (manager *QueueManager) prepareStreaming(streamName string, size float64) error {
	if _, exists := manager.queues[streamName]; exists {
		return errors.New("Stream " + streamName + " already exists")
	}
	manager.queues[streamName] = NewStream(streamName, size)
	fmt.Println("Stream", streamName, "created successfully.")
	return nil
}
