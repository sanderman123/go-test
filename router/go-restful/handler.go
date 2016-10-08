package gorestful

import (
	"github.com/emicklei/go-restful"
	"github.com/sanderman123/user-service/model"
	"net/http"
	"github.com/sanderman123/user-service/controller"
	"github.com/dgrijalva/jwt-go"
	"github.com/sanderman123/user-service/util"
)

type pathParameterController func(parameter string) controller.Response
type pathParameterStringController func(parameter string, string string) controller.Response
type entityController func(entity interface{}) controller.Response
type entityRequestController func(entity interface{}, request *http.Request) controller.Response
type entityResponseWriterController func(entity interface{}, request http.ResponseWriter) controller.Response

type authenticatedPathParameterController func(claims jwt.MapClaims, parameter string) controller.Response


//Generic Handlers
func PathParameterHander(request *restful.Request, response *restful.Response, fn pathParameterController, parameterName string) {
	result := fn(request.PathParameter(parameterName))
	response.WriteHeaderAndEntity(result.Status, result.Body)
}

func AuthenticatedPathParameterHandler(request *restful.Request, response *restful.Response, fn authenticatedPathParameterController, parameterName string) {
	claims, err := util.IsAuthenticated(request.Request)
	if err == nil {
		result := fn(claims, request.PathParameter(parameterName))
		response.WriteHeaderAndEntity(result.Status, result.Body)
	} else {
		response.WriteHeader(http.StatusUnauthorized)
	}
}

func PathParameterStringHander(request *restful.Request, response *restful.Response, fn pathParameterStringController, parameterName string) {
	str := ""
	err := request.ReadEntity(&str)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		result := fn(request.PathParameter(parameterName), str)
		response.WriteHeaderAndEntity(result.Status, result.Body);
	}
}

func EntityHander(request *restful.Request, response *restful.Response, fn entityController, factory interface{}) {
	b := factory.(model.Factory)
	entity := b.NewEntity()
	err := request.ReadEntity(&entity)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		result := fn(entity)
		response.WriteHeaderAndEntity(result.Status, result.Body);
	}
}

func EntityRequestHander(request *restful.Request, response *restful.Response, fn entityRequestController, factory interface{}) {
	b := factory.(model.Factory)
	entity := b.NewEntity()
	err := request.ReadEntity(&entity)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		result := fn(entity, request.Request)
		response.WriteHeaderAndEntity(result.Status, result.Body);
	}
}

func EntityResponseWriterHandler(request *restful.Request, response *restful.Response, fn entityResponseWriterController, factory interface{}) {
	b := factory.(model.Factory)
	entity := b.NewEntity()
	err := request.ReadEntity(&entity)
	if err == nil {
		result := fn(entity, response.ResponseWriter)
		response.WriteHeaderAndEntity(result.Status, result.Body);
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}
