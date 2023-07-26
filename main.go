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
	r := gin.Default()
	routers.IdeaRouter(r)
	docs.SwaggerInfo.BasePath = "/"

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
	os.Exit(0)
}
