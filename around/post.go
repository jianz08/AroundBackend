package main

import (
	"reflect"
	"github.com/olivere/elastic/v7"
	"mime/multipart"
)

const (
	POST_INDEX = "post"
)

type Post struct {
	Id string `json:"id"`
	User string `json:"user"`
	Message string `json:"message"`
	Url string `json:"url"`
	Type string `json:"type"`
}

func searchPostsByUser(user string) ([]Post, error) {
	query := elastic.NewTermQuery("user", user)//select user = user
	searchResult, err := readFromES(query, POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

func searchPostsByKeywords(keywords string) ([]Post, error) {
	query := elastic.NewMatchQuery("message", keywords)//keyword match query, fuzzy matching。在message里match
	query.Operator("AND")//'AND'表示好几个关键词搜索结果做交集,还可以"OR"做并集
	if keywords == "" {//如果keywords为空
		query.ZeroTermsQuery("all")//返回所有结果
	}
	searchResult, err := readFromES(query, POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

func getPostFromSearchResult(searchResult *elastic.SearchResult) []Post {
	var ptype Post
	var posts []Post
	for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {//sanity check。reflect.TypeOf判断数据类型是否正确，是就操作，不是就跳过
		p := item.(Post)//cast成Post类型，type assertion
		posts = append(posts, p)
	}
	return posts
}

func savePost(post *Post, file multipart.File) error {
	medialink, err := saveToGCS(file, post.Id)
	if err != nil {
		return err
	}
	post.Url = medialink
	return saveToES(post, POST_INDEX, post.Id)
}
