package capture

import (
	"errors"
	"time"
)

// ClarifiedData represents the processed information from a Thing ready for Action creation
type ClarifiedData struct {
	Title       string
	Description string
	Priority    string
	DueDate     *time.Time
	Context     string
	SourceID    int // Original Thing ID
}

type ThingService interface {
	AddThing(thing Thing) (*Thing, error)
	GetThings() []Thing
	ClarifyThing(thingID int) (*ClarifiedData, error)
	MarkThingAsProcessed(thingID int) error
}

type InmemoryThingService struct {
	things   []Thing
	sequence int
}

func NewInmemoryThingService() ThingService {
	return &InmemoryThingService{
		things:   make([]Thing, 0),
		sequence: 0,
	}
}

func (s *InmemoryThingService) AddThing(thing Thing) (*Thing, error) {
	if thing.Title == "" {
		return nil, errors.New("thing title cannot be empty")
	}
	
	s.sequence++
	thing.ID = s.sequence
	s.things = append(s.things, thing)
	return &s.things[len(s.things)-1], nil
}

func (s *InmemoryThingService) GetThings() []Thing {
	copiedThings := make([]Thing, len(s.things))
	copy(copiedThings, s.things)
	return copiedThings
}

func (s *InmemoryThingService) ClarifyThing(thingID int) (*ClarifiedData, error) {
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
	
	// Process the thing and extract clarified data
	clarified := &ClarifiedData{
		Title:       targetThing.Title,
		Description: targetThing.Description,
		Priority:    "normal", // Default priority, can be enhanced with AI or user input
		DueDate:     nil,      // Could be extracted from description text
		Context:     "inbox",  // Default context
		SourceID:    thingID,
	}
	
	return clarified, nil
}

func (s *InmemoryThingService) MarkThingAsProcessed(thingID int) error {
	// Find and mark the thing as done
	for i, thing := range s.things {
		if thing.ID == thingID {
			s.things[i].Status = Done
			return nil
		}
	}
	
	return errors.New("thing not found")
}
