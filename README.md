# MD5 Calculator

Source code for a MD5 calculator application.

A simple web application calculates a MD5 Hash code of a file located on the Internet.

This application accepts two requests:
1. POST request on `/submit` has a URL parameter (the address of a file). It creates a task that is responsible for the md5 hash code calculation, gives back the ID of the task and starts calculating the hash code as a background process.
2. GET request on `/check/task_id` provides a user with the information of the task: only status if the process is still running or the md5 hash code itself if the calculation has already been finished.

## Building an application locally

1. Make sure you have MongoDB database running either locally or remotely. [MongoDB installation process](https://docs.mongodb.com/manual/installation/)
2. Create an environmental variable DB_CONNECTION_STRING and assign it to an url of your MongoDB database. For example: `mongodb://localhost/md5`
3. Run `make` in a command line in a working directory.
4. Start "playing" with the application.

## How to use:
```$ curl -X POST -d "url=http://site.com/file.txt" http://localhost:8080/submit```

{"id":"5cfbb4561b699b0e9a633fe3"}   

```$ curl -X GET http://localhost:8080/check/5cfbb4561b699b0e9a633fe3```    

{"id":"5cfbb4561b699b0e9a633fe3", "status":"running"}  

```$ curl -X GET http://localhost:8080/check/5cfbb4561b699b0e9a633fe3```    

{"id":"5cfbb4561b699b0e9a633fe3", "status":"done", "md5":"74454cfc8f400c7d32a70d1f1eda6ce8", "url":"http://site.com/file.txt"}  
