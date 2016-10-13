package controller

import (
	"github.com/sanderman123/user-service/model"
	"net/http"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"log"
	"github.com/sanderman123/user-service/util"
	"github.com/sanderman123/user-service/dao"
	"golang.org/x/net/xsrftoken"
)

func FindUser(claims jwt.MapClaims, userName string) Response {
	result := Response{}
	if claims[util.SUB] == userName {
		usr, err := dao.FindUserWithUserName(userName)
		usr.Password = ""
		result = ProduceResponse(usr, err)
	} else {
		result = ProduceStatusResponse("You can only retrieve your own information", nil, http.StatusUnauthorized)
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

func UpdateUser(claims jwt.MapClaims, entity interface{}) Response {
	result := Response{}
	usr := entity.(model.User)
	if claims[util.SUB] == usr.UserName {
		err := util.SetPassword(&usr, usr.Password)
		err = dao.UpdateUser(usr)
		usr.Password = ""
		result = ProduceResponse(usr, err)
	} else {
		result = ProduceStatusResponse("You can only update your own information", nil, http.StatusUnauthorized)
	}
	return result
}

func DeleteUser(claims jwt.MapClaims, userName string) Response {
	var result Response
	if claims[util.SUB] == userName {
		err := dao.RemoveUser(userName)
		result = ProduceResponse(nil, err)
	} else {
		result = ProduceStatusResponse("You can only delete your own information", nil, http.StatusUnauthorized)
	}
	return result
}

func AuthenticateUser(entity interface{}, w http.ResponseWriter) Response {
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
	xsrfToken := xsrftoken.Generate("my-xsrf-secret", usr.UserName, "ANY")

	// Create a new token object with claims
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.UserName,
		"iat": now,
		"exp": now.Add(time.Hour * 24),
		"xsrfToken": xsrfToken,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("my-secret"))

	fmt.Println(tokenString, err)

	http.SetCookie(w, &http.Cookie{
		Name: "jwt",
		Value: tokenString,
		Expires: now.Add(time.Hour * 24 * 7),
		Path: "/",
		Secure: true,
		HttpOnly: true,
	})
	jwtCookie := w.Header().Get("Set-Cookie")
	jwtCookie += "; SameSite=strict"
	w.Header().Set("Set-Cookie", jwtCookie)

	http.SetCookie(w, &http.Cookie{
		Name: "xsrfToken",
		Value: xsrfToken,
		Expires: now.Add(time.Hour * 24 * 7),
		Path: "/",
		Secure: true,
	})

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
