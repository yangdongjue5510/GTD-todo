package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"yangdongju/gtd-todo/workflow"
	"yangdongju/gtd-todo/capture"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// Mock implementations for unit testing handlers
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

func TestNewThingHandler(t *testing.T) {
	t.Parallel()
	
	// given
	mockThingService := &MockThingService{}
	mockActionService := &MockActionService{}
	
	// when
	handler := NewThingHandler(mockThingService, mockActionService)
	
	// then
	if handler == nil {
		t.Error("NewThingHandler() should return non-nil handler")
	}
	if handler.thingService != mockThingService {
		t.Error("NewThingHandler() should set thingService correctly")
	}
	if handler.actionService != mockActionService {
		t.Error("NewThingHandler() should set actionService correctly")
	}
}

func TestSetupRoutes(t *testing.T) {
	t.Parallel()
	
	// given
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockThingService := &MockThingService{}
	mockActionService := &MockActionService{}
	handler := NewThingHandler(mockThingService, mockActionService)
	
	// when
	SetupRoutes(router, handler)
	
	// then
	// Verify routes are registered by checking route info
	routes := router.Routes()
	
	expectedRoutes := map[string]string{
		"POST /things/":          "POST",
		"GET /things/":           "GET", 
		"POST /things/:id/clarify": "POST",
	}
	
	actualRoutes := make(map[string]string)
	for _, route := range routes {
		actualRoutes[route.Method+" "+route.Path] = route.Method
	}
	
	for expectedRoute, expectedMethod := range expectedRoutes {
		if actualMethod, exists := actualRoutes[expectedRoute]; !exists {
			t.Errorf("SetupRoutes() missing route: %s", expectedRoute)
		} else if actualMethod != expectedMethod {
			t.Errorf("SetupRoutes() wrong method for %s: got %s, expected %s", expectedRoute, actualMethod, expectedMethod)
		}
	}
}

func TestThingHandler_AddThing_Unit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockThingService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid thing should call service and return created thing",
			requestBody: capture.Thing{
				Title:       "Test Thing",
				Description: "Test Description",
				Status:      capture.Pending,
			},
			setupMock: func(m *MockThingService) {
				createdThing := &capture.Thing{
					ID:          1,
					Title:       "Test Thing",
					Description: "Test Description",
					Status:      capture.Pending,
				}
				m.On("AddThing", mock.AnythingOfType("capture.Thing")).Return(createdThing, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Service error should return bad request",
			requestBody: capture.Thing{
				Title:       "",
				Description: "Test Description",
				Status:      capture.Pending,
			},
			setupMock: func(m *MockThingService) {
				m.On("AddThing", mock.AnythingOfType("capture.Thing")).Return((*capture.Thing)(nil), errors.New("thing title cannot be empty"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "thing title cannot be empty",
		},
		{
			name:           "Invalid JSON should return bad request",
			requestBody:    "invalid json",
			setupMock:      func(m *MockThingService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			// given
			gin.SetMode(gin.TestMode)
			mockThingService := &MockThingService{}
			mockActionService := &MockActionService{}
			tt.setupMock(mockThingService)
			
			handler := NewThingHandler(mockThingService, mockActionService)
			
			var reqBody bytes.Buffer
			if str, ok := tt.requestBody.(string); ok {
				reqBody.WriteString(str)
			} else {
				json.NewEncoder(&reqBody).Encode(tt.requestBody)
			}
			
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "/things/", &reqBody)
			c.Request.Header.Set("Content-Type", "application/json")
			
			// when
			handler.AddThing(c)
			
			// then
			if w.Code != tt.expectedStatus {
				t.Errorf("AddThing() status = %d, expected %d", w.Code, tt.expectedStatus)
			}
			
			if tt.expectedError != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				if response["error"] != tt.expectedError {
					t.Errorf("AddThing() error = %v, expected %v", response["error"], tt.expectedError)
				}
			}
			
			if tt.expectedStatus == http.StatusCreated {
				var response capture.Thing
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("AddThing() failed to parse response: %v", err)
				}
				if response.ID == 0 {
					t.Error("AddThing() should return thing with assigned ID")
				}
			}
			
			mockThingService.AssertExpectations(t)
		})
	}
}

func TestThingHandler_GetThings_Unit(t *testing.T) {
	t.Parallel()
	
	// given
	gin.SetMode(gin.TestMode)
	mockThingService := &MockThingService{}
	mockActionService := &MockActionService{}

	expectedThings := []capture.Thing{
		{ID: 1, Title: "Thing 1", Description: "Desc 1", Status: capture.Pending},
		{ID: 2, Title: "Thing 2", Description: "Desc 2", Status: capture.Done},
	}
	
	mockThingService.On("GetThings").Return(expectedThings)
	
	handler := NewThingHandler(mockThingService, mockActionService)
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/things/", nil)
	
	// when
	handler.GetThings(c)
	
	// then
	if w.Code != http.StatusOK {
		t.Errorf("GetThings() status = %d, expected %d", w.Code, http.StatusOK)
	}

	var response []capture.Thing
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("GetThings() failed to parse response: %v", err)
	}
	
	if len(response) != len(expectedThings) {
		t.Errorf("GetThings() returned %d things, expected %d", len(response), len(expectedThings))
	}
	
	mockThingService.AssertExpectations(t)
}

func TestThingHandler_ClarifyThing_Unit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		thingID         string
		setupThingMock  func(*MockThingService)
		setupActionMock func(*MockActionService)
		expectedStatus  int
		expectedError   string
	}{
		{
			name:    "Valid thing ID should clarify and create action",
			thingID: "1",
			setupThingMock: func(m *MockThingService) {
				clarifiedData := &capture.ClarifiedData{
					Title:       "Test Thing",
					Description: "Test Description",
					Priority:    "normal",
					Context:     "inbox",
					SourceID:    1,
				}
				m.On("ClarifyThing", 1).Return(clarifiedData, nil)
				m.On("MarkThingAsProcessed", 1).Return(nil)
			},
			setupActionMock: func(m *MockActionService) {
				createdAction := &workflow.Action{
					ID:          1,
					Title:       "Test Thing",
					Description: "Test Description",
					Status:      workflow.ToDo,
				}
				m.On("CreateActionFromClarified", mock.AnythingOfType("workflow.ClarifiedData")).Return(createdAction, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:    "Invalid thing ID format should return bad request",
			thingID: "invalid",
			setupThingMock: func(m *MockThingService) {},
			setupActionMock: func(m *MockActionService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid thing ID",
		},
		{
			name:    "Non-existent thing should return not found",
			thingID: "999",
			setupThingMock: func(m *MockThingService) {
				m.On("ClarifyThing", 999).Return((*capture.ClarifiedData)(nil), errors.New("thing not found"))
			},
			setupActionMock: func(m *MockActionService) {},
			expectedStatus: http.StatusNotFound,
			expectedError:  "thing not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			// given
			gin.SetMode(gin.TestMode)
			mockThingService := &MockThingService{}
			mockActionService := &MockActionService{}
			
			tt.setupThingMock(mockThingService)
			tt.setupActionMock(mockActionService)
			
			handler := NewThingHandler(mockThingService, mockActionService)
			
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "/things/"+tt.thingID+"/clarify", nil)
			c.Params = gin.Params{{Key: "id", Value: tt.thingID}}
			
			// when
			handler.ClarifyThing(c)
			
			// then
			if w.Code != tt.expectedStatus {
				t.Errorf("ClarifyThing() status = %d, expected %d", w.Code, tt.expectedStatus)
			}
			
			if tt.expectedError != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				if response["error"] != tt.expectedError {
					t.Errorf("ClarifyThing() error = %v, expected %v", response["error"], tt.expectedError)
				}
			}
			
			if tt.expectedStatus == http.StatusCreated {
				var response workflow.Action
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("ClarifyThing() failed to parse action response: %v", err)
				}
				if response.ID == 0 {
					t.Error("ClarifyThing() should return action with assigned ID")
				}
			}
			
			mockThingService.AssertExpectations(t)
			mockActionService.AssertExpectations(t)
		})
	}
}