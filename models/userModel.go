package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  string `gorm:"size:100;not null" json:"name"`
	Email string `gorm:"size:150;uniqueIndex;not null" json:"email"`
	// Gender          string    `gorm:"size:10" json:"gender"`
	Password        string    `gorm:"size:255;not null" json:"password"`
	ConfirmPassword string    `gorm:"-" json:"confirmPassword"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UserLoginToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Token     string    `gorm:"size:255;not null" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type JwtClaims struct {
	ID    uint
	Name  string
	Email string
	jwt.RegisteredClaims
}
