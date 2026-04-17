package repositories

import (
	"context"
	"teste-media-movel/internal/models"
)

type MercadoBitcoinRepository interface {
	GetCandles(ctx context.Context, symbol, resolution string, from, to int64) ([]models.Candle, error)
}

type MarketIndicatorRepository interface {
	Save(ctx context.Context, indicator models.MarketIndicator) error
	GetByPairAndTimestamp(ctx context.Context, pair string, timestampUnix int64) (*models.MarketIndicator, error)
	GetByPairAndDateRange(ctx context.Context, pair string, from, to int64) ([]models.MarketIndicator, error)
}
