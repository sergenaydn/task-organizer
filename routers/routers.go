package routers

import (
	"task-organizer/handlers"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// IdeaRouter sets up the routes and handlers for the "task" API endpoints.
// It takes a *gin.Engine as input to add the routes to.
func IdeaRouter(r *gin.Engine) {
	// Setup the route for Swagger documentation.
	// This serves the Swagger UI to visualize and interact with the API documentation.
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Create a new route group for the "/tasks" endpoint.
	iR := r.Group("/tasks")

	// Define the individual API routes and their corresponding handler functions.
	// The handlers are from the "handlers" package, which contains the logic for each endpoint.
	iR.GET("", handlers.GetAllTasks)       // Get a list of all tasks
	iR.GET(":id", handlers.GetTask)        // Get a task by its ID
	iR.POST("", handlers.CreateTask)       // Create a new task
	iR.DELETE(":id", handlers.DeleteTask)  // Delete a task by its ID
	iR.DELETE("", handlers.DeleteAllTasks) // Delete all tasks
	iR.PUT(":id", handlers.UpdateTask)     // Update a task by its ID
}
