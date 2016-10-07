package gorestful

import (
	"github.com/emicklei/go-restful"
	"github.com/sanderman123/user-service/model"
	"net/http"
	"github.com/sanderman123/user-service/controller"
)

type pathParameterController func(parameter string) controller.Response
type pathParameterStringController func(parameter string, string string) controller.Response
type entityController func(entity interface{}) controller.Response
type entityRequestController func(entity interface{}, request *http.Request) controller.Response


//Generic Handlers
func PathParameterHander(request *restful.Request, response *restful.Response, fn pathParameterController, parameterName string) {
	result := fn(request.PathParameter(parameterName))
	response.WriteHeaderAndEntity(result.Status, result.Body)
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
