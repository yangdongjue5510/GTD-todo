package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"yangdongju/gtd-todo/workflow"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func setupActionRouter(actionService workflow.ActionService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewActionHandler(actionService)
	SetupActionRoutes(r, handler)
	return r
}

func TestActionHandler_CreateAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockActionService := workflow.NewMockActionService(ctrl)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid action should be created successfully",
			requestBody: workflow.Action{
				Title:       "Test Action",
				Description: "Test Description",
				Status:      workflow.ToDo,
			},
			setupMock: func() {
				mockActionService.EXPECT().
					Save(gomock.Any()).
					Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"message":"Action created successfully"}`,
		},
		{
			name: "Invalid JSON should return bad request",
			requestBody: `{"invalid": json}`,
			setupMock: func() {
				// No mock call expected
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid JSON format"}`,
		},
		{
			name: "Service error should return bad request",
			requestBody: workflow.Action{
				Title: "", // Empty title will cause service error
			},
			setupMock: func() {
				mockActionService.EXPECT().
					Save(gomock.Any()).
					Return(fmt.Errorf("action title cannot be empty"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"action title cannot be empty"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			
			router := setupActionRouter(mockActionService)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req, _ := http.NewRequest("POST", "/actions/", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body %s, got %s", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestActionHandler_GetActions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockActionService := workflow.NewMockActionService(ctrl)

	expectedActions := []workflow.Action{
		{ID: 1, Title: "Action 1", Status: workflow.ToDo},
		{ID: 2, Title: "Action 2", Status: workflow.InProgress},
	}

	mockActionService.EXPECT().
		GetActions().
		Return(expectedActions)

	router := setupActionRouter(mockActionService)

	req, _ := http.NewRequest("GET", "/actions/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var responseActions []workflow.Action
	err := json.Unmarshal(w.Body.Bytes(), &responseActions)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(responseActions) != 2 {
		t.Errorf("Expected 2 actions, got %d", len(responseActions))
	}
	if responseActions[0].Title != "Action 1" {
		t.Errorf("Expected first action title 'Action 1', got %s", responseActions[0].Title)
	}
}

func TestActionHandler_UpdateAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockActionService := workflow.NewMockActionService(ctrl)

	tests := []struct {
		name           string
		actionID       string
		requestBody    workflow.Action
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "Valid update should succeed",
			actionID: "1",
			requestBody: workflow.Action{
				Title:       "Updated Action",
				Description: "Updated Description",
				Status:      workflow.InProgress,
			},
			setupMock: func() {
				mockActionService.EXPECT().
					UpdateAction(1, gomock.Any()).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Action updated successfully"}`,
		},
		{
			name:     "Invalid action ID should return bad request",
			actionID: "invalid",
			requestBody: workflow.Action{
				Title: "Updated Action",
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid action ID"}`,
		},
		{
			name:     "Action not found should return not found",
			actionID: "999",
			requestBody: workflow.Action{
				Title: "Updated Action",
			},
			setupMock: func() {
				mockActionService.EXPECT().
					UpdateAction(999, gomock.Any()).
					Return(workflow.ErrActionNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Action not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			
			router := setupActionRouter(mockActionService)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("PUT", "/actions/"+tt.actionID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			if w.Body.String() != tt.expectedBody {
				t.Errorf("Expected body %s, got %s", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestActionHandler_UpdateActionStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockActionService := workflow.NewMockActionService(ctrl)

	statusRequest := struct {
		Status workflow.Status `json:"status"`
	}{
		Status: workflow.InProgress,
	}

	mockActionService.EXPECT().
		UpdateActionStatus(1, workflow.InProgress).
		Return(nil)

	router := setupActionRouter(mockActionService)

	body, _ := json.Marshal(statusRequest)
	req, _ := http.NewRequest("PUT", "/actions/1/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	expectedBody := `{"message":"Action status updated successfully"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestActionHandler_DeleteAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockActionService := workflow.NewMockActionService(ctrl)

	mockActionService.EXPECT().
		DeleteAction(1).
		Return(nil)

	router := setupActionRouter(mockActionService)

	req, _ := http.NewRequest("DELETE", "/actions/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}