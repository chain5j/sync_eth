// description: sync_eth
//
// @author: xwc1125
// @date: 2020/10/05
package chain

import (
	"encoding/json"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
)

// BlockHeader block header
type BlockHeader struct {
	Number     hexutil.Uint64 `json:"number"`
	Hash       string         `json:"hash"`
	ParentHash string         `json:"parentHash"`
	Coinbase   string         `json:"coinbase"`
	Size       hexutil.Uint64 `json:"size"`
	Timestamp  hexutil.Uint64 `json:"timestamp"`
	GasLimit   hexutil.Uint64 `json:"gasLimit"`
	GasUsed    hexutil.Uint64 `json:"gasUsed"`
}

// Block block with transactions
type Block struct {
	BlockHeader
	Transactions []*Transaction `json:"transactions"`
}

// BlockTxHashes block with txHashes
type BlockTxHashes struct {
	BlockHeader
	Transactions []string `json:"transactions"`
}

// Block2BlockTxHashes block to blockTxHashes
func Block2BlockTxHashes(b *Block) *BlockTxHashes {
	if b == nil {
		return nil
	}
	b2 := &BlockTxHashes{
		BlockHeader: b.BlockHeader,
	}
	if b.Transactions != nil && len(b.Transactions) > 0 {
		for _, tx := range b.Transactions {
			b2.Transactions = append(b2.Transactions, tx.Hash)
		}
	}
	return b2
}

// Bytes blockTxHashes json bytes
func (b *BlockTxHashes) Bytes() ([]byte, error) {
	return json.Marshal(b)
}
