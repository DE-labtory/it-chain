package rabbitmq

import (
	"github.com/streadway/amqp"
)

var (
	EX_CHANGE_NAME = "it-chain"
	QUEUE_NAME     = "it-chain-queue"
)

type MessageQueue struct {
	url  string
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func Connect(rabbitmqUrl string) *MessageQueue {

	mq := &MessageQueue{url: rabbitmqUrl}
	mq.Start()

	return mq
}

func (m *MessageQueue) Start() {

	//connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		panic("Failed to connect to RabbitMQ" + err.Error())
	}

	//channel
	ch, err := conn.Channel()

	if err != nil {
		panic("Failed to open a channel" + err.Error())
	}

	err = ch.ExchangeDeclare(
		EX_CHANGE_NAME, // name
		"topic",        // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)

	if err != nil {
		panic("Failed to open a channel" + err.Error())
	}

	m.conn = conn
	m.ch = ch
}

func (m *MessageQueue) Close() {
	m.conn.Close()
	m.ch.Close()
}

func (m *MessageQueue) Publish(topic string, data []byte) error {

	//exchange
	err := m.ch.ExchangeDeclare(
		EX_CHANGE_NAME, // name
		"topic",        // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)

	if err != nil {
		panic("Failed to open exchange" + err.Error())
	}

	err = m.ch.Publish(
		EX_CHANGE_NAME, // exchange
		topic,          // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})

	if err != nil {
		return err
	}

	return nil
}

func (m *MessageQueue) consume(topic string) (<-chan amqp.Delivery, error) {

	q, err := m.ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		panic("Failed to open a channel" + err.Error())
	}

	err = m.ch.QueueBind(
		q.Name,         // queue name
		topic,          // routing key
		EX_CHANGE_NAME, // exchange
		false,
		nil)

	if err != nil {
		return nil, err
	}

	msgs, err := m.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	if err != nil {
		return nil, err
	}

	return msgs, nil
}

type Handler func(delivery amqp.Delivery)

func (m *MessageQueue) Consume(topic string, handler Handler) error {

	chanDelivery, err := m.consume(topic)

	if err != nil {
		return err
	}

	go func() {
		for delivery := range chanDelivery {
			handler(delivery)
		}
	}()

	return nil
}
