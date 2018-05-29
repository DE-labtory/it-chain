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

// 새로운 메세지 큐 만들어서 반환
func NewMessageQueue(url string, conn *amqp.Connection, ch *amqp.Channel, q amqp.Queue) *MessageQueue {
	return &MessageQueue{
		url: url,
		conn: conn,
		ch: ch,
		q: q,
	}
}

// Connect 만 호출하면 아래의 start 등의 귀찮은 내용을 다 처리해줌
func Connect(rabbitmqUrl string) *MessageQueue {

	mq := &MessageQueue{url: rabbitmqUrl}
	mq.Start()

	return mq
}


// amqp 서버에 dial 하고 연결하여 channel을 형성한다.
func (m *MessageQueue) Start() {

	// rabbitMQ 서버에 dial 하여 connection 객체 혹은 err를 받아온다.
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	// 에러 발생 시 rabbitMQ 접속 실패 에러를 띄운다.
	if err != nil {
		panic("Failed to connect to RabbitMQ" + err.Error())
	}

	// connection을 통해 channel을 형성한다.
	ch, err := conn.Channel()

	// 채널 형성 시패 시 오류를 보낸다.
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

// 연결을 끊어주는 메소이드이다.
func (m *MessageQueue) Close() {
	m.conn.Close()
	m.ch.Close()
}


// message que 에 publish 한다.
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

// message que 에서 pop 하여 message 를 수신한다.
func (m *MessageQueue) consume(topic string) (<-chan amqp.Delivery, error) {

	// 사용할 que를 선언한다.
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

	// message를 receive 하여 msgs에 담음
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

	// 읽은 메세지 반환
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
