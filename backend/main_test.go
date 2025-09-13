package main

import (
	"testing"
	"yangdongju/gtd-todo/capture"
	"yangdongju/gtd-todo/web"
	"yangdongju/gtd-todo/workflow"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// Mock implementations for testing main functions
type MockThingService struct {
	mock.Mock
}

func (m *MockThingService) AddThing(thing capture.Thing) (*capture.Thing, error) {
	args := m.Called(thing)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*capture.Thing), args.Error(1)
}

func (m *MockThingService) GetThings() []capture.Thing {
	args := m.Called()
	return args.Get(0).([]capture.Thing)
}

func (m *MockThingService) ClarifyThing(thingID int) (*capture.ClarifiedData, error) {
	args := m.Called(thingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*capture.ClarifiedData), args.Error(1)
}

func (m *MockThingService) MarkThingAsProcessed(thingID int) error {
	args := m.Called(thingID)
	return args.Error(0)
}

type MockActionService struct {
	mock.Mock
}

func (m *MockActionService) Save(action workflow.Action) error {
	args := m.Called(action)
	return args.Error(0)
}

func (m *MockActionService) GetActions() []workflow.Action {
	args := m.Called()
	return args.Get(0).([]workflow.Action)
}

func (m *MockActionService) CreateActionFromClarified(data workflow.ClarifiedData) (*workflow.Action, error) {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*workflow.Action), args.Error(1)
}

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
		"POST /things/:id/clarify": "POST",
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
	actionService := workflow.NewInmemoryActionService()
	thingService := capture.NewInmemoryThingService()
	
	// when & then
	if actionService == nil {
		t.Error("NewInmemoryActionService() should return non-nil service")
	}
	if thingService == nil {
		t.Error("NewInmemoryThingService() should return non-nil service")
	}
	
	// Verify services can be used to create handler
	handler := web.NewThingHandler(thingService, actionService)
	if handler == nil {
		t.Error("NewThingHandler() should return non-nil handler")
	}
}