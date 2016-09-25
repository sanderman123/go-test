package util

import (
	"log"
	"golang.org/x/crypto/bcrypt"
	"github.com/sanderman123/user-service/model"
)

func SetPassword(usr *model.User, password string) error {
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		message := "Error comparing password to hash"
		log.Println(message, err)
		return err
	}

	usr.Password = string(hashedPassword)
	return nil
}