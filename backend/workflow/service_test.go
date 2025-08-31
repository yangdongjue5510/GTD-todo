package workflow

import (
	"testing"
	"time"
)

func TestInmemoryActionService_Save(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		action      Action
		wantError   bool
		errorMsg    string
	}{
		{
			name: "Valid action should be saved successfully",
			action: Action{
				Title:       "Test Action",
				Description: "Test Description",
				Status:      ToDo,
			},
			wantError: false,
		},
		{
			name: "Action with empty title should return error",
			action: Action{
				Title:       "",
				Description: "Test Description",
				Status:      ToDo,
			},
			wantError: true,
			errorMsg:  "action title cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			// given
			service := NewInmemoryActionService()
			
			// when
			err := service.Save(tt.action)
			
			// then
			if tt.wantError {
				if err == nil {
					t.Errorf("Save() expected error but got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Save() error = %v, expected %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Save() unexpected error = %v", err)
				}
				
				// Verify action was saved with correct ID
				actions := service.GetActions()
				if len(actions) != 1 {
					t.Errorf("Expected 1 action, got %d", len(actions))
				}
				if actions[0].ID != 1 {
					t.Errorf("Expected ID 1, got %d", actions[0].ID)
				}
			}
		})
	}
}

func TestInmemoryActionService_GetActions(t *testing.T) {
	t.Parallel()
	
	// given
	service := NewInmemoryActionService()
	action1 := Action{Title: "Action 1", Description: "Desc 1", Status: ToDo}
	action2 := Action{Title: "Action 2", Description: "Desc 2", Status: Completed}
	
	service.Save(action1)
	service.Save(action2)
	
	// when
	actions := service.GetActions()
	
	// then
	if len(actions) != 2 {
		t.Errorf("Expected 2 actions, got %d", len(actions))
	}
	
	// Verify immutability (should return copy, not reference)
	actions[0].Title = "Modified"
	originalActions := service.GetActions()
	if originalActions[0].Title == "Modified" {
		t.Error("GetActions() should return a copy, not reference")
	}
}

func TestInmemoryActionService_CreateActionFromClarified(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		clarifiedData ClarifiedData
		wantError    bool
		errorMsg     string
		expectedContext string
	}{
		{
			name: "Valid clarified data should create action successfully",
			clarifiedData: ClarifiedData{
				Title:       "Test Action",
				Description: "Test Description",
				Priority:    "normal",
				Context:     "work",
				SourceID:    1,
			},
			wantError: false,
			expectedContext: "work",
		},
		{
			name: "High priority should map to urgent context",
			clarifiedData: ClarifiedData{
				Title:       "Urgent Action",
				Description: "Urgent Description",
				Priority:    "high",
				Context:     "work",
				SourceID:    1,
			},
			wantError: false,
			expectedContext: "urgent",
		},
		{
			name: "Low priority should map to someday context",
			clarifiedData: ClarifiedData{
				Title:       "Low Priority Action",
				Description: "Low Description",
				Priority:    "low",
				Context:     "work",
				SourceID:    1,
			},
			wantError: false,
			expectedContext: "someday",
		},
		{
			name: "Empty title should return error",
			clarifiedData: ClarifiedData{
				Title:       "",
				Description: "Test Description",
				Priority:    "normal",
				Context:     "work",
				SourceID:    1,
			},
			wantError: true,
			errorMsg:  "clarified data title cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			// given
			service := NewInmemoryActionService()
			
			// when
			createdAction, err := service.CreateActionFromClarified(tt.clarifiedData)
			
			// then
			if tt.wantError {
				if err == nil {
					t.Errorf("CreateActionFromClarified() expected error but got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("CreateActionFromClarified() error = %v, expected %v", err.Error(), tt.errorMsg)
				}
				if createdAction != nil {
					t.Error("CreateActionFromClarified() should return nil on error")
				}
			} else {
				if err != nil {
					t.Errorf("CreateActionFromClarified() unexpected error = %v", err)
				}
				if createdAction == nil {
					t.Error("CreateActionFromClarified() should return created action on success")
				}
				
				// Verify created action properties
				if createdAction.Title != tt.clarifiedData.Title {
					t.Errorf("CreateActionFromClarified() Title = %v, expected %v", createdAction.Title, tt.clarifiedData.Title)
				}
				if createdAction.Description != tt.clarifiedData.Description {
					t.Errorf("CreateActionFromClarified() Description = %v, expected %v", createdAction.Description, tt.clarifiedData.Description)
				}
				if createdAction.Status != ToDo {
					t.Errorf("CreateActionFromClarified() Status = %v, expected %v", createdAction.Status, ToDo)
				}
				if createdAction.Context != tt.expectedContext {
					t.Errorf("CreateActionFromClarified() Context = %v, expected %v", createdAction.Context, tt.expectedContext)
				}
				if createdAction.ID == 0 {
					t.Error("CreateActionFromClarified() should assign ID to created action")
				}
				
				// Verify action was saved
				actions := service.GetActions()
				if len(actions) != 1 {
					t.Errorf("Expected 1 action in service, got %d", len(actions))
				}
			}
		})
	}
}

func TestInmemoryActionService_CreateActionFromClarified_WithDueDate(t *testing.T) {
	t.Parallel()
	
	// given
	service := NewInmemoryActionService()
	dueDate := time.Now().Add(24 * time.Hour)
	clarifiedData := ClarifiedData{
		Title:       "Action with Due Date",
		Description: "Description",
		Priority:    "normal",
		DueDate:     &dueDate,
		Context:     "work",
		SourceID:    1,
	}
	
	// when
	createdAction, err := service.CreateActionFromClarified(clarifiedData)
	
	// then
	if err != nil {
		t.Errorf("CreateActionFromClarified() unexpected error = %v", err)
	}
	if createdAction.DueDate == nil {
		t.Error("CreateActionFromClarified() should preserve due date")
	}
	if !createdAction.DueDate.Equal(dueDate) {
		t.Errorf("CreateActionFromClarified() DueDate = %v, expected %v", *createdAction.DueDate, dueDate)
	}
}