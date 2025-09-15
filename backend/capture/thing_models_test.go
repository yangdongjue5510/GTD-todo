package capture

import "testing"

func TestStatus_String(t *testing.T) {
	t.Parallel()
	
	// given, when, then
	tests := []struct {
		name     string
		status   Status
		expected string
	}{
		{
			name:     "Active status returns correct string",
			status:   Active,
			expected: "Active",
		},
		{
			name:     "Done status returns correct string",
			status:   Done,
			expected: "Done",
		},
		{
			name:     "Unknown status returns Unknown",
			status:   Status(999), // Invalid status
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			// when
			result := tt.status.String()
			
			// then
			if result != tt.expected {
				t.Errorf("Status.String() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestThingProcess(t *testing.T) {
	t.Parallel()

	// given
	thing := &Thing{
		ID:          1,
		Title: 	 "Sample Thing",
		Description: "This is a sample thing",
		Status: Active,
	}
	
	// when
	thing.Process()

	// then
	if thing.Status != Done {
		t.Errorf("Expected status to be Done, got %v", thing.Status)
	}
}