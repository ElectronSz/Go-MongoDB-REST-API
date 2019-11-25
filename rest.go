/*
Author: ElectronSz
Github: https://github.com/ElectronSz
Date: 2019-11-23

*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type todo struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Status      bool   `json:"Status"`
}

type allTodos []todo

var todos = allTodos{}

func createTodo(w http.ResponseWriter, r *http.Request) {
	credential := options.Credential{
		Username: "api",
		Password: "api1008",
	}
	clientOptions := options.Client().ApplyURI("mongodb://167.71.60.107:27017").SetAuth(credential)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("Api").Collection("todos")

	var newTodo todo

	//reuest body and error
	reqBody, err := ioutil.ReadAll(r.Body)

	//check if we have an error
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newTodo)

	av := newTodo

	insertResult, err := collection.InsertOne(context.TODO(), av)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	//header status['200 ok]
	w.WriteHeader(http.StatusCreated)

	//json response
	json.NewEncoder(w).Encode(insertResult)

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func getOneTodo(w http.ResponseWriter, r *http.Request) {
	credential := options.Credential{
		Username: "api",
		Password: "api1008",
	}
	clientOptions := options.Client().ApplyURI("mongodb://167.71.60.107:27017").SetAuth(credential)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("Api").Collection("todos")

	/******************************************/
	todoID := mux.Vars(r)["id"]

	var oneTodo todo
	filter := bson.D{{"id", todoID}}

	err = collection.FindOne(context.TODO(), filter).Decode(&oneTodo)
	if err != nil {
		json.NewEncoder(w).Encode("Di not find any todo")
	}

	fmt.Printf("Found a single document: %+v\n", oneTodo)
	json.NewEncoder(w).Encode(oneTodo)

	/*******************************************************/
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

}
func getAllTodos(w http.ResponseWriter, r *http.Request) {
	credential := options.Credential{
		Username: "api",
		Password: "api1008",
	}
	clientOptions := options.Client().ApplyURI("mongodb://167.71.60.107:27017").SetAuth(credential)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	//
	fmt.Println("Connected to MongoDB!")

	collection := client.Database("Api").Collection("todos")

	/******************************************/
	findOptions := options.Find()
	//findOptions.SetLimit(2)

	var results []*todo

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem todo
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)

	json.NewEncoder(w).Encode(results)

	/*******************************************************/
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
func updateTodo(w http.ResponseWriter, r *http.Request) {
	credential := options.Credential{
		Username: "api",
		Password: "api1008",
	}
	clientOptions := options.Client().ApplyURI("mongodb://167.71.60.107:27017").SetAuth(credential)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	//
	fmt.Println("Connected to MongoDB!")

	collection := client.Database("Api").Collection("todos")

	/************************************************************/

	todoID := mux.Vars(r)["id"]

	var updatedTodo todo

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &updatedTodo)

	//opts := options.Update().SetUpsert(true)
	filter := bson.D{{"id", todoID}}
	update := bson.D{{"$set", bson.D{{"title", updatedTodo.Title}, {"description", updatedTodo.Description}, {"status", updatedTodo.Status}}}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document with ID %v\n", result.ModifiedCount)
		json.NewEncoder(w).Encode(result)
		return
	}
	if result.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
		json.NewEncoder(w).Encode(result)
	}

	if result.MatchedCount == 0 {
		fmt.Printf("We can not find that document yur want to update=> Result: %v\n", result.MatchedCount)
		json.NewEncoder(w).Encode(result)
	}

	/*******************************************************/
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	credential := options.Credential{
		Username: "api",
		Password: "api1008",
	}
	clientOptions := options.Client().ApplyURI("mongodb://167.71.60.107:27017").SetAuth(credential)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	//
	fmt.Println("Connected to MongoDB!")

	collection := client.Database("Api").Collection("todos")

	/******************************************/
	todoID := mux.Vars(r)["id"]

	filter := bson.D{{"id", todoID}}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the events collection\n", deleteResult.DeletedCount)

	json.NewEncoder(w).Encode(deleteResult)

	/*******************************************************/
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/todo", createTodo).Methods("POST")
	router.HandleFunc("/todos", getAllTodos).Methods("GET")
	router.HandleFunc("/todos/{id}", getOneTodo).Methods("GET")
	router.HandleFunc("/todos/{id}", updateTodo).Methods("PATCH")
	router.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
