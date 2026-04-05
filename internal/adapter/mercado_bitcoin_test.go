package adapter

import (
	"context"
	"testing"
	"time"

	"github.com/moises-ba/mb-crypto-sma-api/internal/domain"
	appErrors "github.com/moises-ba/mb-crypto-sma-api/internal/errors"
	"github.com/moises-ba/mb-crypto-sma-api/internal/http"
	"github.com/stretchr/testify/assert"
)

type mockHTTPClient struct {
	getFunc func(ctx context.Context, url string) (*http.Response, appErrors.ApiError)
}

func (m *mockHTTPClient) Get(ctx context.Context, url string) (*http.Response, appErrors.ApiError) {
	return m.getFunc(ctx, url)
}

func TestListLastClosedPrices(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		mock        mockHTTPClient
		expectError bool
		expected    []string
	}{
		{
			name: "success",
			mock: mockHTTPClient{
				getFunc: func(ctx context.Context, url string) (*http.Response, appErrors.ApiError) {
					return &http.Response{
						Response: []byte(`{"c":["10","20","30"]}`),
					}, nil
				},
			},
			expectError: false,
			expected:    []string{"10", "20", "30"},
		},
		{
			name: "http client error",
			mock: mockHTTPClient{
				getFunc: func(ctx context.Context, url string) (*http.Response, appErrors.ApiError) {
					return nil, appErrors.NewApiError("http error")
				},
			},
			expectError: true,
		},
		{
			name: "invalid json",
			mock: mockHTTPClient{
				getFunc: func(ctx context.Context, url string) (*http.Response, appErrors.ApiError) {
					return &http.Response{
						Response: []byte(`invalid-json`),
					}, nil
				},
			},
			expectError: true,
		},
		{
			name: "empty response",
			mock: mockHTTPClient{
				getFunc: func(ctx context.Context, url string) (*http.Response, appErrors.ApiError) {
					return &http.Response{
						Response: []byte(`{"c":[]}`),
					}, nil
				},
			},
			expectError: false,
			expected:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewMercadoBitcoin(&tt.mock)

			res, err := adapter.ListLastClosedPrices(
				context.Background(),
				domain.BTC,
				now.AddDate(0, 0, -10),
				now,
			)

			if tt.expectError {
				assert.NotNil(t, err)
				assert.Nil(t, res)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, tt.expected, res)
		})
	}
}
