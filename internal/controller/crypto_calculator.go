package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moises-ba/mb-crypto-mms-api/internal/service"
)

type CriptorCalculatorController interface {
	ListLastClosedPrices(c *gin.Context)
}

type criptorCalculatorController struct {
	srv service.CryptoCalculatorService
}

func NewCriptorCalculatorController(srv service.CryptoCalculatorService) CriptorCalculatorController {
	return &criptorCalculatorController{srv: srv}
}

// ListLastClosedPrices godoc
// @Summary Lista os últimos preços fechados
// @Description Retorna os últimos preços fechados com base na quantidade de dias e data de referência
// @Tags LastClosedPrices
// @Accept json
// @Produce json
// @Param days path int true "Quantidade de dias"
// @Param dt_reference path string true "Data de referência (formato: YYYY-MM-DD)"
// @Success 200 {object} {}
// @Failure 400 {object} ApiError
// @Failure 500 {object} ApiError
// @Router /v1/lastclosedprices/{days}/{dt_reference} [get]
func (ctrl *criptorCalculatorController) ListLastClosedPrices(c *gin.Context) {

	days, err := strconv.ParseInt(c.Param("days"), 10, 32)
	if err != nil {
		///
	}

	dtReference, err := time.Parse("2006-01-02", c.Param("dt_reference"))
	if err != nil {
		///
	}

	res, apiErr := ctrl.srv.GetLastAVGPrices(c, int(days), dtReference)
	if err != nil {
		///
	}

	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":     id,
		"nome":   "Teclado Mecânico",
		"status": "disponível",
	})
}
