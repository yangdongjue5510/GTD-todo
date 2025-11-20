package capture

import (
	"testing"
)

func TestNewInmemoryThingRepository(t *testing.T) {
	repo := NewInmemoryThingRepository()
	
	if repo == nil {
		t.Fatal("NewInmemoryThingRepository should return a non-nil repository")
	}
}

func TestAddThing_validateIdAssignment(t *testing.T) {
	repo := NewInmemoryThingRepository()
	
	tests := []struct {
		name        string
		thing       *Thing
		expectedID  int
	}{
		{
			name: "첫 번째 Thing 추가",
			thing: &Thing{
				Title:       "첫 번째 일",
				Description: "첫 번째 일 설명",
				Status:      Active,
			},
			expectedID: 1,
		},
		{
			name: "두 번째 Thing 추가",
			thing: &Thing{
				Title:       "두 번째 일",
				Description: "두 번째 일 설명", 
				Status:      Active,
			},
			expectedID: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.AddThing(tt.thing)
			
			if err != nil {
				t.Fatalf("AddThing() error = %v", err)
			}
			
			if result.ID != tt.expectedID {
				t.Errorf("AddThing() ID = %d, expected %d", result.ID, tt.expectedID)
			}
		})
	}
}

func TestGetThings(t *testing.T) {
	repo := NewInmemoryThingRepository()
	
	
	// Thing 추가 후 조회
	thing1 := &Thing{Title: "첫 번째", Description: "설명1", Status: Active}
	thing2 := &Thing{Title: "두 번째", Description: "설명2", Status: Done}
	
	repo.AddThing(thing1)
	repo.AddThing(thing2)

	things, err := repo.GetThings()
	if err != nil {
		t.Fatalf("GetThings() error = %v", err)
	}
	
	if len(things) != 2 {
		t.Errorf("GetThings() should return 2 items, got %d", len(things))
	}
	
	// ID 순서대로 정렬되어 있는지 확인하지 않음 (맵 순회는 순서 보장 안함)
	foundThing1 := false
	foundThing2 := false
	
	for _, thing := range things {
		if thing.ID == 1 {
			foundThing1 = true
		}
		if thing.ID == 2 {
			foundThing2 = true
		}
	}
	
	if !foundThing1 {
		t.Error("GetThings() should include first thing")
	}
	if !foundThing2 {
		t.Error("GetThings() should include second thing")
	}
}

func TestGetThingByID(t *testing.T) {
	repo := NewInmemoryThingRepository()
	
	// Thing 추가 후 조회
	originalThing := &Thing{
		Title:       "테스트 Thing",
		Description: "테스트 설명",
		Status:      Active,
	}
	
	addedThing, err := repo.AddThing(originalThing)
	if err != nil {
		t.Fatalf("AddThing() error = %v", err)
	}
	
	// 추가된 Thing 조회
	foundThing, err := repo.GetThingByID(addedThing.ID)
	if err != nil {
		t.Fatalf("GetThingByID() error = %v", err)
	}

	if foundThing.ID != addedThing.ID {
		t.Fatalf("GetThingByID() ID = %d, expected %d", foundThing.ID, addedThing.ID)
	}
}

func TestGetThingById_notExist(t *testing.T) {
	repo := NewInmemoryThingRepository()

	thing, err := repo.GetThingByID(999)
	if err != ErrThingNotFound {
		t.Errorf("GetThingByID() should return ErrThingNotFound for non-existent ID, got %v", err)
	}
	if thing != nil {
		t.Error("GetThingByID() should return nil Thing for non-existent ID")
	}
}
