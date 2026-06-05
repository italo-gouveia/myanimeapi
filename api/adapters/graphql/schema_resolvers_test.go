package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseID(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedID  uint
		expectError bool
	}{
		{
			name:        "valid positive integer string",
			input:       "42",
			expectedID:  42,
			expectError: false,
		},
		{
			name:        "zero returns error",
			input:       "0",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "negative integer returns error",
			input:       "-5",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "non-numeric string returns error",
			input:       "abc",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "empty string returns error",
			input:       "",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "very large number within int range",
			input:       "2147483647",
			expectedID:  2147483647,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := parseID(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, uint(0), id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}
