package util

import (
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"log"
)

const XSRF_KEY = "xsrfToken"
const SUB = "sub"

func IsAuthenticated(request *http.Request) (jwt.MapClaims, error) {
	xsrfToken := request.Header.Get(XSRF_KEY)
	tokenString := ""
	cookie, err := request.Cookie("jwt")

	if err == nil {
		tokenString = cookie.Value
		token, errr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Unexpected signing method")
			}
			return []byte("my-secret"), nil
		})

		if (errr == nil) {
			return GetClaims(token, xsrfToken)
		}
		log.Println(errr)
		err = errr
	}
	log.Println(err)
	return nil, err
}

func GetClaims(token *jwt.Token, xsrfToken string) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid && xsrfToken == claims[XSRF_KEY] {
		return claims, nil
	}
	return nil, errors.New("Invalid token")
}