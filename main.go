package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mariaefi29/md5_calculator/config"
)

func main() {
	router := httprouter.New()
	router.GET("/check/:id", check)
	router.POST("/submit", submit)
	log.Fatal(http.ListenAndServe(":8080", router))
	defer config.Session.Close()
}
