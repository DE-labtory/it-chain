/*
 * Copyright 2018 DE-labtory
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

package rpc_t

import (
	"encoding/json"
	"log"

	"github.com/DE-labtory/it-chain/common/rabbitmq"
	"github.com/streadway/amqp"
)

type Handler func(data []byte) ([]byte, error)

type Server struct {
	rabbitmq.Session
}

func NewServer(rabbitmqUrl string) Server {

	return Server{
		Session: rabbitmq.CreateSession(rabbitmqUrl),
	}
}

func (s Server) Register(queue string, handler Handler) error {

	q, err := s.Ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return err
	}

	err = s.Ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		return err
	}

	msgs, err := s.Ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	go func() {
		for d := range msgs {
			s.handleCall(d, handler)
		}
	}()

	return nil
}

func (s Server) handleCall(d amqp.Delivery, handler Handler) {
	response, err := handler(d.Body)

	errStr := ""

	if err != nil {
		errStr = err.Error()
	}

	result := Result{
		Err:  errStr,
		Data: response,
	}

	data, err := json.Marshal(result)

	if err != nil {
		log.Println(err.Error())
	}

	err = s.Ch.Publish(
		"",        // exchange
		d.ReplyTo, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: d.CorrelationId,
			Body:          data,
		})

	if err != nil {
		log.Println(err.Error())
	}

	d.Ack(false)
}
