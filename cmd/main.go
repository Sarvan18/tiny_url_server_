package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Sarvan18/tiny_url_server_.git/config"
	userroutes "github.com/Sarvan18/tiny_url_server_.git/router/userRoutes"
	"github.com/gin-contrib/cors"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Cache-Control"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
