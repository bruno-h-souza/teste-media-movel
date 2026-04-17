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
						indicator := models.MarketIndicator{
							Pair:          symbol,
							TimestampUnix: candle.Timestamp,
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

func (s *marketIndicatorService) GetIndicator(ctx context.Context, pair string, timestampUnix int64) (*models.MarketIndicator, error) {
	return s.miRepo.GetByPairAndTimestamp(ctx, pair, timestampUnix)
}
