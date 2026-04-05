package timeutil

import (
	"time"

	"github.com/moises-ba/mb-crypto-sma-api/internal/errors"
)

const YearMonthDayPattern = "2006-01-02"

func Parse(dateStr, datePattern string) (*time.Time, errors.ApiError) {
	t, err := time.Parse(datePattern, dateStr)
	if err != nil {
		return nil, errors.NewApiError("invalid date",
			errors.WithKind(errors.Invalid), errors.WithError(err))
	}

	return &t, nil
}
