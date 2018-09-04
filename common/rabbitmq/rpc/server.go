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

package rpc

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"

	"github.com/it-chain/engine/common/rabbitmq"
	"github.com/streadway/amqp"
)

type Server struct {
	rabbitmq.Session
}

func NewServer(rabbitmqUrl string) Server {

	return Server{
		Session: rabbitmq.CreateSession(rabbitmqUrl),
	}
}

//todo need handler params and return value check logic
func (s Server) Register(queue string, handler interface{}) error {

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

			sourceValue := reflect.ValueOf(handler)
			sourceType := reflect.TypeOf(handler)

			numOfParam := sourceType.NumIn()

			if numOfParam != 1 {
				log.Println(err.Error())
			}

			callbackParam := sourceType.In(0)
			v := reflect.New(callbackParam)
			initializeStruct(callbackParam, v.Elem())
			paramInterface := v.Interface()

			err := json.Unmarshal(d.Body, paramInterface)

			if err != nil {
				log.Println(err.Error())
			}

			paramValue := reflect.ValueOf(paramInterface).Elem().Interface()
			values := sourceValue.Call([]reflect.Value{reflect.ValueOf(paramValue)})

			r, err := buildResult(values)

			err = s.Ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          r,
				})

			if err != nil {
				log.Println(err.Error())
			}
		}
	}()

	return nil
}

type Result struct {
	Data []byte
	Err  Error
}

func buildResult(values []reflect.Value) ([]byte, error) {

	if len(values) != 2 {
		return []byte{}, errors.New("return should 2")
	}

	d, err := json.Marshal(values[0].Interface())

	if err != nil {
		return []byte{}, err
	}

	e, ok := values[1].Interface().(Error)

	if !ok {
		return []byte{}, err
	}

	return json.Marshal(Result{
		Data: d,
		Err:  e,
	})
}
