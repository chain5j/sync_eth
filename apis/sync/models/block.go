// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/10/05
package models

import (
	"encoding/json"
	"github.com/chain5j/chain5j-pkg/util/dateutil"
	"github.com/chain5j/sync_eth/params/global"
	"github.com/shopspring/decimal"
	"time"
)

type Block struct {
	BlockNumber uint64          `json:"block_number" es:"long"`
	BlockHash   string          `json:"block_hash" es:"text"`
	ParentHash  string          `json:"parent_hash" es:"text"`
	Miner       string          `json:"miner" es:"text"`
	BlockTime   uint64          `json:"block_time" es:"long"`
	BlockAward  decimal.Decimal `json:"block_award" es:"double"`
	TxCount     int             `json:"tx_count" es:"int"`
	BlockSize   uint64          `json:"block_size" es:"long"`
	Txs         []string        `json:"txs,omitempty" es:"text"`
	Timestamp   time.Time       `json:"timestamp" es:"date"`
}

func (a Block) TableName() string {
	return global.Config.ChainConfig.ChainName + "_block"
}

func (a *Block) String() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}

func (a *Block) MarshalJSON() ([]byte, error) {
	a.Timestamp = dateutil.SecondToTime(int64(a.BlockTime))
	return json.Marshal(*a)
}
