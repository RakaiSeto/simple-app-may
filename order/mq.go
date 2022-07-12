package main

import (
	"bytes"
	"encoding/json"

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

func Produce(data *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	uuid := uuid.New().String()
	data.RequestBody.QueueUUID = &uuid

	reqBytes := new(bytes.Buffer)
	json.NewEncoder(reqBytes).Encode(data)
	err := Rabchan.Publish(
		"",          // exchange
		"order_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "application/json",
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

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: "Order processed, payment link will be shown in your order details"}}}, nil
}