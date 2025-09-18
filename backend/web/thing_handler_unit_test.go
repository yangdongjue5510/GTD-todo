package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"yangdongju/gtd-todo/capture"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)


func TestNewThingHandler(t *testing.T) {
	t.Parallel()
	
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockThingService := capture.NewMockThingService(ctrl)
	
	// when
	handler := NewThingHandler(mockThingService)
	
	// then
	if handler == nil {
		t.Error("NewThingHandler() should return non-nil handler")
	}
	if handler.thingService != mockThingService {
		t.Error("NewThingHandler() should set thingService correctly")
	}
}

func TestSetupRoutes(t *testing.T) {
	t.Parallel()
	
	// given
	gin.SetMode(gin.TestMode)
	router := gin.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockThingService := capture.NewMockThingService(ctrl)
	handler := NewThingHandler(mockThingService)
	
	// when
	SetupRoutes(router, handler)
	
	// then
	// Verify routes are registered by checking route info
	routes := router.Routes()
	
	expectedRoutes := map[string]string{
		"POST /things/":          "POST",
		"GET /things/":           "GET", 
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
		setupMock      func(*capture.MockThingService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid thing should call service and return created thing",
			requestBody: capture.Thing{
				Title:       "Test Thing",
				Description: "Test Description",
				Status:      capture.Active,
			},
			setupMock: func(m *capture.MockThingService) {
				createdThing := &capture.Thing{
					ID:          1,
					Title:       "Test Thing",
					Description: "Test Description",
					Status:      capture.Active,
				}
				m.EXPECT().AddThing(gomock.Any()).Return(createdThing, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Service error should return bad request",
			requestBody: capture.Thing{
				Title:       "",
				Description: "Test Description",
				Status:      capture.Active,
			},
			setupMock: func(m *capture.MockThingService) {
				m.EXPECT().AddThing(gomock.Any()).Return((*capture.Thing)(nil), errors.New("thing title cannot be empty"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "thing title cannot be empty",
		},
		{
			name:           "Invalid JSON should return bad request",
			requestBody:    "invalid json",
			setupMock:      func(m *capture.MockThingService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			// given
			gin.SetMode(gin.TestMode)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockThingService := capture.NewMockThingService(ctrl)
			tt.setupMock(mockThingService)
			
			handler := NewThingHandler(mockThingService)
			
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
		})
	}
}

func TestThingHandler_GetThings_Unit(t *testing.T) {
	t.Parallel()
	
	// given
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockThingService := capture.NewMockThingService(ctrl)

	expectedThings := []*capture.Thing{
		{ID: 1, Title: "Thing 1", Description: "Desc 1", Status: capture.Active},
		{ID: 2, Title: "Thing 2", Description: "Desc 2", Status: capture.Done},
	}
	
	mockThingService.EXPECT().GetThings().Return(expectedThings, nil)
	
	handler := NewThingHandler(mockThingService)
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/things/", nil)
	
	// when
	handler.GetThings(c)
	
	// then
	if w.Code != http.StatusOK {
		t.Errorf("GetThings() status = %d, expected %d", w.Code, http.StatusOK)
	}

	var response []*capture.Thing
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("GetThings() failed to parse response: %v", err)
	}
	
	if len(response) != len(expectedThings) {
		t.Errorf("GetThings() returned %d things, expected %d", len(response), len(expectedThings))
	}
}

