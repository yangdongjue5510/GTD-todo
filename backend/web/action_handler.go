package web

import (
	"net/http"
	"strconv"
	"yangdongju/gtd-todo/workflow"

	"github.com/gin-gonic/gin"
)

type ActionHandler struct {
	actionService workflow.ActionService
}

func NewActionHandler(actionService workflow.ActionService) *ActionHandler {
	return &ActionHandler{
		actionService: actionService,
	}
}

func SetupActionRoutes(r *gin.Engine, handler *ActionHandler) {
	actionsGroup := r.Group("/actions")
	{
		actionsGroup.POST("/", handler.CreateAction)
		actionsGroup.GET("/", handler.GetActions)
		actionsGroup.PUT("/:id", handler.UpdateAction)
		actionsGroup.PUT("/:id/status", handler.UpdateActionStatus)
		actionsGroup.DELETE("/:id", handler.DeleteAction)
	}
}

// CreateAction handles POST /actions
func (h *ActionHandler) CreateAction(c *gin.Context) {
	var action workflow.Action
	if err := c.ShouldBindJSON(&action); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	err := h.actionService.Save(action)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Action created successfully"})
}

// GetActions handles GET /actions
func (h *ActionHandler) GetActions(c *gin.Context) {
	actions := h.actionService.GetActions()
	c.JSON(http.StatusOK, actions)
}

// UpdateAction handles PUT /actions/:id
func (h *ActionHandler) UpdateAction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action ID"})
		return
	}

	var action workflow.Action
	if err := c.ShouldBindJSON(&action); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	err = h.actionService.UpdateAction(id, action)
	if err != nil {
		if err == workflow.ErrActionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Action not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Action updated successfully"})
}

// UpdateActionStatus handles PUT /actions/:id/status
func (h *ActionHandler) UpdateActionStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action ID"})
		return
	}

	var statusRequest struct {
		Status workflow.Status `json:"status"`
	}
	if err := c.ShouldBindJSON(&statusRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	err = h.actionService.UpdateActionStatus(id, statusRequest.Status)
	if err != nil {
		if err == workflow.ErrActionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Action not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Action status updated successfully"})
}

// DeleteAction handles DELETE /actions/:id
func (h *ActionHandler) DeleteAction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action ID"})
		return
	}

	err = h.actionService.DeleteAction(id)
	if err != nil {
		if err == workflow.ErrActionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Action not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}