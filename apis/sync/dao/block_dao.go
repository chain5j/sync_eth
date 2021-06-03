// description: sync_eth
//
// @author: xwc1125
// @date: 2020/10/05
package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/chain5j/log15"
	"github.com/chain5j/sync_eth/apis/sync/models"
	"github.com/chain5j/sync_eth/pkg/database/es"
	"github.com/chain5j/sync_eth/pkg/util/convutil"
	"github.com/olivere/elastic/v7"
	"sync"
)

var BlockDao = newBlockDao()
var lock sync.RWMutex

type blockDao struct {
}

func newBlockDao() *blockDao {
	return &blockDao{}
}

// LatestBlock get latest block
func (d *blockDao) LatestBlock(es *es.ES) (*models.Block, error) {
	searchResult, err := es.Client().Search().
		Index(models.Block{}.TableName()).
		//Query(elastic.NewMatchAllQuery()).
		Sort("block_number", false). // sort
		From(0).Size(1).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		log.Error("LatestBlock is err", "err", err)
		return nil, err
	}
	if searchResult.TotalHits() > 0 {
		for _, hit := range searchResult.Hits.Hits {
			latestBlock := new(models.Block)
			err := json.Unmarshal(hit.Source, latestBlock)
			if err != nil {
				return nil, err
			}
			return latestBlock, nil
		}
	}
	return nil, nil
}

// GetBlockByHeight get block by height
func (d *blockDao) GetBlockByHeight(es *es.ES, blockNumber uint64) (*models.Block, error) {
	result, err := es.Client().Get().
		Index(models.Block{}.TableName()).
		Id(convutil.ToString(blockNumber)).
		Do(context.Background())
	if err != nil {
		switch {
		case elastic.IsNotFound(err):
			log.Error("NotFound by GetBlockByHeight", "blockNumber", blockNumber, "err", fmt.Sprintf("Document not found: %v", err))
			return nil, err
		case elastic.IsTimeout(err):
			log.Error("Timeout by GetBlockByHeight", "blockNumber", blockNumber, "err", fmt.Sprintf("Timeout retrieving document: %v", err))
			return nil, err
		case elastic.IsConnErr(err):
			log.Error("ConnErr by GetBlockByHeight", "blockNumber", blockNumber, "err", fmt.Sprintf("Connection problem: %v", err))
			return nil, err
		default:
			return nil, err
		}
	}
	block := new(models.Block)
	err = json.Unmarshal(result.Source, block)
	if err != nil {
		return nil, err
	}
	return block, nil
}

// BatchSave batch save block and txs
func (d *blockDao) BatchSave(es *es.ES, blockInfo *models.Block, contract *models.ContractMapSet, tempRawTxList []*models.Transaction) error {
	lock.Lock()
	defer lock.Unlock()
	bulk := es.Client().Bulk()
	// ===============================================
	// start insert to db
	// save block
	request := es.BulkIndexRequest(blockInfo.TableName()).
		Id(convutil.ToString(blockInfo.BlockNumber)).
		Doc(blockInfo)
	bulk.Add(request)

	// save contract
	if contract != nil && contract.Len() > 0 {
		for _, v := range contract.List() {
			request := es.BulkIndexRequest(v.TableName()).
				Id(convutil.ToString(v.Address)).
				Doc(v)
			bulk.Add(request)
		}
	}

	// save transaction
	if tempRawTxList != nil && len(tempRawTxList) > 0 {
		for _, tx := range tempRawTxList {
			request := es.BulkIndexRequest(tx.TableName()).
				Id(tx.TxHash).
				Doc(tx)
			bulk.Add(request)
		}
	}
	// ===============================================
	bulkResponse, err := bulk.Refresh("true").Do(context.Background())
	if err != nil {
		return err
	}
	if bulkResponse == nil {
		err = errors.New("expected bulkResponse to be != nil; got nil")
		return err
	}
	if bulkResponse.Errors {
		var buffer bytes.Buffer
		for _, item := range bulkResponse.Items {
			for _, i := range item {
				if i.Error != nil {
					fmt.Println(i.Error)
					bytes, _ := json.Marshal(i.Error)
					buffer.WriteString(string(bytes) + "\n")
				}
			}
		}
		// TODO need to rollback
		d.Rollback(es, blockInfo.BlockNumber)
		return errors.New(buffer.String())
	}
	return nil
}

// Rollback rollback
func (d *blockDao) Rollback(es *es.ES, blockNumber uint64) error {
	lock.Lock()
	defer lock.Unlock()
	// ===============================================
	// del block
	deleteResponse, err := es.Client().Delete().Index(models.Block{}.TableName()).
		Id(convutil.ToString(blockNumber)).
		Refresh("true").
		Do(context.Background())
	if err != nil {
		switch {
		case elastic.IsNotFound(err):
			log.Info("NotFound by Rollback", "blockNumber", blockNumber, "err", fmt.Sprintf("Document not found: %v", err))
		default:
			return err
		}
	} else {
		if deleteResponse == nil {
			log.Info("NotResponse by Rollback", "blockNumber", blockNumber)
		} else {
			marshal, _ := json.Marshal(deleteResponse)
			log.Info("rollback block", "blockNumber", blockNumber, "deleteResponse", string(marshal))
		}
	}

	// del transaction
	q := elastic.NewTermQuery("block_number", blockNumber)
	res, err := es.Client().DeleteByQuery().
		Index(models.Transaction{}.TableName()).
		Query(q).
		Slices("auto").
		Refresh("true").
		Pretty(true).
		Do(context.Background())
	if err != nil {
		log.Error("del transaction by rollback", "blockNumber", blockNumber, "err", err)
		return err
	}
	if res == nil {
		log.Info("del transaction no resp by rollback", "blockNumber", blockNumber)
	} else {
		marshal, _ := json.Marshal(res)
		log.Info("rollback transaction", "blockNumber", blockNumber, "res", string(marshal))
	}
	return nil
}

// VerifyBlock verify block
func (d *blockDao) VerifyBlock(es *es.ES, sysStartBlockNumber uint64, latestBlock *models.Block) (blockInfo *models.Block, err error) {
	var (
		latestBlockNumber = latestBlock.BlockNumber
	)
	if latestBlock == nil {
		return nil, errors.New("latestBlock is empty")
	}
	if sysStartBlockNumber >= latestBlockNumber {
		return latestBlock, nil
	}
	for {
		block, err := d.GetBlockByHeight(es, latestBlockNumber-1)
		if err != nil {
			if elastic.IsNotFound(err) {
			} else {
				return latestBlock, err
			}
		}
		if block == nil {
			d.Rollback(es, latestBlockNumber)
			latestBlockNumber = latestBlockNumber - 1
			latestBlock = block
			continue
		}
		if latestBlock == nil {
			latestBlockNumber = latestBlockNumber - 1
			latestBlock = block
			continue
		}
		if latestBlock.ParentHash != block.BlockHash {
			d.Rollback(es, latestBlockNumber)
			latestBlockNumber = latestBlockNumber - 1
			latestBlock = block
			if sysStartBlockNumber >= latestBlockNumber {
				return latestBlock, nil
			}
		} else {
			return latestBlock, nil
		}
	}
	return nil, nil
}
