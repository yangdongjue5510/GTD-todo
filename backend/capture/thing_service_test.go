package capture

import (
	"testing"
)

func TestInmemoryThingService_AddThing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		thing       Thing
		wantError   bool
		errorMsg    string
	}{
		{
			name: "Valid thing should be added successfully",
			thing: Thing{
				Title:       "Test Thing",
				Description: "Test Description",
				Status:      Active,
			},
			wantError: false,
		},
		{
			name: "Thing with empty title should return error",
			thing: Thing{
				Title:       "",
				Description: "Test Description",
				Status:      Active,
			},
			wantError: true,
			errorMsg:  "thing title cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			// given
			service := NewInmemoryThingService()
			
			// when
			createdThing, err := service.AddThing(tt.thing)
			
			// then
			if tt.wantError {
				if err == nil {
					t.Errorf("AddThing() expected error but got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("AddThing() error = %v, expected %v", err.Error(), tt.errorMsg)
				}
				if createdThing != nil {
					t.Error("AddThing() should return nil on error")
				}
			} else {
				if err != nil {
					t.Errorf("AddThing() unexpected error = %v", err)
				}
				if createdThing == nil {
					t.Error("AddThing() should return created thing on success")
				}
				if createdThing.ID != 1 {
					t.Errorf("Expected ID 1, got %d", createdThing.ID)
				}
				
				// Verify thing was added to service
				things := service.GetThings()
				if len(things) != 1 {
					t.Errorf("Expected 1 thing, got %d", len(things))
				}
			}
		})
	}
}

func TestInmemoryThingService_GetThings(t *testing.T) {
	t.Parallel()
	
	// given
	service := NewInmemoryThingService()
	thing1 := Thing{Title: "Thing 1", Description: "Desc 1", Status: Active}
	thing2 := Thing{Title: "Thing 2", Description: "Desc 2", Status: Done}
	
	service.AddThing(thing1)
	service.AddThing(thing2)
	
	// when
	things := service.GetThings()
	
	// then
	if len(things) != 2 {
		t.Errorf("Expected 2 things, got %d", len(things))
	}
	
	// Verify immutability (should return copy, not reference)
	things[0].Title = "Modified"
	originalThings := service.GetThings()
	if originalThings[0].Title == "Modified" {
		t.Error("GetThings() should return a copy, not reference")
	}
}

func TestInmemoryThingService_ClarifyThing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupThing *Thing
		thingID   int
		wantError bool
		errorMsg  string
	}{
		{
			name: "Valid thing ID should return clarified data",
			setupThing: &Thing{
				Title:       "Test Thing",
				Description: "Test Description",
				Status:      Active,
			},
			thingID:   1,
			wantError: false,
		},
		{
			name:      "Non-existent thing ID should return error",
			thingID:   999,
			wantError: true,
			errorMsg:  "thing not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			// given
			service := NewInmemoryThingService()
			if tt.setupThing != nil {
				service.AddThing(*tt.setupThing)
			}
			
			// when
			clarified, err := service.ClarifyThing(tt.thingID)
			
			// then
			if tt.wantError {
				if err == nil {
					t.Errorf("ClarifyThing() expected error but got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("ClarifyThing() error = %v, expected %v", err.Error(), tt.errorMsg)
				}
				if clarified != nil {
					t.Error("ClarifyThing() should return nil on error")
				}
			} else {
				if err != nil {
					t.Errorf("ClarifyThing() unexpected error = %v", err)
				}
				if clarified == nil {
					t.Error("ClarifyThing() should return clarified data on success")
				}
				
				// Verify clarified data structure
				if clarified.Title != tt.setupThing.Title {
					t.Errorf("ClarifyThing() Title = %v, expected %v", clarified.Title, tt.setupThing.Title)
				}
				if clarified.Priority != "normal" {
					t.Errorf("ClarifyThing() Priority = %v, expected normal", clarified.Priority)
				}
				if clarified.Context != "inbox" {
					t.Errorf("ClarifyThing() Context = %v, expected inbox", clarified.Context)
				}
				if clarified.SourceID != tt.thingID {
					t.Errorf("ClarifyThing() SourceID = %v, expected %v", clarified.SourceID, tt.thingID)
				}
			}
		})
	}
}

func TestInmemoryThingService_MarkThingAsProcessed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupThing *Thing
		thingID   int
		wantError bool
		errorMsg  string
	}{
		{
			name: "Valid thing ID should be marked as Done",
			setupThing: &Thing{
				Title:       "Test Thing",
				Description: "Test Description",
				Status:      Active,
			},
			thingID:   1,
			wantError: false,
		},
		{
			name:      "Non-existent thing ID should return error",
			thingID:   999,
			wantError: true,
			errorMsg:  "thing not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			// given
			service := NewInmemoryThingService()
			if tt.setupThing != nil {
				service.AddThing(*tt.setupThing)
			}
			
			// when
			err := service.MarkThingAsProcessed(tt.thingID)
			
			// then
			if tt.wantError {
				if err == nil {
					t.Errorf("MarkThingAsProcessed() expected error but got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("MarkThingAsProcessed() error = %v, expected %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("MarkThingAsProcessed() unexpected error = %v", err)
				}
				
				// Verify thing status was updated
				things := service.GetThings()
				if len(things) == 0 {
					t.Error("Expected at least 1 thing")
					return
				}
				if things[0].Status != Done {
					t.Errorf("Expected status Done, got %v", things[0].Status)
				}
			}
		})
	}
}