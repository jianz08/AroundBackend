package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"github.com/pborman/uuid"
)

var (
	mediaTypes = map[string] string {
		".jpeg": "image",
		".jpg": "image",
		".gif": "iamge",
		".png": "iamge",
		".mov": "video",
		".mp4": "video",
		".avi": "video",
		".flv": "video",
		".wmv": "video",
	}
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	//parse from body of request to get a json object
	fmt.Println("Received one upload request")
	
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		return
	}	
	
	p := Post{
		Id: uuid.New(),
		User: r.FormValue("user"),
		Message: r.FormValue("message"),
	}
	file, header, err := r.FormFile("media_file")
	if err != nil {
		http.Error(w, "Media file is not available", http.StatusBadRequest)
		fmt.Printf("Media file is not available %v\n", err)
		return
	}
	suffix := filepath.Ext(header.Filename)
	if t, ok := mediaTypes[suffix]; ok {
		p.Type = t
	} else {
		p.Type = "unknown"
	}
	err = savePost(&p, file)
	if err != nil {
		http.Error(w, "Failed to save post to GCS or Elasticsearch", http.StatusInternalServerError)
		fmt.Printf("Failed to save post to GCS or Elasticsearch %v\n", err)
		return
	}
	fmt.Println("Post is saved successfully.")

	// decoder := json.NewDecoder(r.Body)//读出r.Body
	// var p Post//post.go 也在main package，所以不需要import
	// if err := decoder.Decode(&p); err != nil {//把r.Body convert成Post. Decode: json->go object
	// 	panic(err)//throw exception
	// }
	// fmt.Fprintf(w, "Post received: %s\n", p.Message)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {//response是interface，没有实例化，所以没有*。 request是struct，所以前面可以有*
	fmt.Println("Received one request for search")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		return
	}

	user := r.URL.Query().Get("user")//读取request参数
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

	js, err := json.Marshal(posts)//Marshal，go object -> json
	if err != nil {
		http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
		fmt.Printf("Failed to parse posts into JSON format %v. \n", err)
	}
	w.Write(js)
}