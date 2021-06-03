// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/10/05
package dao

import (
	"database/sql"
	"github.com/chain5j/sync_eth/apis/sync/models"
	"github.com/chain5j/sync_eth/params/global"
	"github.com/chain5j/sync_eth/pkg/database/db_xorm"
)

var AccountsDao = newAccountsDao()

type accountsDao struct {
}

func newAccountsDao() *accountsDao {
	return &accountsDao{}
}

func (d accountsDao) GetAccount(address string, contract string) (*models.Accounts, error) {
	result := make([]*models.Accounts, 0)
	engine := db_xorm.MasterEngine(global.Config.Database.Master)
	err := engine.Where("address = ? and contract = ?", address, contract).Limit(1).Find(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if result == nil || len(result) == 0 {
		return nil, nil
	}
	return result[0], nil
}

func (d accountsDao) Insert(account *models.Accounts) error {
	engine := db_xorm.MasterEngine(global.Config.Database.Master)
	_, err := engine.InsertOne(account)
	return err
}
