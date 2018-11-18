package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/transactions", GetTransactions).Methods("GET")
	log.Print(http.ListenAndServe(":8000", router))
}

func GetTransactions(writer http.ResponseWriter, request *http.Request)  {
	writer.WriteHeader(200)
	writer.Write([]byte("Here are your transactions"))
}