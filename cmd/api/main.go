package main

import (
	"context"
	"flag"
	"log"

	"github.com/gin-gonic/gin"

	"teste-media-movel/configs"
	httphandler "teste-media-movel/internal/handler/http"
	"teste-media-movel/internal/interfaces/services"
	"teste-media-movel/internal/repository"
	"teste-media-movel/internal/service"
	"teste-media-movel/internal/utils"
)

func main() {
	// Flag para definir o modo de execução da aplicação
	mode := flag.String("mode", "api", "Modo de execução: 'api' ou 'job'")
	flag.Parse()

	db, err := configs.NewMySqlConnection()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Inicializando os repositórios
	mbRepo := repository.NewMercadoBitcoinRepository(nil)
	miRepo := repository.NewMarketIndicatorRepository(db)

	// Inicializando o serviço
	miService := service.NewMarketIndicatorService(mbRepo, miRepo)

	// Roteamento de execução baseado na flag
	if *mode == "job" {
		runJob(miService)
	} else {
		runAPI(miService)
	}
}

// runJob executa a rotina de buscar dados da API do Mercado Bitcoin e salvá-los no banco
func runJob(svc services.MarketIndicatorService) {
	ctx := context.Background()
	if err := svc.SyncIndicators(ctx); err != nil {
		log.Printf("Erro na execução do Job: %v", err)
	}
}

// runAPI inicia um servidor web simples para consultar os dados
func runAPI(svc services.MarketIndicatorService) {
	router := gin.Default()

	// Inicializa e registra os endpoints do handler de indicadores
	handler := httphandler.NewIndicatorHandler(svc)
	handler.RegisterRoutes(router)

	port := utils.GetEnv("PORT", "8080")
	log.Printf("Servidor API iniciado e escutando na porta %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
