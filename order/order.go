package main

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/RakaiSeto/simple-app-may/db"
	proto "github.com/RakaiSeto/simple-app-may/service"
	"github.com/google/uuid"
)

var dbconn *sql.DB

func init() {
	dbconn = db.Db
}

var funcCtx = context.TODO()

func AllOrder(input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	creden, err := proto.ParseJWT(funcCtx, input.GetString_())
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	// get the order
	ordersRows, err := dbconn.Query("SELECT * FROM public.order where userid=$1", creden["userid"].(float64))
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown error",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, err 
	}
	defer ordersRows.Close()

	orders := make([]*proto.Order, 0)
	for ordersRows.Next() {
		var order proto.Order
		var created int64
		var updated int64
		err := ordersRows.Scan(&order.Id, &order.UserId, &order.ProductId, &order.Quantity, &order.PaymentMethod, &order.PaymentUrl, &order.OrderValue, &order.PaymentStatus, &order.MidtransStatus, &created, &updated)
		if err != nil {
			return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: "*cricket..."}}}, nil
		}

		createdString := time.Unix(created, 0).In(proto.WIB_TIME).Format(proto.TIME_LAYOUT_ALL)
		updatedString := time.Unix(updated, 0).In(proto.WIB_TIME).Format(proto.TIME_LAYOUT_ALL)

		order.Created = &createdString
		order.Updated = &updatedString

		order.OrderValue /= 100

		orders = append(orders, &order)
	}

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{Orders: &proto.Orders{Order: orders}}}, nil
}

func OneOrder(input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	creden, err := proto.ParseJWT(funcCtx, input.GetString_())
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	// get the order
	orderRow := dbconn.QueryRow("SELECT * FROM public.order WHERE userid = $1 AND id = $2", creden["userid"].(float64), input.GetId())

	var order proto.Order
	var created int64
	var updated int64
	err = orderRow.Scan(&order.Id, &order.UserId, &order.ProductId, &order.Quantity, &order.PaymentMethod, &order.PaymentUrl, &order.OrderValue, &order.PaymentStatus, &order.MidtransStatus, &created, &updated)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			var errString string = "order not found"
			return &proto.ResponseWrapper{Code: 404, Message: "not found", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
		}
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}	

	createdString := time.Unix(created, 0).In(proto.WIB_TIME).Format(proto.TIME_LAYOUT_ALL)
	updatedString := time.Unix(updated, 0).In(proto.WIB_TIME).Format(proto.TIME_LAYOUT_ALL)

	order.Created = &createdString
	order.Updated = &updatedString

	order.OrderValue /= 100	

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{Order: &order}}, nil
}

func AddOrder(input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	uuid := uuid.New().String()

	// update user
	_, err := dbconn.Exec("UPDATE public.user SET updated_at = $1 WHERE id = $2",time.Now().Unix(), input.User.GetId()) 
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	var totalprice float64
	for i := 0; i < len(input.Order.Quantity); i++ {
		totalprice += (float64(input.Products.Product[i].Price) * float64(input.Order.Quantity[i]))
	}

	totalprice *= 100

	req := &proto.RequestWrapper{Method: "", Url: "", RequestBody: &proto.RequestBody{Order: &proto.Order{Id: uuid, UserId:input.User.Id,  PaymentMethod: input.Order.GetPaymentMethod(), OrderValue:int64(totalprice), ProductId: input.Order.ProductId, Quantity: input.Order.Quantity}}}
	_, err = Produce(req)
	if err != nil {
		errString := err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown",
			ResponseBody: &proto.ResponseBody{Error: &errString},
		}, nil
	}

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{Order: input.Order}}, nil
}

func DeleteOrder(input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	creden, err := proto.ParseJWT(funcCtx, input.GetString_())
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	// delete order
	_, err = dbconn.Exec("DELETE FROM public.order WHERE id=$1 AND userid=$2", input.Order.Id, creden["userid"].(float64))
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	success := "success deleting order " + input.Order.GetId()
	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: success}}}, nil
}

func GetOrder(order *proto.Order, id string) (*proto.Order, error) {
	row := dbconn.QueryRow("SELECT * FROM public.order WHERE id = $1", id)
	var created int
	var updated int
	err := row.Scan(&order.Id, &order.UserId, &order.ProductId, &order.Quantity, &order.PaymentMethod, &order.PaymentUrl, &order.OrderValue, &order.PaymentStatus, &order.MidtransStatus, &created, &updated)
	if err != nil {
		return nil, err
	}

	return order, nil
}