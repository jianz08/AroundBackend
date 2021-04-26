package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"github.com/pborman/uuid"

	"regexp"
	"time"
	jwt "github.com/form3tech-oss/jwt-go"//import jwt-go as jwt
	"github.com/gorilla/mux"

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
	
	w.Header().Set("Access-Control-Allow-Origin", "*")//后端设置允许跨域
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		return
	}	

	//通过token获取username
	user := r.Context().Value("user")//raw token string. requet header里user key保存的是token
	claims := user.(*jwt.Token).Claims//cast raw token string to jwt.token
	username := claims.(jwt.MapClaims)["username"]
	
	p := Post{
		Id: uuid.New(),
		User: username.(string),
		Message: r.FormValue("message"),
	}
	file, header, err := r.FormFile("media_file")
	if err != nil {
		http.Error(w, "Media file is not available", http.StatusBadRequest)
		fmt.Printf("Media file is not available %v\n", err)
		return
	}
	suffix := filepath.Ext(header.Filename)//读取文件后缀
	if t, ok := mediaTypes[suffix]; ok {
		p.Type = t
	} else {
		p.Type = "unknown"
	}
	err = savePost(&p, file)//saveToGCS and saveToES
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

func signinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signin request")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}

	//Get User information from client
	decoder := json.NewDecoder(r.Body)
	var user User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode user data from client %v\n", err)
		return
	}
	exists, err := checkUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, "Failed to read user from Elasticsearch", http.StatusInternalServerError)
		fmt.Printf("Failed to read user from Elasticsearch %v/n", err)
		return
	}
	if !exists {
		http.Error(w, "User doesn't exists or wrong password", http.StatusUnauthorized)
		fmt.Printf("User doesn't exists or wrong password\n")
		return
	}

	//claim ~= payload
	//三个参数：加密算法，payload，private key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp": time.Now().Add(time.Hour*24).Unix(),//Unix()代表unix timestamp. 1970/1/1到现在多少秒
	})
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		fmt.Printf("Failed to generate token %v\n", err)
		return
	}
	w.Write([]byte(tokenString))
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signup request")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTION" {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var user User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
		fmt.Println("Cannot decode user data from client %v\n", err)
		return
	}
	//sanity check
	//username 只能小写字母和数字。^开头 $结尾
	if user.Username == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9]$`).MatchString(user.Username) {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		fmt.Printf("Invalid username or password")
		return
	}
	success, err := addUser(&user)
	if err != nil {
		http.Error(w, "Failed to save user to Elasticsearch", http.StatusInternalServerError)
		fmt.Printf("Failed to save user to Elasticsearch %v\n", err)
		return
	}
	if !success {
		http.Error(w, "User already exists", http.StatusBadRequest)
		fmt.Println("User already exists")
		return
	}
	fmt.Printf("User added successfully: %s\n", user.Username)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one delete for search")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Conrtrol-Allow-Headers", "Content-Type, Authorization")
	if r.Method == "OPTIONS" {
		return
	}

	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"].(string)
	id := mux.Vars(r)["id"]//通过mux解析url request里的 post Id

	if err := deletePost(id, username); err != nil {
		http.Error(w, "Failed to delete post from Elasticsearch", http.StatusInternalServerError)
		fmt.Printf("Failed to delete post from Elasticsearch %v\n", err)
		return
	}
	fmt.Println("Post is deleted successfully")
}