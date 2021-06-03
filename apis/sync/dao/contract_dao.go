// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/10/05
package dao

import (
	"context"
	"encoding/json"
	"github.com/chain5j/sync_eth/apis/sync/models"
	"github.com/chain5j/sync_eth/pkg/database/es"
)

var ContractDao = newContractDao()

type contractDao struct {
}

func newContractDao() *contractDao {
	return &contractDao{}
}

func (d contractDao) Info(es *es.ES, contractAddress string) (*models.Contract, error) {
	searchResult, err := es.Client().Get().
		Index(models.Contract{}.TableName()).
		Id(contractAddress).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	if searchResult == nil {
		return nil, nil
	}
	contract := new(models.Contract)
	err = json.Unmarshal(searchResult.Source, contract)
	if err != nil {
		return nil, err
	}
	return contract, nil
}
