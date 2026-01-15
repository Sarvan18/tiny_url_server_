package helpers

import (
	"os"
	"strconv"
	"time"

	"github.com/Sarvan18/tiny_url_server_.git/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	bcryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bcryptedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJwt(user *models.User) (string, error) {
	userSigning := models.JwtClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.FormatUint(uint64(user.ID), 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userSigning)
	secret := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
