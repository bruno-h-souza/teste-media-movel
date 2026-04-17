package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"teste-media-movel/internal/interfaces/repositories"
	"teste-media-movel/internal/models"
	"time"
)

type MercadoBitcoinRepository struct {
	client  *http.Client
	baseURL string
}

func NewMercadoBitcoinRepository(client *http.Client) repositories.MercadoBitcoinRepository {
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	return &MercadoBitcoinRepository{
		client:  client,
		baseURL: "https://api.mercadobitcoin.net/api/v4",
	}
}

// GetCandles busca o histórico de barras (candles) para um par, resolução e período específicos.
func (r *MercadoBitcoinRepository) GetCandles(ctx context.Context, symbol, resolution string, from, to int64) ([]models.Candle, error) {
	url := fmt.Sprintf("%s/candles?symbol=%s&resolution=%s&from=%d&to=%d", r.baseURL, symbol, resolution, from, to)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar request: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha ao executar request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("a API retornou status de erro: %d", resp.StatusCode)
	}

	var apiResp models.MercadoBitcoinResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("falha ao decodificar JSON: %w", err)
	}

	return mapToCandles(apiResp)
}

// mapToCandles converte as listas orientadas a colunas da API para uma lista de objetos do tipo Candle.
func mapToCandles(apiResp models.MercadoBitcoinResponse) ([]models.Candle, error) {
	length := len(apiResp.Time)

	// Validação de segurança para garantir que todos os arrays tenham o mesmo tamanho
	if len(apiResp.Close) != length || len(apiResp.High) != length ||
		len(apiResp.Low) != length || len(apiResp.Open) != length ||
		len(apiResp.Volume) != length {
		return nil, fmt.Errorf("tamanhos de array divergentes na resposta da API")
	}

	candles := make([]models.Candle, 0, length)

	// Percorremos os arrays simultaneamente, transformando strings para float64
	for i := 0; i < length; i++ {
		open, _ := strconv.ParseFloat(apiResp.Open[i], 64)
		high, _ := strconv.ParseFloat(apiResp.High[i], 64)
		low, _ := strconv.ParseFloat(apiResp.Low[i], 64)
		closePrice, _ := strconv.ParseFloat(apiResp.Close[i], 64)
		volume, _ := strconv.ParseFloat(apiResp.Volume[i], 64)

		candles = append(candles, models.Candle{
			Timestamp: apiResp.Time[i],
			Open:      open,
			High:      high,
			Low:       low,
			Close:     closePrice,
			Volume:    volume,
		})
	}

	return candles, nil
}
