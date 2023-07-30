package handlers

import (
	"context"
	"net/http"
	"task-organizer/models"
	"time"

	"github.com/gin-gonic/gin"
)

// DeleteTask godoc
// @Summary Delete a task by ID
// @Description Deletes a task with the specified ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID" Format(int64)
// @Success 200 {object} models.Task
// @Failure 400 {object} nil
// @Failure 404 {object} nil
// @Failure 500 {object} nil
// @Router /tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
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
	taskID := c.Param("id")

	// Check if the task exists before attempting to delete it
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := h.Client.Get(ctx, "tasks/"+taskID)
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
	_, err = h.Client.Delete(ctx, "tasks/"+taskID)
	cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
