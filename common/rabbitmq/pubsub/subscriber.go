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

	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/common/rabbitmq"
	"github.com/rs/xid"
)

type Message struct {
	MatchingValue string
	Data          []byte
}

type TopicSubscriber struct {
	rabbitmq.Session
	exchange string
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

	return &TopicSubscriber{
		Session:  session,
		exchange: exchange,
	}
}

func (t *TopicSubscriber) SubscribeTopic(topic string, source interface{}) error {
	q, err := t.Ch.QueueDeclare(
		xid.New().String(), // name
		false,              // durable
		true,               // delete when usused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
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

	p, _ := NewParamBasedRouter()
	p.SetHandler(q.Name, source)

	if err != nil {
		return err
	}

	go func(queueName string, router *ParamBasedRouter) {
		for delivery := range msgs {

			message := &Message{}
			data := delivery.Body

			if err := json.Unmarshal(data, message); err != nil {
				logger.Errorf(nil, "[Common] Fail to unmarshal rabbitmq message - Err: [%s]", err.Error())
			}

			p.Route(queueName, message.Data, message.MatchingValue) //해당 event를 처리하기 위한 matching value 에는 structName이 들어간다.
		}
	}(q.Name, p)

	return nil
}

func (t *TopicSubscriber) Close() {
	t.Session.Close()
}
