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

// GetMMS godoc
// @Summary      Retorna a Média Móvel Simples (MMS)
// @Description  Calcula e retorna a MMS (20, 50 ou 200) para um par de criptomoeda específico dentro de um período.
// @Tags         mms
// @Produce      json
// @Param        pair   path      string  true  "Par de moedas (BRLBTC ou BRLETH)"
// @Param        from   query     int     true  "Timestamp inicial (Unix)"
// @Param        to     query     int     true  "Timestamp final (Unix)"
// @Param        range  query     int     true  "Período da MMS (20, 50, 200)"
// @Success      200    {array}   models.MMSResponse
// @Failure      400    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /{pair}/mms [get]
func (h *IndicatorHandler) GetMMS(c *gin.Context) {
	pairStr := c.Param("pair")
	toStr := c.Query("to")
	fromStr := c.Query("from")
	rangeStr := c.Query("range")

	if pairStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Os parâmetros 'pair' é obrigatório"})
		return
	}

	if pairStr != "BRLBTC" && pairStr != "BRLETH" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'pair' inválido"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'to' inválido"})
		return
	}

	fromTimestamp, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'from' inválido"})
		return
	}

	var pair string
	switch pairStr {
	case "BRLBTC":
		pair = "BTC-BRL"
	case "BRLETH":
		pair = "ETH-BRL"
	}

	indicators, err := h.svc.GetIndicator(c.Request.Context(), pair, fromTimestamp, toTimestamp, rangeStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao buscar indicador"})
		return
	}
	if len(indicators) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nenhum indicador encontrado para o período"})
		return
	}
	c.JSON(http.StatusOK, indicators)
}
