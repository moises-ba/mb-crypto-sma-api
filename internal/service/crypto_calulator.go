package service

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/moises-ba/mb-crypto-mms-api/internal/adapter"
	"github.com/moises-ba/mb-crypto-mms-api/internal/domain"
	"github.com/moises-ba/mb-crypto-mms-api/internal/errors"
	"github.com/moises-ba/mb-crypto-mms-api/internal/strategy"
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
	if !slices.Contains(allowedQtDays, qtDays) {
		return nil, errors.NewApiError(fmt.Sprintf("invalid qt_days: %v", qtDays), errors.WithKind(errors.Internal))
	}
	return s.findAverages(ctx, qtDays, referenceDate)
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
	res, err := createAverengeResponse(digitalCoin, dtEnd, qtDays, lastPriceRes, strategy.SimpleMovingAverage)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func createAverengeResponse(digitalCoin domain.DigitalCoin, referenceDate time.Time, qtDays int, lastClosed []string, calcStrategy strategy.AvgCalculator) (*domain.AverageResponse, errors.ApiError) {
	values, err := convert(lastClosed)
	if err != nil {
		return nil, err
	}

	avg, err := calcStrategy(values)
	if err != nil {
		return nil, err
	}

	return &domain.AverageResponse{
		ReferenceDate: referenceDate.Format("2006-01-02"),
		AvgType:       strconv.Itoa(qtDays),
		DigitalCoin:   digitalCoin,
		Value:         avg.String(),
	}, nil
}

func convert(priceValues []string) ([]decimal.Decimal, errors.ApiError) {
	res := make([]decimal.Decimal, len(priceValues))
	for _, priceValue := range priceValues {
		v, err := decimal.NewFromString(priceValue)
		if err != nil {
			return nil, errors.NewApiError("unexpected response", errors.WithKind(errors.Unexpected), errors.WithError(err))
		}
		res = append(res, v)
	}
	return res, nil
}
