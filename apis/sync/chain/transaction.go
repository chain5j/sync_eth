// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/10/05
package chain

import (
	"encoding/json"
	"github.com/chain5j/chain5j-pkg/types"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
)

type Transaction struct {
	Hash             string         `json:"hash"`
	Nonce            hexutil.Uint64 `json:"nonce"`
	BlockHash        string         `json:"blockHash"`
	BlockNumber      hexutil.Uint64 `json:"blockNumber"`
	TransactionIndex hexutil.Uint64 `json:"transactionIndex"`
	From             types.Address  `json:"from"`
	To               *types.Address `json:"to"`
	Value            *hexutil.Big   `json:"value"`
	GasPrice         *hexutil.Big   `json:"gasPrice"`
	GasLimit         hexutil.Uint64 `json:"gas"`
	Input            hexutil.Bytes  `json:"input"`
}

func (tx *Transaction) Bytes() ([]byte, error) {
	return json.Marshal(tx)
}
