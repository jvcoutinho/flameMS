package server

import (
	"container/list"

	"github.com/golang-collections/go-datastructures/queue"
)

type Topic struct {
	name        string
	queue       *queue.Queue
	subscribers *list.List
}

func NewTopic(name string) *Topic {
	return &Topic{
		name:        name,
		queue:       queue.New(20),
		subscribers: list.New(),
	}
}

func (topic *Topic) push(item interface{}) error {
	return topic.queue.Put(item)
}

func (topic *Topic) subscribe(subscriber string) {
	topic.subscribers.PushBack(subscriber)
}

func (topic *Topic) isSubscriber(subscriber string) bool {
	for e := topic.subscribers.Front(); e != nil; e = e.Next() {
		if e.Value == subscriber {
			return true
		}
	}
	return false
}
