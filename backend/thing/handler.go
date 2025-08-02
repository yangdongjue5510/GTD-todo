package thing

import "github.com/gin-gonic/gin"

type ThingHandler struct {
	service ThingService
}

func NewThingHandler(service ThingService) *ThingHandler {
	return &ThingHandler{service: service}
}

func SetupRoutes(r *gin.Engine, handler *ThingHandler) {
	thingsGroup := r.Group("/things")
	{
		thingsGroup.POST("/", handler.AddThing)
		thingsGroup.GET("/", handler.GetThings)
	}
}

func (h *ThingHandler) AddThing(c *gin.Context) {
	var thing Thing
	if err := c.ShouldBindJSON(&thing); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	h.service.AddThing(thing)
	c.JSON(201, thing)
}

func (h *ThingHandler) GetThings(c *gin.Context) {
	things := h.service.GetThings()
	c.JSON(200, things)
}
