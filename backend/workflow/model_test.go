package workflow

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
			name:     "ToDo status returns correct string",
			status:   ToDo,
			expected: "ToDo",
		},
		{
			name:     "InProgress status returns correct string",
			status:   InProgress,
			expected: "InProgress",
		},
		{
			name:     "Completed status returns correct string",
			status:   Completed,
			expected: "Completed",
		},
		{
			name:     "Delayed status returns correct string",
			status:   Delayed,
			expected: "Delayed",
		},
		{
			name:     "Delegated status returns correct string",
			status:   Delegated,
			expected: "Delegated",
		},
		{
			name:     "Planned status returns correct string",
			status:   Planned,
			expected: "Planned",
		},
		{
			name:     "Someday status returns correct string",
			status:   Someday,
			expected: "Someday",
		},
		{
			name:     "Removed status returns correct string",
			status:   Removed,
			expected: "Removed",
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