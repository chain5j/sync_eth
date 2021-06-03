// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/10/05
package models

import (
	"encoding/json"
	"github.com/chain5j/chain5j-pkg/util/dateutil"
	"github.com/chain5j/sync_eth/params/global"
	"github.com/shopspring/decimal"
	"sort"
	"sync"
	"time"
)

var (
	ContractTypeUnknown = 0
	ContractTypeERC20   = 1
	ContractTypeERC721  = 2
)

type Contract struct {
	Address          string          `json:"address" es:"keyword"`          // contract
	Name             string          `json:"name" es:"text"`                // contract name
	Symbol           string          `json:"symbol" es:"text"`              // contract symbol
	Creator          string          `json:"creator" es:"text"`             // creator
	Icon             string          `json:"icon" es:"text"`                // icon
	UnitLength       int64           `json:"unit_length" es:"integer"`      // unit_length
	Introduction     string          `json:"introduction" es:"text"`        // introduction
	WebSite          string          `json:"web_site" es:"text"`            // web_site
	BlockTime        uint64          `json:"block_time" es:"long"`          // block_time
	BlockNumber      uint64          `json:"block_number" es:"long"`        // block_number
	Total            decimal.Decimal `json:"total" es:"double"`             // total
	TxNum            int64           `json:"tx_num" es:"integer"`           // tx_num
	TxHash           string          `json:"tx_hash" es:"text"`             // tx_hash
	Type             int             `json:"type" es:"integer"`             // type，0：unknown，1：ERC20，2：ERC721
	AccountNum       int             `json:"account_num" es:"integer"`      // account num
	TokenCirculation decimal.Decimal `json:"token_circulation" es:"double"` // token_circulation
	TokenTxNum       int64           `json:"token_tx_num" es:"integer"`     // token_tx_num
	IsShow           bool            `json:"is_show" es:"boolean"`          // is_show
	Timestamp        time.Time       `json:"timestamp" es:"date"`           // timestamp
}

func (a Contract) TableName() string {
	return global.Config.ChainConfig.ChainName + "_contract"
}

func (c *Contract) Bytes() ([]byte, error) {
	return json.Marshal(c)
}

func (a *Contract) MarshalJSON() ([]byte, error) {
	if a.BlockTime > 0 {
		a.Timestamp = dateutil.SecondToTime(int64(a.BlockTime))
	} else {
		a.Timestamp = time.Now()
	}
	return json.Marshal(*a)
}

type ContractMap map[string]*Contract

type ContractMapSet struct {
	m ContractMap
	sync.RWMutex
}

func NewContractMapSet() *ContractMapSet {
	return &ContractMapSet{
		m: ContractMap{},
	}
}

func (s *ContractMapSet) Update(items ...*Contract) {
	s.Lock()
	defer s.Unlock()
	if len(items) == 0 {
		return
	}

	for _, item := range items {
		s.m[item.Address] = item
	}
}

func (s *ContractMapSet) Remove(items ...string) {
	s.Lock()
	defer s.Unlock()
	if len(items) == 0 {
		return
	}
	for _, item := range items {
		delete(s.m, item)
	}
}

func (s *ContractMapSet) Get(item string) (*Contract, bool) {
	s.RLock()
	defer s.RUnlock()
	revenue, ok := s.m[item]
	return revenue, ok
}

func (s *ContractMapSet) Len() int {
	return len(s.List())
}

func (s *ContractMapSet) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = ContractMap{}
}

func (s *ContractMapSet) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func (s *ContractMapSet) List() ContractMap {
	s.RLock()
	defer s.RUnlock()
	return s.m
}

func (s *ContractMapSet) SortList() []string {
	s.RLock()
	defer s.RUnlock()
	list := []string{}
	for item := range s.m {
		list = append(list, item)
	}
	sort.Strings(list)
	return list
}
