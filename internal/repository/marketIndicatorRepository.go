package repository

import (
	"context"
	"database/sql"
	"fmt"
	"teste-media-movel/internal/interfaces/repositories"
	"teste-media-movel/internal/models"
)

type MarketIndicatorRepository struct {
	db *sql.DB
}

func NewMarketIndicatorRepository(db *sql.DB) repositories.MarketIndicatorRepository {
	return &MarketIndicatorRepository{
		db: db,
	}
}

func (r *MarketIndicatorRepository) Save(ctx context.Context, indicator models.MarketIndicator) error {
	query := `
		INSERT INTO market_indicators (pair, timestamp_unix, mms_20, mms_50, mms_200)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		indicator.Pair,
		indicator.TimestampUnix,
		indicator.MMS20,
		indicator.MMS50,
		indicator.MMS200,
	)
	if err != nil {
		return fmt.Errorf("falha ao salvar indicador de mercado: %w", err)
	}

	return nil
}

func (r *MarketIndicatorRepository) GetByPairAndTimestamp(ctx context.Context, pair string, timestampUnix int64) (*models.MarketIndicator, error) {
	query := `
		SELECT id, pair, timestamp_unix, mms_20, mms_50, mms_200
		FROM market_indicators
		WHERE pair = ? AND timestamp_unix = ?
	`

	var indicator models.MarketIndicator
	err := r.db.QueryRowContext(ctx, query, pair, timestampUnix).Scan(
		&indicator.ID,
		&indicator.Pair,
		&indicator.TimestampUnix,
		&indicator.MMS20,
		&indicator.MMS50,
		&indicator.MMS200,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("falha ao buscar indicador de mercado: %w", err)
	}

	return &indicator, nil
}

func (r *MarketIndicatorRepository) GetByPairAndDateRange(ctx context.Context, pair string, from, to int64) ([]models.MarketIndicator, error) {
	query := `
		SELECT id, pair, timestamp_unix, mms_20, mms_50, mms_200
		FROM market_indicators
		WHERE pair = ? AND timestamp_unix >= ? AND timestamp_unix <= ?
		ORDER BY timestamp_unix ASC
	`

	rows, err := r.db.QueryContext(ctx, query, pair, from, to)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar indicadores por range: %w", err)
	}
	defer rows.Close()

	var indicators []models.MarketIndicator
	for rows.Next() {
		var indicator models.MarketIndicator
		if err := rows.Scan(&indicator.ID, &indicator.Pair, &indicator.TimestampUnix, &indicator.MMS20, &indicator.MMS50, &indicator.MMS200); err != nil {
			return nil, fmt.Errorf("falha ao escanear indicador: %w", err)
		}
		indicators = append(indicators, indicator)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das linhas: %w", err)
	}

	return indicators, nil
}
