// Package handlers provides HTTP handlers for managing tasks using an etcd database.
package handlers

// Import necessary packages
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"task-organizer/models"
	"time"

	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// cli is the global etcd client used to interact with the etcd database.
var cli *clientv3.Client

// init function initializes the connection to the etcd database.
// This function will be automatically executed when the package is imported.
func init() {
	Connection()
}

// Connection function establishes a connection to the etcd database.
func Connection() {
	var err error
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("Failed to connect to etcd:", err)
	}
}

// GetAllTasks function retrieves all tasks from the etcd database and returns them as a JSON array.
// It handles the HTTP GET request for the route "/tasks".
func GetAllTasks(c *gin.Context) {
	// Fetch all tasks from the database with the "tasks/" prefix
	resp, err := cli.Get(context.Background(), "tasks/", clientv3.WithPrefix())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	var tasks []models.Task
	for _, kv := range resp.Kvs {
		var task models.Task
		// Unmarshal the task data from etcd's byte format to the Task struct
		if err := json.Unmarshal(kv.Value, &task); err != nil {
			// If there is an issue parsing the task data, return an error response
			fmt.Printf("Failed to parse task: %s\n", string(kv.Value))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse task"})
			return
		}
		tasks = append(tasks, task)
	}

	// Sort the tasks by ID before returning the JSON response
	sort.Slice(tasks, func(i, j int) bool {
		return *tasks[i].ID < *tasks[j].ID
	})

	c.JSON(http.StatusOK, tasks)
}

// GetTask function retrieves a task with the specified ID from the etcd database.
// It handles the HTTP GET request for the route "/tasks/{id}".
func GetTask(c *gin.Context) {
	var task models.Task

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Fetch the task with the provided ID from the database
	resp, err := cli.Get(ctx, "tasks/"+c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		return
	}

	// Check if the task with the specified ID exists in the database
	if len(resp.Kvs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Unmarshal the task data from etcd's byte format to the Task struct
	if err := json.Unmarshal(resp.Kvs[0].Value, &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse task"})
		return
	}

	// Return the JSON response with the retrieved task data
	c.IndentedJSON(http.StatusOK, task)
}

// CreateTask function creates a new task in the etcd database.
// It handles the HTTP POST request for the route "/tasks".
func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		// If there is an issue with the JSON request body, return an error response
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the request contains an ID (IDs are auto-generated and not allowed in the request)
	if task.ID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Manual ID entry is not allowed"})
		return
	}

	// Generate a unique ID for the task using the GenerateUniqueID function
	id := models.GenerateUniqueID()
	task.ID = &id

	// Check if the Title field is empty
	if task.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title cannot be empty"})
		return
	}

	// Marshal the task data to JSON format
	data, err := json.Marshal(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Store the task in the etcd database with the generated ID as the key
	_, err = cli.Put(ctx, "tasks/"+strconv.Itoa(*task.ID), string(data))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

// UpdateTask function updates an existing task in the etcd database.
// It handles the HTTP PUT request for the route "/tasks/{id}".
func UpdateTask(c *gin.Context) {
	// Get the task ID from the URL path
	taskID := c.Param("id")

	// Fetch the existing task from the database
	resp, err := cli.Get(context.Background(), "tasks/"+taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Check if the task with the specified ID exists in the database
	if len(resp.Kvs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Unmarshal the existing task data into a temporary task variable
	var existingTask models.Task
	err = json.Unmarshal(resp.Kvs[0].Value, &existingTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Bind the JSON request body to the update request
	var updateReq models.UpdateReq
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate if the title field is empty
	if updateReq.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title cannot be empty"})
		return
	}

	// Update the existing task with the new title
	existingTask.Title = updateReq.Title

	// Marshal the updated task into JSON format
	taskJSON, err := json.Marshal(existingTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save the updated task back to the database
	_, err = cli.Put(context.Background(), "tasks/"+taskID, string(taskJSON))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updateReq)
}

// DeleteTask function deletes a task with the specified ID from the etcd database.
// It handles the HTTP DELETE request for the route "/tasks/{id}".
func DeleteTask(c *gin.Context) {
	taskID := c.Param("id")

	// Check if the task exists before attempting to delete it
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := cli.Get(ctx, "tasks/"+taskID)
	cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if the task with the specified ID exists in the database
	if len(resp.Kvs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Perform the delete operation for the task
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	_, err = cli.Delete(ctx, "tasks/"+taskID)
	cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

// DeleteAllTasks function deletes all tasks from the etcd database.
// It handles the HTTP DELETE request for the route "/tasks/".
func DeleteAllTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get all keys with the "tasks/" prefix from the database
	resp, err := cli.Get(ctx, "tasks/", clientv3.WithPrefix())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	// Check if there are any tasks to delete
	if len(resp.Kvs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No tasks to delete"})
		return
	}

	// Delete all tasks one by one
	for _, kv := range resp.Kvs {
		_, err := cli.Delete(ctx, string(kv.Key))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tasks"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "All tasks deleted"})
}
