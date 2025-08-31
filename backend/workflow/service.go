package workflow

import (
	"errors"
	"time"
)

// ClarifiedData represents processed data from Capture Context
type ClarifiedData struct {
	Title       string
	Description string
	Priority    string
	DueDate     *time.Time
	Context     string
	SourceID    int
}

type ActionService interface {
	Save(action Action) error
	GetActions() []Action
	CreateActionFromClarified(data ClarifiedData) (*Action, error)
} 

type InmemoryActionService struct {
	actions  []Action
	sequence int
}

func NewInmemoryActionService() ActionService {
	return &InmemoryActionService{
		actions:  make([]Action, 0),
		sequence: 0,
	}
}

func (s *InmemoryActionService) Save(action Action) error {
	if action.Title == "" {
		return errors.New("action title cannot be empty")
	}
	
	s.sequence++
	action.ID = s.sequence
	s.actions = append(s.actions, action)
	return nil
}

func (s *InmemoryActionService) GetActions() []Action {
	copiedActions := make([]Action, len(s.actions))
	copy(copiedActions, s.actions)
	return copiedActions
}

func (s *InmemoryActionService) CreateActionFromClarified(data ClarifiedData) (*Action, error) {
	if data.Title == "" {
		return nil, errors.New("clarified data title cannot be empty")
	}
	
	// Apply Workflow Context business rules for Action creation
	action := Action{
		Title:       data.Title,
		Description: data.Description,
		Status:      ToDo, // Default status in workflow
		DueDate:     data.DueDate,
	}
	
	// Apply priority mapping (Workflow Context logic)
	switch data.Priority {
	case "high":
		action.Context = "urgent"
	case "low":
		action.Context = "someday"
	default:
		action.Context = data.Context
	}
	
	// Save the action
	if err := s.Save(action); err != nil {
		return nil, err
	}
	
	// Return the created action with assigned ID
	return &s.actions[len(s.actions)-1], nil
}