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

type Transaction struct {
	//Id       int64           `json:"id"`
	TxHash   string          `json:"tx_hash" es:"keyword"`
	From     string          `json:"from" es:"keyword"`
	To       string          `json:"to" es:"keyword"`
	Value    decimal.Decimal `json:"value" es:"double"`
	Nonce    decimal.Decimal `json:"nonce" es:"long"`
	GasPrice decimal.Decimal `json:"gas_price" es:"long"`
	GasLimit uint64          `json:"gas_limit" es:"long"`

	Contract string `json:"contract" es:"keyword"`
	TxType   int    `json:"tx_type" es:"integer"`
	Status   bool   `json:"status" es:"boolean"`

	BlockNumber      decimal.Decimal `json:"block_number" es:"long"`
	BlockTime        uint64          `json:"block_time" es:"long"`
	TransactionIndex uint64          `json:"transaction_index" es:"integer"`

	GasUsed decimal.Decimal `json:"gas_used" es:"long"`
	Fee     decimal.Decimal `json:"fee" es:"double"`

	Input     string    `json:"input" es:"input"`
	Timestamp time.Time `json:"timestamp" es:"date"`
}

func (a Transaction) TableName() string {
	return global.Config.ChainConfig.ChainName + "_transaction"
}

func (a *Transaction) MarshalJSON() ([]byte, error) {
	a.Timestamp = dateutil.SecondToTime(int64(a.BlockTime))
	return json.Marshal(*a)
}
