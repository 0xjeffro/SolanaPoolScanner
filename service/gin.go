package service

import (
	"SolanaPoolScanner/workers"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func StartGin() {
	port := "8080"
	if os.Getenv("SERVICE_PORT") != "" {
		port = os.Getenv("SERVICE_PORT")
	}
	debug := false
	if os.Getenv("DEBUG") == "true" {
		debug = true
	}

	router := gin.New()
	if debug {
		router.Use(gin.Logger())
	}
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, workers.WorkerStatus)
	})

	err := router.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
