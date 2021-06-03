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

type TransactionReceipt struct {
	TransactionHash   string          `json:"transactionHash"`
	TransactionIndex  hexutil.Uint64  `json:"transactionIndex"`
	BlockHash         string          `json:"blockHash"`
	BlockNumber       hexutil.Uint64  `json:"blockNumber"`
	From              types.Address   `json:"from"`
	To                *types.Address  `json:"to"`
	CumulativeGasUsed hexutil.Uint64  `json:"cumulativeGasUsed"`
	GasUsed           hexutil.Uint64  `json:"gasUsed"`
	ContractAddress   *types.Address  `json:"contractAddress"`
	Logs              json.RawMessage `json:"logs"`
	Status            hexutil.Uint64  `json:"status"`
}

func (t *TransactionReceipt) Bytes() ([]byte, error) {
	return json.Marshal(t)
}
