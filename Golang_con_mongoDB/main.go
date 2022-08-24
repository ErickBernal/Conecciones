package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// desabilitar cors, lo necesitan GET,POST,PUT,DELETE
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// desabilitar cors, lo necesitan TODOS MENOS GET
func createUser_opt(response http.ResponseWriter, request *http.Request) { // POST
	response.Header().Add("content-type", "application/json") //aceptar json en el request
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	response.Header().Set("Access-Control-Allow-Methods", "POST")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(response).Encode("<repuesta POST, insertado>")
}

type User struct {
	Name   string `json:"name,omitempty"`
	Carnet string `json:"carnet,omitempty"`
	Time   string `json:"time,omitempty"`
}

var client *mongo.Client

// endPoints
func createUser(response http.ResponseWriter, request *http.Request) { // POST
	//corss
	enableCors(&response)

	response.Header().Add("content-type", "application/json") //aceptar json en el request
	var usr User
	json.NewDecoder(request.Body).Decode(&usr)
	usr.Time = time.Now().String()
	collection := client.Database("goDB").Collection("user")            //creamos la base de datos y la coleccion
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) //canal
	result, _ := collection.InsertOne(ctx, usr)
	json.NewEncoder(response).Encode(result)

}

func getUser(response http.ResponseWriter, request *http.Request) { // POST
	enableCors(&response)
	response.Header().Add("content-type", "application/json") //aceptar json en el request
	var listUsr []User
	collection := client.Database("goDB").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(ctx)
	//lectura de datos
	for cursor.Next(ctx) {
		var usr User
		cursor.Decode(&usr)
		listUsr = append(listUsr, usr)
	}
	json.NewEncoder(response).Encode(listUsr)

}

func main() {
	fmt.Println("Coneccion Golang y MongoDB")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) //canal
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()

	router.HandleFunc("/create", createUser).Methods("POST")
	router.HandleFunc("/create", createUser_opt).Methods("OPTIONS")

	router.HandleFunc("/get", getUser).Methods("GET")

	fmt.Println("server on port: ", 8000)
	http.ListenAndServe(":8000", router)

	defer client.Disconnect(ctx)
}
