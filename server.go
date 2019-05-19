package main

import (
	"./server"
)

func main() {
	queueManager := server.NewQueueManager()
	queueManager.Manage("localhost", 2020)
}
