package models

import (
	"github.com/globalsign/mgo/bson"
	"github.com/mariaefi29/md5_calculator/config"
	"github.com/pkg/errors"
)

//Task Struct
type Task struct {
	ID     bson.ObjectId `json:"id" bson:"_id"`
	IDstr  string        `json:"idstr" bson:"idstr"`
	Status string        `json:"status" bson:"status"`
	MD5    string        `json:"md5" bson:"md5,omitempty"`
	URL    string        `json:"url" bson:"url,omitempty"`
}

//CreateTask creates a task in a database for calculation a md5 hash code of a file
func CreateTask(url string) (Task, error) {

	//copy the global session to get access to a database
	currentSession := config.Session.Copy()
	defer currentSession.Close()

	//initialize the task
	task := Task{}

	//creates initial parametes
	task.ID = bson.NewObjectId()
	task.IDstr = task.ID.Hex()
	task.URL = url
	task.Status = "running"

	//creates a corresponding document in a mongo database
	err4 := config.Tasks.Insert(task)
	if err4 != nil {
		return task, errors.Wrap(err4, "Database error: fail to insert a task")
	}

	return task, nil
}

//FindTask fins a corresponding task by an id
func FindTask(taskIDstr string) (Task, error) {

	//копируем сессию для доступа к базе данных
	currentSession := config.Session.Copy()
	defer currentSession.Close()

	//copy the global session to get access to a database
	task := Task{}

	err := config.Tasks.Find(bson.M{"idstr": taskIDstr}).One(&task)
	if err != nil {
		return task, errors.Wrap(err, "Database error: fail to find a task")
	}
	return task, nil
}

//FailedTask updates the task's status to "failed"
func FailedTask(task Task) error {

	//копируем сессию для доступа к базе данных
	currentSession := config.Session.Copy()
	defer currentSession.Close()

	//updates the task's status to "failed"
	task.Status = "failed"

	err := config.Tasks.Update(bson.M{"_id": task.ID}, &task)
	if err != nil {
		return errors.Wrap(err, "Database error: failed to update a task")
	}
	return nil
}
