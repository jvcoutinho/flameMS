package server

import (
	"github.com/golang-collections/go-datastructures/queue"
)

type Topic struct {
	name  string
	queue *queue.Queue
}

func NewTopic(name string) *Topic {
	return &Topic{
		name:  name,
		queue: queue.New(20),
	}
}

func (topic *Topic) push(item interface{}) error {
	return topic.queue.Put(item)
}

func (topic *Topic) peek() (interface{}, error) {
	return topic.queue.Get(1)
}
