// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/10/05
package chain

import (
	"encoding/json"
	"fmt"
	"github.com/chain5j/chain5j-pkg/math"
	"github.com/chain5j/chain5j-pkg/types"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestEth_GetLatestBlock(t *testing.T) {
	eth, _ := NewETH("http://127.0.0.1:7545", "eth", 1, false)
	var blockInfo Block
	eth.GetLatestBlock(true, &blockInfo)
	spew.Dump(blockInfo)
	var block2 BlockTxHashes
	eth.GetLatestBlock(false, &block2)
	spew.Dump(block2)
}

func TestEth_GetBlockByNumber(t *testing.T) {
	eth, _ := NewETH("http://127.0.0.1:8545", "eth", 1587545710, true)
	var raw json.RawMessage
	err := eth.GetBlockByNumber(0, true, &raw)
	if err != nil {
		panic(err)
	}
	spew.Dump(raw)

	var blockInfo Block
	eth.GetBlockByNumber(0, true, &blockInfo)
	spew.Dump(blockInfo)
}

func TestEth_GetTransactionByHash(t *testing.T) {
	eth, _ := NewETH("http://127.0.0.1:7545", "eth", 1, false)
	var tx Transaction
	eth.GetTransactionByHash("0x97ca16d41c59c80e2ff7a2d7f8bde1a4f5a3abcf79a71f98fec2d2999ae9ccb9", &tx)
	spew.Dump(tx)
}

func TestEth_GetTransactionReceipt(t *testing.T) {
	eth, _ := NewETH("http://127.0.0.1:7545", "eth", 1587545710, true)
	var tx TransactionReceipt
	eth.GetTransactionReceipt("0x4f98a47f864aec943abdf91b735fa720a401bfbee1ae1d0b43de4544199f3765", &tx)
	spew.Dump(tx)
}

func TestEth_GetBalance(t *testing.T) {
	eth, _ := NewETH("http://127.0.0.1:7545", "eth", 1587545710, true)
	var balance *math.HexOrDecimal256
	eth.GetBalance("0xaeff996f0efb374fcf95eb6b38fd4aa5e4bbc1b1", &balance)
	spew.Dump(balance)
	eth.GetBalance("0x2f4bb54f039ebd5e3476cb82c6357e43d3080c37", &balance)
	spew.Dump(balance)
	eth.GetBalance("0x9254e62fbca63769dfd4cc8e23f630f0785610ce", &balance)
	spew.Dump(balance)
}

func TestInput(t *testing.T) {
	input := "0xa9059cbb" +
		"000000000000000000000000aeff996f0efb374fcf95eb6b38fd4aa5e4bbc1b1" +
		"0000000000000000000000000000000000000000000000000000000005f5e100"
	method := input[:10]
	fmt.Println(method)
	to := input[10:74]
	fmt.Println(to)
	value := input[74:]
	fmt.Println(value)
	if method == "0xa9059cbb" {
		fmt.Println("transfer(address _to, uint256 _value)")
	}
	toAddress := types.HexToAddress(to)
	fmt.Println(toAddress.Hex())
}

func TestEth_GetTokenBalance(t *testing.T) {
	eth, _ := NewETH("http://127.0.0.1:9545", "eth", 1587545710, true)
	var balance *math.HexOrDecimal256
	eth.GetTokenBalance("0xa9f168bd6ef64a5a3788f04d18132e7ca60e1f5a", "0xb127c4a462a23f6f2b52027422fc64291cdd568b", nil, &balance)
	fmt.Println(balance)
}

func TestEth_GetTokenTotalSupply(t *testing.T) {
	eth, _ := NewETH("http://127.0.0.1:9545", "eth", 1587545710, true)
	totalSupply, err := eth.GetTokenTotalSupply("0x71c83e8e6b581a50614346c2a77cb1606ca9aaf8")
	if err != nil {
		panic(err)
	}
	fmt.Println(totalSupply)
}

func TestEth_GetTokenDecimals(t *testing.T) {
	eth, _ := NewETH("http://127.0.0.1:9545", "eth", 1587545710, true)
	decimals, err := eth.GetTokenDecimals("0x71c83e8e6b581a50614346c2a77cb1606ca9aaf8")
	if err != nil {
		panic(err)
	}
	fmt.Println(decimals)
}

func TestEth_GetTokenName(t *testing.T) {
	eth, _ := NewETH("http://127.0.0.1:9545", "eth", 1587545710, true)
	name, err := eth.GetTokenName("0x71c83e8e6b581a50614346c2a77cb1606ca9aaf8")
	if err != nil {
		panic(err)
	}
	fmt.Println(name)
}

func TestEth_GetTokenSymbol(t *testing.T) {
	eth, _ := NewETH("http://127.0.0.1:9545", "eth", 1587545710, true)
	name, err := eth.GetTokenSymbol("0x210e31dedeab9eb1b280e8b60b5e5ff81c0136c2")
	if err != nil {
		panic(err)
	}
	fmt.Println(name)
}
