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
var cont = 0

// #coneccion con "dbaver"
// 	#entrar a conectar
// 	#seleccionar mysql
// 	#ingresar los datos
// 		server host: localhost
// 		port: 3306
// 		nombre de usuario: root
// 		contrase;a: 1234
// 	#probar coneccion

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

// get  ---------------------------------------------------------------------
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getUser(response http.ResponseWriter, request *http.Request) {
	enableCors(&response)

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

// post  ---------------------------------------------------------------------
func postUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	enableCors(&response)
	var usr User
	json.NewDecoder(request.Body).Decode(&usr)
	cont = cont + 1
	fmt.Println(cont, "(name:", usr.Name, " carnet:", usr.Carnet, ")")
	query := `INSERT INTO prueba(name, carnet) VALUES (?,?);`
	result, err := conn.Exec(query, usr.Name, usr.Carnet)
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(response).Encode(result)
}

// quitando los cors
func optPostCors(toAnsw *http.ResponseWriter) {
	(*toAnsw).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	(*toAnsw).Header().Set("Access-Control-Allow-Methods", "POST")
	(*toAnsw).Header().Set("Access-Control-Allow-Origin", "*")
}

func postUser_options(response http.ResponseWriter, request *http.Request) { // POST
	response.Header().Add("content-type", "application/json") //aceptar json en el request
	optPostCors(&response)
	json.NewEncoder(response).Encode("<repuesta POST>")
}

//*******************     main 		**********************************

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/postUser", postUser).Methods("POST")
	router.HandleFunc("/postUser", postUser_options).Methods("OPTIONS")

	//http://localhost:8000/getUsers
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
