package strategy

import (
	"github.com/moises-ba/mb-crypto-sma-api/internal/domain"
	"github.com/moises-ba/mb-crypto-sma-api/internal/errors"
	"github.com/shopspring/decimal"
)

var mapAvgStrategy = map[domain.AvgType]AvgCalculator{
	domain.SMA: simpleMovingAverage,
}

type AvgCalculator func(values []decimal.Decimal) (*decimal.Decimal, errors.ApiError)

func simpleMovingAverage(values []decimal.Decimal) (*decimal.Decimal, errors.ApiError) {
	totalElements := len(values)
	if totalElements == 0 {
		return nil, errors.NewApiError("values are required", errors.WithKind(errors.Invalid))
	}

	var total decimal.Decimal
	for _, val := range values {
		total = total.Add(val)
	}

	avg := total.Div(decimal.NewFromInt(int64(totalElements)))

	return &avg, nil
}

func GetStrategy(avgType domain.AvgType) AvgCalculator {
	if f, ok := mapAvgStrategy[avgType]; ok {
		return f
	}
	return nil
}
