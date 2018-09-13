/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pubsub

import (
	"encoding/json"
	"log"

	"regexp"

	"github.com/it-chain/engine/common/rabbitmq"
	"github.com/streadway/amqp"
)

type Message struct {
	MatchingValue string
	Data          []byte
}

type TopicSubscriber struct {
	rabbitmq.Session
	exchange string
	router   Router
	topics   []string
}

func NewTopicSubscriber(rabbitmqUrl string, exchange string) *TopicSubscriber {

	session := rabbitmq.CreateSession(rabbitmqUrl)

	err := session.Ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		panic(err.Error())
	}

	p, _ := NewParamBasedRouter()

	return &TopicSubscriber{
		Session:  session,
		exchange: exchange,
		router:   p,
		topics:   make([]string, 0),
	}
}

func (t *TopicSubscriber) SubscribeTopic(topic string, source interface{}) error {

	var q amqp.Queue
	var err error
	if t.checkRegex(topic) {
		q, err = t.Ch.QueueDeclare(
			"",    // name
			false, // durable
			true,  // delete when usused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)

		if err != nil {
			return err
		}

		err = t.Ch.QueueBind(
			q.Name,     // queue name
			topic,      // routing key
			t.exchange, // exchange
			false,
			nil)

		if err != nil {
			return err
		}
	}

	msgs, err := t.Ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	if err != nil {
		return err
	}

	err = t.router.SetHandler(source)

	if err != nil {
		return err
	}

	go func() {
		for delivery := range msgs {
			data := delivery.Body

			message := &Message{}
			err := json.Unmarshal(data, message)

			if err != nil {
				log.Print(err.Error())
			}

			t.router.Route(message.Data, message.MatchingValue) //해당 event를 처리하기 위한 matching value 에는 structName이 들어간다.
		}
	}()

	return nil
}

func (t *TopicSubscriber) checkRegex(topic string) bool {
	r, _ := regexp.Compile(topic)

	var index = -1
	for i, topic := range t.topics {
		if r.MatchString(topic) {
			index = i
		}
	}

	if index >= 0 {
		t.topics[index] = topic
		return false
	} else {
		t.topics = append(t.topics, topic)
		return true
	}
}

func (t *TopicSubscriber) Close() {
	t.Session.Close()
}
