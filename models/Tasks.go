package models

import (
	"time"

	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Task represents a task with an ID, title, and completion status.
type Task struct {
	ID        string `json:"id"`        // ID of the task (string format)
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

// Init initializes two etcd clients for different endpoints and returns their handlers.
// It creates and configures etcd clients for the specified endpoints.
// Returns handler1, handler2, and any error encountered during client creation.
func Init() (*Handler, *Handler, error) {
	// Create a client for the first etcd member (using port 2379)
	cli1, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		cli1.Close()
		return nil, nil, err
	}

	// Create a client for the second etcd member (using port 2380)
	cli2, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2380"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		cli2.Close()
		return nil, nil, err
	}

	// Create handlers for each client
	handler1 := &Handler{Client: cli1}
	handler2 := &Handler{Client: cli2}

	// Return the handlers for the two etcd clients and no error
	return handler1, handler2, nil
}

// GenerateUniqueID generates a unique ID for a new task using UUID.
func GenerateUniqueID() string {
	return uuid.New().String()
}
