package main

import (
	"testing"
	"yangdongju/gtd-todo/capture"
	"yangdongju/gtd-todo/web"
	"yangdongju/gtd-todo/workflow"

	"github.com/gin-gonic/gin"
)

func TestSetupRoutes(t *testing.T) {
	t.Parallel()
	
	// given
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// when
	setupRoutes(router)
	
	// then
	// Verify routes are registered by checking route info
	routes := router.Routes()
	
	expectedRoutes := map[string]string{
		"POST /things/":           "POST",
		"GET /things/":            "GET",
		"POST /actions/":          "POST",
		"GET /actions/":           "GET",
		"PUT /actions/:id":        "PUT",
		"PUT /actions/:id/status": "PUT",
		"DELETE /actions/:id":     "DELETE",
	}
	
	actualRoutes := make(map[string]string)
	for _, route := range routes {
		actualRoutes[route.Method+" "+route.Path] = route.Method
	}
	
	for expectedRoute, expectedMethod := range expectedRoutes {
		if actualMethod, exists := actualRoutes[expectedRoute]; !exists {
			t.Errorf("setupRoutes() missing route: %s", expectedRoute)
		} else if actualMethod != expectedMethod {
			t.Errorf("setupRoutes() wrong method for %s: got %s, expected %s", expectedRoute, actualMethod, expectedMethod)
		}
	}
}

func TestSetupRoutes_ServiceCreation(t *testing.T) {
	t.Parallel()
	
	// given
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// when
	setupRoutes(router)
	
	// then
	// Verify that function completes without panic
	// This tests the service creation and handler setup logic
	if len(router.Routes()) < 7 {
		t.Errorf("setupRoutes() should register at least 7 routes, got %d", len(router.Routes()))
	}
}

func TestMain_ServicesInitialization(t *testing.T) {
	t.Parallel()
	
	// This test verifies that services can be created without errors
	// given
	thingRepository := capture.NewInmemoryThingRepository()
	thingService := capture.NewThingService(thingRepository)
	
	actionRepository := workflow.NewInmemoryActionRepository()
	actionService := workflow.NewActionService(actionRepository)
	
	// when & then
	if thingRepository == nil {
		t.Error("NewInmemoryThingRepository() should return non-nil repository")
	}
	if thingService == nil {
		t.Error("NewThingService() should return non-nil service")
	}
	if actionRepository == nil {
		t.Error("NewInmemoryActionRepository() should return non-nil repository")
	}
	if actionService == nil {
		t.Error("NewActionService() should return non-nil service")
	}
	
	// Verify services can be used to create handlers
	thingHandler := web.NewThingHandler(thingService)
	if thingHandler == nil {
		t.Error("NewThingHandler() should return non-nil handler")
	}
	
	actionHandler := web.NewActionHandler(actionService)
	if actionHandler == nil {
		t.Error("NewActionHandler() should return non-nil handler")
	}
}