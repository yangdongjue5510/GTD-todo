package web

import (
	"yangdongju/gtd-todo/capture"
	"github.com/gin-gonic/gin"
)

type ThingHandler struct {
	thingService  capture.ThingService
}

func NewThingHandler(thingService capture.ThingService) *ThingHandler {
	return &ThingHandler{
		thingService:  thingService,
	}
}

func SetupRoutes(r *gin.Engine, handler *ThingHandler) {
	thingsGroup := r.Group("/things")
	{
		thingsGroup.POST("/", handler.AddThing)
		thingsGroup.GET("/", handler.GetThings)
		
	}
}

func (h *ThingHandler) AddThing(c *gin.Context) {
	var thing capture.Thing
	if err := c.ShouldBindJSON(&thing); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}
	
	createdThing, err := h.thingService.AddThing(&thing)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(201, createdThing)
}

func (h *ThingHandler) GetThings(c *gin.Context) {
	things, err := h.thingService.GetThings()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, things)
}