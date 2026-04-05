package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moises-ba/mb-crypto-sma-api/internal/domain"
	appErrors "github.com/moises-ba/mb-crypto-sma-api/internal/errors"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	getFunc func(ctx context.Context, days int, ref time.Time) ([]domain.AverageResponse, appErrors.ApiError)
}

func (m *mockService) GetLastAVGPrices(ctx context.Context, days int, ref time.Time) ([]domain.AverageResponse, appErrors.ApiError) {
	if m.getFunc != nil {
		return m.getFunc(ctx, days, ref)
	}
	return nil, nil
}

func TestListLastClosedPrices(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		daysParam      string
		dateParam      string
		mock           mockService
		expectedStatus int
	}{
		{
			name:           "invalid days",
			daysParam:      "abc",
			dateParam:      "2024-01-01",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid date",
			daysParam:      "20",
			dateParam:      "invalid-date",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "service error - invalid",
			daysParam: "20",
			dateParam: "2024-01-01",
			mock: mockService{
				getFunc: func(ctx context.Context, days int, ref time.Time) ([]domain.AverageResponse, appErrors.ApiError) {
					return nil, appErrors.NewApiError("invalid request", appErrors.WithKind(appErrors.Invalid))
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "service error - generic (fallback 500)",
			daysParam: "20",
			dateParam: "2024-01-01",
			mock: mockService{
				getFunc: func(ctx context.Context, days int, ref time.Time) ([]domain.AverageResponse, appErrors.ApiError) {
					return nil, appErrors.NewApiError("internal error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:      "success",
			daysParam: "20",
			dateParam: "2024-01-01",
			mock: mockService{
				getFunc: func(ctx context.Context, days int, ref time.Time) ([]domain.AverageResponse, appErrors.ApiError) {
					return []domain.AverageResponse{
						{DigitalCoin: domain.BTC, Value: "10"},
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			ctrl := NewCriptorCalculatorController(&tt.mock)
			router.GET("/v1/lastclosedprices/:days/:dt_reference", ctrl.ListLastClosedPrices)

			req := httptest.NewRequest(
				http.MethodGet,
				"/v1/lastclosedprices/"+tt.daysParam+"/"+tt.dateParam,
				nil,
			)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
