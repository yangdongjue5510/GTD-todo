//go:generate mockgen -source=action_repository.go -destination=mock_action_repository.go -package=workflow

package workflow

import (
	"errors"
)

var (
	ErrActionNotFound = errors.New("action not found")
)

type ActionRepository interface {
	AddAction(action *Action) (*Action, error)
	GetActions() ([]*Action, error)
	GetActionByID(id int) (*Action, error)
	UpdateAction(action *Action) error
	DeleteAction(id int) error
}

type InmemoryActionRepository struct {
	actions  map[int]*Action
	sequence int
}

func NewInmemoryActionRepository() ActionRepository {
	return &InmemoryActionRepository{
		actions:  make(map[int]*Action),
		sequence: 0,
	}
}

func (r *InmemoryActionRepository) AddAction(action *Action) (*Action, error) {
	r.sequence++
	action.ID = r.sequence
	r.actions[action.ID] = action
	return action, nil
}

func (r *InmemoryActionRepository) GetActions() ([]*Action, error) {
	actionList := make([]*Action, 0, len(r.actions))
	for _, a := range r.actions {
		actionList = append(actionList, a)
	}
	return actionList, nil
}

func (r *InmemoryActionRepository) GetActionByID(id int) (*Action, error) {
	foundAction, exists := r.actions[id]
	if exists {
		return foundAction, nil
	}
	return nil, ErrActionNotFound
}

func (r *InmemoryActionRepository) UpdateAction(action *Action) error {
	_, exists := r.actions[action.ID]
	if !exists {
		return ErrActionNotFound
	}
	r.actions[action.ID] = action
	return nil
}

func (r *InmemoryActionRepository) DeleteAction(id int) error {
	_, exists := r.actions[id]
	if !exists {
		return ErrActionNotFound
	}
	delete(r.actions, id)
	return nil
}