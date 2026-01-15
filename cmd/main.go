package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Sarvan18/tiny_url_server_.git/config"
	userroutes "github.com/Sarvan18/tiny_url_server_.git/router/userRoutes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "4757"
	}

	if err := config.LoadConfig(); err != nil {
		log.Fatal("Could not load config: ", err)
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

	userroutes.UserRoutes(r)

	if err := r.Run(":" + port); err != nil {
		panic(err)
	}

}
