package ginrouter

import (
	"github.com/gin-gonic/gin"
	"github.com/sanderman123/user-service/controller"
	"github.com/sanderman123/user-service/model"
)

const USER_NAME = "user-name"
const TOKEN = "token"

func Init(r *gin.Engine) {
	userFactory := model.UserFactory{}

	users := r.Group("/users")
	{
		users.GET("/:" + USER_NAME, func(c *gin.Context) {
			AuthenticatedPathParameterHandler(c, controller.FindUser, USER_NAME)
		})
		users.POST("", func(c *gin.Context) {
			EntityRequestHandler(c, controller.CreateUser, model.UserFactory{})
		})
		users.PUT("", func(c *gin.Context) {
			AuthenticatedEntityHandler(c, controller.UpdateUser, userFactory)
		})
		users.DELETE("/:" + USER_NAME, func(c *gin.Context) {
			AuthenticatedPathParameterHandler(c, controller.DeleteUser, USER_NAME)
		})
		users.POST("/login", func(c *gin.Context) {
			EntityResponseWriterHandler(c, controller.AuthenticateUser, userFactory)
		})
		users.POST("/activate/:" + TOKEN, func(c *gin.Context) {
			PathParameterHandler(c, controller.ActivateUser, TOKEN)
		})
		users.POST("/forgot", func(c *gin.Context) {
			EntityRequestHandler(c, controller.ForgotPassword, userFactory)
		})
		users.POST("/reset/:" + TOKEN, func(c *gin.Context) {
			PathParameterStringHandler(c, controller.ResetPassword, TOKEN)
		})
	}
}

