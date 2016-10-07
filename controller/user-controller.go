package controller

import (
	"github.com/sanderman123/user-service/model"
	"net/http"
	"gopkg.in/mgo.v2"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"log"
	"github.com/sanderman123/user-service/util"
	"github.com/sanderman123/user-service/dao"
)

func FindUser(userName string) Response {
	result := Response{}

	usr, err := dao.FindUserWithUserName(userName)
	usr.Password = ""

	if err == nil {
		result = Response{Status: http.StatusOK, Body: usr}
	} else if err == mgo.ErrNotFound {
		result = Response{Status: http.StatusNotFound}
	} else {
		log.Print("Error for user with userName ", userName, ": ", err)
		result = Response{Status: http.StatusInternalServerError, Body: err}
	}
	return result
}

func CreateUser(entity interface{}, request *http.Request) Response {
	result := Response{}
	usr := entity.(model.User)
	err := util.SetPassword(&usr, usr.Password)

	token, err := util.GenerateRandomString(32)
	if err != nil {
		message := "Generating activation token failed"
		log.Println(message, err)
		result = ProduceResponse(nil, err)
	} else {
		usr.ActivationToken = token

		err = dao.InsertUser(usr)
		usr.Password = ""

		SendActivationEmail(usr, request.Host)
		result = ProduceStatusResponse(usr, err, http.StatusCreated)
	}
	return result
}

func UpdateUser(entity interface{}) Response {
	result := Response{}
	usr := entity.(model.User)
	err := util.SetPassword(&usr, usr.Password)
	err = dao.UpdateUser(usr)
	usr.Password = ""

	if err == nil {
		result = Response{Status: http.StatusOK, Body: usr}
	} else if err == mgo.ErrNotFound {
		result = Response{Status: http.StatusNotFound}
	} else {
		result = Response{Status: http.StatusInternalServerError, Body: err}
	}
	return result
}

func DeleteUser(userName string) Response {
	err := dao.RemoveUser(userName)
	return ProduceResponse(nil, err)
}

func AuthenticateUser(entity interface{}) Response {
	usr := entity.(model.User)
	unHashedPassword := usr.Password
	usr, err := dao.FindUserWithUserName(usr.UserName)
	if err != nil {
		return ProduceResponse(nil, err)
	}

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(unHashedPassword))
	if err != nil {
		log.Println("[ERROR] Error comparing password to hash:", err)
		return ProduceStatusResponse(nil, nil, http.StatusUnauthorized)
	}

	usr.Password = ""

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userName": usr.UserName,
		"nbf": time.Date(2016, 9, 24, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("my-secret"))

	fmt.Println(tokenString, err)

	return ProduceResponse(tokenString, err)
}

func ActivateUser(token string) Response {
	result := Response{}
	usr, err := dao.FindUserWithActivationToken(token)

	if err == nil {
		usr.ActivationToken = ""
		err = dao.UpdateUser(usr)
		if err == nil {
			result = ProduceResponse(nil, nil)
		} else {
			log.Print("[ERROR] Error updating user with userName", usr.UserName, ": ", err)
			result = ProduceResponse(nil, err)
		}
	} else {
		result = ProduceResponse(nil, err)
	}

	return result
}

func ForgotPassword(entity interface{}, request *http.Request) Response {
	usr := entity.(model.User)
	usr, err := dao.FindUserWithEmail(usr.Email)
	message := ""

	if err == nil {

		token, err := util.GenerateRandomString(32)
		usr.ResetToken = token

		err = dao.UpdateUser(usr)

		if err == nil {
			err = SendPasswordResetEmail(&usr, request.Host)
			message = "A password reset email has been sent"
		}
	}

	return ProduceResponse(message, err)
}

func ResetPassword(token string, password string) Response {
	usr, err := dao.FindUserWithResetToken(token)

	if err == nil {
		err = util.SetPassword(&usr, password)

		if err == nil {
			usr.ResetToken = ""
			err = dao.UpdateUser(usr)
		}
	}
	return ProduceResponse(nil, err)
}
