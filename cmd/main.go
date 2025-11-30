package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "4757"
	}

	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "App is Running on Port : " + port,
		})
	})

	if err := r.Run(":" + port); err != nil {
		panic(err)
	}

}
