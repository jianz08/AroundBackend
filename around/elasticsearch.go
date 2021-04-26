package main

import (
	"context"
	"github.com/olivere/elastic/v7"
)

func readFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
	client, err := elastic.NewClient(//ES实际是建立了一个connection pool
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth(ES_Username, ES_Password))
	if err != nil {
		return nil, err
	}
	//可以在query里加pagination的设置
	searchResult, err := client.Search().
						Index(index).Query(query).Pretty(true).Do(context.Background())//searchResult 是一个pointer
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func saveToES(i interface{}, index string, id string) error {
	//这里用i interface{} 而不用post *Post， 是为了扩展性，比如也可以save user index
	//Hibernate is also an implementation of the Java Persistence API (JPA)
	//Hibernate ORM is an object–relational mapping（ORM） tool for the Java programming language. 
	//It provides a framework for mapping an object-oriented domain model to a relational database
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth(ES_Username, ES_Password))
	if err != nil {
		return err
	}
	
	_, err = client.Index().//Index 代表 Insert
			Index(index).Id(id).BodyJson(i).Do(context.Background())
	return err
}

func deleteFromES(query elastic.Query, index string) error{
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth(ES_Username, ES_Password))
	if err != nil {
		return err
	}
	_, err = client.DeleteByQuery().
			Index(index).Query(query).Pretty(true).Do(context.Background())
	return err
}