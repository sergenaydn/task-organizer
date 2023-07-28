package routers

import (
	"task-organizer-copy/handlers"
	"task-organizer-copy/models"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// IdeaRouter sets up the routes and handlers for the "task" API endpoints.
// It takes a *gin.Engine as input to add the routes to.
func IdeaRouter(r *gin.Engine) {
	// Create a client for the first etcd member (using port 2379)
	client1, err := models.Connection2379()
	if err != nil {
		panic(err)
	}

	// Create a client for the second etcd member (using port 2380)
	client2, err := models.Connection2380()
	if err != nil {
		panic(err)
	}

	// Create a client for the third etcd member (using port 2381)
	client3, err := models.Connection2381()
	if err != nil {
		panic(err)
	}

	// Create a client for the fourth etcd member (using port 2382)
	client4, err := models.Connection2382()
	if err != nil {
		panic(err)
	}

	// Create a client for the fifth etcd member (using port 2383)
	client5, err := models.Connection2383()
	if err != nil {
		panic(err)
	}

	// Create a client for the sixth etcd member (using port 2384)
	client6, err := models.Connection2384()
	if err != nil {
		panic(err)
	}

	// Create handlers for each client
	handler1 := &models.Handler{Client: client1}
	handler2 := &models.Handler{Client: client2}
	handler3 := &models.Handler{Client: client3}
	handler4 := &models.Handler{Client: client4}
	handler5 := &models.Handler{Client: client5}
	handler6 := &models.Handler{Client: client6}

	// Setup the route for Swagger documentation.
	// This serves the Swagger UI to visualize and interact with the API documentation.
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Create a new route group for the "/tasks" endpoint.
	iR := r.Group("/tasks")

	// Define the individual API routes and their corresponding handler functions.
	// The handlers are from the "handlers" package, which contains the logic for each endpoint.
	iR.GET("", func(c *gin.Context) {
		c.Set("handler", handler1) // Set the handler for port 2379
		handlers.GetAllTasks(c)
	}) // Get a list of all tasks from port 2379

	iR.GET(":id", func(c *gin.Context) {
		c.Set("handler", handler2) // Set the handler for port 2380
		handlers.GetTask(c)
	}) // Get a task by its ID from port 2380

	iR.POST("", func(c *gin.Context) {
		c.Set("handler", handler3) // Set the handler for port 2381
		handlers.CreateTask(c)
	}) // Create a new task from port 2381

	iR.DELETE(":id", func(c *gin.Context) {
		c.Set("handler", handler4) // Set the handler for port 2382
		handlers.DeleteTask(c)
	}) // Delete a task by its ID from port 2382

	iR.DELETE("", func(c *gin.Context) {
		c.Set("handler", handler5) // Set the handler for port 2383
		handlers.DeleteAllTasks(c)
	}) // Delete all tasks from port 2383

	iR.PUT(":id", func(c *gin.Context) {
		c.Set("handler", handler6) // Set the handler for port 2384
		handlers.UpdateTask(c)
	}) // Update a task by its ID from port 2384
}
