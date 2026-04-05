package timeutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		dateStr     string
		pattern     string
		expectError bool
	}{
		{
			name:        "valid date",
			dateStr:     "2024-01-01",
			pattern:     "2006-01-02",
			expectError: false,
		},
		{
			name:        "invalid date format",
			dateStr:     "01-01-2024",
			pattern:     "2006-01-02",
			expectError: true,
		},
		{
			name:        "invalid pattern",
			dateStr:     "2024-01-01",
			pattern:     "invalid-pattern",
			expectError: true,
		},
		{
			name:        "empty date string",
			dateStr:     "",
			pattern:     "2006-01-02",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Parse(tt.dateStr, tt.pattern)

			if tt.expectError {
				assert.NotNil(t, err)
				assert.Nil(t, res)

				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, res)
			expected, _ := time.Parse(tt.pattern, tt.dateStr)

			assert.True(t, expected.Equal(*res))
		})
	}
}
