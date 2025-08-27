package main

import (
	"github.com/gin-gonic/gin"
	"yangdongju/gtd-todo/capture"
	"yangdongju/gtd-todo/workflow"
	"github.com/gin-contrib/cors"
	"time"
)
func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins for simplicity))
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}))
	thingRoutes(r) // Setup the thing routes
	r.Run(":8080") // Start the server on port 8080
}

func thingRoutes(r *gin.Engine) {
	// Create independent domain services
	actionService := workflow.NewInmemoryActionService()
	thingService := capture.NewInmemoryThingService()
	
	// Create handlers
	handler := capture.NewThingHandler(thingService, actionService)
	
	capture.SetupRoutes(r, handler)
}