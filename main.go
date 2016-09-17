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
	//err = c.Insert(&User{"ID1", "Ale"},
	//	&User{"ID2", "Cla"})
	//if err != nil {
	//	log.Fatal(err)
	//}

	restful.Add(New())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func New() *restful.WebService {
	service := new(restful.WebService)
	service.
	Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_XML, restful.MIME_JSON)

	service.Route(service.GET("/{user-id}").To(FindUser))
	//service.Route(service.POST("").To(UpdateUser))
	//service.Route(service.PUT("/{user-id}").To(CreateUser))
	//service.Route(service.DELETE("/{user-id}").To(RemoveUser))

	return service
}

type User struct {
	Id, Name string
}

func FindUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	result := User{}
	err := collection.Find(bson.M{"id": id}).One(&result)
	if err != nil {
		log.Print("Error for user with id ", id, ": ", err)
		response.WriteHeader(404)
		return
	}
	response.WriteEntity(result)
}
