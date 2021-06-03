// @author: xwc1125
// @date: 2020/10/05
package rpc

import (
	"github.com/chain5j/chain5j-pkg/network/rpc"
)

func NewRpc(host string) (JsonRpc, error) {
	return rpc.Dial(host)
}
