package adapter

import (
	"context"
	"time"

	"github.com/moises-ba/mb-crypto-mms-api/internal/domain"
	"github.com/moises-ba/mb-crypto-mms-api/internal/errors"
)

type mercadoBitcoin struct {
}

func (m *mercadoBitcoin) ListLastPrices(ctx context.Context, digitalCurrency domain.DigitalCoin, referenceDate time.Time) ([]string, errors.ApiError) {
	//todo implemenar
	return nil, nil
}
