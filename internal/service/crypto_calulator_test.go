package service

import (
	"context"
	"testing"
	"time"

	"github.com/moises-ba/mb-crypto-sma-api/internal/domain"
	"github.com/moises-ba/mb-crypto-sma-api/internal/dto"
	appErrors "github.com/moises-ba/mb-crypto-sma-api/internal/errors"
	"github.com/stretchr/testify/assert"
)

type mockExchange struct {
	listFunc func(ctx context.Context, coin domain.DigitalCoin, start, end time.Time) ([]string, appErrors.ApiError)
}

func (m *mockExchange) ListLastClosedPrices(ctx context.Context, coin domain.DigitalCoin, start, end time.Time) ([]string, appErrors.ApiError) {
	return m.listFunc(ctx, coin, start, end)
}

func TestGetLastAVGPrices(t *testing.T) {
	refDate := time.Now()

	tests := []struct {
		name        string
		qtDays      int
		mock        mockExchange
		expectError bool
	}{
		{
			name:        "should return error when qtDays invalid",
			qtDays:      10,
			mock:        mockExchange{},
			expectError: true,
		},
		{
			name:   "should calculate averages successfully",
			qtDays: 20,
			mock: mockExchange{
				listFunc: func(ctx context.Context, coin domain.DigitalCoin, start, end time.Time) ([]string, appErrors.ApiError) {
					return generatePrices(20, "10"), nil
				},
			},
			expectError: false,
		},
		{
			name:   "should return error from adapter",
			qtDays: 20,
			mock: mockExchange{
				listFunc: func(ctx context.Context, coin domain.DigitalCoin, start, end time.Time) ([]string, appErrors.ApiError) {
					return nil, appErrors.NewApiError("adapter error")
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCryptoCalculatorService(&tt.mock)

			res, err := svc.GetLastAVGPrices(context.Background(), tt.qtDays, refDate)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, res)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, res, 2) // BTC + ETH
		})
	}
}

func TestFindAverage(t *testing.T) {
	req := dummyRequest(20)

	tests := []struct {
		name        string
		mock        mockExchange
		expectError bool
	}{
		{
			name: "should return not found when empty prices",
			mock: mockExchange{
				listFunc: func(ctx context.Context, coin domain.DigitalCoin, start, end time.Time) ([]string, appErrors.ApiError) {
					return []string{}, nil
				},
			},
			expectError: true,
		},
		{
			name: "should return error when insufficient data",
			mock: mockExchange{
				listFunc: func(ctx context.Context, coin domain.DigitalCoin, start, end time.Time) ([]string, appErrors.ApiError) {
					return generatePrices(10, "10"), nil
				},
			},
			expectError: true,
		},
		{
			name: "should calculate successfully",
			mock: mockExchange{
				listFunc: func(ctx context.Context, coin domain.DigitalCoin, start, end time.Time) ([]string, appErrors.ApiError) {
					return generatePrices(20, "10"), nil
				},
			},
			expectError: false,
		},
		{
			name: "should propagate adapter error",
			mock: mockExchange{
				listFunc: func(ctx context.Context, coin domain.DigitalCoin, start, end time.Time) ([]string, appErrors.ApiError) {
					return nil, appErrors.NewApiError("fail")
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := cryptoCalculatorService{exchangeAdapter: &tt.mock}

			res, err := svc.findAverage(context.Background(), domain.BTC, req)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, res)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, res)
			assert.Equal(t, "10", res.Value)
		})
	}
}

func TestConvert(t *testing.T) {
	tests := []struct {
		name        string
		input       []string
		expectError bool
	}{
		{
			name:        "should convert successfully",
			input:       []string{"10", "20"},
			expectError: false,
		},
		{
			name:        "should return error on invalid value",
			input:       []string{"invalid"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := convert(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, res)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, res, len(tt.input))
		})
	}
}

func generatePrices(n int, value string) []string {
	res := make([]string, n)
	for i := range res {
		res[i] = value
	}
	return res
}

func dummyRequest(days int) dto.AverageRequest {
	now := time.Now()
	return dto.AverageRequest{
		DigitalCoins: []domain.DigitalCoin{domain.BTC},
		AvgType:      domain.SMA,
		Days:         days,
		StartDate:    now.AddDate(0, 0, -(days - 1)),
		EndDate:      now,
	}
}
