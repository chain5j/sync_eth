// description: sync_eth 
// 
// @author: xwc1125
// @date: 2020/10/05
package models

import (
	"encoding/json"
	"github.com/chain5j/sync_eth/params/global"
)

type ListenAddress struct {
	Address string `json:"address" es:"keyword"` // address
	Remark  string `json:"remark" es:"text"`     // memo
}

func (a ListenAddress) TableName() string {
	if global.Config != nil {
		return global.Config.ChainConfig.ChainName + "_listen_address"
	} else {
		return "dev_eth" + "_listen_address"
	}
}

func (a *ListenAddress) String() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}
