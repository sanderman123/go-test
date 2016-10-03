package router

import (
	"github.com/emicklei/go-restful"
	"github.com/sanderman123/user-service/controller"
	"github.com/sanderman123/user-service/model"
)

const USER_NAME = "user-name"
const TOKEN = "token"

func Init(service *restful.WebService) {
	service.
	Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_XML, restful.MIME_JSON)

	service.Route(service.GET("/{" + USER_NAME + "}").To(FindUser))
	service.Route(service.POST("").To(CreateUser))
	service.Route(service.PUT("").To(UpdateUser))
	service.Route(service.DELETE("/{" + USER_NAME + "}").To(DeleteUser))
	service.Route(service.POST("/login").To(AuthenticateUser))
	service.Route(service.GET("/activate/{" + TOKEN + "}").To(ActivateUser))
	service.Route(service.POST("/forgot").To(ForgotPassword))
	service.Route(service.POST("/reset/{" + TOKEN + "}").To(ResetPassword))
}

func FindUser(request *restful.Request, rp *restful.Response) {
	PathParameterHander(request, rp, controller.FindUser, USER_NAME)
}

func CreateUser(request *restful.Request, response *restful.Response) {
	EntityRequestHander(request, response, controller.CreateUser, model.UserFactory{})
}

func UpdateUser(request *restful.Request, response *restful.Response) {
	EntityHander(request, response, controller.UpdateUser, model.UserFactory{})
}

func DeleteUser(request *restful.Request, response *restful.Response) {
	PathParameterHander(request, response, controller.DeleteUser, USER_NAME)
}

func AuthenticateUser(request *restful.Request, response *restful.Response) {
	EntityHander(request, response, controller.AuthenticateUser, model.UserFactory{})
}

func ActivateUser(request *restful.Request, response *restful.Response) {
	PathParameterHander(request, response, controller.ActivateUser, TOKEN)
}

func ForgotPassword(request *restful.Request, response *restful.Response) {
	EntityRequestHander(request, response, controller.ForgotPassword, model.UserFactory{})
}

func ResetPassword(request *restful.Request, response *restful.Response) {
	PathParameterStringHander(request, response, controller.ResetPassword, TOKEN)
}
