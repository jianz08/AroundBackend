package main

import (
	"context"
	"github.com/olivere/elastic/v7"
)

const (
	ES_URL = "http://10.128.0.2:9200"
)

func readFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth("elastic", "Pass123!"))
	if err != nil {
		return nil, err
	}
	searchResult, err := client.Search().Index(index).Query(query).Pretty(true).Do(context.Background())
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}