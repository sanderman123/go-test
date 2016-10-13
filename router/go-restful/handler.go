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
type authenticatedEntityController func(claims jwt.MapClaims, entity interface{}) controller.Response

//Generic Handlers
func PathParameterHandler(request *restful.Request, response *restful.Response, fn pathParameterController, parameterName string) {
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

func PathParameterStringHandler(request *restful.Request, response *restful.Response, fn pathParameterStringController, parameterName string) {
	str := ""
	err := request.ReadEntity(&str)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		result := fn(request.PathParameter(parameterName), str)
		response.WriteHeaderAndEntity(result.Status, result.Body);
	}
}

func EntityHandler(request *restful.Request, response *restful.Response, fn entityController, factory interface{}) {
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

func AuthenticatedEntityHandler(request *restful.Request, response *restful.Response, fn authenticatedEntityController, factory interface{}) {
	claims, err := util.IsAuthenticated(request.Request)
	if err == nil {
		b := factory.(model.Factory)
		entity := b.NewEntity()
		err := request.ReadEntity(&entity)
		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
		} else {
			result := fn(claims, entity)
			response.WriteHeaderAndEntity(result.Status, result.Body);
		}
	} else {
		response.WriteHeader(http.StatusUnauthorized)
	}
}

func EntityRequestHandler(request *restful.Request, response *restful.Response, fn entityRequestController, factory interface{}) {
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
