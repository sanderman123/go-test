package main

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

var d *gomail.Dialer

func Setup(host string, port int, user string, password string) {
	d = gomail.NewDialer(host, port, user, password)
}

func SendActivationEmail(usr *User) {
	m := gomail.NewMessage()
	m.SetHeader("From", "activation@example.com")
	m.SetHeader("To", usr.Email)
	m.SetHeader("Subject", "Welcome to example!")
	m.SetBody("text/html", fmt.Sprintf(`Hello <b>%s</b>, </br>
	Welcome to example! </br>
	Please click the following link to activate your account: </br>
	<a href=%s>%s</a>`, usr.UserName, "LINK", "LINK"))

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}