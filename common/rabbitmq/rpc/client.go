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
	"math/rand"
	"reflect"

	"github.com/it-chain/engine/common/rabbitmq"
	"github.com/streadway/amqp"
)

type Client struct {
	rabbitmq.Session
}

func NewClient(rabbitmqUrl string) Client {

	return Client{
		Session: rabbitmq.CreateSession(rabbitmqUrl),
	}
}

//todo need to implement timeout
func (c Client) Call(queue string, params interface{}, callback interface{}) error {

	if !hasConsumer(c.Ch, queue) {
		return errors.New("no consumer")
	}

	data, err := json.Marshal(params)

	if err != nil {
		return err
	}

	replyQ, err := c.Ch.QueueDeclare(
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

	err = c.Ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		return err
	}

	msgs, err := c.Ch.Consume(
		replyQ.Name, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)

	if err != nil {
		return err
	}

	corrId := randomString(32)

	err = c.Ch.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       replyQ.Name,
			Body:          data,
		})

	if err != nil {
		return err
	}

	for d := range msgs {

		if corrId == d.CorrelationId {

			err := handleResponse(d.Body, callback)

			if err != nil {
				log.Fatal(err)
				return err
			}

			c.Ch.QueueDelete(replyQ.Name, false, false, true)
			break
		}
	}

	return nil
}

func handleResponse(data []byte, callback interface{}) error {

	sourceValue := reflect.ValueOf(callback)
	sourceType := reflect.TypeOf(callback)

	len := sourceType.NumIn()

	if len != 2 {
		return errors.New("callback function parameter should have only one struct")
	}

	callbackParam := sourceType.In(0)
	v, err := toValues(data, callbackParam)

	if err != nil {
		return err
	}

	sourceValue.Call(v)

	return nil
}

func toValues(data []byte, paramType reflect.Type) ([]reflect.Value, error) {

	v := reflect.New(paramType)
	initializeStruct(paramType, v.Elem())
	paramInterface := v.Interface()

	r := Result{}
	err := json.Unmarshal(data, &r)

	if err != nil {
		return []reflect.Value{}, err
	}

	err = json.Unmarshal(r.Data, paramInterface)

	if err != nil {
		return []reflect.Value{}, err
	}

	paramValue := reflect.ValueOf(paramInterface).Elem().Interface()

	return []reflect.Value{reflect.ValueOf(paramValue), reflect.ValueOf(r.Err)}, nil
}

func hasConsumer(channel *amqp.Channel, queueName string) bool {

	q, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if q.Consumers == 0 || err != nil {
		return false
	}

	return true
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

//build empty struct from struct type
func initializeStruct(t reflect.Type, v reflect.Value) {

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)

		if !f.CanSet() {
			continue
		}

		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			initializeStruct(ft.Type, f)
		case reflect.Ptr:
			fv := reflect.New(ft.Type.Elem())
			initializeStruct(ft.Type.Elem(), fv.Elem())
			f.Set(fv)
		default:
		}
	}
}
