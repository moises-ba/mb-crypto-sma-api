package strategy

import (
	"github.com/moises-ba/mb-crypto-mms-api/internal/errors"
	"github.com/shopspring/decimal"
)

type AvgCalculator func(values []decimal.Decimal) (*decimal.Decimal, errors.ApiError)

func SimpleMovingAverage(values []decimal.Decimal) (*decimal.Decimal, errors.ApiError) {
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
