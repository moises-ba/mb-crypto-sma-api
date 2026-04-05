package adapter

import (
	"context"
	"encoding/json"
	"time"

	"github.com/moises-ba/mb-crypto-mms-api/internal/domain"
	"github.com/moises-ba/mb-crypto-mms-api/internal/errors"
	"github.com/moises-ba/mb-crypto-mms-api/internal/http"
)

const (
	maercadobitcoinFQDN       = "https://api.mercadobitcoin.net/api/v4"
	maercadobitcoincandlesURL = "/candles"
)

type mercadoBitCoinResponse struct {
	Closeds []string `json:"c"` // close
}

type MercadoBitcoin interface {
	ListLastClosedPrices(ctx context.Context, digitalCurrency domain.DigitalCoin, dtIni, dtEnd time.Time) ([]string, errors.ApiError)
}

type mercadoBitcoin struct {
	client http.Client
}

func NewMercadoBitcoin(client http.Client) MercadoBitcoin {
	return &mercadoBitcoin{client: client}
}

func (m *mercadoBitcoin) ListLastClosedPrices(ctx context.Context, digitalCurrency domain.DigitalCoin, dtIni, dtEnd time.Time) ([]string, errors.ApiError) {
	resp, apiErr := m.client.Get(ctx, maercadobitcoinFQDN+maercadobitcoincandlesURL)
	if apiErr != nil {
		return nil, apiErr
	}

	var mercadoBitcoinResp mercadoBitCoinResponse
	if err := json.Unmarshal(resp.Response, &mercadoBitcoinResp); err != nil {
		return nil, errors.NewApiError("error on unmarshalling last closed prices", errors.WithError(err), errors.WithKind(errors.Unexpected))
	}

	return mercadoBitcoinResp.Closeds, nil
}
