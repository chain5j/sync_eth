// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/10/05
package chain

import (
	"github.com/chain5j/chain5j-pkg/math"
	"github.com/chain5j/chain5j-pkg/types"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	log "github.com/chain5j/log15"
	"github.com/shopspring/decimal"
	"strings"
)

type TokenTransferInput struct {
	Method     string `json:"method"`
	MethodName string `json:"methodName"`
	From       string `json:"from,omitempty"`
	To         string `json:"to"`
	Value      string `json:"value"`
}

type TokenTransferParams struct {
	MethodName string          `json:"methodName"`
	From       string          `json:"from,omitempty"`
	To         string          `json:"to"`
	Value      decimal.Decimal `json:"value"`
}

func ParseInputBytes(inputBytes hexutil.Bytes) *TokenTransferInput {
	return ParseInput(inputBytes.String())
}

func ParseInput(input string) *TokenTransferInput {
	if len(input) != 136 && len(input) != 138 && len(input) != 202 && len(input) != 200 {
		log.Debug("input len err", "inputLen", len(input))
		return nil
	}
	if !strings.HasPrefix(input, "0x") {
		input = "0x" + input
	}
	method := input[:10]
	if method != "0xa9059cbb" && method != "0x23b872dd" {
		log.Debug("the input is not erc20 transfer(address _to, uint256 _value) or transferFrom(address _from, address _to, uint256 _value)", "method", method)
		return nil
	}
	var (
		fromStr  string
		toStr    string
		valueStr string
		mathName string
	)
	if method == "0xa9059cbb" {
		toStr = input[10:74]
		valueStr = input[74:]
		mathName = "transfer(address _to, uint256 _value)"
	} else if method == "0x23b872dd" {
		fromStr = input[10:74]
		toStr = input[74:138]
		valueStr = input[138:]
		mathName = "transferFrom(address _from, address _to, uint256 _value)"
	}
	return &TokenTransferInput{
		Method:     method,
		MethodName: mathName,
		From:       fromStr,
		To:         toStr,
		Value:      valueStr,
	}
}

func ParseTransferInput(input *TokenTransferInput) *TokenTransferParams {
	if input == nil {
		return nil
	}
	if input.Method != "0xa9059cbb" && input.Method != "0x23b872dd" {
		return nil
	}
	tokenParams := &TokenTransferParams{
		MethodName: input.MethodName,
	}
	toAddress := types.HexToAddress(input.To)
	tokenParams.To = toAddress.Hex()
	if input.From != "" {
		tokenParams.From = types.HexToAddress(input.From).Hex()
	}

	bigInt, _ := math.ParseBig256(input.Value)
	tokenParams.Value = decimal.NewFromBigInt(bigInt, 0)
	return tokenParams
}
