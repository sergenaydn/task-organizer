package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"task-organizer/models"
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
// @Failure 500 {object} nil
// @Failure 400 {object} nil
// @Router /tasks [post]
func CreateTask(c *gin.Context) {
	client, ok := c.Get("handler")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch handler"})
		return
	}

	h, ok := client.(*models.Handler)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid handler type"})
		return
	}

	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the request contains an ID
	if task.ID != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Manual ID entry is not allowed"}) // Reject request with manual ID entry
		return
	}

	// Generate a unique ID using the GenerateUniqueID function
	task.ID = models.GenerateUniqueID()

	// Check if the Title is empty
	if task.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title cannot be empty"}) // Reject request with empty title
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
	_, err = h.Client.Put(ctx, "tasks/"+task.ID, string(data))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}
