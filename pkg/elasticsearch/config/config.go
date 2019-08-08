package es

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/spf13/viper"
)

type ESClient struct {
	*elasticsearch.Client
}

func NewESClient() (*ESClient, error) {
	cfg := elasticsearch.Config{
		Addresses: viper.GetStringSlice("elasticsearch.nodes"),
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   50,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				// ...
			},
		},
	}
	es, err := elasticsearch.NewClient(cfg)

	esClient := &ESClient{Client: es}

	return esClient, err
}

func (es *ESClient) CreateIndex() {
	res, err := es.Indices.Exists([]string{"test"})
	if err != nil {

	}
	if res.StatusCode == http.StatusOK {

	}
	es.Index.WithRefresh("true")
	es.Indices.Create("test")
}

func (es *ESClient) IndexDocument(index, documentId, data string) error {
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: documentId,
		Body:       strings.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es)
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
