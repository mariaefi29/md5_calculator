package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/julienschmidt/httprouter"
	"github.com/mariaefi29/md5_calculator/config"
	"github.com/mariaefi29/md5_calculator/models"
	"github.com/pkg/errors"
)

//POST request /submit with url parameter.
//Creates a task with an id that lets the user find the status of its progress.
//Responds with a corresponding task id.
func submit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, http.StatusText(400), http.StatusBadGateway)
		error := errors.New("Failed to obtain post request parameter")
		log.Println(error)
		return
	}

	task, err := models.CreateTask(url)
	if err != nil {
		http.Error(w, http.StatusText(500)+" "+err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	//Writes an task's id.
	// fmt.Println("id:", task.IDstr)
	fmt.Fprintln(w, "{\"id\":"+"\""+task.IDstr+"\"}")

	//Starts a goroutine that retrives data from the file and calculate a md5 hash code.
	go calculateMD5(url, task)

	return
}

func calculateMD5(url string, task models.Task) (string, error) {

	md5Code := ""

	//Send a GET request to the corresponding URL
	resp, err := http.Get(url)
	if err != nil {
		models.FailedTask(task)
		return md5Code, errors.Wrap(err, "Failed to access file through the URL")
	}
	defer resp.Body.Close()

	h := md5.New()

	//Copy the data from the response Body into
	if _, err := io.Copy(h, resp.Body); err != nil {
		models.FailedTask(task)
		log.Println(errors.Wrap(err, "Failed to copy data from the response"))
		return md5Code, errors.Wrap(err, "Failed to copy data from the response")
	}

	md5Code = hex.EncodeToString(h.Sum(nil))

	//Update hash code and a task's status
	task.MD5 = md5Code
	task.Status = "done"

	//Update data in the database

	if err := config.Tasks.Update(bson.M{"_id": task.ID}, &task); err != nil {
		log.Println(errors.Wrap(err, "Failed to update data in a database"))
		return md5Code, errors.Wrap(err, "Failed to update data in a database")
	}

	return md5Code, nil
}

// GET request /check with id parameter. The check handler responds woth a task id
// and status.

func check(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Obtain the id from the URL using the request parameters
	id := ps.ByName("id")
	if id == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		error := errors.New("Failed to parse url and find a task id")
		log.Println(error)
		return
	}

	//Find the task
	task, err := models.FindTask(id)
	if err != nil {
		http.Error(w, http.StatusText(404)+" "+err.Error(), http.StatusNotFound)
		log.Println(err)
		return
	}
	//Print out all the data to the user. If the task is still running or was failed, prints out the data parcially.
	if task.Status == "done" {
		fmt.Fprintf(w, "{\"id\":\"%s\", \"status\":\"%s\", \"md5\":\"%s\", \"url\":\"%s\"}\n", task.IDstr, task.Status, task.MD5, task.URL)
	} else {
		fmt.Fprintf(w, "{\"id\":\"%s\", \"status\":\"%s\"}\n", task.IDstr, task.Status)
	}
}
