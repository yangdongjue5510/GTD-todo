package capture

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewThingService(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockThingRepository(ctrl)
	service := NewThingService(repo)

	if service == nil {
		t.Fatal("Expected non-nil ThingService")
	}
}

func TestThingService_AddThing_Success(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockRepo := NewMockThingRepository(ctrl)
	service := NewThingService(mockRepo)
	
	inputThing := &Thing{
		Title:       "Test Thing",
		Description: "Test Description",
		Status:      Active,
	}
	
	expectedThing := &Thing{
		ID:          1,
		Title:       "Test Thing",
		Description: "Test Description",
		Status:      Active,
	}
	
	// when
	mockRepo.EXPECT().
		AddThing(inputThing).
		Return(expectedThing, nil).
		Times(1)
	
	// then
	result, err := service.AddThing(inputThing)
	
	assert.NoError(t, err)
	assert.Equal(t, expectedThing, result)
}

func TestThingService_AddThing_RepositoryError(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockRepo := NewMockThingRepository(ctrl)
	service := NewThingService(mockRepo)
	
	inputThing := &Thing{
		Title:       "Test Thing",
		Description: "Test Description",
		Status:      Active,
	}
	
	expectedError := errors.New("repository error")
	
	// when
	mockRepo.EXPECT().
		AddThing(inputThing).
		Return(nil, expectedError).
		Times(1)
	
	// then
	result, err := service.AddThing(inputThing)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestThingService_GetThings_Success(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockRepo := NewMockThingRepository(ctrl)
	service := NewThingService(mockRepo)
	
	expectedThings := []*Thing{
		{ID: 1, Title: "Thing 1", Description: "Desc 1", Status: Active},
		{ID: 2, Title: "Thing 2", Description: "Desc 2", Status: Done},
	}
	
	// when
	mockRepo.EXPECT().
		GetThings().
		Return(expectedThings, nil).
		Times(1)
	
	// then
	result, err := service.GetThings()
	
	assert.NoError(t, err)
	assert.Equal(t, expectedThings, result)
	assert.Len(t, result, 2)
}

func TestThingService_GetThings_RepositoryError(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockRepo := NewMockThingRepository(ctrl)
	service := NewThingService(mockRepo)
	
	expectedError := errors.New("database connection error")
	
	// when
	mockRepo.EXPECT().
		GetThings().
		Return(nil, expectedError).
		Times(1)
	
	// then
	result, err := service.GetThings()
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestThingService_MarkThingAsProcessed_Success(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockRepo := NewMockThingRepository(ctrl)
	service := NewThingService(mockRepo)
	
	thingID := 1
	foundThing := &Thing{
		ID:          thingID,
		Title:       "Test Thing",
		Description: "Test Description",
		Status:      Active,
	}
	
	// when
	mockRepo.EXPECT().
		GetThingByID(thingID).
		Return(foundThing, nil).
		Times(1)
	
	// then
	err := service.MarkThingAsProcessed(thingID)
	
	assert.NoError(t, err)
	assert.Equal(t, Done, foundThing.Status) // Process() 호출로 상태 변경 확인
}

func TestThingService_MarkThingAsProcessed_ThingNotFound(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockRepo := NewMockThingRepository(ctrl)
	service := NewThingService(mockRepo)
	
	thingID := 999
	
	// when
	mockRepo.EXPECT().
		GetThingByID(thingID).
		Return(nil, ErrThingNotFound).
		Times(1)
	
	// then
	err := service.MarkThingAsProcessed(thingID)
	
	assert.Error(t, err)
	assert.Equal(t, ErrThingNotFound, err)
}
