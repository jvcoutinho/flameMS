package server

import (
	"fmt"
	"strconv"

	"../marshaller"
	"../message"
)

// QueueManager ffifui
type QueueManager struct {
	queues map[string]*Topic
}

// NewQueueManager Creates new Queue Manager.
func NewQueueManager() *QueueManager {
	return &QueueManager{make(map[string]*Topic)}
}

func (manager *QueueManager) Manage(host string, port int) error {
	listener := NewListener(host, strconv.Itoa(port))
	marshaller := &marshaller.JSONMarshaller{}
	for {
		handler := listener.Accept()
		go manager.handleRequests(marshaller, handler)
	}
}

func (manager *QueueManager) handleRequests(marshaller *marshaller.JSONMarshaller, handler *RequestHandler) {
	for {
		message := receiveRequest(marshaller, handler)
		switch message.GetOperation() {
		case "init":
			manager.initializeQueue(message.GetQueueName())
		case "push":
			manager.push(message.GetQueueName(), message.GetBody())
		case "peek":
			manager.peek(message.GetQueueName())
		}
	}
}

func receiveRequest(marshaller *marshaller.JSONMarshaller, handler *RequestHandler) message.Message {
	serializedData := handler.Receive()
	return marshaller.Unmarshal(serializedData)
}

/*
 * REQUESTS
 */

func (manager *QueueManager) initializeQueue(queueName string) {
	//TODO: checar se já existe.
	manager.queues[queueName] = NewTopic(queueName)
	fmt.Println("Queue", queueName, "created successfully.")
}

func (manager *QueueManager) push(queueName string, item interface{}) {
	topic := manager.queues[queueName]
	if err := topic.push(item); err != nil {
		fmt.Println(err.Error())
		// TODO.
	}
	fmt.Println("Element added to topic", queueName)
}

func (manager *QueueManager) peek(queueName string) {
	topic := manager.queues[queueName]
	topic.peek()
}
