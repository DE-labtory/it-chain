package messaging

import (
	"github.com/streadway/amqp"
)

var (
	EX_CHANGE_NAME = "it-chain"
	QUEUE_NAME = "it-chain-queue"
)

type Messaging struct {
	url  string
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func NewMessaging(rabbitmqUrl string) *Messaging{

	return &Messaging{url: rabbitmqUrl}
}

func (m *Messaging) Start(){

	//connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil{
		panic("Failed to connect to RabbitMQ" + err.Error())
	}

	//channel
	ch, err := conn.Channel()

	if err != nil{
		panic("Failed to open a channel" + err.Error())
	}

	//exchange
	err = ch.ExchangeDeclare(
		EX_CHANGE_NAME, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil{
		panic("Failed to open exchange" + err.Error())
	}

	//queue
	q, err := ch.QueueDeclare(
		QUEUE_NAME,    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	m.conn = conn
	m.ch = ch
	m.q = q
}

func (m *Messaging) Close(){

}

func (m *Messaging) Publish(topic string, data []byte) error{

	err := m.ch.Publish(
		EX_CHANGE_NAME,          // exchange
		topic, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})

	if err != nil{
		return err
	}

	return nil
}

func (m *Messaging) Consume(topic string) (<-chan amqp.Delivery, error){

	err := m.ch.QueueBind(
		m.q.Name,       // queue name
		topic,            // routing key
		EX_CHANGE_NAME, // exchange
		false,
		nil)

	if err != nil{
		return nil, err
	}

	msgs, err := m.ch.Consume(
		m.q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	if err != nil{
		return nil, err
	}

	return msgs,nil
}


