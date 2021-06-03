// description: sync_eth 
// 
// @author: xwc1125
// @date: 2020/10/05
package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chain5j/sync_eth/apis/sync/chain"
	"github.com/chain5j/sync_eth/apis/sync/dao"
	"github.com/chain5j/sync_eth/apis/sync/models"
	"github.com/chain5j/sync_eth/pkg/util/convutil"
	"github.com/olivere/elastic/v7"
	"github.com/shopspring/decimal"
	"testing"
)

var (
	syncClient, _ = NewSync([]string{"http://127.0.0.1:9200"})
)

func TestAddBlock(t *testing.T) {
	block := &models.Block{
		BlockNumber: 7,
		BlockHash:   "0xa911a79f7112fb4574a2584bfe74e168d4dced22a2996d6785a98988ace1050b",
		ParentHash:  "0xb8214b7254ba008240cfaddf92822070460710f0d1c8045f18de1573a91ee203",
		Miner:       "0x6e60f5243e1a3f0be3f407b5afe9e5395ee82aa2",
		BlockTime:   1597817607,
		BlockAward:  decimal.NewFromFloat(1.0),
		TxCount:     0,
		BlockSize:   539,
		Txs: []string{
			"0x2637fae015f832b962720692d9c5276f8c1b09ed7d811d06ffdfd0263e75c648",
			"0x2637fae015f832b962720692d9c5276f8c1b09ed7d811d06ffdfd0263e75c648",
		},
	}
	put, err := syncClient.es.Client().Index().
		Index(block.TableName()).
		Id(convutil.ToString(block.BlockNumber)).
		BodyJson(block).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("insert block:blockHeight: %d, %s to index %s, type %s\n", block.BlockNumber, put.Id, put.Index, put.Type)
}

func TestDelete(t *testing.T) {
	deleteResponse, err := syncClient.es.Client().Delete().
		Index(models.Block{}.TableName()).
		Id("6").
		Refresh("true").
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(deleteResponse)
}

func TestBulk(t *testing.T) {
	bulk := syncClient.es.Client().Bulk()

	block := &models.Block{
		BlockNumber: 1,
		BlockHash:   "0x2637fae015f832b962720692d9c5276f8c1b09ed7d811d06ffdfd0263e75c648",
		ParentHash:  "0xb8214b7254ba008240cfaddf92822070460710f0d1c8045f18de1573a91ee203",
		Miner:       "0x6e60f5243e1a3f0be3f407b5afe9e5395ee82aa2",
		BlockTime:   1597817607,
		BlockAward:  decimal.NewFromFloat(1.0),
		TxCount:     0,
		BlockSize:   539,
		Txs: []string{
			"0x2637fae015f832b962720692d9c5276f8c1b09ed7d811d06ffdfd0263e75c648",
			"0x2637fae015f832b962720692d9c5276f8c1b09ed7d811d06ffdfd0263e75c648",
		},
	}
	request := syncClient.es.BulkIndexRequest(block.TableName()).
		Id(block.BlockHash).
		Doc(block)
	bulk.Add(request)

	txs := make([]*models.Transaction, 0)
	txs = append(txs, &models.Transaction{
		TxHash:           "0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b",
		From:             "0x407d73d8a49eeb85d32cf465507dd71d507100c1",
		To:               "0x85h43d8a49eeb85d32cf465507dd71d507100c1",
		Value:            decimal.NewFromFloat(1),
		Nonce:            decimal.NewFromFloat(1),
		GasPrice:         decimal.NewFromFloat(100000000),
		GasLimit:         21000,
		GasUsed:          decimal.NewFromFloat32(3000000),
		Input:            "0x603880600c6000396000f300603880600c6000396000f3603880600c6000396000f360",
		BlockNumber:      decimal.NewFromFloat32(1),
		BlockTime:        1597817607,
		Status:           true,
		TransactionIndex: 0,
	}, &models.Transaction{
		TxHash:           "0x14ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b",
		From:             "0x407d73d8a49eeb85d32cf465507dd71d507100c1",
		To:               "0x85h43d8a49eeb85d32cf465507dd71d507100c1",
		Value:            decimal.NewFromFloat(1),
		Nonce:            decimal.NewFromFloat(1),
		GasPrice:         decimal.NewFromFloat(100000000),
		GasLimit:         21000,
		GasUsed:          decimal.NewFromFloat32(3000000),
		Input:            "0x603880600c6000396000f300603880600c6000396000f3603880600c6000396000f360",
		BlockNumber:      decimal.NewFromFloat32(1),
		BlockTime:        1597817607,
		Status:           true,
		TransactionIndex: 0,
	})
	for _, tx := range txs {
		request := syncClient.es.BulkIndexRequest(models.Transaction{}.TableName()).
			Id(tx.TxHash).
			Doc(tx)
		bulk.Add(request)
	}
	bulkResponse, err := bulk.Do(context.Background())
	if err != nil {
		panic(err)
	}
	if bulkResponse == nil {
		err = errors.New("expected bulkResponse to be != nil; got nil")
		panic(err)
	}
	if bulkResponse.Errors {
		for _, item := range bulkResponse.Items {
			for _, i := range item {
				fmt.Println(i.Error)
			}
		}
	}

}

func TestLatestBlock(t *testing.T) {
	searchResult, err := syncClient.es.Client().Search().
		Index(models.Block{}.TableName()).
		Query(elastic.NewMatchAllQuery()).
		Sort("block", false). // sort
		From(0).Size(1).
		Pretty(true).
		Do(context.Background())

	if err != nil {
		panic(err)
	}
	fmt.Println(searchResult)
	if searchResult.TotalHits() > 0 {
		for _, hit := range searchResult.Hits.Hits {
			block := new(models.Block)
			err := json.Unmarshal(hit.Source, block)
			if err != nil {
				panic(err)
			}
			fmt.Println(block)
		}
	}
}

func TestGetBlockByHeight(t *testing.T) {
	// 9208,1586
	block, err := dao.BlockDao.GetBlockByHeight(syncClient.es, 9208)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(block)
}

func TestRollback(t *testing.T) {
	err := dao.BlockDao.Rollback(syncClient.es, 9208)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("rollback success")
}

func TestVerifyBlock(t *testing.T) {
	latestBlock, err := dao.BlockDao.LatestBlock(syncClient.es)
	if err != nil {
		fmt.Println(err)
		return
	}
	block, err := dao.BlockDao.VerifyBlock(syncClient.es, 3, latestBlock)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(block)
}

func TestBigInt(t *testing.T) {
	a := decimal.NewFromFloat(1)
	div := a.Div(decimal.NewFromFloat(10000000000000000))
	fmt.Println(div)
}

func TestGetContractInfo(t *testing.T) {
	eth, err := chain.NewETH("http://127.0.0.1:7545", "eth", 5777, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	syncClient.eth = eth
	isContract, err := syncClient.eth.IsContract("0xCb10efC721268a27b521b8604ab069f56078e663")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("isContract", isContract)
}
