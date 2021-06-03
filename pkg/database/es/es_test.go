// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/8/14 0014
package es

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/chain5j/log15"
	"github.com/olivere/elastic/v7"
	"reflect"
	"testing"
	"time"
)

func TestNewES(t *testing.T) {
	es, err := NewES([]string{"http://127.0.0.1:9200"})
	if err != nil {
		panic(err)
	}

	indexName := "twitter"
	exists, err := es.Client().IndexExists(indexName).Do(context.Background())
	if err != nil {
		log.Error("ES IndexExists is err", "err", err)
	}
	fmt.Println("exists", exists)
	if !exists {
		// table is not exist
		//		mapping := `
		//{
		//	"settings":{
		//		"number_of_shards":1,
		//		"number_of_replicas":0
		//	},
		//	"mappings":{
		//		"properties":{
		//				"user":{
		//					"type":"keyword"
		//				},
		//				"message":{
		//					"type":"text",
		//					"store": true,
		//					"fielddata": true
		//				},
		//                "retweets":{
		//                    "type":"long"
		//                },
		//				"tags":{
		//					"type":"keyword"
		//				},
		//				"location":{
		//					"type":"geo_point"
		//				},
		//				"suggest_field":{
		//					"type":"completion"
		//				}
		//			}
		//	}
		//}
		//`
		// create table
		createIndex, err := es.Client().CreateIndex(indexName).Body(Mapping(Tweet{})).Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
			fmt.Println("createIndex.Acknowledged", createIndex.Acknowledged)
		}
	}

	// write
	tweet1 := Tweet{
		User:     "A",
		Message:  "a",
		Retweets: 0,
		Image:    "http://127.0.0.1:22122/group1/big/upload/f61433fe81ab23917b58cca177f1f10b",
		Created:  time.Now(),
		Tags:     []string{"11"},
	}
	put1, err := es.Client().Index().
		Index(indexName).
		Id("2").
		BodyJson(tweet1).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	// get
	result, err := es.Client().Get().Index(indexName).Id("2").
		Do(context.Background())
	if err != nil {
		switch {
		case elastic.IsNotFound(err):
			panic(fmt.Sprintf("Document not found: %v", err))
		case elastic.IsTimeout(err):
			panic(fmt.Sprintf("Timeout retrieving document: %v", err))
		case elastic.IsConnErr(err):
			panic(fmt.Sprintf("Connection problem: %v", err))
		default:
			// Some other kind of error
			panic(err)
		}
	}
	fmt.Printf("Got document %s in version %d from index %s, type %s\n", result.Id, result.Version, result.Index, result.Type)

	//
	_, err = es.Client().Refresh().Index(indexName).Do(context.Background())
	if err != nil {
		panic(err)
	}

	//
	termQuery := elastic.NewTermQuery("user", "A")
	searchResult, err := es.Client().Search().
		Index(indexName).
		Query(termQuery). // specify the query
		Sort("user", true). // sort by "user" field, ascending
		From(0).Size(10). // take documents 0-9
		Pretty(true). // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	var ttyp Tweet
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		t := item.(Tweet)
		fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
	}

	//
	if searchResult.TotalHits() > 0 {
		fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

		for _, hit := range searchResult.Hits.Hits {
			var t Tweet
			err := json.Unmarshal(hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	} else {
		fmt.Print("Found no tweets\n")
	}

	//
	script := elastic.NewScript("ctx._source.retweets += params.num").Param("num", 1)
	update, err := es.Client().Update().Index(indexName).Id("1").
		Script(script).
		Upsert(map[string]interface{}{"retweets": 0}).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("New version of tweet %q is now %d", update.Id, update.Version)
	//
	deleteIndex, err := es.Client().DeleteIndex(indexName).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
	}
}
