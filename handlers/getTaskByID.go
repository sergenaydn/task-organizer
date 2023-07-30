package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"task-organizer/models"

	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/gin-swagger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// GetTask godoc
// @Summary Get a task by ID
// @Description Retrieves a task with the specified ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID" Format(int64)
// @Success 200 {object} models.Task
// @Failure 404 {object} nil
// @Failure 500 {object} nil
// @Router /tasks/{id} [get]
func GetTask(c *gin.Context) {
	// Fetch the task ID from the URL path parameter
	var task models.Task
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

	resp, err := h.Client.Get(context.Background(), "tasks/", clientv3.WithPrefix())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
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
