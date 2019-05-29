package queue

import (
	"sort"
	"sync"

	"github.com/golang-collections/go-datastructures/queue"
)

type priorityItems []Item

// Item is an item that can be added to the priority queue.
type Item struct {
	item     interface{}
	priority int
}

func (item Item) Compare(other queue.Item) int {
	if item.priority > other.(Item).priority {
		return 1
	} else if item.priority == other.(Item).priority {
		return 0
	}
	return -1
}

type PriorityQueue struct {
	empty chan bool
	items priorityItems
	lock  sync.Mutex
}

func New(size int) *PriorityQueue {
	return &PriorityQueue{
		empty: make(chan bool),
		items: make([]Item, 0, size),
	}
}

func (items *priorityItems) get(number int) []Item {
	returnItems := make([]Item, 0, number)
	index := 0
	for i := 0; i < number; i++ {
		if i >= len(*items) {
			break
		}

		returnItems = append(returnItems, (*items)[i])
		(*items)[i] = Item{}
		index++
	}

	*items = (*items)[index:]
	return returnItems
}

func (items *priorityItems) insert(item Item) {
	if len(*items) == 0 {
		*items = append(*items, item)
		return
	}

	equalFound := false
	i := sort.Search(len(*items), func(i int) bool {
		result := (*items)[i].Compare(item)
		if result == 0 {
			equalFound = true
		}
		return result >= 0
	})

	if equalFound {
		return
	}

	if i == len(*items) {
		*items = append(*items, item)
		return
	}

	*items = append(*items, Item{})
	copy((*items)[i+1:], (*items)[i:])
	(*items)[i] = item
}

func (priorityQueue *PriorityQueue) Peek(number int) []Item {

	priorityQueue.lock.Lock()
	defer priorityQueue.lock.Unlock()

	queueSize := len(priorityQueue.items)
	for queueSize == 0 {
		<-priorityQueue.empty
	}
	if number > queueSize {
		return priorityQueue.items[:queueSize-1]
	}
	return priorityQueue.items[:number-1]
}

func (priorityQueue *PriorityQueue) Pull(number int) []Item {
	priorityQueue.lock.Lock()
	defer priorityQueue.lock.Unlock()

	for priorityQueue.IsEmpty() {
		<-priorityQueue.empty
	}
	return priorityQueue.items.get(number)
}

func (priorityQueue *PriorityQueue) Push(item interface{}, priority int) error {

	priorityQueue.lock.Lock()
	defer priorityQueue.lock.Unlock()

	priorityQueue.items.insert(Item{item, priority})

	priorityQueue.empty <- false
	return nil
}

func (priorityQueue *PriorityQueue) Size() int {
	priorityQueue.lock.Lock()
	defer priorityQueue.lock.Unlock()

	return len(priorityQueue.items)
}

func (priorityQueue *PriorityQueue) IsEmpty() bool {
	priorityQueue.lock.Lock()
	defer priorityQueue.lock.Unlock()

	return len(priorityQueue.items) == 0
}
