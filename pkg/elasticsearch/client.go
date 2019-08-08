package elasticsearch

import (
	"context"
	"encoding/json"
	"github.com/chandresh-pancholi/edith/model"
	es "github.com/chandresh-pancholi/edith/pkg/elasticsearch/config"
	"log"
	"strings"


	"github.com/elastic/go-elasticsearch/v7/esapi"
)

//IndexDocument is to index document in Elasticsearch
func IndexDocument(esClient *es.ESClient, index, documentType, documentID string, data interface{}) error {
	//https://github.com/elastifeed/es-pusher/blob/8fc8c086439a5e668ade1cd2d6ed71b68102cc78/poc.go
	//TODO Use a Waitgroup to ensure everything was sent to Elasticsearch before exiting the program. follow above link
	//var wg sync.WaitGroup

	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:        index,
		DocumentType: documentType,
		DocumentID:   documentID,
		Body:         strings.NewReader(string(out)),
		Refresh:      "true",
	}

	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), 101)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}
	return nil
}

//Search in elasticsearch index
func Search(esClient *es.ESClient, index, query, value string) (*model.ESResponse, error) {
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()), esClient.Search.WithIndex(index),
		esClient.Search.WithBody(buildQuery(query, value)),
		esClient.Search.WithTrackTotalHits(true),
		esClient.Search.WithPretty(),
		esClient.Search.WithHuman(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result model.ESResponse
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
