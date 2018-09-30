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
	"errors"
	"reflect"
	"sync"
)

var ErrTypeNotFound = errors.New("Type of handler not found")

//Route data depends on type
type Router interface {

	//route data depends on matching value
	Route(key string, data []byte, matchingValue string) error

	SetHandler(key string, handler interface{}) error
}

type Handler struct {
	Func func(param interface{})
	Type reflect.Type
}

//ParamBasedRouter routes data through the structure and structure name(matching value) of the parameter
type ParamBasedRouter struct {
	handlerMap map[string]Handler
	sync.Mutex
}

func NewParamBasedRouter() (*ParamBasedRouter, error) {

	return &ParamBasedRouter{
		handlerMap: make(map[string]Handler),
	}, nil
}

////handler should be a struct pointer which has handler method
func (c *ParamBasedRouter) SetHandler(key string, handler interface{}) error {

	if reflect.TypeOf(handler).Kind() != reflect.Ptr {
		return errors.New("handler should be ptr type")
	}

	sourceType := reflect.TypeOf(handler)
	methodCount := sourceType.NumMethod() // NumMethod returns the number of exported methods in the value's method set.

	// register every method that has only one input parameter to the handlerMap
	for i := 0; i < methodCount; i++ {

		// Method returns a function value corresponding to v's i'th method.
		// The arguments to a Call on the returned function should not include a receiver;
		// the returned function will always use v as the receiver.
		// Method panics if i is out of range or if v is a nil interface value.
		method := sourceType.Method(i)

		// 핸들러 메소드의 인풋 파라미터가 1개가 아니면 에러 반환
		// 객체의 함수는 기본 인풋으로 자기 자신이 들어오기 때문에 인자가 최소 하나가 있다.
		if method.Type.NumIn() != 2 {
			return errors.New("number of parameter of handler is not 2")
		}

		paramType := method.Type.In(1) //returns the type of a function type's i'th input parameter.
		handler := createEventHandler(method, handler)
		// 핸들러맵에 핸들러 매소드 등록
		c.handlerMap[key+paramType.Name()] = Handler{Type: paramType, Func: handler}

	}

	return nil
}

// 핸들러의 특정 메소드의 입력으로 핸들러 자체와 파라미터인자를 받아 처리하는 핸들러 반환
func createEventHandler(method reflect.Method, handler interface{}) func(interface{}) {
	return func(param interface{}) {
		sourceValue := reflect.ValueOf(handler)
		eventValue := reflect.ValueOf(param)

		// Call actual event handling method.
		method.Func.Call([]reflect.Value{sourceValue, eventValue})
	}
}

func (c ParamBasedRouter) Route(key string, data []byte, structName string) (err error) {
	paramType, handler, err := c.findTypeOfHandlers(key, structName)

	if paramType == nil {
		//logger.Errorf(nil, "No handler found for struct [%s]", structName)
		return nil
	}

	v := reflect.New(paramType)
	initializeStruct(paramType, v.Elem())
	paramInterface := v.Interface()

	err = json.Unmarshal(data, paramInterface)
	if err != nil {
		return err
	}

	paramValue := reflect.ValueOf(paramInterface).Elem().Interface()
	handler(paramValue)

	return nil
}

//find type of handler by struct name
func (c ParamBasedRouter) findTypeOfHandlers(key string, matchingValue string) (reflect.Type, func(param interface{}), error) {
	for k, handlers := range c.handlerMap {
		if k == key+matchingValue {
			return handlers.Type, handlers.Func, nil
		}
	}

	return nil, nil, nil
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
