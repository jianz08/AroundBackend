package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

const (
	POST_INDEX = "post"
	USER_INDEX = "user"
	ES_URL = "http://10.128.0.2:9200"//internal IP
)

func main() {
	client, err := elastic.NewClient(//建立连接
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth("elastic", "Pass123!"))//如何不在这里显示密码，把密码写在一个config文件里
	if err != nil {
		panic(err)
	}

	exists, err := client.IndexExists(POST_INDEX).Do(context.Background())//判断index是否存在
	//Do代表执行，context.Background()表示没有额外参数，直到等待结果结束。可以加的参数比如deadline, cancel，callback func 
	if err != nil {
		panic(err)
	}

	if !exists {
		//mapping = schema
		//properties = column
		//keyword = 单词，要求完全匹配，full match
		//text = 句子，不要求完全匹配，fuzzy match
		//index 默认是true，建立索引
		//user 和 USER_INDEX 没有联系，NoSQL 非关系型数据库。不像在SQL里user是foreign key
		mapping := `{
			"mappings": {
				"properties": {
					"id": {"type": "keyword" },
					"user": {"type": "keyword" },
					"message": {"type": "text" },
					"url": {"type": "keyword", "index": false },
					"type": {"type": "keyword", "index": false }
				}	
			}
		}`
		_, err := client.CreateIndex(POST_INDEX).Body(mapping).Do(context.Background())//创建index
		if err != nil {
			panic(err)
		}
	}

	if !exists {
		mapping := `{
			"mappings": {
				"properties": {
					"username": {"type": "keyword" },
					"passoword": {"type": "keyword" },
					"age": {"type": "long", "index": false },
					"gender": {"type": "keyword", "index": false }
				}	
			}
		}`
		_, err := client.CreateIndex(USER_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Indexes are created.")
}