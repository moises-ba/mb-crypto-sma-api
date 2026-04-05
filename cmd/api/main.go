package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moises-ba/mb-crypto-mms-api/internal/adapter"
	"github.com/moises-ba/mb-crypto-mms-api/internal/controller"
	h "github.com/moises-ba/mb-crypto-mms-api/internal/http"
	"github.com/moises-ba/mb-crypto-mms-api/internal/service"
)

func main() {
	client := h.NewClient()
	exchangeAdapter := adapter.NewMercadoBitcoin(client)
	lastClosedPricesSrv := service.NewCryptoCalculatorService(exchangeAdapter)
	lastClosedPricesCtrl := controller.NewCriptorCalculatorController(lastClosedPricesSrv)

	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("/ping", func(c *gin.Context) { //para health check
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		v1.GET("/lastclosedprices/:days/:dt_reference", lastClosedPricesCtrl.ListLastClosedPrices) //escolhido path para porque os parametros sao obrigatorios data no formato YYYY-MM-DD
	}

	r.Run(":8080")
}
