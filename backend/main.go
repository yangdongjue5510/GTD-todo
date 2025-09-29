package main

import (
	"time"
	"yangdongju/gtd-todo/capture"
	"yangdongju/gtd-todo/web"
	"yangdongju/gtd-todo/workflow"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	setupRoutes(r) // Setup all routes
	r.Run(":8080") // Start the server on port 8080
}

func setupRoutes(r *gin.Engine) {
	// Create independent domain services
	thingRepository := capture.NewInmemoryThingRepository()
	thingService := capture.NewThingService(thingRepository)
	
	actionRepository := workflow.NewInmemoryActionRepository()
	actionService := workflow.NewActionService(actionRepository)
	
	// Create handlers
	thingHandler := web.NewThingHandler(thingService)
	actionHandler := web.NewActionHandler(actionService)
	
	// Setup routes
	web.SetupRoutes(r, thingHandler)
	web.SetupActionRoutes(r, actionHandler)
}