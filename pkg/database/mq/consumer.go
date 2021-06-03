package mq

import (
	"fmt"
	log "github.com/chain5j/log15"
	"time"

	"github.com/streadway/amqp"
)

// ConsumeHandler .
type ConsumeHandler func(d amqp.Delivery) error

func (w *Wrapper) Consume(consumer string, autoAck, exclusive, noLocal, noWait bool,
	args amqp.Table, handler ConsumeHandler) {
	var (
		ch       *amqp.Channel
		delivery <-chan amqp.Delivery
		err      error
	)

	w.hasConsumer = true

	for {
		select {
		case <-w.changeConn:
			log.Debug("evt 'changeConn' triggered.")
			if ch, err = w.Channel(5 * time.Second); err != nil {
				log.Error("could not get channel for now with error: ", "err", err)
				break
			}
			if delivery, err = ch.Consume(
				w.queue,
				consumer,
				autoAck,
				exclusive,
				noLocal,
				noWait,
				args,
			); err != nil {
				log.Error("could not start consuming with error: ", "err", err)
				break
			}
			log.Debug("initial consumer finished")
		default:
			if !w.isConnected || delivery == nil {
				// true: wrapper has not connected or consumer has not initialized
				// must to wait `changeConn` evt
				time.Sleep(1 * time.Second)
				break
			}
			// delivery will be closed, then this `range` will be finished
			for d := range delivery {
				if err := handler(d); err != nil {
					log.Error(fmt.Sprintf("could not consume message: %v with error: %v", d, err))
				}
			}
		}
	}
}
