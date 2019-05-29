package server

import (
	"container/list"
	"math"

	"github.com/golang-collections/go-datastructures/queue"
)

type Item struct {
	item      interface{}
	topicName string
	priority  int
}

func (item Item) Compare(other queue.Item) int {
	if item.priority > other.(Item).priority {
		return 1
	} else if item.priority == other.(Item).priority {
		return 0
	}
	return -1
}

type Topic struct {
	name                string
	queue               *queue.PriorityQueue
	subscribers         *list.List
	//hasActiveSubscriber chan bool
}

func NewTopic(name string) *Topic {
	return &Topic{
		name:        name,
		queue:       queue.NewPriorityQueue(20),
		subscribers: list.New(),
	}
}

func NewStream(name string, size float64) *Topic {
	return &Topic{
		name:        name,
		queue:       queue.NewPriorityQueue(int(math.Ceil(size / 1024))),
		subscribers: list.New(),
	}
}

func (topic *Topic) push(item interface{}, priority int) error {
	priorityItem := Item{item, topic.name, priority}
	topic.alertSubscribers(priorityItem)
	return topic.queue.Put(priorityItem)
}

func (topic *Topic) subscribe(subscriber *User) {
	topic.subscribers.PushBack(subscriber)
}

func (topic *Topic) isSubscriber(subscriber *User) bool {
	for e := topic.subscribers.Front(); e != nil; e = e.Next() {
		if e.Value.(*User).getIdentifier() == subscriber.getIdentifier() {
			return true
		}
	}
	return false
}

func (topic *Topic) alertSubscribers(item Item) {
	for subscriber := topic.subscribers.Front(); subscriber != nil; subscriber = subscriber.Next() {
		subscriber.Value.(*User).addTopicItem(item)
	}
}
