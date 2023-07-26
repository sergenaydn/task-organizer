package handlers

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
	_ "github.com/swaggo/gin-swagger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var cli *clientv3.Client

func init() {
	Connection()
}

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

// GetAllTasks godoc
// @Summary Get a list of all tasks
// @Description Returns the data of all the tasks
// @Tags Tasks
// @Accept json
// @Produce json
// @Success 200 {array} models.Task
// @Failure 500 {object} nil
// @Router /tasks [get]

func GetAllTasks(c *gin.Context) {
	resp, err := cli.Get(context.Background(), "tasks/", clientv3.WithPrefix())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	var tasks []models.Task
	for _, kv := range resp.Kvs {
		var task models.Task
		if err := json.Unmarshal(kv.Value, &task); err != nil {
			// Print the value that failed to parse to help diagnose the issue
			fmt.Printf("Failed to parse task: %s\n", string(kv.Value))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse task"})
			return
		}
		tasks = append(tasks, task)
	}

	// Sort the tasks by ID before returning
	sort.Slice(tasks, func(i, j int) bool {
		return *tasks[i].ID < *tasks[j].ID
	})

	c.JSON(http.StatusOK, tasks)
}

// GetTask godoc
// @Summary Get a task by ID
// @Description Retrieves a task with the specified ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID" Format(int64)
// @Success 200 {object} models.Task
// @Failure 404 {object} nil
// @Failure 500 {object} nil
// @Router /tasks/{id} [get]
func GetTask(c *gin.Context) {
	var task models.Task

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := cli.Get(ctx, "tasks/"+c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		return
	}

	// Check if the key exists in the response
	if len(resp.Kvs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse task"})
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

// CreateTask godoc
// @Summary Create a new task
// @Description Creates a new task with a unique ID if not provided
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task object to be created"
// @Success 201 {object} models.Task
// @Failure 400 {object} nil
// @Router /tasks [post]
func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the request contains an ID
	if task.ID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Manual ID entry is not allowed"})
		return
	}

	// Generate a unique ID using the GenerateUniqueID function
	id := models.GenerateUniqueID()
	task.ID = &id

	// Check if the Title is empty
	if task.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Idea and owner cannot be empty"})
		return
	}

	data, err := json.Marshal(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Store the task in the database with the generated ID
	_, err = cli.Put(ctx, "tasks/"+strconv.Itoa(*task.ID), string(data))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

// CreateTask godoc
// @Summary Create a new task
// @Description Creates a new task with a unique ID if not provided
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body models.Task true "Task object to be created"
// @Success 201 {object} models.Task
// @Failure 400 {object} nil
// @Router /tasks [put]
func UpdateTask(c *gin.Context) {
	// Get the task ID from the URL path
	taskID := c.Param("id")

	// Fetch the existing task from the database
	resp, err := cli.Get(context.Background(), "tasks/"+taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

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

	// Validate if the idea field is empty
	if updateReq.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title cannot be empty"})
		return
	}

	// Update the existing task with the new title
	existingTask.Title = updateReq.Title

	// Marshal the updated task into JSON
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

// DeleteTask godoc
// @Summary Delete a task by ID
// @Description Deletes a task with the specified ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID" Format(int64)
// @Success 200 {object} models.Task
// @Failure 400 {object} nil
// @Failure 404 {object} nil
// @Failure 500 {object} nil
// @Router /tasks/{id} [DELETE]
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

	if len(resp.Kvs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Perform the delete operation
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	_, err = cli.Delete(ctx, "tasks/"+taskID)
	cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
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
