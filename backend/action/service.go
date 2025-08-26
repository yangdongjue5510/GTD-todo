package action

import "errors"

type ActionService interface {
	Save(action Action) error
	GetActions() []Action
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