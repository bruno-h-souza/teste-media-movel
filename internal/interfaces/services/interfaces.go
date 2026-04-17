package services

import (
	"context"
	"teste-media-movel/internal/models"
)

type MarketIndicatorService interface {
	SyncIndicators(ctx context.Context) error
	GetIndicator(ctx context.Context, pair string, timestampUnix int64) (*models.MarketIndicator, error)
}
