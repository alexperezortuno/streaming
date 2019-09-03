package main

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "log"
    "net/http"
)

type Book struct {
    ID  string `json:"id"`
}

var books []Book

func GetBooks(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    json.NewEncoder(w).Encode(books)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/books", GetBooks).Methods("GET")
    log.Fatal(http.ListenAndServe(":8081", r))
}
