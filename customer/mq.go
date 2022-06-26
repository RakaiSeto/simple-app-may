package main

import (
	"bytes"
	"encoding/json"
	"math/rand"

	proto "github.com/RakaiSeto/simple-app-may/service"
	"github.com/google/uuid"
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
	uuid := uuid.New().String()
	data.RequestBody.String_ = &uuid

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

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: "Request key : " + uuid}}}, nil
}