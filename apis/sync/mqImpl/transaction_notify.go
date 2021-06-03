// description: sync_eth 
// 
// @author: xwc1125
// @date: 2020/10/05
package mqImpl

import (
	"encoding/json"
	"github.com/chain5j/sync_eth/apis/sync/models"
	"github.com/chain5j/sync_eth/params/global"
	"github.com/chain5j/sync_eth/pkg/database/mq"
)

const (
	MqTxNotifyEx       = "tx.notify.ex"
	MqTxNotifyQueue    = "tx.notify.queue"
	MqTxNotifyRouteKey = "tx.notify.routekey"
)

type TransactionNotifyMq struct {
	w         *mq.Wrapper
	chainName string
	chainType int64
}

func NewTxNotifyMq(w *mq.Wrapper) *TransactionNotifyMq {
	chainName := global.Config.ChainConfig.ChainName
	chainType := global.Config.ChainConfig.ChainType
	return &TransactionNotifyMq{
		w:         w,
		chainType: chainType,
		chainName: chainName,
	}
}

type NotifyInfo struct {
	ChainName string
	ChainType int64
	Tx        *models.Transaction
}

func (m *TransactionNotifyMq) Produce(tx *models.Transaction) error {
	// (1:BTC 2:ETH)
	info := NotifyInfo{
		ChainName: m.chainName,
		ChainType: m.chainType,
		Tx:        tx,
	}
	bytes, _ := json.Marshal(info)
	return m.w.Produce(false, false, bytes)
}
