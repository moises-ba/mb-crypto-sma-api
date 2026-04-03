package service

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/moises-ba/mb-crypto-mms-api/internal/adapter"
	"github.com/moises-ba/mb-crypto-mms-api/internal/domain"
	"github.com/moises-ba/mb-crypto-mms-api/internal/errors"
	"github.com/shopspring/decimal"
)

var allowedQtDays = []int{20, 50, 200}

type CryptoCalculatorService interface {
	GetLastAVGPrices(ctx context.Context, qtDays int, referenceDate time.Time) ([]domain.AverageResponse, errors.ApiError)
}

type cryptoCalculatorService struct {
	exchangeAdapter adapter.Exchange
}

func (s *cryptoCalculatorService) GetLastAVGPrices(ctx context.Context, qtDays int, referenceDate time.Time) ([]domain.AverageResponse, errors.ApiError) {
	if !isValidaQtDays(qtDays) {
		return nil, errors.NewApiError(fmt.Sprintf("invalid qt_days: %v", qtDays), errors.WithKind(errors.Internal))
	}

	return s.findAverages(ctx, qtDays, referenceDate)
}

func isValidaQtDays(qtDays int) bool {
	for _, qtAllowed := range allowedQtDays {
		if qtAllowed == qtDays {
			return true
		}
	}
	return false
}

func (s *cryptoCalculatorService) findAverages(ctx context.Context, qtDays int, referenceDate time.Time) ([]domain.AverageResponse, errors.ApiError) {
	digitalCoinss := []domain.DigitalCoin{domain.BTC, domain.ETH}
	dtEnd := referenceDate
	dtIni := referenceDate.AddDate(0, 0, -qtDays)
	chanAverageResponse := make(chan domain.AverageResponse, len(digitalCoinss))
	chanErr := make(chan errors.ApiError, len(digitalCoinss))
	wg := sync.WaitGroup{}
	for _, digtitalCoin := range digitalCoinss {
		dc := digtitalCoin //versoes de versoes do go antes de 1.22 precisam disso para evitar race condition
		wg.Add(1)
		go func(dcParam domain.DigitalCoin) {
			defer close(chanAverageResponse)
			average, err := s.findAverage(ctx, dcParam, qtDays, dtIni, dtEnd)
			if err != nil {
				chanErr <- err
				return
			}
			chanAverageResponse <- *average
		}(dc)
	}

	go func() {
		wg.Wait()
		close(chanAverageResponse)
		close(chanErr)
	}()

	var resLastPrices []domain.AverageResponse

	for {
		select {
		case res, ok := <-chanAverageResponse:
			if ok {
				resLastPrices = append(resLastPrices, res)
			} else {
				chanAverageResponse = nil
			}
		case err, ok := <-chanErr:
			if ok {
				return nil, err
			} else {
				chanErr = nil
			}
		}

		if chanAverageResponse == nil && chanErr == nil {
			break
		}
	}

	return resLastPrices, nil
}

func (s *cryptoCalculatorService) findAverage(ctx context.Context, digitalCoin domain.DigitalCoin, qtDays int, dtIni, dtEnd time.Time) (*domain.AverageResponse, errors.ApiError) {
	lastPriceRes, err := s.exchangeAdapter.ListLastPrices(ctx, digitalCoin, dtIni, dtEnd)
	if err != nil {
		return nil, err
	}
	res, err := createAverengeResponse(qtDays, lastPriceRes)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func createAverengeResponse(qtDays int, lastClosed *domain.LastPrices) (*domain.AverageResponse, errors.ApiError) {
	//TODO IMPLEMENTAR
}

func getDatesDescending(lastCloseds map[string]decimal.Decimal) []string {
	lastDatesOrderDesc := make([]string, 0, len(lastCloseds))
	for k := range lastCloseds {
		lastDatesOrderDesc = append(lastDatesOrderDesc, k)
	}
	sort.Slice(lastDatesOrderDesc, func(i, j int) bool {
		return lastDatesOrderDesc[i] > lastDatesOrderDesc[j]
	})
	return lastDatesOrderDesc
}
