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

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for search")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		return
	}

	user := r.URL.Query().Get("user")
	keywords := r.URL.Query().Get("keywords")

	var posts []Post
	var err error

	if user != "" {
		posts, err = searchPostsByUser(user)
	} else {
		posts, err = searchPostsByKeywords(keywords)
	}

	if err != nil {
		http.Error(w, "Failed to read post from Elasticsearch", http.StatusInternalServerError)
		fmt.Printf("Failed to read post from Elasticsearch %v. \n", err)
	}

	js, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
		fmt.Printf("Failed to parse posts into JSON format %v. \n", err)
	}
	w.Write(js)
}