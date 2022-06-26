package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"

	proto "github.com/RakaiSeto/simple-app-may/service"
	"github.com/streadway/amqp"
)

var Rabconn *amqp.Connection
var Rabchan *amqp.Channel
var Q amqp.Queue
var Msgs <-chan amqp.Delivery

func init() {
	Rabconn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	
	Rabchan, err = Rabconn.Channel()
	if err != nil {
		panic(err)
	}
	
	Q, err = Rabchan.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		panic(err)
	}

	Msgs, err = Rabchan.Consume(
		Q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("rabbitmq init done")
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

func Produce(data *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	reqBytes := new(bytes.Buffer)
	json.NewEncoder(reqBytes).Encode(data)
	corrId := randomString(32)
	err := Rabchan.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			ReplyTo:       Q.Name,
			Body:          reqBytes.Bytes(),
		})
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown",
			ResponseBody: &proto.ResponseBody{Error: &errString},
		}, err
	}

	var res proto.ResponseWrapper
	for d := range Msgs {
		// fmt.Println("get message")
		if corrId == d.CorrelationId {
			err = json.Unmarshal(d.Body, &res)
			// fmt.Println("unmarshal message")
			if err != nil {
				var errString string = err.Error()
				return &proto.ResponseWrapper{Code: 500, Message: "unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, err
			}
			return &res, nil
			// break
		}
	}

	return &res, nil
}