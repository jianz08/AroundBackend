package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    fmt.Println("started-service")

    r := mux.NewRouter()
    r.Handle("/upload", http.HandlerFunc(uploadHandler)).Methods("POST", "OPTIONS")
    log.Fatal(http.ListenAndServe(":8080", r))
    //http.HandleFunc("/upload", uploadHandler)
    //log.Fatal(http.ListenAndServe(":8080", nil))//nil is the default http router,不能区分post, get    
}