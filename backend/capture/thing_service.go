package capture

import (
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

type AddThingUseCase interface {
	AddThing(thing Thing) (*Thing, error)
}

type GetThingUseCase interface {
	GetThingByID(id int) (*Thing, error)
}

type GetThingsUseCase interface {
	GetThings() ([]Thing, error)
}

type UpdateThingStatusUseCase interface {
	UpdateThingStatus(id int, status Status) error
}

type ThingService struct {
	thingRepository ThingRepository
}

func NewThingService(repo ThingRepository) *ThingService {
	return &ThingService{thingRepository: repo}
}

func (s *ThingService) AddThing(thing *Thing) (*Thing, error) {
	return s.thingRepository.AddThing(thing)
}

func (s *ThingService) GetThings() ([]*Thing, error) {
	return s.thingRepository.GetThings()
}

func (s *ThingService) MarkThingAsProcessed(thingID int) error {
	foundThing, err := s.thingRepository.GetThingByID(thingID)
	if err != nil {
		return err
	}
	foundThing.Process()
	return nil
}
