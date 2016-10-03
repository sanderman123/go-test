package controller

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"github.com/sanderman123/user-service/model"
)

var d *gomail.Dialer

func Init(host string, port int, user string, password string) {
	d = gomail.NewDialer(host, port, user, password)
}

func SendActivationEmail(usr model.User, host string) error {
	url := fmt.Sprintf("http://%s/users/activate/%s", host, usr.ActivationToken)

	m := gomail.NewMessage()
	m.SetHeader("From", "activation@example.com")
	m.SetHeader("To", usr.Email)
	m.SetHeader("Subject", "Welcome to example!")
	m.SetBody("text/html", fmt.Sprintf(`Hello <b>%s</b>, </br>
	Welcome to example! </br>
	Please click the following link to activate your account: </br>
	<a href=%s>%s</a>`, usr.UserName, url, url))

	// Send the email
	return d.DialAndSend(m)
}

func SendPasswordResetEmail(usr *model.User, host string) error {
	url := fmt.Sprintf("http://%s/users/reset/%s", host, usr.ResetToken)

	m := gomail.NewMessage()
	m.SetHeader("From", "reset@example.com")
	m.SetHeader("To", usr.Email)
	m.SetHeader("Subject", "Password reset")
	m.SetBody("text/html", fmt.Sprintf(`Hello <b>%s</b>, </br>
	A password reset was requested. </br>
	If the reset was indeed requested by you, please click the following link to set a new password: </br>
	<a href=%s>%s</a>`, usr.UserName, url, url))

	// Send the email
	return d.DialAndSend(m);
}