package main

import (
	"fmt"
	"log"
	"time"
	"flag"
	"gopkg.in/mgo.v2"
	"github.com/emicklei/go-restful"
	"github.com/magiconair/properties"
	"github.com/sanderman123/user-service/controller"
	"github.com/sanderman123/user-service/dao"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/sanderman123/user-service/router/ginrouter"
	"github.com/sanderman123/user-service/router/go-restful"
)

var (
	props          *properties.Properties
	propertiesFile = flag.String("config", "user-service.properties", "the configuration file")
)

var serverKey = ""
var serverCertificate = ""

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

	serverKey = props.GetString("https.serverKey", "")
	serverCertificate = props.GetString("https.serverCertificate", "")

	//GoRestful()
	Gin()
}

func GoRestful() {
	service := new(restful.WebService)
	gorestful.Init(service)
	restful.Add(service)
	//log.Fatal(http.ListenAndServe(":8080", nil))
	log.Fatal(http.ListenAndServeTLS(":8080", serverCertificate, serverKey, nil))
}

func Gin() {
	r := gin.Default()
	ginrouter.Init(r)
	//r.Run(":8080")
	r.RunTLS(":8080", serverCertificate, serverKey)
}

func Log(l *log.Logger, msg string) {
	//Format("2006-01-02 15:04:05:006")
	l.SetPrefix(fmt.Sprint(time.Now().UnixNano() / 1000000) + " [AAA] ")
	l.Print(msg)
}
