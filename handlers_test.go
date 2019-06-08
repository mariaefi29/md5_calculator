package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/mariaefi29/md5_calculator/config"
	"github.com/mariaefi29/md5_calculator/models"
)

var ts *httptest.Server
var router *httprouter.Router

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	os.Exit(code)
}

func setUp() {
	router := httprouter.New()
	router.GET("/check/:id", check)
	router.POST("/submit", submit)
	ts := httptest.NewServer(router)
	defer ts.Close()
}

func TestCheck(t *testing.T) {

	//Adds a test task to a database
	url := "https/test.com"
	task, err := models.CreateTask(url)
	if err != nil {
		t.Errorf("Database error: %v", err)
	}
	id := task.IDstr

	//Creates a new writer and a request we are going to test
	writer := httptest.NewRecorder()

	//Creates a new request to the corresponding URL
	req := httptest.NewRequest("GET", "/check/"+id, nil)

	//Creates router parameters to transfer to the handler
	ps1 := httprouter.Param{Key: "id", Value: id}
	ps := []httprouter.Param{ps1}

	//Start the handler we are testing
	check(writer, req, ps)

	if writer.Code != 200 {
		t.Errorf("Response code: %v", writer.Code)
	}
	//Gets the response
	resp := writer.Result()

	//Gets the response Body and converts it to a string represantation
	body, err0 := ioutil.ReadAll(resp.Body)
	if err0 != nil {
		t.Errorf("Cannot read the response body: %s", err0)
	}
	bodyNew := string(body)

	defer resp.Body.Close()

	//Creates a response that the user is supposed to revieve
	resp0 := fmt.Sprintf("{\"id\":\"%s\", \"status\":\"%s\"}\n", task.IDstr, task.Status)

	//Compares two results
	if bodyNew != resp0 {
		t.Errorf("Expected %s, got %s", resp0, bodyNew)
	}

	//Cleans test data after testing
	if err := models.DeleteTask(task); err != nil {
		t.Error(err)
	}

	//Tests a request to get the data of a task that doesn't exist

	writer1 := httptest.NewRecorder()
	req1 := httptest.NewRequest("GET", "/check/"+id, nil)

	check(writer1, req1, ps)

	if writer1.Code != http.StatusNotFound {
		t.Errorf("Response Code is %v, expected %v", writer1.Code, http.StatusNotFound)
	}

	//Tests a request to with an empty id
	writer2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/check", nil)

	check(writer2, req2, nil)

	if writer2.Code != http.StatusBadRequest {
		t.Errorf("Response Code is %v, expected %v", writer2.Code, http.StatusBadRequest)
	}
}

func TestSubmit(t *testing.T) {
	//Tests the handler to handle a real url where we have a text file

	urlReal := "https://www.w3.org/TR/PNG/iso_8859-1.txt"

	//constracts a test form for a post Request
	form := url.Values{}
	form.Add("url", urlReal)

	testForm := strings.NewReader(form.Encode())

	writer := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/submit", testForm)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	submit(writer, req, nil)

	if writer.Code != http.StatusOK {
		t.Errorf("The Responce Code is %v, expected %v", writer.Code, http.StatusOK)
	}

	resp := writer.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Cannot read the response body: %s", err)
	}
	bodyNew := string(body)
	defer resp.Body.Close()

	task := models.Task{}
	// err := config.Tasks.Find(bson.M{"idstr": taskIDstr}).One(&task)
	dbSize, _ := config.Tasks.Count()
	err0 := config.Tasks.Find(nil).Skip(dbSize - 1).One(&task)
	if err0 != nil {
		t.Errorf("Cannot find a required task: %s", err0)
	}

	expected := fmt.Sprintf("{\"id\":\"%s\"}\n", task.IDstr)

	if bodyNew != expected {
		t.Errorf("got %v, expected %v", bodyNew, expected)
	}

	//Cleans test data after testing
	err1 := models.DeleteTask(task)
	if err1 != nil {
		t.Error(err1)
	}

	//test a post request with empty parametes
	//constracts a test form for a post request
	form1 := url.Values{}
	form1.Add("url", "")

	testForm1 := strings.NewReader(form1.Encode())

	writer1 := httptest.NewRecorder()
	req1 := httptest.NewRequest("POST", "/submit", testForm1)
	req1.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	submit(writer1, req1, nil)

	if writer1.Code != http.StatusBadGateway {
		t.Errorf("The Responce Code is %v, expected %v", writer.Code, http.StatusOK)
	}

}

func TestCalculateMD5(t *testing.T) {
	urlReal := "https://www.w3.org/TR/PNG/iso_8859-1.txt"

	urlUnreal := "https://www.site.com/file.txt"

	task1, err1 := models.CreateTask(urlReal)
	if err1 != nil {
		t.Error(err1)
	}

	code1, err2 := calculateMD5(urlReal, task1)
	if err2 != nil {
		t.Errorf("Failed to calcilate md5 hash code: %v", err2)
	}
	//updating the task struct
	task1, err3 := models.FindTask(task1.IDstr)
	if err3 != nil {
		t.Error(err3)
	}
	//We are testing here that if we calculated md5, it was properly assigned to the task
	if task1.MD5 != code1 {
		t.Errorf("MD5 hash code mismatch: should be %s, got %s", code1, task1.MD5)
	}
	//and status should be equal to "done"
	if task1.Status != "done" {
		t.Errorf("Status is %s, expected \"done\"", task1.Status)
	}

	//
	task2, err4 := models.CreateTask(urlUnreal)
	if err4 != nil {
		t.Error(err4)
	}
	code2, err5 := calculateMD5(urlUnreal, task2)
	if err2 != nil {
		t.Errorf("Failed to calcilate md5 hash code: %v", err5)
	}

	if code2 != "" {
		t.Errorf("Expected an empty string, got %s", code2)
	}

	//updating the task struct
	task2, err6 := models.FindTask(task2.IDstr)
	if err6 != nil {
		t.Error(err6)
	}

	if task2.Status != "failed" {
		t.Errorf("Expected the status \"failed\", got %s", task2.Status)
	}

	if task2.MD5 != "" {
		t.Errorf("Expected an empty string, got %s", task2.MD5)
	}

	//Clean up after testing
	if err := models.DeleteTask(task1); err != nil {
		t.Error(err)
	}

	if err := models.DeleteTask(task2); err != nil {
		t.Error(err)
	}
}
