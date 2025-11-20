//go:generate mockgen -source=service.go -destination=mock_action_service.go -package=workflow

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
	GetActionByID(id int) (*Action, error)
	UpdateAction(id int, action Action) error
	UpdateActionStatus(id int, status Status) error
	DeleteAction(id int) error
	CreateActionFromClarified(data ClarifiedData) (*Action, error)
} 

type ActionServiceImpl struct {
	actionRepository ActionRepository
}

func NewActionService(repo ActionRepository) ActionService {
	return &ActionServiceImpl{
		actionRepository: repo,
	}
}

func (s *ActionServiceImpl) Save(action Action) error {
	if action.Title == "" {
		return errors.New("action title cannot be empty")
	}
	
	_, err := s.actionRepository.AddAction(&action)
	return err
}

func (s *ActionServiceImpl) GetActions() []Action {
	actions, err := s.actionRepository.GetActions()
	if err != nil {
		return []Action{}
	}
	
	// Convert []*Action to []Action
	result := make([]Action, len(actions))
	for i, action := range actions {
		result[i] = *action
	}
	return result
}

func (s *ActionServiceImpl) GetActionByID(id int) (*Action, error) {
	return s.actionRepository.GetActionByID(id)
}

func (s *ActionServiceImpl) UpdateAction(id int, action Action) error {
	if action.Title == "" {
		return errors.New("action title cannot be empty")
	}
	
	// Preserve the ID
	action.ID = id
	return s.actionRepository.UpdateAction(&action)
}

func (s *ActionServiceImpl) UpdateActionStatus(id int, status Status) error {
	existingAction, err := s.actionRepository.GetActionByID(id)
	if err != nil {
		return err
	}
	
	// Update only the status
	existingAction.Status = status
	return s.actionRepository.UpdateAction(existingAction)
}

func (s *ActionServiceImpl) DeleteAction(id int) error {
	return s.actionRepository.DeleteAction(id)
}

func (s *ActionServiceImpl) CreateActionFromClarified(data ClarifiedData) (*Action, error) {
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
	
	// Save the action through repository
	return s.actionRepository.AddAction(&action)
}