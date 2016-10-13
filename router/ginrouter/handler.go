package ginrouter

import (
	"github.com/sanderman123/user-service/controller"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/sanderman123/user-service/model"
	"github.com/sanderman123/user-service/util"
	"github.com/dgrijalva/jwt-go"
)

type pathParameterController func(parameter string) controller.Response
type pathParameterStringController func(parameter string, string string) controller.Response
type entityController func(entity interface{}) controller.Response
type entityRequestController func(entity interface{}, request *http.Request) controller.Response
type entityResponseWriterController func(entity interface{}, request http.ResponseWriter) controller.Response

type authenticatedPathParameterController func(claims jwt.MapClaims, parameter string) controller.Response
type authenticatedEntityController func(claims jwt.MapClaims, entity interface{}) controller.Response

//Generic Handlers
func PathParameterHandler(c *gin.Context, fn pathParameterController, parameterName string) {
	result := fn(c.Param(parameterName))
	respond(c, result)
}

func AuthenticatedPathParameterHandler(c *gin.Context, fn authenticatedPathParameterController, parameterName string) {
	claims, err := util.IsAuthenticated(c.Request)
	if err == nil {
		respond(c, fn(claims, c.Param(parameterName)))
	} else {
		respond(c, controller.Response{Status: http.StatusUnauthorized})
	}
}

func PathParameterStringHandler(c *gin.Context, fn pathParameterStringController, parameterName string) {
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

func AuthenticatedEntityHandler(c *gin.Context, fn authenticatedEntityController, factory interface{}) {
	claims, err := util.IsAuthenticated(c.Request)
	if err == nil {
		b := factory.(model.Factory)
		entity := b.NewEntity()
		err := c.Bind(&entity)
		if err == nil {
			respond(c, fn(claims, entity))
		} else {
			respond(c, controller.Response{Status: http.StatusInternalServerError})
		}
	} else {
		respond(c, controller.Response{Status: http.StatusUnauthorized})
	}
}

func EntityRequestHandler(c *gin.Context, fn entityRequestController, factory interface{}) {
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

func EntityResponseWriterHandler(c *gin.Context, fn entityResponseWriterController, factory interface{}) {
	b := factory.(model.Factory)
	entity := b.NewEntity()
	err := c.Bind(&entity)
	if err == nil {
		respond(c, fn(entity, c.Writer.(http.ResponseWriter)))
	} else {
		respond(c, controller.Response{Status: http.StatusInternalServerError})
	}
}

func respond(c *gin.Context, response controller.Response) {
	if response.Body == nil {
		c.Writer.WriteHeader(response.Status)
	} else {
		c.JSON(response.Status, response.Body)
	}
}
