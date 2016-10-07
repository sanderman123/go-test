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
		users.GET("/:name", func(c *gin.Context) {
			name := c.Param("name")
			response := controller.FindUser(name)
			c.JSON(response.Status, response.Body)
		})
		users.POST("", func(c *gin.Context) {
			EntityRequestHander(c, controller.CreateUser, model.UserFactory{})
		})
		users.PUT("", func(c *gin.Context) {
			EntityHandler(c, controller.UpdateUser, userFactory)
		})
		users.DELETE("/:user-name", func(c *gin.Context) {
			PathParameterHander(c, controller.DeleteUser, USER_NAME)
		})
		users.POST("/login", func(c *gin.Context) {
			EntityHandler(c, controller.AuthenticateUser, userFactory)
		})
		users.POST("/activate/:token", func(c *gin.Context) {
			PathParameterHander(c, controller.ActivateUser, TOKEN)
		})
		users.POST("/forgot", func(c *gin.Context) {
			EntityRequestHander(c, controller.ForgotPassword, userFactory)
		})
		users.POST("/reset/:token", func(c *gin.Context) {
			PathParameterStringHander(c, controller.ResetPassword, TOKEN)
		})
	}
}

