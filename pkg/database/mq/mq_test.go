// description: sync_eth 
// 
// @author: xwc1125
// @date: 2020/10/05
package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"testing"
	"time"
)

func TestNewMq(t *testing.T) {
	var (
		ex         = "t-ex"
		queue      = "t-queue"
		routingKey = "t-routing"
	)
	config := &Config{
		Host:         "127.0.0.1",
		Port:         5672,
		Username:     "guest",
		Password:     "guest",
		//ExchangeName: "test.rabbit.mq",
		//ExchangeType: "direct",
	}
	wrapper := NewDefaultMq(config, ex, queue, routingKey)
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for {
			select {
			case <-ticker.C:
				json:=`hello`
				if err := wrapper.Produce(false, false, []byte(json)); err != nil {
					log.Println("could not produce: ", err)
				}
				log.Println("send success")
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()

	go wrapper.Consume("", false, false, false, false, nil, consume)

	quit := make(chan bool)
	<-quit
}

func consume(d amqp.Delivery) error {
	fmt.Printf("data: %s\n", d.Body)
	d.Ack(true)
	return nil
}
