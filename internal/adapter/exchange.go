package adapter

import (
	"context"
	"time"

	"github.com/moises-ba/mb-crypto-sma-api/internal/domain"
	"github.com/moises-ba/mb-crypto-sma-api/internal/errors"
)

type Exchange interface {
	ListLastClosedPrices(ctx context.Context, digitalCurrency domain.DigitalCoin, dateIni, dateTo time.Time) ([]string, errors.ApiError)
}
