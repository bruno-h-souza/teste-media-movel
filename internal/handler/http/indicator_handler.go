package httphandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"teste-media-movel/internal/interfaces/services"
)

type IndicatorHandler struct {
	svc services.MarketIndicatorService
}

func NewIndicatorHandler(svc services.MarketIndicatorService) *IndicatorHandler {
	return &IndicatorHandler{svc: svc}
}

func (h *IndicatorHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/api/indicator", h.GetIndicator)
}

func (h *IndicatorHandler) GetIndicator(c *gin.Context) {
	pair := c.Query("pair")
	timestampStr := c.Query("timestamp")

	if pair == "" || timestampStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Os parâmetros 'pair' e 'timestamp' são obrigatórios"})
		return
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'timestamp' inválido"})
		return
	}

	indicator, err := h.svc.GetIndicator(c.Request.Context(), pair, timestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao buscar indicador"})
		return
	}
	if indicator == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Indicador não encontrado"})
		return
	}
	c.JSON(http.StatusOK, indicator)
}
