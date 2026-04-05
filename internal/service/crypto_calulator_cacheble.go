package service

import (
	"context"
	"fmt"
	"time"

	"github.com/moises-ba/mb-crypto-mms-api/internal/cache"
	"github.com/moises-ba/mb-crypto-mms-api/internal/domain"
	"github.com/moises-ba/mb-crypto-mms-api/internal/errors"
	"github.com/moises-ba/mb-crypto-mms-api/internal/log"
)

// Decorator que faz cache dos resultados
type cryptorCalculatorCachable struct {
	CryptoCalculatorService
	c cache.Cache
}

func NewCryptorCalculatorCachable(srv CryptoCalculatorService, c cache.Cache) CryptoCalculatorService {
	return &cryptorCalculatorCachable{
		CryptoCalculatorService: srv,
		c:                       c,
	}
}

func (s *cryptorCalculatorCachable) GetLastAVGPrices(ctx context.Context, qtDays int, referenceDate time.Time) ([]domain.AverageResponse, errors.ApiError) {
	key := fmt.Sprintf("%s_%v", referenceDate.Format("01-01-2026"), qtDays)

	var response []domain.AverageResponse
	ok, err := s.c.Get(ctx, key, &response)
	if err != nil {
		log.Error("erro trying get from cache", err)
		//ignora o erro para poder ir pelo fluxo normal
	}

	if ok {
		log.Info("obtido do cache")
		return response, nil
	}

	response, err = s.CryptoCalculatorService.GetLastAVGPrices(ctx, qtDays, referenceDate)
	if err != nil {
		return nil, err
	}

	if len(response) > 0 {
		if errCache := s.c.Set(ctx, key, response, 24*time.Hour); errCache != nil { //aqui podemos configurar o ttl dependendo das dadas mais acessadas, considerei aqui um dia.
			log.Error("erro trying set on cache", err)
			//ignora o erro
		}
	}

	return response, err
}
