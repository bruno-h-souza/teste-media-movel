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
	r.GET("/:pair/mms", h.GetMMS)
}

func (h *IndicatorHandler) GetMMS(c *gin.Context) {
	pair := c.Param("pair")
	toStr := c.Query("to")
	fromStr := c.Query("from")
	rangeStr := c.Query("range")

	if pair == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Os parâmetros 'pair' é obrigatório"})
		return
	}

	if toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Os parâmetros 'toStr' é obrigatório"})
		return
	}

	if fromStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Os parâmetros 'fromStr' é obrigatório"})
		return
	}

	if rangeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Os parâmetros 'rangeStr' é obrigatório"})
		return
	}

	toTimestamp, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'timestamp' inválido"})
		return
	}

	fromTimestamp, err := strconv.ParseInt(fromStr, 10, 64)
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
