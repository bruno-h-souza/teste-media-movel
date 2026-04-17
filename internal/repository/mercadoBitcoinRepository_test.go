package repository

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMercadoBitcoinRepository_GetCandles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("symbol") != "BTC-BRL" {
			t.Errorf("esperava symbol BTC-BRL, recebido %s", r.URL.Query().Get("symbol"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"timestamp": [1620000000],
			"t": [1620000000],
			"o": ["100.0"],
			"h": ["110.0"],
			"l": ["90.0"],
			"c": ["105.5"],
			"v": ["10.0"]
		}`))
	}))
	defer server.Close()

	repo := &MercadoBitcoinRepository{
		client:  server.Client(),
		baseURL: server.URL,
	}

	candles, err := repo.GetCandles(context.Background(), "BTC-BRL", "1d", 1610000000, 1620000000)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if len(candles) != 1 {
		t.Fatalf("esperava 1 candle na resposta, obteve %d", len(candles))
	}

	if candles[0].Close != 105.5 {
		t.Errorf("esperava candle.Close igual a 105.5, obteve %f", candles[0].Close)
	}
}
