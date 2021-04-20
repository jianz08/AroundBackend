package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	//parse from body of request to get a json object
	fmt.Println("Received one post request")
	
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		return
	}	
	
	decoder := json.NewDecoder(r.Body)//读出r.Body
	var p Post//post.go 也在main package，所以不需要import
	if err := decoder.Decode(&p); err != nil {//把r.Body convert成Post
		panic(err)//throw exception
	}
	fmt.Fprintf(w, "Post received: %s\n", p.Message)
}