package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
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

	resp, err := h.Client.Get(context.Background(), "tasks/", clientv3.WithPrefix())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	var tasks []models.Task
	for _, KV := range resp.Kvs {
		var task models.Task
		if err := json.Unmarshal(KV.Value, &task); err != nil {
			// Print the value that failed to parse to help diagnose the issue
			fmt.Printf("Failed to parse task: %s\n", string(KV.Value))
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
