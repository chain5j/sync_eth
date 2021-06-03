// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/10/05
package chain

import (
	"github.com/chain5j/chain5j-pkg/codec/rlp"
	"github.com/chain5j/chain5j-pkg/crypto/keccak"
	"github.com/chain5j/chain5j-pkg/types"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	log "github.com/chain5j/log15"
	"math/big"
)

type RawTransaction struct {
	Nonce   uint64         `json:"nonce"    gencodec:"required"`
	GaPrice *big.Int       `json:"gasPrice" gencodec:"required"`
	Gas     uint64         `json:"gas"      gencodec:"required"`
	To      *types.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Value   *big.Int       `json:"value"    gencodec:"required"`
	Input   []byte         `json:"input"    gencodec:"required"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`
}

func (t *RawTransaction) GetRawTx() []byte {
	bytes, err := rlp.EncodeToBytes(t)
	if err != nil {
		log.Error("GetRawTx err", "err", err)
	}
	return bytes
}
func (t *RawTransaction) GetTxHash() string {
	return hexutil.Bytes2Hex(keccak.Keccak256(t.GetRawTx()))
}
