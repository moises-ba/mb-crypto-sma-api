package dto

import (
	"time"

	"github.com/moises-ba/mb-crypto-mms-api/internal/domain"
)

type AverageRequest struct {
	DigitalCoins []domain.DigitalCoin
	AvgType      domain.AvgType
	Days         int
	StartDate    time.Time
	EndDate      time.Time
}
