package models

import (
	"log"
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

// GenerateUniqueID generates a unique ID for a new task using UUID.
func GenerateUniqueID() string {
	return uuid.New().String()
}

// Connection2379 creates a client for the first etcd member using port 2379
func Connection2379() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("Failed to connect to etcd", err)
	}
	return cli, nil
}

// Connection2380 creates a client for the second etcd member using port 2380
func Connection2380() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2380"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}
