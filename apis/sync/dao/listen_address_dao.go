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
	"github.com/olivere/elastic/v7"
	"strings"
	"sync"
)

var ListenAddressDao = newListenAddressDaoDao()

type listenAddressDao struct {
}

func newListenAddressDaoDao() *listenAddressDao {
	return &listenAddressDao{}
}

func (d *listenAddressDao) GetListenAddress(es *es.ES, address string) (*models.ListenAddress, error) {
	result, err := es.Client().Get().
		Index(models.ListenAddress{}.TableName()).
		Id(strings.ToLower(address)).
		Do(context.Background())
	if err != nil {
		switch {
		case elastic.IsNotFound(err):
			log.Debug("NotFound by GetListenAddress", "address", address, "err", fmt.Sprintf("Document not found: %v", err))
			return nil, nil
		case elastic.IsTimeout(err):
			log.Error("Timeout by GetListenAddress", "address", address, "err", fmt.Sprintf("Timeout retrieving document: %v", err))
			return nil, err
		case elastic.IsConnErr(err):
			log.Error("ConnErr by GetListenAddress", "address", address, "err", fmt.Sprintf("Connection problem: %v", err))
			return nil, err
		default:
			return nil, err
		}
	}
	block := new(models.ListenAddress)
	err = json.Unmarshal(result.Source, block)
	if err != nil {
		return nil, err
	}
	return block, nil
}

var listenAddressLock sync.RWMutex

// BatchAddListenAddress batch add listen address
func (d *listenAddressDao) BatchAddListenAddress(es *es.ES, addressList []*models.ListenAddress) error {
	listenAddressLock.Lock()
	defer listenAddressLock.Unlock()
	bulk := es.Client().Bulk()
	// ===============================================
	// save address
	if addressList != nil && len(addressList) > 0 {
		for _, address := range addressList {
			request := es.BulkIndexRequest(address.TableName()).
				Id(strings.ToLower(address.Address)).
				Doc(address)
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
		return errors.New(buffer.String())
	}
	return nil
}
