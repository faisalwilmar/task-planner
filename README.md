# Project Title

The most common way to practice building RESTful API is a Task Planner. This is it.

## Lessons Learned

This project is based on Nic Raboy's Article in [here]("https://www.thepolyglotdeveloper.com/2019/02/developing-restful-api-golang-mongodb-nosql-database/")

When following his article using Visual Studio Code in windows, the intellisense said that required module is missing. I fix that by initiating go module before started to working on the project.
```bash
go mod init task-planner
```
Everytime the intellisense says the module is missing, or the required dependency is incomplete, just run go mod tidy.
```bash
go mod tidy
```

  
## Installation

Clone this project.
Make sure you have Golang >1.6.6 installed and a MongoDB Cluster Database

    