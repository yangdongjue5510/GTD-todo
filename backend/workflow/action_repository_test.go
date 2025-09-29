package workflow

import (
	"testing"
	"time"
)

func TestInmemoryActionRepository_AddAction(t *testing.T) {
	repo := NewInmemoryActionRepository()

	action := &Action{
		Title:       "Test Action",
		Description: "Test Description",
		Status:      ToDo,
	}

	result, err := repo.AddAction(action)

	if err != nil {
		t.Errorf("AddAction should not return error, got: %v", err)
	}
	if result.ID == 0 {
		t.Error("AddAction should assign ID to action")
	}
	if result.Title != "Test Action" {
		t.Errorf("Expected title 'Test Action', got: %s", result.Title)
	}
}

func TestInmemoryActionRepository_AddAction_AssignsIncrementalID(t *testing.T) {
	repo := NewInmemoryActionRepository()

	action1 := &Action{Title: "Action 1", Status: ToDo}
	action2 := &Action{Title: "Action 2", Status: ToDo}

	result1, _ := repo.AddAction(action1)
	result2, _ := repo.AddAction(action2)

	if result1.ID != 1 {
		t.Errorf("Expected first action ID to be 1, got: %d", result1.ID)
	}
	if result2.ID != 2 {
		t.Errorf("Expected second action ID to be 2, got: %d", result2.ID)
	}
}

func TestInmemoryActionRepository_GetActions_EmptyRepository(t *testing.T) {
	repo := NewInmemoryActionRepository()

	actions, err := repo.GetActions()

	if err != nil {
		t.Errorf("GetActions should not return error, got: %v", err)
	}
	if len(actions) != 0 {
		t.Errorf("Expected empty list, got %d actions", len(actions))
	}
}

func TestInmemoryActionRepository_GetActions_MultipleActions(t *testing.T) {
	repo := NewInmemoryActionRepository()

	action1 := &Action{Title: "Action 1", Status: ToDo}
	action2 := &Action{Title: "Action 2", Status: InProgress}

	repo.AddAction(action1)
	repo.AddAction(action2)

	actions, err := repo.GetActions()

	if err != nil {
		t.Errorf("GetActions should not return error, got: %v", err)
	}
	if len(actions) != 2 {
		t.Errorf("Expected 2 actions, got: %d", len(actions))
	}
}

func TestInmemoryActionRepository_GetActionByID_ExistingAction(t *testing.T) {
	repo := NewInmemoryActionRepository()

	action := &Action{Title: "Test Action", Status: ToDo}
	addedAction, _ := repo.AddAction(action)

	foundAction, err := repo.GetActionByID(addedAction.ID)

	if err != nil {
		t.Errorf("GetActionByID should not return error for existing action, got: %v", err)
	}
	if foundAction.Title != "Test Action" {
		t.Errorf("Expected title 'Test Action', got: %s", foundAction.Title)
	}
	if foundAction.ID != addedAction.ID {
		t.Errorf("Expected ID %d, got: %d", addedAction.ID, foundAction.ID)
	}
}

func TestInmemoryActionRepository_GetActionByID_NonExistingAction(t *testing.T) {
	repo := NewInmemoryActionRepository()

	foundAction, err := repo.GetActionByID(999)

	if err == nil {
		t.Error("GetActionByID should return error for non-existing action")
	}
	if err != ErrActionNotFound {
		t.Errorf("Expected ErrActionNotFound, got: %v", err)
	}
	if foundAction != nil {
		t.Error("GetActionByID should return nil action for non-existing ID")
	}
}

func TestInmemoryActionRepository_UpdateAction_ExistingAction(t *testing.T) {
	repo := NewInmemoryActionRepository()

	action := &Action{Title: "Original Title", Status: ToDo}
	addedAction, _ := repo.AddAction(action)

	dueDate := time.Now().Add(24 * time.Hour)
	updatedAction := &Action{
		ID:          addedAction.ID,
		Title:       "Updated Title",
		Description: "Updated Description",
		Status:      InProgress,
		DueDate:     &dueDate,
		Context:     "work",
	}

	err := repo.UpdateAction(updatedAction)

	if err != nil {
		t.Errorf("UpdateAction should not return error for existing action, got: %v", err)
	}

	foundAction, _ := repo.GetActionByID(addedAction.ID)
	if foundAction.Title != "Updated Title" {
		t.Errorf("Expected updated title 'Updated Title', got: %s", foundAction.Title)
	}
	if foundAction.Status != InProgress {
		t.Errorf("Expected updated status InProgress, got: %v", foundAction.Status)
	}
}

func TestInmemoryActionRepository_UpdateAction_NonExistingAction(t *testing.T) {
	repo := NewInmemoryActionRepository()

	action := &Action{ID: 999, Title: "Non-existing Action", Status: ToDo}

	err := repo.UpdateAction(action)

	if err == nil {
		t.Error("UpdateAction should return error for non-existing action")
	}
	if err != ErrActionNotFound {
		t.Errorf("Expected ErrActionNotFound, got: %v", err)
	}
}

func TestInmemoryActionRepository_DeleteAction_ExistingAction(t *testing.T) {
	repo := NewInmemoryActionRepository()

	action := &Action{Title: "To be deleted", Status: ToDo}
	addedAction, _ := repo.AddAction(action)

	err := repo.DeleteAction(addedAction.ID)

	if err != nil {
		t.Errorf("DeleteAction should not return error for existing action, got: %v", err)
	}

	_, getErr := repo.GetActionByID(addedAction.ID)
	if getErr != ErrActionNotFound {
		t.Error("Action should be deleted and not found")
	}
}

func TestInmemoryActionRepository_DeleteAction_NonExistingAction(t *testing.T) {
	repo := NewInmemoryActionRepository()

	err := repo.DeleteAction(999)

	if err == nil {
		t.Error("DeleteAction should return error for non-existing action")
	}
	if err != ErrActionNotFound {
		t.Errorf("Expected ErrActionNotFound, got: %v", err)
	}
}