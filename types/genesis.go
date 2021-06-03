// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/10/05
package types

import (
	"github.com/chain5j/chain5j-pkg/math"
)

type GenesisAlloc map[string]string

type GenesisAccount struct {
	//Code       []byte                    `json:"code,omitempty"`
	//Storage    map[types.Hash]types.Hash `json:"storage,omitempty"`
	Balance *math.HexOrDecimal256 `json:"balance" gencodec:"required"`
	//Nonce      uint64                    `json:"nonce,omitempty"`
	//PrivateKey []byte                    `json:"secretKey,omitempty"` // for tests
}

//func (ga *GenesisAlloc) UnmarshalJSON(data []byte) error {
//	m := make(map[types.UnprefixedAddress]GenesisAccount)
//	if err := json.Unmarshal(data, &m); err != nil {
//		return err
//	}
//	*ga = make(GenesisAlloc)
//	for addr, a := range m {
//		(*ga)[types.Address(addr)] = a
//	}
//	return nil
//}
