package main

import (
	"testing"
	"yangdongju/gtd-todo/capture"
	"yangdongju/gtd-todo/web"

	"github.com/gin-gonic/gin"
)

func TestThingRoutes(t *testing.T) {
	t.Parallel()
	
	// given
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// when
	thingRoutes(router)
	
	// then
	// Verify routes are registered by checking route info
	routes := router.Routes()
	
	expectedRoutes := map[string]string{
		"POST /things/":           "POST",
		"GET /things/":            "GET",
	}
	
	actualRoutes := make(map[string]string)
	for _, route := range routes {
		actualRoutes[route.Method+" "+route.Path] = route.Method
	}
	
	for expectedRoute, expectedMethod := range expectedRoutes {
		if actualMethod, exists := actualRoutes[expectedRoute]; !exists {
			t.Errorf("thingRoutes() missing route: %s", expectedRoute)
		} else if actualMethod != expectedMethod {
			t.Errorf("thingRoutes() wrong method for %s: got %s, expected %s", expectedRoute, actualMethod, expectedMethod)
		}
	}
}

func TestThingRoutes_ServiceCreation(t *testing.T) {
	t.Parallel()
	
	// given
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// when
	thingRoutes(router)
	
	// then
	// Verify that function completes without panic
	// This tests the service creation and handler setup logic
	if len(router.Routes()) == 0 {
		t.Error("thingRoutes() should register at least one route")
	}
}

func TestMain_ServicesInitialization(t *testing.T) {
	t.Parallel()
	
	// This test verifies that services can be created without errors
	// given
	thingRepository := capture.NewInmemoryThingRepository()
	thingService := capture.NewThingService(thingRepository)
	
	// when & then
	if thingRepository == nil {
		t.Error("NewInmemoryThingRepository() should return non-nil repository")
	}
	if thingService == nil {
		t.Error("NewThingService() should return non-nil service")
	}
	
	// Verify services can be used to create handler
	handler := web.NewThingHandler(thingService)
	if handler == nil {
		t.Error("NewThingHandler() should return non-nil handler")
	}
}