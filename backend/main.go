package main

import (
	"github.com/gin-gonic/gin"
	"yangdongju/gtd-todo/thing"
	
)
func main() {
	r := gin.Default()
	thingRoutes(r) // Setup the thing routes
	r.Run(":8080") // Start the server on port 8080
}

func thingRoutes(r *gin.Engine) {
	service := &thing.InmemoryThingService{}
	handler := thing.NewThingHandler(service)
	
	thing.SetupRoutes(r, handler)
}