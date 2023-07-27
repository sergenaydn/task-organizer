package models

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// counter is a private variable used to generate unique IDs for tasks.
var counter = 1

// Task represents a task with an ID, title, and completion status.
type Task struct {
	ID        *int   `json:"id,string"` // ID of the task (pointer to an integer to allow nil)
	Title     string `json:"title"`     // Title of the task
	Completed bool   `json:"completed"` // Completion status of the task
}

// UpdateReq represents a request to update a task with a new title and completion status.
type UpdateReq struct {
	Title     string `json:"title"`     // New title for the task update
	Completed bool   `json:"completed"` // New completion status for the task update
}

type Handler struct {
	Client *clientv3.Client
}

// GenerateUniqueID generates a unique ID for a new task.
// It uses a simple counter approach to increment the ID for each new task.
func GenerateUniqueID() int {
	id := counter
	counter++
	return id
}

func Connection() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}
