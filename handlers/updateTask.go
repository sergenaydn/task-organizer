package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"task-organizer/models"

	"github.com/gin-gonic/gin"
)

// UpdateTask godoc
// @Summary Update a task by ID
// @Description Updates a task with the specified ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID" Format(int64)
// @Param task body models.UpdateReq true "Task object with fields to be updated"
// @Success 200 {object} models.UpdateReq
// @Failure 400 {object} nil
// @Failure 404 {object} nil
// @Failure 500 {object} nil
// @Router /tasks/{id} [put]
func UpdateTask(c *gin.Context) {
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
	// Get the task ID from the URL path
	taskID := c.Param("id")

	// Fetch the existing task from the database
	resp, err := h.Client.Get(context.Background(), "tasks/"+taskID)
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
	_, err = h.Client.Put(context.Background(), "tasks/"+taskID, string(taskJSON))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updateReq)
}
