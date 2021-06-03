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

type TransactionDetail struct {
	*Transaction

	CumulativeGasUsed hexutil.Uint64  `json:"cumulativeGasUsed"`
	GasUsed           hexutil.Uint64  `json:"gasUsed"`
	ContractAddress   *types.Address  `json:"contractAddress"`
	Logs              json.RawMessage `json:"logs"`
	//LogsBloom         hexutil.Bytes  `json:"logsBloom"`
	Status      bool                `json:"status"`
	InputFormat *TokenTransferInput `json:"inputFormat"`
}
