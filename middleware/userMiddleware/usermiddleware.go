package usermiddleware

import (
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"unicode"

	"github.com/Sarvan18/tiny_url_server_.git/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var User *models.User

func UserRegisterMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")

		if ctx.Request.Method != "POST" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Method not allowed"})
		}

		defer ctx.Request.Body.Close()

		if err := ctx.Request.ParseForm(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse form", "message": err.Error()})
		}
		User = &models.User{
			Name:            ctx.Request.FormValue("name"),
			Email:           ctx.Request.FormValue("email"),
			Gender:          ctx.Request.FormValue("gender"),
			Password:        ctx.Request.FormValue("password"),
			ConfirmPassword: ctx.Request.FormValue("confirmPassword"),
		}

		if strings.TrimSpace(string(User.Name)) == "" || strings.TrimSpace(string(User.Email)) == "" || strings.TrimSpace(string(User.Password)) == "" || strings.TrimSpace(string(User.ConfirmPassword)) == "" {

			ctx.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
			return
		}

		_, err := mail.ParseAddress(User.Email)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			return
		}

		if strings.TrimSpace(string(User.Gender)) == "Male" || strings.TrimSpace(string(User.Gender)) == "Female" || strings.TrimSpace(string(User.Gender)) == "Other" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Gender"})
		}

		hasNum, hasUpper, hasSpecial := validatePassword(User.Password)

		if len(User.Password) <= 6 || !hasNum || !hasUpper || !hasSpecial {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long and include at least one number, one uppercase letter, and one special character"})
			return
		}

		if User.Password != User.ConfirmPassword {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password and Confirm Password do not match"})
			return
		}

		ctx.Next()
	}
}

func UserLoginMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")

		if ctx.Request.Method != "POST" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Method not allowed"})
			return
		}

		defer ctx.Request.Body.Close()

		if err := ctx.Request.ParseForm(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse form", "message": err.Error()})
			return
		}

		userEmail := ctx.Request.FormValue("email")
		userPassword := ctx.Request.FormValue("password")
		if strings.TrimSpace(string(userEmail)) == "" || strings.TrimSpace(string(userPassword)) == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email and Password are required"})
			return
		}

		_, err := mail.ParseAddress(userEmail)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			return
		}

		User = &models.User{
			Email:    userEmail,
			Password: userPassword,
		}

		ctx.Next()

	}
}

func GetUserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Id Not Found"})
			return
		}

		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Token Not Found"})
			return
		}
		tokenString = tokenString[7:]
		ctx.Set("BearerToken", tokenString)
		ctx.Set("id", id)
		ctx.Next()
	}
}

func UpdateUserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Token Not Found"})
			return
		}
		var user models.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Invalid request body"})
			return
		}

		climerEmail, err := Auth(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "message": "Unauthorized"})
			return
		}

		ctx.Set("climerEmail", climerEmail)
		ctx.Set("user", &user)
		ctx.Next()
	}
}

func DeleteUserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id := ctx.Param("id")

		if id == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Id Not Found"})
			return
		}

		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Token Not Found"})
			return
		}

		climerEmail, err := Auth(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Set("climerEmail", climerEmail)
		ctx.Set("id", id)
		ctx.Next()
	}
}

func Auth(token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("Missing or Invalid Authorization header")
	}

	token = token[7:]
	claims := models.JwtClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return "", err
	}

	if !jwtToken.Valid {
		return "", fmt.Errorf("Invalid Token")
	}
	return claims.Email, nil
}

func validatePassword(s string) (hasNum, hasUpper, hasSpecial bool) {
	for _, value := range s {
		switch {
		case unicode.IsNumber(value):
			hasNum = true
		case unicode.IsUpper(value):
			hasUpper = true

		case unicode.IsSymbol(value) || unicode.IsPunct(value):
			hasSpecial = true
		}

	}
	return hasNum, hasUpper, hasSpecial
}
