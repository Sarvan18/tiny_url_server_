package userroutes

import (
	usercontroller "github.com/Sarvan18/tiny_url_server_.git/controller/userController"
	usermiddleware "github.com/Sarvan18/tiny_url_server_.git/middleware/userMiddleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(_router *gin.Engine) {
	api := _router.Group("/api/users")
	{
		api.GET("/:id", usercontroller.GetUserHandler())
		api.POST("/register", usermiddleware.UserRegisterMiddleware(), usercontroller.RegisterUserController())
		api.POST("/login", usermiddleware.UserLoginMiddleware(), usercontroller.LoginUserController())
		api.POST("/update/:id", usermiddleware.UpdateUserMiddleware(), usercontroller.UpdateUserController())
		api.DELETE("/delete/:id", usermiddleware.DeleteUserMiddleware(), usercontroller.DeleteUserController())
	}
}
