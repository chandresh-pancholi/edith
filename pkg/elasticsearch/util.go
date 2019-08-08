package elasticsearch

import (
	"encoding/json"
	"fmt"
	"github.com/chandresh-pancholi/edith/model"
	"io"
	"strings"

)

const searchAll = `
	"query" : { "match_all" : {} },
	"size" : 25,
	"sort" : { "published" : "desc", "_doc" : "asc" }`

const searchMatch = `
	"query" : {
		"match" : {
			%q : %q
		}
	}`

func buildQuery(query string, value ...string) io.Reader {
	var b strings.Builder

	b.WriteString("{\n")

	//if query == "" {
	//	b.WriteString(searchAll)
	//} else {
	//
	//}

	if len(value) > 0 && value[0] != "" && value[0] != "null" {
		b.WriteString(fmt.Sprintf(searchMatch, query, value[0]))
		//b.WriteString(",\n")
		//b.WriteString(fmt.Sprintf(`	"search_after": %s`, after))
	}

	b.WriteString("\n}")

	// fmt.Printf("%s\n", b.String())
	return strings.NewReader(b.String())
}

func ConsumerSearchResult(esResponse *model.ESResponse) ([]model.Consumer, error) {
	var consumerList []model.Consumer
	data := esResponse.Hits.Hits
	for _, hits := range data {

		byt, err := json.Marshal(hits.Source)
		if err != nil {
			//logger.Error("es result marshalling failed", zap.String("id", hits.ID), zap.Error(err))

			return nil, err
		}

		var client model.Consumer
		err = json.Unmarshal(byt, &client)
		if err != nil {
			//logger.Error("es source to MSG client mapping failed", zap.String("source", string(byt)), zap.Error(err))
			return nil, err
		}

		consumerList = append(consumerList, client)
	}

	return consumerList, nil
}
