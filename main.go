package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"task/handlers"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Setup gin router
	r := gin.Default()

	// Rute users
	r.POST("/users", handlers.CreateUser)
	r.GET("/users", handlers.GetUsers)
	r.DELETE("/users/:id", handlers.DeleteUser)

	// Nested rute user -> task
	r.GET("/users/:id/tasks", handlers.GetUserTasks)
	r.POST("/users/:id/tasks", handlers.CreateTaskByUser)

	// Rute tasks
	r.POST("/tasks", handlers.CreateTask)
	r.GET("/tasks", handlers.GetTasks)
	r.PATCH("/tasks/:id", handlers.UpdateTask)
	r.DELETE("/tasks/:id", handlers.DeleteTask)

	// Run server
	r.Run(":8080")
}

