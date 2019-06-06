package config

import (
	"fmt"
	"log"
	"os"

	"github.com/globalsign/mgo"
)

// DB instance of MongoDB
var DB *mgo.Database

// Tasks are tasks for calculation a md5 sum for a file
var Tasks *mgo.Collection

// Session is a mongo session
var Session *mgo.Session

func init() {
	// get a mongo session
	//DB_CONNECTION_STRING = mongodb://localhost/md5 (env variable)
	var err error

	DbConnectionString := os.Getenv("DB_CONNECTION_STRING")

	if DbConnectionString == "" {
		log.Println("env variable DB_CONNECTION_STRING is not defined")
	}

	Session, err = mgo.Dial(DbConnectionString)
	if err != nil {
		log.Fatal("cannot dial mongo:", err)
	}

	if err = Session.Ping(); err != nil {
		log.Fatal("cannot ping mongo:", err)
	}

	mgo.SetStats(true)

	DB = Session.DB("md5")
	Tasks = DB.C("tasks")

	fmt.Println("You connected to your mongo database.")
}
