//go:build acceptance

package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"yangdongju/gtd-todo/workflow"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	
	// Create services
	thingService := NewInmemoryThingService()
	actionService := workflow.NewInmemoryActionService()
	
	// Create handler
	handler := NewThingHandler(thingService, actionService)
	
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
			requestBody: Thing{
				Title:       "Test Thing",
				Description: "Test Description",
				Status:      Pending,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Thing with empty title should return bad request",
			requestBody: Thing{
				Title:       "",
				Description: "Test Description",
				Status:      Pending,
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
				var response Thing
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
	addThingReq := Thing{Title: "Test Thing 1", Description: "Desc 1", Status: Pending}
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
	
	var response []Thing
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("GetThings() failed to parse response: %v", err)
	}
	
	if len(response) == 0 {
		t.Error("GetThings() should return at least 1 thing")
	}
}

func TestThingHandler_ClarifyThing_API(t *testing.T) {
	tests := []struct {
		name           string
		setupThing     *Thing
		thingID        string
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid thing ID should create action and mark thing as processed",
			setupThing: &Thing{
				Title:       "Test Thing",
				Description: "Test Description",
				Status:      Pending,
			},
			thingID:        "1",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid thing ID format should return bad request",
			thingID:        "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid thing ID",
		},
		{
			name:           "Non-existent thing ID should return not found",
			thingID:        "999",
			expectedStatus: http.StatusNotFound,
			expectedError:  "thing not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			router := setupTestRouter()
			
			if tt.setupThing != nil {
				// Add test thing first
				reqBody, _ := json.Marshal(*tt.setupThing)
				postReq, _ := http.NewRequest("POST", "/things/", bytes.NewBuffer(reqBody))
				postReq.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(httptest.NewRecorder(), postReq)
			}
			
			// when
			req, _ := http.NewRequest("POST", "/things/"+tt.thingID+"/clarify", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
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
				// Verify action was created
				var response workflow.Action
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("ClarifyThing() failed to parse action response: %v", err)
				}
				if response.ID == 0 {
					t.Error("ClarifyThing() should return action with assigned ID")
				}
				if response.Title != tt.setupThing.Title {
					t.Errorf("ClarifyThing() action title = %v, expected %v", response.Title, tt.setupThing.Title)
				}
			}
		})
	}
}

func TestThingHandler_ClarifyThing_Integration_API(t *testing.T) {
	// given
	router := setupTestRouter()
	
	// Step 1: Add a thing
	thingReq := Thing{Title: "Integration Test Thing", Description: "Integration Description", Status: Pending}
	reqBody, _ := json.Marshal(thingReq)
	
	postReq, _ := http.NewRequest("POST", "/things/", bytes.NewBuffer(reqBody))
	postReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), postReq)
	
	// Step 2: Clarify the thing
	clarifyReq, _ := http.NewRequest("POST", "/things/1/clarify", nil)
	clarifyW := httptest.NewRecorder()
	router.ServeHTTP(clarifyW, clarifyReq)
	
	// then
	if clarifyW.Code != http.StatusCreated {
		t.Errorf("ClarifyThing() status = %d, expected %d", clarifyW.Code, http.StatusCreated)
	}
	
	// Step 3: Verify thing was marked as processed
	getReq, _ := http.NewRequest("GET", "/things/", nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)
	
	var things []Thing
	json.Unmarshal(getW.Body.Bytes(), &things)
	
	if len(things) == 0 {
		t.Error("Integration test: should have at least 1 thing")
		return
	}
	
	if things[0].Status != Done {
		t.Errorf("Integration test: thing status = %v, expected Done", things[0].Status)
	}
}