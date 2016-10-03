package main

import (
	"fmt"
	//"net/http"
	"log"
	"time"
	"flag"
	"gopkg.in/mgo.v2"
	"github.com/emicklei/go-restful"
	"github.com/magiconair/properties"
	"github.com/sanderman123/user-service/controller"
	"github.com/sanderman123/user-service/dao"
	"github.com/sanderman123/user-service/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	props          *properties.Properties
	propertiesFile = flag.String("config", "user-service.properties", "the configuration file")
)

func main() {
	message := "Hello World!"
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

	database := session.DB("gotest")
	//userDao := dao.UserDaoImpl{}
	//userDao.Init(database)
	dao.Init(database)

	controller.Init(props.GetString("mail.host", ""), props.GetInt("mail.port", 0), props.GetString("mail.username", ""), props.GetString("mail.password", ""))

	restful.Add(New())
	log.Fatal(http.ListenAndServe(":8080", nil))
	//Gin()
}

func New() *restful.WebService {
	service := new(restful.WebService)
	router.Init(service)
	return service
}

func Gin() {
	r := gin.Default()
	//_ = r.Group("/users")
	//{
		r.GET("/users/:name", func(c *gin.Context) {
			name := c.Param("name")
			response := controller.FindUser(name)
			c.JSON(response.Status, response.Body)
		})
		//r.POST("", controller.CreateUser)
		//r.PUT("", controller.UpdateUser)
		//r.DELETE("/{user-name}", controller.DeleteUser)
		//r.POST("/login", controller.AuthenticateUser)
		//r.GET("/activate/{token}", controller.ActivateUser)
		//r.POST("/forgot", controller.ForgotPassword)
		//r.POST("/reset/{token}", controller.ResetPassword)
	//}

	r.Run(":8080")
}

func Log(l *log.Logger, msg string) {
	//Format("2006-01-02 15:04:05:006")
	l.SetPrefix(fmt.Sprint(time.Now().UnixNano() / 1000000) + " [AAA] ")
	l.Print(msg)
}
