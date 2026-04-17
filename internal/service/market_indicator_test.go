package service

import (
	"context"
	"errors"
	"testing"

	"teste-media-movel/internal/models"
)

type mockMarketIndicatorRepo struct {
	MockIndicators            []models.MarketIndicator
	SaveFunc                  func(ctx context.Context, indicator models.MarketIndicator) error
	GetByPairAndTimestampFunc func(ctx context.Context, pair string, timestampUnix int64) (*models.MarketIndicator, error)
}

func (m *mockMarketIndicatorRepo) Save(ctx context.Context, indicator models.MarketIndicator) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, indicator)
	}
	return nil
}

func (m *mockMarketIndicatorRepo) GetByPairAndTimestamp(ctx context.Context, pair string, timestampUnix int64) (*models.MarketIndicator, error) {
	if m.GetByPairAndTimestampFunc != nil {
		return m.GetByPairAndTimestampFunc(ctx, pair, timestampUnix)
	}
	return nil, nil
}

func (m *mockMarketIndicatorRepo) GetByPairAndDateRange(ctx context.Context, pair string, from, to int64) ([]models.MarketIndicator, error) {
	return m.MockIndicators, nil
}

type mockMercadoBitcoinRepo struct {
	GetCandlesFunc func(ctx context.Context, symbol, resolution string, from, to int64) ([]models.Candle, error)
}

func (m *mockMercadoBitcoinRepo) GetCandles(ctx context.Context, symbol, resolution string, from, to int64) ([]models.Candle, error) {
	if m.GetCandlesFunc != nil {
		return m.GetCandlesFunc(ctx, symbol, resolution, from, to)
	}
	return nil, nil
}

func TestMarketIndicatorService_GetIndicator(t *testing.T) {
	mms20Value := 50000.50

	mockRepo := &mockMarketIndicatorRepo{
		MockIndicators: []models.MarketIndicator{
			{Pair: "BTC-BRL", TimestampUnix: 1620000000, MMS20: &mms20Value},
		},
	}

	svc := NewMarketIndicatorService(&mockMercadoBitcoinRepo{}, mockRepo)

	responses, err := svc.GetIndicator(context.Background(), "BTC-BRL", 1610000000, 1620000000, "20")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if len(responses) != 1 {
		t.Fatalf("esperava 1 resultado, obteve %d", len(responses))
	}

	if *responses[0].MMS != mms20Value {
		t.Errorf("esperava MMS igual a %f, obteve %f", mms20Value, *responses[0].MMS)
	}
}

func TestMarketIndicatorService_SyncIndicators_Success(t *testing.T) {
	mbRepo := &mockMercadoBitcoinRepo{
		GetCandlesFunc: func(ctx context.Context, symbol, resolution string, from, to int64) ([]models.Candle, error) {
			return []models.Candle{
				{Timestamp: 1620000000, Close: 100.0},
			}, nil
		},
	}

	saveCalled := false
	miRepo := &mockMarketIndicatorRepo{
		SaveFunc: func(ctx context.Context, indicator models.MarketIndicator) error {
			saveCalled = true
			return nil
		},
	}

	svc := NewMarketIndicatorService(mbRepo, miRepo)

	err := svc.SyncIndicators(context.Background())
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if !saveCalled {
		t.Errorf("esperava que o método Save fosse chamado ao menos uma vez")
	}
}

func TestMarketIndicatorService_SyncIndicators_Error(t *testing.T) {
	mbRepo := &mockMercadoBitcoinRepo{
		GetCandlesFunc: func(ctx context.Context, symbol, resolution string, from, to int64) ([]models.Candle, error) {
			return nil, errors.New("erro simulado na API do MB")
		},
	}

	miRepo := &mockMarketIndicatorRepo{}

	svc := NewMarketIndicatorService(mbRepo, miRepo)

	err := svc.SyncIndicators(context.Background())
	if err == nil {
		t.Fatalf("esperava um erro, mas obteve sucesso")
	}

	if err.Error() != "erro simulado na API do MB" {
		t.Errorf("esperava 'erro simulado na API do MB', obteve '%v'", err)
	}
}
