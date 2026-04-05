package service

import (
	"context"
	"testing"
	"time"

	"github.com/moises-ba/mb-crypto-sma-api/internal/domain"
	appErrors "github.com/moises-ba/mb-crypto-sma-api/internal/errors"
	"github.com/stretchr/testify/assert"
)

type mockCache struct {
	getFunc func(ctx context.Context, key string, dest any) (bool, appErrors.ApiError)
	setFunc func(ctx context.Context, key string, value any, ttl time.Duration) appErrors.ApiError
}

func (m *mockCache) Get(ctx context.Context, key string, dest interface{}) (bool, appErrors.ApiError) {
	return m.getFunc(ctx, key, dest)
}

func (m *mockCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) appErrors.ApiError {
	return m.setFunc(ctx, key, value, ttl)
}

type mockService struct {
	getFunc func(ctx context.Context, qtDays int, ref time.Time) ([]domain.AverageResponse, appErrors.ApiError)
}

func (m *mockService) GetLastAVGPrices(ctx context.Context, qtDays int, ref time.Time) ([]domain.AverageResponse, appErrors.ApiError) {
	return m.getFunc(ctx, qtDays, ref)
}

func TestCacheDecorator_GetLastAVGPrices(t *testing.T) {
	refDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	expectedResponse := []domain.AverageResponse{
		{DigitalCoin: domain.BTC, Value: "10"},
	}

	tests := []struct {
		name        string
		cache       mockCache
		service     mockService
		expectError bool
		expectCall  bool
	}{
		{
			name: "should return from cache",
			cache: mockCache{
				getFunc: func(ctx context.Context, key string, dest interface{}) (bool, appErrors.ApiError) {
					ptr := dest.(*[]domain.AverageResponse)
					*ptr = expectedResponse
					return true, nil
				},
				setFunc: func(ctx context.Context, key string, value interface{}, ttl time.Duration) appErrors.ApiError {
					return nil
				},
			},
			service: mockService{
				getFunc: func(ctx context.Context, qtDays int, ref time.Time) ([]domain.AverageResponse, appErrors.ApiError) {
					t.Fatal("service should not be called")
					return nil, nil
				},
			},
			expectError: false,
		},
		{
			name: "should fallback to service when cache get fails",
			cache: mockCache{
				getFunc: func(ctx context.Context, key string, dest interface{}) (bool, appErrors.ApiError) {
					return false, appErrors.NewApiError("error")
				},
				setFunc: func(ctx context.Context, key string, value interface{}, ttl time.Duration) appErrors.ApiError {
					return nil
				},
			},
			service: mockService{
				getFunc: func(ctx context.Context, qtDays int, ref time.Time) ([]domain.AverageResponse, appErrors.ApiError) {
					return expectedResponse, nil
				},
			},
			expectError: false,
		},
		{
			name: "should call service and cache result",
			cache: mockCache{
				getFunc: func(ctx context.Context, key string, dest interface{}) (bool, appErrors.ApiError) {
					return false, nil
				},
				setFunc: func(ctx context.Context, key string, value interface{}, ttl time.Duration) appErrors.ApiError {
					assert.Equal(t, expectedResponse, value)
					return nil
				},
			},
			service: mockService{
				getFunc: func(ctx context.Context, qtDays int, ref time.Time) ([]domain.AverageResponse, appErrors.ApiError) {
					return expectedResponse, nil
				},
			},
			expectError: false,
		},
		{
			name: "should return error from service",
			cache: mockCache{
				getFunc: func(ctx context.Context, key string, dest interface{}) (bool, appErrors.ApiError) {
					return false, nil
				},
				setFunc: func(ctx context.Context, key string, value interface{}, ttl time.Duration) appErrors.ApiError {
					return nil
				},
			},
			service: mockService{
				getFunc: func(ctx context.Context, qtDays int, ref time.Time) ([]domain.AverageResponse, appErrors.ApiError) {
					return nil, appErrors.NewApiError("service error")
				},
			},
			expectError: true,
		},
		{
			name: "should ignore cache set error",
			cache: mockCache{
				getFunc: func(ctx context.Context, key string, dest interface{}) (bool, appErrors.ApiError) {
					return false, nil
				},
				setFunc: func(ctx context.Context, key string, value interface{}, ttl time.Duration) appErrors.ApiError {
					return appErrors.NewApiError("error")
				},
			},
			service: mockService{
				getFunc: func(ctx context.Context, qtDays int, ref time.Time) ([]domain.AverageResponse, appErrors.ApiError) {
					return expectedResponse, nil
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decorator := NewCryptorCalculatorCachable(&tt.service, &tt.cache)

			res, err := decorator.GetLastAVGPrices(context.Background(), 20, refDate)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, res)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, expectedResponse, res)
		})
	}
}
