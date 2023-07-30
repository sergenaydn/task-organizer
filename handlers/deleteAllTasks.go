package handlers

import (
	"context"
	"net/http"
	"task-organizer/models"

	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// DeleteAllTasks godoc
// @Summary Deletes All Tasks
// @Description Deletes All Data
// @Tags Tasks
// @Accept json
// @Produce json
// @Success 201 {object} models.Task
// @Failure 500 {object} nil
// @Router /tasks/ [delete]
func DeleteAllTasks(c *gin.Context) {

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

	// Get all keys with the "tasks/" prefix from the database
	resp, err := h.Client.Get(context.Background(), "tasks/", clientv3.WithPrefix())
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
		_, err := h.Client.Delete(context.Background(), string(kv.Key))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tasks"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "All tasks deleted"})
}
