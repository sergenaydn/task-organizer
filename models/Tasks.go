package models

import (
	"log"
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

// Connection2381 creates a client for the third etcd member using port 2381
func Connection2381() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2381"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// Connection2382 creates a client for the fourth etcd member using port 2382
func Connection2382() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2382"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// Connection2383 creates a client for the fifth etcd member using port 2383
func Connection2383() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2383"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// Connection2384 creates a client for the sixth etcd member using port 2384
func Connection2384() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2384"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// Connection2385 creates a client for the seventh etcd member using port 2385
func Connection2385() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2385"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}
