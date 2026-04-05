package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/moises-ba/mb-crypto-mms-api/internal/domain"
	"github.com/moises-ba/mb-crypto-mms-api/internal/errors"
	"github.com/moises-ba/mb-crypto-mms-api/internal/http"
)

const (
	mercadobitcoinFQDN       = "https://api.mercadobitcoin.net/api/v4"
	mercadobitcoincandlesURL = "/candles"
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
	url := fmt.Sprintf("%s%s?symbol=%s-BRL&resolution=1d&from=%v&to=%v", mercadobitcoinFQDN, mercadobitcoincandlesURL,
		digitalCurrency, //poderiamos ter passado por parametro a moeda, mas o tempo foi curto pra fazer o refactory
		dtIni.Unix(), dtEnd.Unix())

	resp, apiErr := m.client.Get(ctx, url)
	if apiErr != nil {
		return nil, apiErr
	}

	var mercadoBitcoinResp mercadoBitCoinResponse
	if err := json.Unmarshal(resp.Response, &mercadoBitcoinResp); err != nil {
		return nil, errors.NewApiError("error on unmarshalling last closed prices", errors.WithError(err), errors.WithKind(errors.Unexpected))
	}

	return mercadoBitcoinResp.Closeds, nil
}
