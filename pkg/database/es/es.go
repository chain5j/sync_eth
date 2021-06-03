// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/10/05
package es

import (
	"context"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"time"
)

type ES struct {
	es *elastic.Client
}

func NewES(hosts []string) (*ES, error) {
	if hosts == nil || len(hosts) == 0 {
		return nil, errors.New("es hosts is empty")
	}
	client, err := elastic.NewClient(
		elastic.SetURL(hosts...),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetMaxRetries(5),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		//elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		return nil, err
	}
	info, code, err := client.Ping(hosts[0]).Do(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esVersion, err := client.ElasticsearchVersion(hosts[0])
	if err != nil {
		return nil, err
	}
	fmt.Printf("Elasticsearch version %s\n", esVersion)

	return &ES{
		es: client,
	}, nil
}

func (es *ES) Client() *elastic.Client {
	return es.es
}

func (es *ES) BulkIndexRequest(index string) *elastic.BulkIndexRequest {
	return elastic.NewBulkIndexRequest().Index(index)
}
