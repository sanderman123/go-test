package dao

//import (
//	"gopkg.in/mgo.v2/bson"
//	"gopkg.in/mgo.v2"
//	"github.com/sanderman123/user-service/model"
//)
//
//type UserDao interface {
//	Init()
//	FindUserWithUserName() (model.User, error)
//	FindUserWithEmail() (model.User, error)
//	FindUserWithActivationToken(token string) (model.User, error)
//	FindUserWithResetToken(token string) (model.User, error)
//	InsertUser(usr model.User) error
//	UpdateUser(usr model.User) error
//	RemoveUser(userName string) error
//}
//
//type UserDaoImpl struct {
//	collection mgo.Collection
//}
//
//func (d UserDaoImpl) Init(database *mgo.Database) {
//	d.collection = database.C("users")
//	d.EnsureIndices()
//}
//
//func (d UserDaoImpl) FindUserWithUserName(userName string) (model.User, error) {
//	usr := model.User{}
//	err := d.collection.Find(bson.M{"username": userName}).One(&usr)
//	return usr, err
//}
//
//func (d UserDaoImpl) FindUserWithEmail(email string) (model.User, error) {
//	usr := model.User{}
//	err := d.collection.Find(bson.M{"email": email}).One(&usr)
//	return usr, err
//}
//
//func (d UserDaoImpl) FindUserWithActivationToken(token string) (model.User, error) {
//	usr := model.User{}
//	err := d.collection.Find(bson.M{"activationtoken": token}).One(&usr)
//	return usr, err
//}
//
//func (d UserDaoImpl) FindUserWithResetToken(token string) (model.User, error) {
//	usr := model.User{}
//	err := d.collection.Find(bson.M{"resettoken": token}).One(&usr)
//	return usr, err
//}
//
//func (d UserDaoImpl) InsertUser(usr model.User) error {
//	return d.collection.Insert(usr)
//}
//
//func (d UserDaoImpl) UpdateUser(usr model.User) error {
//	return d.collection.Update(bson.M{"username": usr.UserName}, usr)
//}
//
//func (d UserDaoImpl) RemoveUser(userName string) error {
//	return d.collection.Remove(bson.M{"username": userName});
//}
//
//func (d UserDaoImpl) EnsureIndices() {
//	index := mgo.Index{
//		Key:        []string{"username"},
//		Unique:     true,
//		DropDups:   true,
//		Background: true,
//		Sparse:     true,
//	}
//
//	err := d.collection.EnsureIndex(index)
//	if err != nil {
//		panic(err)
//	}
//
//	index = mgo.Index{
//		Key:        []string{"email"},
//		Unique:     true,
//		DropDups:   true,
//		Background: true,
//		Sparse:     true,
//	}
//
//	err = d.collection.EnsureIndex(index)
//	if err != nil {
//		panic(err)
//	}
//}
