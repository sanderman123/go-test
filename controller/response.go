package controller

import (
	"net/http"
	"gopkg.in/mgo.v2"
	"fmt"
	"log"
)

type Response struct {
	Status int
	Body   interface{}
}

func ProduceResponse(entity interface{}, err error) Response {
	result := Response{Status: http.StatusNotImplemented}
	if err == nil {
		result = Response{Status: http.StatusOK, Body: entity}
	} else if err == mgo.ErrNotFound {
		result = Response{Status: http.StatusNotFound}
	} else {
		log.Println(fmt.Sprintf("[ERROR] An error occurred: %s", err))
		result = Response{Status: http.StatusInternalServerError, Body: "An error occurred"}
	}
	return result
}

func ProduceStatusResponse(entity interface{}, err error, status int) Response {
	result := Response{Status: http.StatusNotImplemented}
	if err == nil {
		result = Response{Status: status, Body: entity}
	} else if err == mgo.ErrNotFound {
		result = Response{Status: http.StatusNotFound}
	} else {
		log.Println(fmt.Sprintf("[ERROR] An error occurred: %s", err))
		result = Response{Status: http.StatusInternalServerError, Body: "An error occurred"}
	}
	return result
}
