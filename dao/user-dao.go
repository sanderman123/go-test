package dao

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/sanderman123/user-service/model"
)

var collection *mgo.Collection

func Init(database *mgo.Database) {
	collection = database.C("users")
	EnsureIndices()
}

func FindUserWithUserName(userName string) (model.User, error) {
	usr := model.User{}
	err := collection.Find(bson.M{"username": userName}).One(&usr)
	return usr, err
}

func FindUserWithEmail(email string) (model.User, error) {
	usr := model.User{}
	err := collection.Find(bson.M{"email": email}).One(&usr)
	return usr, err
}

func FindUserWithActivationToken(token string) (model.User, error) {
	usr := model.User{}
	err := collection.Find(bson.M{"activationtoken": token}).One(&usr)
	return usr, err
}

func FindUserWithResetToken(token string) (model.User, error) {
	usr := model.User{}
	err := collection.Find(bson.M{"resettoken": token}).One(&usr)
	return usr, err
}

func InsertUser(usr model.User) error {
	return collection.Insert(usr)
}

func UpdateUser(usr model.User) error {
	return collection.Update(bson.M{"username": usr.UserName}, usr)
}

func RemoveUser(userName string) error {
	return collection.Remove(bson.M{"username": userName});
}


func EnsureIndices() {
	index := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	index = mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}
