package thing

type ThingService interface {
	AddThing(thing Thing)
	GetThings() []Thing
}

type InmemoryThingService struct {
	things   []Thing
	sequence int
}

func (s *InmemoryThingService) AddThing(thing Thing) {
	s.sequence++          // Increment sequence for unique ID
	thing.ID = s.sequence // Simple ID generation
	s.things = append(s.things, thing)
}

func (s *InmemoryThingService) GetThings() []Thing {
	copiedThings := make([]Thing, len(s.things))
	copy(copiedThings, s.things)
	return copiedThings
}
