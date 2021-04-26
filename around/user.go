package main

import (
	"fmt"
	"reflect"
	"github.com/olivere/elastic/v7"
)

const (
	USER_INDEX = "user"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Age int64 `json:"age"`
	Gender string `json:"gender"`
}

func checkUser(username, password string) (bool, error) {//判断user是否存在
	
	query := elastic.NewBoolQuery()//select * from users where username = ? AND password = ?
	query.Must(elastic.NewTermQuery("username", username))
	query.Must(elastic.NewTermQuery("password", password))
	searchResult, err := readFromES(query, USER_INDEX)

	if err != nil {
		return false, err
	}
	// fmt.Printf("username = %v\n",username)
	// fmt.Printf("password = %v\n",password)
	// fmt.Printf("result size = %v\n",searchResult.TotalHits())

	var utype User
	for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
		u := item.(User)	
		if u.Password == password {
			fmt.Printf("Login as %s\n", username)
			return true, nil
		}		
	}
	return false, nil
}

func addUser(user *User) (bool, error) {
	query := elastic.NewTermQuery("username", user.Username)
	searchResult, err := readFromES(query, USER_INDEX)
	if err != nil {
		return false, err
	}
	if searchResult.TotalHits() > 0 {//如果不check有重复的话，会直接覆盖
		return false, nil
	}
	err = saveToES(user, USER_INDEX, user.Username)
	if err != nil {
		return false, err
	}
	fmt.Printf("User is added: %s\n", user.Username)
	return true, nil
}