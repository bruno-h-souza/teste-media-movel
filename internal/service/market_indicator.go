package service

import (
	"context"
	"log"
	"time"

	"teste-media-movel/internal/interfaces/repositories"
	"teste-media-movel/internal/interfaces/services"
	"teste-media-movel/internal/models"
)

var Symbols = []string{"BTC-BRL", "ETH-BRL"}

const Resolution string = "1d"

type marketIndicatorService struct {
	mbRepo repositories.MercadoBitcoinRepository
	miRepo repositories.MarketIndicatorRepository
}

func NewMarketIndicatorService(mbRepo repositories.MercadoBitcoinRepository, miRepo repositories.MarketIndicatorRepository) services.MarketIndicatorService {
	return &marketIndicatorService{
		mbRepo: mbRepo,
		miRepo: miRepo,
	}
}

func (s *marketIndicatorService) SyncIndicators(ctx context.Context) error {
	log.Println("Iniciando execução do Job de sincronização...")
	now := time.Now()

	// Define data limte de 1 ano atrás
	endTarget := now.AddDate(-1, 0, 0)

	for _, symbol := range Symbols {
		log.Printf("Buscando candles para o par %s...", symbol)
		salvos := 0

		// Retroage mês a mês a partir da data atual até 1 ano atrás
		currentEnd := now
		for currentEnd.After(endTarget) {
			currentStart := currentEnd.AddDate(0, -1, 0)
			if currentStart.Before(endTarget) {
				currentStart = endTarget
			}

			from := currentStart.Unix()
			to := currentEnd.Unix()

			log.Printf("  Buscando dados de %s até %s...", currentStart.Format(time.DateOnly), currentEnd.Format(time.DateOnly))
			candles, err := s.mbRepo.GetCandles(ctx, symbol, Resolution, from, to)
			if err != nil {
				log.Printf("  Falha ao buscar candles para %s neste período: %v", symbol, err)
				return err
			} else {
				log.Printf("  %d registros encontrados. Iniciando salvamento...", len(candles))
				for _, candle := range candles {
					existing, err := s.miRepo.GetByPairAndTimestamp(ctx, symbol, candle.Timestamp)
					if err != nil {
						log.Printf("Erro ao verificar duplicidade para %s no timestamp %d: %v", symbol, candle.Timestamp, err)
						continue
					}

					if existing == nil {
						mms20, err := s.calculateMMS(ctx, symbol, candle.Timestamp, 20)
						if err != nil {
							log.Printf("Erro ao calcular MMS20 para %s no timestamp %d: %v", symbol, candle.Timestamp, err)
						}

						mms50, err := s.calculateMMS(ctx, symbol, candle.Timestamp, 50)
						if err != nil {
							log.Printf("Erro ao calcular MMS50 para %s no timestamp %d: %v", symbol, candle.Timestamp, err)
						}

						mms200, err := s.calculateMMS(ctx, symbol, candle.Timestamp, 200)
						if err != nil {
							log.Printf("Erro ao calcular MMS200 para %s no timestamp %d: %v", symbol, candle.Timestamp, err)
						}

						indicator := models.MarketIndicator{
							Pair:          symbol,
							TimestampUnix: candle.Timestamp,
							MMS20:         mms20,
							MMS50:         mms50,
							MMS200:        mms200,
						}
						if err := s.miRepo.Save(ctx, indicator); err == nil {
							salvos++
						}
					}
				}
			}
			currentEnd = currentStart
		}
		log.Printf("Sincronização para %s finalizada. %d novos registros salvos.", symbol, salvos)
	}
	log.Println("Job finalizado com sucesso!")
	return nil
}

// calculateMMS busca os últimos dias de um par e calcula a Média Móvel Simples (MMS)
func (s *marketIndicatorService) calculateMMS(ctx context.Context, symbol string, currentTimestamp int64, days int) (*float64, error) {
	// Data de início: dias antes da data atual
	from := time.Unix(currentTimestamp, 0).AddDate(0, 0, -days).Unix()

	var candles []models.Candle
	var err error
	maxRetries := 3

	// Busca os registros no intervalo com retentativas
	for attempt := 1; attempt <= maxRetries; attempt++ {
		candles, err = s.mbRepo.GetCandles(ctx, symbol, Resolution, from, currentTimestamp)
		if err == nil {
			break
		}
		log.Printf("Falha ao buscar candles para MMS%d de %s (tentativa %d/%d): %v", days, symbol, attempt, maxRetries, err)
		if attempt < maxRetries {
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		return nil, err
	}

	if len(candles) == 0 {
		return nil, nil
	}

	var sum float64
	for _, c := range candles {
		sum += c.Close
	}

	mms := sum / float64(len(candles))
	return &mms, nil
}

func (s *marketIndicatorService) GetIndicator(ctx context.Context, pair string, from, to int64, rangeStr string) ([]models.MMSResponse, error) {
	indicators, err := s.miRepo.GetByPairAndDateRange(ctx, pair, from, to)
	if err != nil {
		return nil, err
	}

	var result []models.MMSResponse
	for _, ind := range indicators {
		var mmsValue *float64
		switch rangeStr {
		case "20":
			mmsValue = ind.MMS20
		case "50":
			mmsValue = ind.MMS50
		case "200":
			mmsValue = ind.MMS200
		}
		result = append(result, models.MMSResponse{
			Timestamp: ind.TimestampUnix,
			MMS:       mmsValue,
		})
	}
	return result, nil
}
