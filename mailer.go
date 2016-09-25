package main

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"net/http"
)

var d *gomail.Dialer

func Setup(host string, port int, user string, password string) {
	d = gomail.NewDialer(host, port, user, password)
}

func SendActivationEmail(request *http.Request, usr *User) {
	host := request.Host
	url := fmt.Sprintf("http://%s/users/activate/%s", host, usr.ActivationToken)

	m := gomail.NewMessage()
	m.SetHeader("From", "activation@example.com")
	m.SetHeader("To", usr.Email)
	m.SetHeader("Subject", "Welcome to example!")
	m.SetBody("text/html", fmt.Sprintf(`Hello <b>%s</b>, </br>
	Welcome to example! </br>
	Please click the following link to activate your account: </br>
	<a href=%s>%s</a>`, usr.UserName, url, url))

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}