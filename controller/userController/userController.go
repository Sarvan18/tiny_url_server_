package usercontroller

import (
	"net/http"
	"strconv"

	usermiddleware "github.com/Sarvan18/tiny_url_server_.git/middleware/userMiddleware"
	"github.com/Sarvan18/tiny_url_server_.git/models"
	userservices "github.com/Sarvan18/tiny_url_server_.git/services/userServices"
	"github.com/gin-gonic/gin"
)

func RegisterUserController() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := userservices.UserRegisterService(usermiddleware.User)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, res)
	}
}

func LoginUserController() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := userservices.UserLoginService(usermiddleware.User.Email, usermiddleware.User.Password)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, res)
	}

}

func UpdateUserController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, ok := ctx.Get("climerEmail")

		if !ok {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Error While Getting Jwt Claims Data"})
			return
		}

		userId := ctx.Param("id")

		userIDUint, err := strconv.ParseUint(userId, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		userValue, ok := ctx.Get("user")

		if !ok {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Error While Getting User Data"})
			return
		}

		user, ok := userValue.(*models.User)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Error While Converting User Data"})
			return
		}

		res, err := userservices.UpdateUserService(uint(userIDUint), user)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(200, res)
	}
}

func DeleteUserController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("id")

		userIDUint, err := strconv.ParseUint(userId, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		_, ok := ctx.Get("climerEmail")

		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Error while getting user claims",
			})
			return
		}

		err = userservices.DeleteUserService(uint(userIDUint))
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "Error deleting user",
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message": "User deleted successfully",
		})
	}
}

func GetUserHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("id")
		userIDUint, err := strconv.ParseUint(userId, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "Invalid user ID" + err.Error(),
			})
			return
		}

		user, err := userservices.GetUserByIDService(uint(userIDUint))

		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message": "Okay",
			"user":    user,
		})
	}
}
