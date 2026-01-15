package config

import (
	"fmt"
	"os"

	"github.com/Sarvan18/tiny_url_server_.git/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadConfig() error {
	dbUri := os.Getenv("POSTGRESS")

	if dbUri == "" {
		return fmt.Errorf("DB URI NOT SET")
	}

	db, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{})

	if err != nil {
		// panic(err)
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	DB = db
	fmt.Println("DB Connected Successfully")

	DB.AutoMigrate(&models.User{}, &models.UserLoginToken{})

	return nil
}
