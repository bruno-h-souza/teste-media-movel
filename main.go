package main

import (
	"context"
	"flag"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"teste-media-movel/configs"
	httphandler "teste-media-movel/internal/handler/http"
	"teste-media-movel/internal/interfaces/services"
	"teste-media-movel/internal/repository"
	"teste-media-movel/internal/service"
	"teste-media-movel/internal/utils"

	_ "teste-media-movel/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           API de Médias Móveis (MMS)
// @version         1.0
// @description     API para consulta de Médias Móveis Simples de criptomoedas (BTC e ETH).
// @host            localhost:8080
// @BasePath        /
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

	maxRetriesStr := utils.GetEnv("JOB_MAX_RETRIES", "3")
	timeoutSecondsStr := utils.GetEnv("JOB_RETRY_TIMEOUT_SECONDS", "2")

	maxRetries, err := strconv.Atoi(maxRetriesStr)
	if err != nil {
		maxRetries = 3
	}

	timeoutSeconds, err := strconv.Atoi(timeoutSecondsStr)
	if err != nil {
		timeoutSeconds = 2
	}

	for attempt := 1; attempt <= maxRetries+1; attempt++ {
		if err := svc.SyncIndicators(ctx); err != nil {
			log.Printf("Erro na execução do Job (tentativa %d/%d): %v", attempt, maxRetries+1, err)
			if attempt <= maxRetries {
				log.Printf("Aguardando timeout de %d segundos antes da próxima tentativa...", timeoutSeconds)
				time.Sleep(time.Duration(timeoutSeconds) * time.Second)
			}
		} else {
			return // Executado com sucesso, o loop encerra
		}
	}

	log.Printf("O Job falhou definitivamente após %d tentativas.", maxRetries+1)
}

// runAPI inicia um servidor web simples para consultar os dados
func runAPI(svc services.MarketIndicatorService) {
	router := gin.Default()

	// Inicializa e registra os endpoints do handler de indicadores
	handler := httphandler.NewIndicatorHandler(svc)
	handler.RegisterRoutes(router)

	// Rota do Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := utils.GetEnv("PORT", "8080")
	log.Printf("Servidor API iniciado e escutando na porta %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
