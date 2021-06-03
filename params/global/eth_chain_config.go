// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/10/05
package global

import (
	"github.com/chain5j/sync_eth/types"
)

type EthChainConfig struct {
	ChainName        string             `json:"chainName" mapstructure:"chainName"`
	ChainType        int64              `json:"chainType" mapstructure:"chainType"`
	Host             string             `json:"host" mapstructure:"host"`
	ChainId          int64              `json:"chainId" mapstructure:"chainId"`
	ClientIdentifier string             `json:"clientIdentifier" mapstructure:"clientIdentifier"`
	IsEip155         bool               `json:"isEip155" mapstructure:"isEip155"`
	GasPrice         uint64             `json:"gasPrice" mapstructure:"gasPrice"`
	Gas              uint64             `json:"ga" mapstructure:"gas"`
	To               string             `json:"to" mapstructure:"to"`
	Contract         string             `json:"contract" mapstructure:"contract"`
	Value            uint64             `json:"value" mapstructure:"value"`
	Input            string             `json:"input" mapstructure:"input"`
	Alloc            types.GenesisAlloc `json:"alloc" mapstructure:"alloc"`
	BaseDecimals     int32              `json:"baseDecimals" mapstructure:"baseDecimals"`
	MinConfirms      uint64             `json:"minConfirms" mapstructure:"minConfirms"`       // the mini confirms
	SyncStartBlock   int64              `json:"syncStartBlock" mapstructure:"syncStartBlock"` // start sync block
}
