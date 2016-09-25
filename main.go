package main

import (
	"fmt"
	"net/http"
	"log"
	"time"
	"flag"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/emicklei/go-restful"
	"github.com/dgrijalva/jwt-go"
	"github.com/magiconair/properties"
)

var (
	props          *properties.Properties
	propertiesFile = flag.String("config", "user-service.properties", "the configuration file")

	collection = &mgo.Collection{}
)

func main() {
	message := time.Second;
	fmt.Println(message)

	var err error
	if props, err = properties.LoadFile(*propertiesFile, properties.UTF8); err != nil {
		log.Fatalf("[error] Unable to read properties:%v\n", err)
	}

	session, err := mgo.Dial(fmt.Sprintf("%s:%d", props.GetString("mongod.host", ""), props.GetInt("mongod.port", 0)))
	if err != nil {
		panic(err)
	}
	defer session.Close()

	collection = session.DB("gotest").C("users")

	EnsureIndices()
	Setup(props.GetString("mail.host", ""), props.GetInt("mail.port", 0), props.GetString("mail.username", ""), props.GetString("mail.password", ""))

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
	service.Route(service.POST("/login").To(AuthenticateUser))

	return service
}

type User struct {
	UserName, Email string
	Password        string `json:",omitempty" xml:",omitempty"`
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
	usr := User{}

	err := collection.Find(bson.M{"username": userName}).One(&usr)
	usr.Password = ""

	if err == nil {
		response.WriteEntity(usr)
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

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(usr.Password))
	if err != nil {
		log.Println("Error comparing password to hash ", err)
		response.WriteError(http.StatusInternalServerError, err)
	}
	usr.Password = string(hashedPassword)

	err = collection.Insert(usr)
	usr.Password = ""

	SendActivationEmail(usr)

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
	usr.Password = ""

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

func AuthenticateUser(request *restful.Request, response *restful.Response) {
	usr := new(User)
	err := request.ReadEntity(&usr)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	unHashedPassword := usr.Password
	err = collection.Find(bson.M{"username": usr.UserName}).One(&usr)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(unHashedPassword))
	if err != nil {
		log.Println("Error comparing password to hash ", err)
		response.WriteError(http.StatusUnauthorized, err)
	}

	usr.Password = ""

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userName": usr.UserName,
		"nbf": time.Date(2016, 9, 24, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("my-secret"))

	fmt.Println(tokenString, err)

	if err == nil {
		response.WriteEntity(tokenString)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}