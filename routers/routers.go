package routers

import (
	"task-organizer/handlers"
	"task-organizer/models"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// IdeaRouter sets up the routes and handlers for the "task" API endpoints.
// It takes a *gin.Engine as input to add the routes to.
func IdeaRouter(r *gin.Engine) {
	// Initialize the handlers for the etcd clients
	handler1, handler2, err := models.Init()
	if err != nil {
		panic(err)
	}

	// Setup the route for Swagger documentation.
	// This serves the Swagger UI to visualize and interact with the API documentation.
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Create a new route group for the "/tasks" endpoint.
	iR := r.Group("/tasks")

	// Define the individual API routes and their corresponding handler functions.
	// The handlers are from the "handlers" package, which contains the logic for each endpoint.

	// Get a list of all tasks from port 2379
	iR.GET("", func(c *gin.Context) {
		c.Set("handler", handler1) // Set the handler for port 2379
		handlers.GetAllTasks(c)
	})

	// Get a task by its ID from port 2380
	iR.GET(":id", func(c *gin.Context) {
		c.Set("handler", handler2) // Set the handler for port 2380
		handlers.GetTask(c)
	})

	// Create a new task from port 2379
	iR.POST("", func(c *gin.Context) {
		c.Set("handler", handler1) // Set the handler for port 2379
		handlers.CreateTask(c)
	})

	// Delete a task by its ID from port 2380
	iR.DELETE(":id", func(c *gin.Context) {
		c.Set("handler", handler2) // Set the handler for port 2380
		handlers.DeleteTask(c)
	})

	// Delete all tasks from port 2379
	iR.DELETE("", func(c *gin.Context) {
		c.Set("handler", handler1) // Set the handler for port 2379
		handlers.DeleteAllTasks(c)
	})

	// Update a task by its ID from port 2380
	iR.PUT(":id", func(c *gin.Context) {
		c.Set("handler", handler2) // Set the handler for port 2380
		handlers.UpdateTask(c)
	})
}
