package mq

import (
	"errors"
	log "github.com/chain5j/log15"
	"time"

	"github.com/streadway/amqp"
)

const (
	reconnectDelay     = 5 * time.Second
	reconnectDetectDur = 5 * time.Second
)

var (
	errNotConnected  = errors.New("not connected to the AMQP server")
	errAlreadyClosed = errors.New("already closed: not connected to the AMQP server")
)

// ApplyTopology to apply resource from MQ server
// eg. QueueDeclare, ExchangeDeclare
type ApplyTopology func(exchange, queue, routingKey string, ch *amqp.Channel) error

// Wrapper .
type Wrapper struct {
	Addr   string
	Config amqp.Config

	applyTopology ApplyTopology
	connection    *amqp.Connection
	channel       *amqp.Channel
	done          chan bool
	changeConn    chan struct{}
	chNotify      chan *amqp.Error // channel notify
	connNotify    chan *amqp.Error // conn notify

	isConnected bool // mark wrapper is connected to server
	hasConsumer bool // mark wrapper is used by a consumer

	exchange, queue, routingKey string
}

// handleReconnect
func (w *Wrapper) handleReconnect() {
	for {
		if !w.isConnected {
			log.Debug("Attempting to connect")
			var (
				connected = false
				err       error
			)

			for cnt := 0; !connected; cnt++ {
				if connected, err = w.connect(); err != nil {
					log.Debug("Failed to connect: ", "err", err)
				}
				if !connected {
					log.Debug("Retrying... ", "cnt", cnt)
				}
				time.Sleep(reconnectDelay)
			}
		}

		select {
		case <-w.done:
			println("evt `w.done` triggered")
			return
		case err := <-w.chNotify:
			log.Debug("channel close notify: ", "err", err)
			w.isConnected = false
		case err := <-w.connNotify:
			log.Debug("conn close notify: ", "err", err)
			w.isConnected = false
		}
		time.Sleep(reconnectDetectDur)
	}
}

// Connect .
func (w *Wrapper) connect() (bool, error) {
	conn, err := amqp.DialConfig(w.Addr, w.Config)
	if err != nil {
		return false, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return false, err
	}

	if err := w.applyTopology(w.exchange, w.queue, w.routingKey, ch); err != nil {
		return false, err
	}
	w.isConnected = true
	w.changeConnection(conn, ch)
	return true, nil
}

func (w *Wrapper) changeConnection(connection *amqp.Connection, channel *amqp.Channel) {
	w.connection = connection
	w.connNotify = make(chan *amqp.Error, 1)
	w.connection.NotifyClose(w.connNotify)

	w.channel = channel
	w.chNotify = make(chan *amqp.Error, 1)
	w.channel.NotifyClose(w.chNotify)

	// TOFIX: only producer will be blocked here
	if w.hasConsumer {
		// true: cause only consumer will be  notify for now.
		w.changeConn <- struct{}{}
	}
}

// Channel . it will blocked
func (w *Wrapper) Channel(timeout time.Duration) (*amqp.Channel, error) {
	timer := time.NewTimer(timeout)
	for !w.isConnected {
		select {
		case <-timer.C:
			return nil, errNotConnected
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	return w.channel, nil
}

// Close .
func (w *Wrapper) Close() error {
	if !w.isConnected {
		return errAlreadyClosed
	}
	err := w.channel.Close()
	if err != nil {
		return err
	}
	err = w.connection.Close()
	if err != nil {
		return err
	}
	close(w.done)
	w.isConnected = false
	return nil
}

func New(addr string, cfg amqp.Config, exchange, queue, routingKey string, f ApplyTopology) *Wrapper {
	w := &Wrapper{
		Addr:          addr,
		applyTopology: f,
		Config:        cfg,
		changeConn:    make(chan struct{}, 1),
		exchange:      exchange,
		queue:         queue,
		routingKey:    routingKey,
	}

	go w.handleReconnect()

	return w
}
