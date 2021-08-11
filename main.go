package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Task struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
}

var client *mongo.Client

func CreateTaskEndpoint(response http.ResponseWriter, req *http.Request) {
	response.Header().Set("content-type", "application/json")
	var task Task

	_ = json.NewDecoder(req.Body).Decode(&task)
	collection := client.Database("task-planner").Collection("task")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := collection.InsertOne(ctx, task)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":" "` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(response).Encode(result)
}

func GetTasksEndpoint(response http.ResponseWriter, req *http.Request) {
	response.Header().Set("content-type", "application/json")
	var tasks []Task

	collection := client.Database("task-planner").Collection("task")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// we need to retrieve all data
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":" "` + err.Error() + `"}`))
		return
	}

	// close the cursor
	defer cursor.Close(ctx)

	// loop for each item in cursor, decode and insert it to Slice
	for cursor.Next(ctx) {
		var task Task
		cursor.Decode(&task)
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":" "` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(response).Encode(tasks)
}

func GetSingleTaskEndpoint(response http.ResponseWriter, req *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(req)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var task Task

	collection := client.Database("task-planner").Collection("task")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err := collection.FindOne(ctx, Task{ID: id}).Decode(&task)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(response).Encode(task)
}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb+srv://" + os.Getenv("MONGO_USER") + ":" + os.Getenv("MONGO_PASS") + "@task-planner.gbyx1.mongodb.net/" + os.Getenv("MONGO_DB") + "?retryWrites=true&w=majority")
	localClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	client = localClient
	router := mux.NewRouter()
	router.HandleFunc("/task", CreateTaskEndpoint).Methods("POST")
	router.HandleFunc("/task/{id}", GetSingleTaskEndpoint).Methods("GET")
	router.HandleFunc("/task", GetTasksEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}

// reference https://www.thepolyglotdeveloper.com/2019/02/developing-restful-api-golang-mongodb-nosql-database/
