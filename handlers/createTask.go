package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"task-organizer-copy/models"
	"time"

	"github.com/gin-gonic/gin"
)

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
	client, ok := c.Get("client")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch handler"})
		return
	}

	h, ok := client.(*models.Handler)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid handler type"})
		return
	}
	// Bind the JSON request body to the task model
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

	// Convert the task to JSON format
	data, err := json.Marshal(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Store the task in the database with the generated ID
	_, err = h.Client.Put(ctx, "tasks/"+strconv.Itoa(*task.ID), string(data))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}
