package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/moises-ba/mb-crypto-sma-api/internal/adapter"
	"github.com/moises-ba/mb-crypto-sma-api/internal/cache"
	"github.com/moises-ba/mb-crypto-sma-api/internal/controller"
	h "github.com/moises-ba/mb-crypto-sma-api/internal/http"
	"github.com/moises-ba/mb-crypto-sma-api/internal/service"
	"github.com/redis/go-redis/v9"
)

func main() {

	client := h.NewClient()
	exchangeAdapter := adapter.NewMercadoBitcoin(client)

	cache := cache.NewRedisCache(createRedisClient())

	lastClosedPricesSrv := service.NewCryptorCalculatorCachable(service.NewCryptoCalculatorService(exchangeAdapter), cache)
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

func createRedisClient() *redis.Client {
	// Pega o endereço das variáveis de ambiente (definidas no docker-compose)
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	return redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
}
