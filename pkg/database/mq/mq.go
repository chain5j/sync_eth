// description: sync_eth 
// 
// @author: xwc1125
// @date: 2020/10/05
package mq

import (
	"github.com/streadway/amqp"
	"time"
)

func NewDefaultMq(config *Config, ex, queue, routingKey string) *Wrapper {
	wrapper := New(config.GetConnUtl(), amqp.Config{
		//Vhost:     "vhost",
		Heartbeat: 2 * time.Second,
	}, ex, queue, routingKey, DefaultApply)
	return wrapper
}

func DefaultApply(exchange, queue, routingKey string, ch *amqp.Channel) (err error) {
	if err = ch.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if _, err = ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if err = ch.QueueBind(queue, routingKey, exchange, false, nil); err != nil {
		return err
	}
	return nil
}
