package routers

import (
	"task-organizer/handlers"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func IdeaRouter(r *gin.Engine) {
	// Setup the route for Swagger documentation.
	// This serves the Swagger UI to visualize and interact with the API documentation.
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Create a new route group for the "/idea" endpoint.
	iR := r.Group("/tasks")

	// Define the individual API routes and their corresponding handler functions.
	// The handlers are from the "handlers" package, which contains the logic for each endpoint.
	iR.GET("", handlers.GetAllTasks)       // Get a list of all ideas
	iR.GET(":id", handlers.GetTask)        // Get an idea by its ID
	iR.POST("", handlers.CreateTask)       // Create a new idea
	iR.DELETE(":id", handlers.DeleteTask)  // Delete an idea by its ID
	iR.DELETE("", handlers.DeleteAllTasks) // Delete an idea by its ID
	iR.PUT(":id", handlers.UpdateTask)     // Update an idea by its ID
}
