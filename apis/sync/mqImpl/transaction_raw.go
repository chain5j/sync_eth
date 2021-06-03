// description: sync_eth 
// 
// @author: xwc1125
// @date: 2020/10/05
package mqImpl

import (
	"encoding/json"
	"fmt"
	log "github.com/chain5j/log15"
	"github.com/chain5j/sync_eth/apis/sync/dao"
	"github.com/chain5j/sync_eth/apis/sync/models"
	"github.com/chain5j/sync_eth/params/global"
	"github.com/chain5j/sync_eth/pkg/database/es"
	"github.com/chain5j/sync_eth/pkg/database/mq"
	"github.com/streadway/amqp"
)

const (
	MqTxRawEx       = "tx.raw.ex"
	MqTxRawQueue    = "tx.raw.queue"
	MqTxRawRouteKey = "tx.raw.routekey"
)

type TransactionRawMq struct {
	w           *mq.Wrapper
	mqNotifyRaw *TransactionNotifyMq
	es          *es.ES
}

func NewTxRawMq(w *mq.Wrapper, es *es.ES) *TransactionRawMq {
	chainName := global.Config.ChainConfig.ChainName
	mqNotifyRaw1 := mq.NewDefaultMq(
		global.Config.Database.Mq,
		chainName+"."+MqTxNotifyEx,
		chainName+"."+MqTxNotifyQueue,
		chainName+"."+MqTxNotifyRouteKey)
	mqNotifyRaw := NewTxNotifyMq(mqNotifyRaw1)

	return &TransactionRawMq{
		w:           w,
		mqNotifyRaw: mqNotifyRaw,
		es:          es,
	}
}

func (m *TransactionRawMq) Produce(tx *models.Transaction) error {
	bytes, _ := json.Marshal(tx)
	fmt.Println(string(bytes))
	return m.w.Produce(false, false, bytes)
}

func (m *TransactionRawMq) Consumer() {
	m.w.Consume("", false, false, false, false, nil, m.onReceive)
}

func (m *TransactionRawMq) onReceive(d amqp.Delivery) error {
	tx := new(models.Transaction)
	err := json.Unmarshal(d.Body, &tx)
	if err != nil {
		return err
	}
	// Whether the address or hash needs to be monitored.
	// If it is not returned directly, otherwise the notify queue needs to be added
	listenAddress, err := dao.ListenAddressDao.GetListenAddress(m.es, tx.From)
	if err != nil {
		log.Error("GetListenAddress err", "from", tx.From, "err", err)
		return err
	}
	if listenAddress == nil {
		listenAddress, err = dao.ListenAddressDao.GetListenAddress(m.es, tx.To)
		if err != nil {
			log.Error("GetListenAddress err", "to", tx.To, "err", err)
			return err
		}
	}
	if listenAddress != nil {
		err = m.mqNotifyRaw.Produce(tx)
		if err != nil {
			return err
		}
	}
	d.Ack(true)
	return nil
}
