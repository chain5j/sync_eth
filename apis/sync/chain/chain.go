// Package chain
// description: sync_eth
//
// @author: xwc1125
// @date: 2020/10/05
package chain

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/chain5j/chain5j-pkg/crypto/prime256v1"
	"github.com/chain5j/chain5j-pkg/math"
	"github.com/chain5j/chain5j-pkg/types"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"github.com/chain5j/log15"
	"github.com/chain5j/sync_eth/pkg/crypto"
	"github.com/chain5j/sync_eth/pkg/rpc"
	"math/big"
	"strings"
)

// Eth eth
type Eth struct {
	rpc              rpc.JsonRpc
	clientIdentifier string
	chainId          int64
	isEip155         bool
}

// NewETH new eth struct
func NewETH(host string, clientIdentifier string, chainId int64, isEip155 bool) (*Eth, error) {
	rpc, err := rpc.NewRpc(host)
	if err != nil {
		return nil, err
	}
	return &Eth{
		rpc:              rpc,
		clientIdentifier: clientIdentifier,
		chainId:          chainId,
		isEip155:         isEip155,
	}, nil
}

// GetBalance get balance
func (eth *Eth) GetBalance(address string, result interface{}) error {
	return eth.rpc.Call(&result, eth.clientIdentifier+"_getBalance", address, "latest")
}

// GetTokenBalance get token balance
func (eth *Eth) GetTokenBalance(contract, address string, blockNumber *big.Int, result interface{}) error {
	paramsMap := make(map[string]interface{})
	paramsMap["from"] = address
	paramsMap["to"] = contract
	input := "0x70a08231000000000000000000000000" + types.HexToAddress(address).Hex()[2:]
	paramsMap["data"] = input

	extraParam := "latest"
	if blockNumber != nil {
		extraParam = hexutil.IntToHex(blockNumber)
	}
	return eth.rpc.Call(&result, eth.clientIdentifier+"_call", paramsMap, extraParam)
}

// GetLatestBlock get latest block
func (eth *Eth) GetLatestBlock(isFullTx bool, result interface{}) error {
	return eth.rpc.Call(&result, eth.clientIdentifier+"_getBlockByNumber", "latest", isFullTx)
}

// GetBlockByNumber get block by height
// isFullTx whether return the full transaction
func (eth *Eth) GetBlockByNumber(height uint64, isFullTx bool, result interface{}) error {
	toHex := hexutil.IntToHex(height)
	return eth.rpc.Call(&result, eth.clientIdentifier+"_getBlockByNumber", toHex, isFullTx)
}

// GetBlockByHash get block by hash
func (eth *Eth) GetBlockByHash(height uint64, isFullTx bool, result interface{}) error {
	toHex := hexutil.IntToHex(height)
	return eth.rpc.Call(&result, eth.clientIdentifier+"_getBlockByHash", toHex, isFullTx)
}

// GetTransactionReceipt get tx receipt by hash
func (eth *Eth) GetTransactionReceipt(hash string, result interface{}) error {
	return eth.rpc.Call(&result, eth.clientIdentifier+"_getTransactionReceipt", hash)
}

// GetTransactionByHash get transaction by hash
func (eth *Eth) GetTransactionByHash(hash string, result interface{}) error {
	return eth.rpc.Call(&result, eth.clientIdentifier+"_getTransactionByHash", hash)
}

// TxStatus get transaction status
func TxStatus(txReceipt *TransactionReceipt, isToken bool) bool {
	if txReceipt == nil {
		return false
	}
	// 1 success or 0 failed
	if txReceipt.Status == 0 {
		return false
	}
	if isToken {
		if txReceipt.Logs == nil || len(txReceipt.Logs) == 0 || string(txReceipt.Logs) == "[]" {
			return false
		} else {
			log.Trace("txReceipt.Logs", "logs", string(txReceipt.Logs))
		}
	}
	return true
}

// SignTx sign rawTx
func (eth *Eth) SignTx(privateKey string, to *types.Address, value *big.Int, nonce uint64, gasPrice *big.Int, gasLimit uint64, input []byte) (*RawTransaction, error) {
	b, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	priKey, err := ToECDSA(b)

	var hashBytes []byte
	if eth.isEip155 {
		hash := crypto.RlpHash([]interface{}{
			nonce,
			gasPrice,
			gasLimit,
			to,
			value,
			input,
			big.NewInt(eth.chainId), // eip155[big.int]
			uint(0),                 // eip155
			uint(0),                 // eip155
		})
		hashBytes = hash[:]
	} else {
		hash := crypto.RlpHash([]interface{}{
			nonce,
			gasPrice,
			gasLimit,
			to,
			value,
			input,
		})
		hashBytes = hash[:]
	}
	signature, err := Sign(hashBytes, priKey)
	if err != nil {
		return nil, err
	}
	fmt.Println(hexutil.Bytes2Hex(signature))
	sig := &prime256v1.Signature{
		R: new(big.Int).SetBytes(signature[:32]),
		S: new(big.Int).SetBytes(signature[32:64]),
	}
	v := new(big.Int).SetBytes([]byte{signature[64] + 27})
	if eth.isEip155 {
		chainIdM := (eth.chainId << 1) + 8
		v = new(big.Int).Add(v, big.NewInt(chainIdM))
	}
	return &RawTransaction{
		Nonce:   nonce,
		GaPrice: gasPrice,
		Gas:     gasLimit,
		To:      to,
		Value:   value,
		Input:   input,
		V:       v,
		R:       sig.R,
		S:       sig.S,
	}, nil
}

// SendRawTransactionMethod send rawTx
func (eth *Eth) SendRawTransactionMethod(tx *RawTransaction) (string, error) {
	var result string
	rawTx := hexutil.Bytes2Hex(tx.GetRawTx())
	fmt.Println("rawTx", rawTx)
	err := eth.rpc.CallContext(context.Background(), &result, eth.clientIdentifier+"_sendRawTransaction", "0x"+rawTx)
	return result, err
}

// GetNonce get nonce
func (eth *Eth) GetNonce(address string) (hexutil.Uint64, error) {
	var result hexutil.Uint64
	err := eth.rpc.CallContext(context.Background(), &result, eth.clientIdentifier+"_getTransactionCount", address, "latest")
	return result, err
}

// GetCode get code
func (eth *Eth) GetCode(contract string) (hexutil.Bytes, error) {
	var result hexutil.Bytes
	err := eth.rpc.CallContext(context.Background(), &result, eth.clientIdentifier+"_getCode", contract, "latest")
	return result, err
}

// IsContract whether address is contract
func (eth *Eth) IsContract(address string) (bool, error) {
	code, err := eth.GetCode(address)
	if err != nil {
		return false, err
	}
	if code == nil || code.String() == "" || code.String() == "0x" {
		return false, nil
	}
	return true, nil
}

// GetTokenTotalSupply get token total supply
func (eth *Eth) GetTokenTotalSupply(contract string) (totalSupply big.Int, err error) {
	paramsMap := make(map[string]interface{})
	paramsMap["from"] = contract
	paramsMap["to"] = contract
	input := "0x" + "18160ddd"
	paramsMap["data"] = input

	extraParam := "latest"
	//var result hexutil.Bytes
	var result *math.HexOrDecimal256
	err = eth.rpc.Call(&result, eth.clientIdentifier+"_call", paramsMap, extraParam)
	if err != nil {
		return *big.NewInt(0), err
	}
	return big.Int(*result), nil
}

// GetTokenDecimals get token decimals
func (eth *Eth) GetTokenDecimals(contract string) (decimals big.Int, err error) {
	paramsMap := make(map[string]interface{})
	paramsMap["from"] = contract
	paramsMap["to"] = contract
	input := "0x313ce567"
	paramsMap["data"] = input

	extraParam := "latest"
	var result *math.HexOrDecimal256
	err = eth.rpc.Call(&result, eth.clientIdentifier+"_call", paramsMap, extraParam)
	if err != nil {
		return *big.NewInt(0), err
	}
	return big.Int(*result), nil
}

// GetTokenName get token name
func (eth *Eth) GetTokenName(contract string) (name string, err error) {
	paramsMap := make(map[string]interface{})
	paramsMap["from"] = contract
	paramsMap["to"] = contract
	input := "0x06fdde03"
	paramsMap["data"] = input

	extraParam := "latest"
	var result hexutil.Bytes
	err = eth.rpc.Call(&result, eth.clientIdentifier+"_call", paramsMap, extraParam)
	if err != nil {
		return "", err
	}
	resultMap := make(map[int]string)
	resultStr := result.String()
	if strings.HasPrefix(resultStr, "0x") {
		resultStr = resultStr[2:]
	}
	times := len(resultStr) / 64
	for i := 0; i < times; i++ {
		resultMap[i] = decodeStr(resultStr[64*i : 64*(i+1)])
	}

	s := string(hexutil.Hex2Bytes(resultMap[times-1]))
	return s, nil
}

// GetTokenSymbol get token symbol
func (eth *Eth) GetTokenSymbol(contract string) (symbol string, err error) {
	paramsMap := make(map[string]interface{})
	paramsMap["from"] = contract
	paramsMap["to"] = contract
	input := "0x95d89b41"
	paramsMap["data"] = input

	extraParam := "latest"
	var result hexutil.Bytes
	err = eth.rpc.Call(&result, eth.clientIdentifier+"_call", paramsMap, extraParam)
	if err != nil {
		return "", err
	}
	resultMap := make(map[int]string)
	resultStr := result.String()

	if strings.HasPrefix(resultStr, "0x") {
		resultStr = resultStr[2:]
	}
	times := len(resultStr) / 64
	for i := 0; i < times; i++ {
		resultMap[i] = decodeStr(resultStr[64*i : 64*(i+1)])
	}

	s := string(hexutil.Hex2Bytes(resultMap[times-1]))
	return s, nil
}

func decodeStr(hex string) string {
	for i := len(hex); i > 0; i = i - 2 {
		s := hex[i-2 : i]
		if s != "00" {
			return hex[:i]
		}
	}
	return hex
}
