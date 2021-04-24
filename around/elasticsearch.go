package main

import (
	"context"
	"github.com/olivere/elastic/v7"
)

const (
	ES_URL = "http://10.128.0.2:9200"
)

func readFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
	client, err := elastic.NewClient(//ES实际是建立了一个connection pool
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth("elastic", "Pass123!"))
	if err != nil {
		return nil, err
	}
	//可以在query里加pagination的设置
	searchResult, err := client.Search().Index(index).Query(query).Pretty(true).Do(context.Background())//searchResult 是一个pointer
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func saveToES(i interface{}, index string, id string) error {
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth("elastic", "Pass123!"))
	if err != nil {
		return err
	}
	_, err = client.Index().Index(index).Id(id).BodyJson(i).Do(context.Background())
	return err
}