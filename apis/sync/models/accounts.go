// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/10/05
package models

import (
	"github.com/chain5j/sync_eth/params/global"
	"github.com/shopspring/decimal"
	"sort"
	"sync"
)

const (
	TokenETH      = 1 // eth
	TokenContract = 2 // Contract
	TokenERC20    = 3 // ERC20
	TokenDeploy   = 4 // deploy contract

	TypeAccountNormal   = 0
	TypeAccountContract = 1
)

type Accounts struct {
	Id               int64           `json:"id"`
	Address          string          `json:"address"`
	Balance          decimal.Decimal `json:"balance"`       // balance
	BalanceInput     decimal.Decimal `json:"balanceInput"`  // input balance
	BalanceOutput    decimal.Decimal `json:"balanceOutput"` // output balance
	AccountType      int64           `json:"type"`          // 0：common，1：contract
	Contract         string          `json:"contract"`      // contract address
	ContractExist    int             `json:"contractExist"` // 0:contract is nil，1: contract is not nil
	TransactionCount int64           `json:"transactionCount"`
}

func (a *Accounts) TableName() string {
	return global.Config.ChainConfig.ChainName + "_accounts"
}

type AccountMap map[string]*Accounts

type MapSet struct {
	m AccountMap
	sync.RWMutex
}

func NewMapSet() *MapSet {
	return &MapSet{
		m: AccountMap{},
	}
}

func (s *MapSet) Update(items ...*Accounts) {
	s.Lock()
	defer s.Unlock()
	if len(items) == 0 {
		return
	}

	for _, item := range items {
		s.m[item.Address+"_"+item.Contract] = item
	}
}

func (s *MapSet) Remove(items ...string) {
	s.Lock()
	defer s.Unlock()
	if len(items) == 0 {
		return
	}
	for _, item := range items {
		delete(s.m, item)
	}
}

func (s *MapSet) Get(item string, contract string) (*Accounts, bool) {
	s.RLock()
	defer s.RUnlock()
	revenue, ok := s.m[item+"_"+contract]
	return revenue, ok
}

func (s *MapSet) Len() int {
	return len(s.List())
}

func (s *MapSet) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = AccountMap{}
}

func (s *MapSet) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func (s *MapSet) List() AccountMap {
	s.RLock()
	defer s.RUnlock()
	return s.m
}

func (s *MapSet) SortList() []string {
	s.RLock()
	defer s.RUnlock()
	list := []string{}
	for item := range s.m {
		list = append(list, item)
	}
	sort.Strings(list)
	return list
}
