package server

import (
	"../marshaller"
	"../message"
	"github.com/golang-collections/go-datastructures/queue"
)

type User struct {
	identifier        string
	handler           *RequestHandler
	isConnected       bool
	queue             *queue.PriorityQueue
	connectionWarning chan bool
}

func newUser(identifier string, handler *RequestHandler) *User {
	connectionWarning := make(chan bool, 1)
	connectionWarning <- true
	return &User{identifier, handler, true, queue.NewPriorityQueue(20), connectionWarning}
}

func (user *User) getIdentifier() string {
	return user.identifier
}

func (user *User) setHandler(handler *RequestHandler) {
	user.handler = handler
}

func (user *User) setConnected(isConnected bool) {
	user.isConnected = isConnected
	if isConnected {
		user.connectionWarning <- true
	} else {
		<-user.connectionWarning
	}
}

func (user *User) addTopicItem(item Item) {
	user.queue.Put(item)
}

func (user *User) handleSubscriptions() {
	for {
		for !user.isConnected {
			<-user.connectionWarning
		}

		if items, err := user.queue.Get(1); err == nil {
			user.sendTopicItem(items[0].(Item))
		}
	}
}

func (user *User) sendTopicItem(item Item) {
	topicName := item.topicName
	data := item.item

	message := message.NewRequest(user.getIdentifier(), message.Topic, topicName, data, 0)
	marshalledData := marshaller.JSONMarshaller{}.MarshalRequest(*message)
	if err := user.handler.Send(marshalledData); err != nil { // disconection
		user.handler.conn.Close()
	}
}
