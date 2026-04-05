package strategy

import (
	"testing"

	"github.com/moises-ba/mb-crypto-sma-api/internal/domain"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestSimpleMovingAverage(t *testing.T) {
	tests := []struct {
		name        string
		input       []decimal.Decimal
		expected    *decimal.Decimal
		expectError bool
	}{
		{
			name: "should calculate average correctly",
			input: []decimal.Decimal{
				decimal.NewFromInt(10),
				decimal.NewFromInt(20),
				decimal.NewFromInt(30),
			},
			expected: func() *decimal.Decimal {
				v := decimal.NewFromInt(20)
				return &v
			}(),
			expectError: false,
		},
		{
			name:        "should return error when empty slice",
			input:       []decimal.Decimal{},
			expected:    nil,
			expectError: true,
		},
		{
			name: "should handle single value",
			input: []decimal.Decimal{
				decimal.NewFromInt(42),
			},
			expected: func() *decimal.Decimal {
				v := decimal.NewFromInt(42)
				return &v
			}(),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := simpleMovingAverage(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.True(t, tt.expected.Equal(*result))
		})
	}
}

func TestGetStrategy(t *testing.T) {
	tests := []struct {
		name     string
		input    domain.AvgType
		expected bool // se deve retornar função ou nil
	}{
		{
			name:     "should return SMA strategy",
			input:    domain.SMA,
			expected: true,
		},
		{
			name:     "should return nil for unknown strategy",
			input:    domain.AvgType("UNKNOWN"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStrategy(tt.input)

			if tt.expected {
				assert.NotNil(t, result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}
