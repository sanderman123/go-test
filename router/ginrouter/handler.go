package ginrouter

import (
	"github.com/sanderman123/user-service/controller"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/sanderman123/user-service/model"
)

type pathParameterController func(parameter string) controller.Response
type pathParameterStringController func(parameter string, string string) controller.Response
type entityController func(entity interface{}) controller.Response
type entityRequestController func(entity interface{}, request *http.Request) controller.Response

//Generic Handlers
func PathParameterHander(c *gin.Context, fn pathParameterController, parameterName string) {
	result := fn(c.Param(parameterName))
	respond(c, result)
}

func PathParameterStringHander(c *gin.Context, fn pathParameterStringController, parameterName string) {
	str := ""
	err := c.Bind(&str)
	if err != nil {
		respond(c, controller.Response{Status: http.StatusInternalServerError})
	} else {
		respond(c, fn(c.Param(parameterName), str))
	}
}

func EntityHandler(c *gin.Context, fn entityController, factory interface{}) {
	b := factory.(model.Factory)
	entity := b.NewEntity()
	err := c.Bind(&entity)
	if err == nil {
		respond(c, fn(entity))
	} else {
		respond(c, controller.Response{Status: http.StatusInternalServerError})
	}
}

func EntityRequestHander(c *gin.Context, fn entityRequestController, factory interface{}) {
	b := factory.(model.Factory)
	entity := b.NewEntity()
	err := c.Bind(&entity)
	if err != nil {
		respond(c, controller.Response{Status: http.StatusInternalServerError})
	} else {
		result := fn(entity, c.Request)
		respond(c, result)
	}
}

func respond(c *gin.Context, response controller.Response) {
	if response.Body == nil {
		c.Writer.WriteHeader(response.Status)
	} else {
		c.JSON(response.Status, response.Body)
	}
}
