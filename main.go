package main

import (
	"github.com/emicklei/go-restful"
	"fmt"
	"net/http"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var collection = &mgo.Collection{}

func main() {
	message := "Hello, world"
	fmt.Println(message)

	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	collection = session.DB("gotest").C("users")

	EnsureIndices()

	restful.Add(New())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func New() *restful.WebService {
	service := new(restful.WebService)
	service.
	Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_XML, restful.MIME_JSON)

	service.Route(service.GET("/{user-name}").To(FindUser))
	service.Route(service.POST("").To(CreateUser))
	service.Route(service.PUT("").To(UpdateUser))
	service.Route(service.DELETE("/{user-name}").To(DeleteUser))

	return service
}

type User struct {
	UserName, Email string
	Password        string `json:"-" xml:"-"`
}

func EnsureIndices() {
	// Index
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

func FindUser(request *restful.Request, response *restful.Response) {
	userName := request.PathParameter("user-name")
	result := User{}

	err := collection.Find(bson.M{"username": userName}).One(&result)

	if err == nil {
		response.WriteEntity(result)
	} else if err == mgo.ErrNotFound {
		response.WriteHeader(http.StatusNotFound)
	} else {
		log.Print("Error for user with userName ", userName, ": ", err)
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func CreateUser(request *restful.Request, response *restful.Response) {
	usr := new(User)
	err := request.ReadEntity(&usr)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	err = collection.Insert(usr)

	if err == nil {
		response.WriteHeaderAndEntity(http.StatusCreated, usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func UpdateUser(request *restful.Request, response *restful.Response) {
	usr := new(User)
	err := request.ReadEntity(&usr)

	err = collection.Update(bson.M{"username": usr.UserName}, usr)

	if err == nil {
		response.WriteEntity(usr)
	} else if err == mgo.ErrNotFound {
		response.WriteHeader(http.StatusNotFound)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func DeleteUser(request *restful.Request, response *restful.Response) {
	// here you would delete the user from some persistence system
	userName := request.PathParameter("user-name")
	err := collection.Remove(bson.M{"username": userName});

	if err == nil {
		response.WriteHeader(http.StatusOK)
	} else if err == mgo.ErrNotFound {
		response.WriteHeader(http.StatusNotFound)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}