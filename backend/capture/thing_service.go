//go:generate mockgen -source=thing_service.go -destination=mock_thing_service.go -package=capture

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
	AddThing(thing *Thing) (*Thing, error)
	GetThings() ([]*Thing, error)
	MarkThingAsProcessed(thingID int) error
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

type ThingServiceImpl struct {
	thingRepository ThingRepository
}

func NewThingService(repo ThingRepository) ThingService {
	return &ThingServiceImpl{thingRepository: repo}
}

func (s *ThingServiceImpl) AddThing(thing *Thing) (*Thing, error) {
	if thing.Title == "" {
		return nil, errors.New("thing title cannot be empty")
	}
	return s.thingRepository.AddThing(thing)
}

func (s *ThingServiceImpl) GetThings() ([]*Thing, error) {
	return s.thingRepository.GetThings()
}

func (s *ThingServiceImpl) MarkThingAsProcessed(thingID int) error {
	foundThing, err := s.thingRepository.GetThingByID(thingID)
	if err != nil {
		return err
	}
	foundThing.Process()
	return nil
}
