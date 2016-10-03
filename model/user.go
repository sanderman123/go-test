package model

type User struct {
	UserName, Email string
	Password        string `json:",omitempty" xml:",omitempty"`
	ActivationToken string `json:"-" xml:"-" bson:",omitempty"`
	ResetToken      string `json:"-" xml:"-" bson:",omitempty"`
}

type Factory interface {
	NewEntity() User
}

type UserFactory struct {}

func (f UserFactory) NewEntity() User {
	return User{}
}