package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"task-organizer-copy/models"

	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
)

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
	var tasks []models.Task
	for _, kv := range resp.Kvs {
		var task models.Task
		if err := json.Unmarshal([]byte(kv.Value), &task); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse taks"})
			return
		}
		tasks = append(tasks, task)
	}

	c.JSON(http.StatusOK, tasks)
}
