package main

import (
	"log"
	"net/http"
	"os"
	"task-organizer/docs"
	"task-organizer/routers"

	_ "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// @title Task Organizator
// @version 1.0
// @description Manage Your Tasks

// @host localhost:8080
// @BasePath /tasks

func main() {
	// Create a new Gin router with default middleware (logger, recovery).
	r := gin.Default()

	// Set up the API routes and handlers using the IdeaRouter function.
	routers.IdeaRouter(r)

	// Set the base path for Swagger documentation.
	docs.SwaggerInfo.BasePath = "/"

	// Start the HTTP server and listen on port 8080.
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Failed to start the server:", err)
	}

	// Exit the application gracefully with status code 0.
	os.Exit(0)
}
