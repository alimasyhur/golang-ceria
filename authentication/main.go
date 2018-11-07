package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

//DBHost,DBPort,DBUser,DBPassword,DBBase, PORT const
const (
	DBHost     = "127.0.0.1"
	DBPort     = "3306"
	DBUser     = "root"
	DBPassword = ""
	DBBase     = "test"
	PORT       = ":9000"
	//hashCost
	hashCost = 8
)

var db *sql.DB

func main() {
	// "Signin" and "Signup" are handler that we will implement
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/signin", SigninHandler)
	http.HandleFunc("/signup", SignupHandler)
	// initialize our database connection
	initDB()
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func initDB() {
	var err error
	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUser, DBPassword, DBHost, DBPort, DBBase)
	db, err = sql.Open("mysql", dbConn)
	if err != nil {
		panic(err)
	}
}
