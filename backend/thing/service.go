package thing

import (
	"errors"
	"yangdongju/gtd-todo/action"
)

type ThingService interface {
	AddThing(thing Thing) error
	GetThings() []Thing
	Clarify(thingID int) (*action.Action, error)
}

type InmemoryThingService struct {
	things        []Thing
	sequence      int
	actionService action.ActionService
}

func NewInmemoryThingService(actionService action.ActionService) ThingService {
	return &InmemoryThingService{
		things:        make([]Thing, 0),
		sequence:      0,
		actionService: actionService,
	}
}

func (s *InmemoryThingService) AddThing(thing Thing) error {
	if thing.Title == "" {
		return errors.New("thing title cannot be empty")
	}
	
	s.sequence++
	thing.ID = s.sequence
	s.things = append(s.things, thing)
	return nil
}

func (s *InmemoryThingService) GetThings() []Thing {
	copiedThings := make([]Thing, len(s.things))
	copy(copiedThings, s.things)
	return copiedThings
}

func (s *InmemoryThingService) Clarify(thingID int) (*action.Action, error) {
	// Find the thing by ID
	var targetThing *Thing
	for i, thing := range s.things {
		if thing.ID == thingID {
			targetThing = &s.things[i]
			break
		}
	}
	
	if targetThing == nil {
		return nil, errors.New("thing not found")
	}
	
	// Create action from thing
	newAction := action.Action{
		Title:       targetThing.Title,
		Description: targetThing.Description,
		Status:      action.ToDo,
	}
	
	// Save the action
	if err := s.actionService.Save(newAction); err != nil {
		return nil, err
	}
	
	// Mark thing as done
	targetThing.Status = Done
	
	return &newAction, nil
}
