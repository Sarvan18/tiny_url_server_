package userservices

import (
	"fmt"
	"time"

	"github.com/Sarvan18/tiny_url_server_.git/config"
	"github.com/Sarvan18/tiny_url_server_.git/helpers"
	"github.com/Sarvan18/tiny_url_server_.git/models"
)

func UserRegisterService(user *models.User) (*models.User, error) {
	var count int64

	if err := config.DB.Model(&models.User{}).Where("name = ?", user.Name).Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, fmt.Errorf("User With Same Name Already Exists")
	}

	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	user.ConfirmPassword = ""

	if err := config.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func UserLoginService(email, password string) (*models.UserLoginToken, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("Invalid Email or Password")
	}

	if !helpers.CheckPasswordHash(password, user.Password) {
		return nil, fmt.Errorf("Invalid Email or Password")
	}

	jwtToken, err := helpers.GenerateJwt(&user)
	if err != nil {
		return nil, fmt.Errorf("Error while Generating JWT\n Error : %v", err)
	}
	userLoginToken := &models.UserLoginToken{
		Token:     jwtToken,
		ExpiresAt: time.Now().Add(time.Minute * 30),
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(userLoginToken).Error; err != nil {
		return nil, err
	}

	return userLoginToken, nil
}

func GetUserByIDService(userID uint) (*models.User, error) {
	var user models.User

	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("User Not Found")
	}

	return &user, nil
}

func UpdateUserService(userID uint, updatedUser *models.User) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("User Not Found \n Error: %v", err)
	}

	updatedUser.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	if err := config.DB.Model(&user).Updates(updatedUser).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteUserService(userID uint) error {
	var user models.User

	if err := config.DB.First(&user, userID).Error; err != nil {
		return fmt.Errorf("User Not Found")
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("User Not Found")
	}

	return &user, nil
}
