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
			name:     "Pending status returns correct string",
			status:   Pending,
			expected: "Pending",
		},
		{
			name:     "Someday status returns correct string",
			status:   Someday,
			expected: "Someday",
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