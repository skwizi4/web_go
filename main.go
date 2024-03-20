package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"main.go/DatabaseInit"
	"main.go/Logs"
	"main.go/WorkWithUsersTask"
	"main.go/handlers"
	"net/http"
)

func main() {
	DatabaseInit.DatabaseInit()
	router := mux.NewRouter()
	router.Use(Logs.LoggingMiddleware)
	router.HandleFunc("/main", handlers.MainPage)
	router.HandleFunc("/register", handlers.Register)
	router.HandleFunc("/login", handlers.Login)
	router.HandleFunc("/main/{id}", WorkWithUsersTask.WorkWithTask).Methods("POST")
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", router)
}
