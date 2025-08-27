package capture

import (
	"strconv"
	"yangdongju/gtd-todo/workflow"
	
	"github.com/gin-gonic/gin"
)

type ThingHandler struct {
	thingService  ThingService
	actionService workflow.ActionService
}

func NewThingHandler(thingService ThingService, actionService workflow.ActionService) *ThingHandler {
	return &ThingHandler{
		thingService:  thingService,
		actionService: actionService,
	}
}

func SetupRoutes(r *gin.Engine, handler *ThingHandler) {
	thingsGroup := r.Group("/things")
	{
		thingsGroup.POST("/", handler.AddThing)
		thingsGroup.GET("/", handler.GetThings)
		thingsGroup.POST("/:id/clarify", handler.ClarifyThing)
	}
}

func (h *ThingHandler) AddThing(c *gin.Context) {
	var thing Thing
	if err := c.ShouldBindJSON(&thing); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}
	
	if err := h.thingService.AddThing(thing); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(201, thing)
}

func (h *ThingHandler) GetThings(c *gin.Context) {
	things := h.thingService.GetThings()
	c.JSON(200, things)
}

func (h *ThingHandler) ClarifyThing(c *gin.Context) {
	// Get thing ID from URL parameter
	idParam := c.Param("id")
	thingID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid thing ID"})
		return
	}
	
	// Step 1: Capture Context - Clarify the thing
	clarifiedData, err := h.thingService.ClarifyThing(thingID)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	
	// Step 2: Workflow Context - Create Action from clarified data
	actionData := workflow.ClarifiedData{
		Title:       clarifiedData.Title,
		Description: clarifiedData.Description,
		Priority:    clarifiedData.Priority,
		DueDate:     clarifiedData.DueDate,
		Context:     clarifiedData.Context,
		SourceID:    clarifiedData.SourceID,
	}
	
	createdAction, err := h.actionService.CreateActionFromClarified(actionData)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create action: " + err.Error()})
		return
	}
	
	// Step 3: Mark original thing as processed
	if err := h.thingService.MarkThingAsProcessed(thingID); err != nil {
		c.JSON(500, gin.H{"error": "Failed to mark thing as processed: " + err.Error()})
		return
	}
	
	// Return the created action
	c.JSON(201, createdAction)
}
