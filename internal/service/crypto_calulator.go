package service

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/moises-ba/mb-crypto-sma-api/internal/adapter"
	"github.com/moises-ba/mb-crypto-sma-api/internal/domain"
	"github.com/moises-ba/mb-crypto-sma-api/internal/dto"
	"github.com/moises-ba/mb-crypto-sma-api/internal/errors"
	"github.com/moises-ba/mb-crypto-sma-api/internal/strategy"
	"github.com/moises-ba/mb-crypto-sma-api/internal/timeutil"
	"github.com/shopspring/decimal"

	"golang.org/x/sync/singleflight"
)

var allowedQtDays = []int{20, 50, 200}

type CryptoCalculatorService interface {
	GetLastAVGPrices(ctx context.Context, qtDays int, referenceDate time.Time) ([]domain.AverageResponse, errors.ApiError)
}

type cryptoCalculatorService struct {
	singleFlightGroup singleflight.Group
	exchangeAdapter   adapter.Exchange
}

func NewCryptoCalculatorService(exchangeAdapter adapter.Exchange) CryptoCalculatorService {
	return &cryptoCalculatorService{exchangeAdapter: exchangeAdapter}
}

func (s *cryptoCalculatorService) GetLastAVGPrices(ctx context.Context, qtDays int, referenceDate time.Time) ([]domain.AverageResponse, errors.ApiError) {
	if err := validate(qtDays, referenceDate); err != nil {
		return nil, err
	}

	return s.findAverages(ctx, dto.AverageRequest{
		DigitalCoins: []domain.DigitalCoin{domain.BTC, domain.ETH},
		AvgType:      domain.SMA,
		Days:         qtDays,
		StartDate:    referenceDate.AddDate(0, 0, -(qtDays - 1)),
		EndDate:      referenceDate,
	})
}

func validate(qtDays int, referenceDate time.Time) errors.ApiError {
	if !slices.Contains(allowedQtDays, qtDays) {
		return errors.NewApiError(fmt.Sprintf("invalid qt_days: %v", qtDays), errors.WithKind(errors.Invalid))
	}

	if referenceDate.After(time.Now().UTC()) {
		return errors.NewApiError("date is greather today", errors.WithKind(errors.Invalid))
	}

	return nil
}

func (s *cryptoCalculatorService) findAverages(ctx context.Context, req dto.AverageRequest) ([]domain.AverageResponse, errors.ApiError) {
	chanAverageResponse := make(chan domain.AverageResponse, len(req.DigitalCoins))
	chanErr := make(chan errors.ApiError, len(req.DigitalCoins))
	wg := sync.WaitGroup{}
	for _, digtitalCoin := range req.DigitalCoins {
		dc := digtitalCoin //versoes de versoes do go antes de 1.22 precisam disso para evitar race condition
		wg.Add(1)
		go func(dcParam domain.DigitalCoin) {
			defer wg.Done()
			average, err := s.findAverage(ctx, dc, req)
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

func (s *cryptoCalculatorService) findAverage(ctx context.Context, digitalCoin domain.DigitalCoin, req dto.AverageRequest) (*domain.AverageResponse, errors.ApiError) {
	var lastPriceRes []string

	//single flight serve para evitar o thundering herd problem,  requisições duplicadas concorrentes, onde multiplas goroutines tentam acessar o mesmo dado a mesmo tempo.
	//com isso vamos na api apenas uma vez e todas as goroutines recebem o resultado para a mesma chave
	//observacao: ele nao guarda cache, so evita que varias goroutines faça a mesma requisicao ao mesmo tempo em q ainda outra esta sendo executada
	res, err, _ := s.singleFlightGroup.Do(fmt.Sprintf("%s_%s_%v", req.EndDate.Format(timeutil.YearMonthDayPattern), digitalCoin, req.Days),
		func() (any, error) {
			return s.exchangeAdapter.ListLastClosedPrices(ctx, digitalCoin, req.StartDate, req.EndDate)
		})
	if err != nil {
		return nil, err.(errors.ApiError)
	}

	lastPriceRes = res.([]string)

	if len(lastPriceRes) == 0 {
		return nil, errors.NewApiError(fmt.Sprintf("last price for %s not found", digitalCoin), errors.WithKind(errors.NotFound))
	}

	if len(lastPriceRes) < req.Days {
		return nil, errors.NewApiError(fmt.Sprintf("it is not posible calculate avg for %s because has no all closed values for %v days", digitalCoin, req.Days),
			errors.WithKind(errors.Unprocessable))
	}

	avgCalculatorF := strategy.GetStrategy(req.AvgType)
	if avgCalculatorF == nil {
		return nil, errors.NewApiError(fmt.Sprintf("calculator for type: %s is invalid", req.AvgType), errors.WithKind(errors.Invalid))
	}

	return createAverengeResponse(digitalCoin, req, lastPriceRes, avgCalculatorF)
}

func createAverengeResponse(digitalCoin domain.DigitalCoin, req dto.AverageRequest, lastClosed []string, calcStrategy strategy.AvgCalculator) (*domain.AverageResponse, errors.ApiError) {
	values, err := convert(lastClosed)
	if err != nil {
		return nil, err
	}

	avg, err := calcStrategy(values)
	if err != nil {
		return nil, err
	}

	return &domain.AverageResponse{
		ReferenceDate: req.EndDate.Format("2006-01-02"),
		AvgType:       string(req.AvgType),
		Days:          req.Days,
		DigitalCoin:   digitalCoin,
		Value:         avg.String(),
	}, nil
}

func convert(priceValues []string) ([]decimal.Decimal, errors.ApiError) {
	res := make([]decimal.Decimal, 0, len(priceValues))
	for _, priceValue := range priceValues {
		v, err := decimal.NewFromString(priceValue)
		if err != nil {
			return nil, errors.NewApiError("unexpected response", errors.WithKind(errors.Unexpected), errors.WithError(err))
		}
		res = append(res, v)
	}
	return res, nil
}
