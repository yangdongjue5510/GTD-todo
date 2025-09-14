package capture

import (
	"errors"
)

var (
	ErrThingNotFound = errors.New("thing not found")
)

type ThingRepository interface {
	AddThing(thing Thing) (*Thing, error)
	GetThings() ([]Thing, error)
	GetThingByID(id int) (*Thing, error)
}

type InmemoryThingRepository struct {
	things map[int]Thing
	sequence int
}

func NewInmemoryThingRepository() ThingRepository {
	return &InmemoryThingRepository{
		things:  make(map[int]Thing),
		sequence: 0,
	}
}
func (r *InmemoryThingRepository) AddThing(thing Thing) (*Thing, error) {
	r.sequence += 1
	thing.ID = r.sequence
	r.things[thing.ID] = thing
	return &thing, nil
}

func (r *InmemoryThingRepository) GetThings() ([]Thing, error) {
    thingList := make([]Thing, 0, len(r.things))
	for _, t := range r.things {
		thingList = append(thingList, t)
	}
	return thingList, nil
}

func (r *InmemoryThingRepository) GetThingByID(id int) (*Thing, error) {
	foundThing, exists := r.things[id]
	if exists {	
		return &foundThing, nil
	}
	return nil, ErrThingNotFound
}
