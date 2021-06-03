package mq

import (
	"time"

	"github.com/streadway/amqp"
)

func (w *Wrapper) Produce(mandatory, immediate bool, dat []byte) error {
	ch, err := w.Channel(5 * time.Second)
	if err != nil {
		return err
	}

	return ch.Publish(
		w.exchange,
		w.routingKey,
		mandatory,
		immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        dat,
		})
}
