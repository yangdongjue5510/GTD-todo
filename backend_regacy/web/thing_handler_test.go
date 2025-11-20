//go:build acceptance

package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"yangdongju/gtd-todo/capture"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	
	// Create services
	thingRepository := capture.NewInmemoryThingRepository()
	thingService := capture.NewThingService(thingRepository)
	
	// Create handler
	handler := NewThingHandler(thingService)
	
	// Setup router
	router := gin.New()
	SetupRoutes(router, handler)
	
	return router
}

func TestThingHandler_AddThing_API(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid thing should be created successfully",
			requestBody: capture.Thing{
				Title:       "Test Thing",
				Description: "Test Description",
				Status:      capture.Active,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Thing with empty title should return bad request",
			requestBody: capture.Thing{
				Title:       "",
				Description: "Test Description",
				Status:      capture.Active,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "thing title cannot be empty",
		},
		{
			name:           "Invalid JSON should return bad request",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			router := setupTestRouter()
			
			var reqBody bytes.Buffer
			if str, ok := tt.requestBody.(string); ok {
				reqBody.WriteString(str)
			} else {
				json.NewEncoder(&reqBody).Encode(tt.requestBody)
			}
			
			req, _ := http.NewRequest("POST", "/things/", &reqBody)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			// when
			router.ServeHTTP(w, req)
			
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

func TestThingHandler_GetThings_API(t *testing.T) {
	// given
	router := setupTestRouter()
	
	// Add some test things
	addThingReq := capture.Thing{Title: "Test Thing 1", Description: "Desc 1", Status: capture.Active}
	reqBody, _ := json.Marshal(addThingReq)
	
	postReq, _ := http.NewRequest("POST", "/things/", bytes.NewBuffer(reqBody))
	postReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), postReq)
	
	// when
	req, _ := http.NewRequest("GET", "/things/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// then
	if w.Code != http.StatusOK {
		t.Errorf("GetThings() status = %d, expected %d", w.Code, http.StatusOK)
	}
	
	var response []*capture.Thing
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("GetThings() failed to parse response: %v", err)
	}
	
	if len(response) == 0 {
		t.Error("GetThings() should return at least 1 thing")
	}
}
