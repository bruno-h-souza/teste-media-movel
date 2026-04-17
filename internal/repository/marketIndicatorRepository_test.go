package repository

import (
	"context"
	"testing"

	"teste-media-movel/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMarketIndicatorRepository_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewMarketIndicatorRepository(db)

	mms20 := 150000.50
	indicator := models.MarketIndicator{
		Pair:          "BTC-BRL",
		TimestampUnix: 1620000000,
		MMS20:         &mms20,
	}

	mock.ExpectExec("INSERT INTO market_indicators").
		WithArgs(indicator.Pair, indicator.TimestampUnix, indicator.MMS20, indicator.MMS50, indicator.MMS200).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Save(context.Background(), indicator)
	if err != nil {
		t.Errorf("erro inesperado ao salvar: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas de query não foram atendidas: %v", err)
	}
}

func TestMarketIndicatorRepository_GetByPairAndDateRange(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewMarketIndicatorRepository(db)
	mms20 := 100.0

	rows := sqlmock.NewRows([]string{"id", "pair", "timestamp_unix", "mms_20", "mms_50", "mms_200"}).
		AddRow(1, "BTC-BRL", 1620000000, &mms20, nil, nil).
		AddRow(2, "BTC-BRL", 1620086400, &mms20, nil, nil)

	mock.ExpectQuery("SELECT id, pair, timestamp_unix, mms_20, mms_50, mms_200 FROM market_indicators").
		WithArgs("BTC-BRL", int64(1620000000), int64(1620086400)).
		WillReturnRows(rows)

	indicators, err := repo.GetByPairAndDateRange(context.Background(), "BTC-BRL", 1620000000, 1620086400)
	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}

	if len(indicators) != 2 {
		t.Errorf("esperava 2 indicadores, obteve %d", len(indicators))
	}

	if indicators[0].ID != 1 || indicators[0].Pair != "BTC-BRL" {
		t.Errorf("dados do primeiro indicador incorretos")
	}

	if indicators[0].MMS20 == nil || *indicators[0].MMS20 != 100.0 {
		t.Errorf("valor de MMS20 incorreto")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas de query não foram atendidas: %v", err)
	}
}
