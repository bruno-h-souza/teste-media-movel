package services

import (
	"context"
	"teste-media-movel/internal/models"
)

type MarketIndicatorService interface {
	SyncIndicators(ctx context.Context) error
	GetIndicator(ctx context.Context, pair string, from, to int64, rangeStr string) ([]models.MMSResponse, error)
}
