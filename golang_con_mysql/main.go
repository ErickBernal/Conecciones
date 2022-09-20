package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//go get github.com/go-sql-driver/mysql
//go get github.com/gorilla/mux

var conn = MySQLConn()

func MySQLConn() *sql.DB {
	//connString := "root:password@tcp(34.135.81.94:3306)/pro1_bases1"
	connString := "root:1234@tcp(localhost:3306)/bases1_pro1"

	conn, err := sql.Open("mysql", connString)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("\"coneccion con MySql\"")
	}
	return conn
}

// *******************     ejemplos **********************************
type User struct {
	Name   string `json:"name"`
	Carnet int    `json:"carnet"`
}

func getUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var listUsr []User
	query := "SELECT * FROM prueba;"
	result, err := conn.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	for result.Next() {
		var usr User
		err = result.Scan(&usr.Name, &usr.Carnet)
		if err != nil {
			fmt.Println(err)
		}
		listUsr = append(listUsr, usr)
	}
	json.NewEncoder(response).Encode(listUsr)
}

func postUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var usr User
	json.NewDecoder((request.Body)).Decode(&usr)
	query := `INSERT INTO prueba(name, carnet) VALUES (?,?);`
	result, err := conn.Exec(query, usr.Name, usr.Carnet)
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(response).Encode(result)
}

//*******************     main 		**********************************

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/postUser", postUser).Methods("POST")
	router.HandleFunc("/getUsers", getUser).Methods("GET")

	fmt.Println("Server on port", 8000)
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		fmt.Println(err)
	}
}

// POST
/*
	#EndPoint
		http://localhost:8000/postUser

	#body
		{
			"name":"Go2",
			"carnet":312
		}
*/
// GET
//	http://localhost:8000/getUsers
