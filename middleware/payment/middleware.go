package main

import (
	// "bytes"
	// "context"
	"encoding/json"
	"strings"
	"time"

	// "errors"
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/RakaiSeto/simple-app-may/db"

	"github.com/lib/pq"
	_ "github.com/lib/pq"

	// server "github.com/RakaiSeto/simple-app-may/server"
	proto "github.com/RakaiSeto/simple-app-may/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

// var justcontext context.Context
var dbconn *sql.DB
var Rabconn *amqp.Connection
var Rabchan *amqp.Channel
var Q amqp.Queue
var Msgs <-chan amqp.Delivery

type RequestBody proto.RequestBody

type Request struct {
	Method      string       `json:"Method,omitempty"`
	Url         string       `json:"Url,omitempty"`
	RequestBody RequestBody `json:"RequestBody,omitempty"`
}

func init() {
	dbconn = db.Db
	Rabconn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// fmt.Println("dial")
	if err != nil {
		panic(err)
	}
	
	
	Rabchan, err = Rabconn.Channel()
	// fmt.Println("channel")
	if err != nil {
		panic(err)
	}
	
	
	Q, err = Rabchan.QueueDeclare(
		"order_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	// fmt.Println("queue")
	if err != nil {
		panic(err)
	}
	
	// fmt.Println("qos")
	if err != nil {
		panic(err)
	}
	

}

func main() {
	Msgs, err := Rabchan.Consume(
		Q.Name, // queue
		"",     // consumer
		true,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	// fmt.Println("consume")
	if err != nil {
		panic(err)
	}
	fmt.Println("rabbitmq init done")
	var forever chan struct{}

	go func() {
		for d := range Msgs {
			var req Request
			err = json.Unmarshal(d.Body, &req)
			fmt.Println("unmarshal")
			if err != nil {
				panic(err)
			}

			_, err = dbconn.Exec(`INSERT INTO public.order (id, user_id, product_id, quantity, payment_method, order_value, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, req.RequestBody.Order.GetId(), req.RequestBody.Order.GetUserId(), pq.Array(req.RequestBody.Order.GetProductId()), pq.Array(req.RequestBody.Order.GetQuantity()), req.RequestBody.Order.GetPaymentMethod(), req.RequestBody.Order.GetOrderValue(), time.Now().Unix(), time.Now().Unix())
			if err != nil {
				if strings.Contains(err.Error(), "duplicate key value violates") {
					continue
				}
				fmt.Println(err)
			}
			fmt.Println("input")
		}
		}()

	fmt.Printf(" [*] Awaiting RPC requests")
	<-forever

	defer Rabconn.Close()
	defer Rabchan.Close()
}

func (r *RequestBody) Value() (driver.Value, error) {
	return json.Marshal(r)
}