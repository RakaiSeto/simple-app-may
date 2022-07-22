package main

import (
	"context"
	"database/sql"
	"strings"

	// "strings"
	"time"

	"github.com/RakaiSeto/simple-app-may/db"
	"github.com/RakaiSeto/simple-app-may/helper"
	proto "github.com/RakaiSeto/simple-app-may/service"
	"github.com/google/uuid"
)

var dbconn *sql.DB

func init() {
	dbconn = db.Db
}

var funcCtx = context.TODO()

func AllOrder(input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	creden, err := helper.ParseJWT(funcCtx, input.GetString_())
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	// get the order
	ordersRows, err := dbconn.Query("SELECT * FROM public.order where userid=$1", creden["userid"].(float64))
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: "*cricket..."}}}, nil
		} else {
			var errString string = err.Error()
			return &proto.ResponseWrapper{
				Code: 500,
				Message: "unknown error",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, err 
		}
	}
	defer ordersRows.Close()

	orders := make([]*proto.Order, 0)
	for ordersRows.Next() {
		var order proto.Order
		var created string
		var updated string
		var byteSlice []byte
		err := ordersRows.Scan(&order.Id, &order.UserId, &byteSlice, &order.PaymentMethod, &order.PaymentUrl, &order.OrderValue, &order.PaymentStatus, &order.MidtransStatus, &created, &updated)

		created_at, err := helper.ParseTimeToWIB(created, helper.TIME_LAYOUT_ALL)
		if err != nil {
			var errString string = err.Error()
			return &proto.ResponseWrapper{Code: 500, Message: "unknown error", ResponseBody: &proto.ResponseBody{Error: &errString,},}, err 
		}
		order.Created = &created_at
		updated_at, err := helper.ParseTimeToWIB(updated, helper.TIME_LAYOUT_ALL)
		if err != nil {
			var errString string = err.Error()
			return &proto.ResponseWrapper{Code: 500, Message: "unknown error", ResponseBody: &proto.ResponseBody{Error: &errString,},}, err 
		}
		order.Updated = &updated_at

		order.OrderValue /= 100
		newOrder, err := helper.ParseOrderProducts(&order, byteSlice)
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

		orders = append(orders, newOrder)
	}

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{Orders: &proto.Orders{Order: orders}}}, nil
}

func OneOrder(input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	creden, err := helper.ParseJWT(funcCtx, input.GetString_())
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	// get the order
	orderRow := dbconn.QueryRow("SELECT * FROM public.order WHERE userid = $1 AND id = $2", creden["userid"].(float64), input.GetId())

	var order proto.Order
	var created string
	var updated string
	var byteSlice []byte
	err = orderRow.Scan(&order.Id, &order.UserId, &byteSlice, &order.PaymentMethod, &order.PaymentUrl, &order.OrderValue, &order.PaymentStatus, &order.MidtransStatus, &created, &updated)

	created_at, err := helper.ParseTimeToWIB(created, helper.TIME_LAYOUT_ALL)
		if err != nil {
			var errString string = err.Error()
			return &proto.ResponseWrapper{Code: 500, Message: "unknown error", ResponseBody: &proto.ResponseBody{Error: &errString,},}, err 
		}
		order.Created = &created_at
		updated_at, err := helper.ParseTimeToWIB(updated, helper.TIME_LAYOUT_ALL)
		if err != nil {
			var errString string = err.Error()
			return &proto.ResponseWrapper{Code: 500, Message: "unknown error", ResponseBody: &proto.ResponseBody{Error: &errString,},}, err 
		}
		order.Updated = &updated_at
	
	order.OrderValue /= 100
	newOrder, err := helper.ParseOrderProducts(&order, byteSlice)
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

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{Order: newOrder}}, nil
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
	for i := 0; i < len(input.Order.OrderProducts.OrderProduct); i++ {
		totalprice += (float64(input.Products.Product[i].GetPrice()) * float64(input.Order.OrderProducts.OrderProduct[i].GetQuantity()))
	}

	totalprice *= 10

	req := &proto.RequestWrapper{Method: "", Url: "", RequestBody: &proto.RequestBody{Order: &proto.Order{Id: uuid, UserId:input.User.Id,  PaymentMethod: input.Order.GetPaymentMethod(), OrderValue:int64(totalprice), OrderProducts: &proto.OrderProducts{OrderProduct: input.Order.OrderProducts.OrderProduct}}}}
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
	creden, err := helper.ParseJWT(funcCtx, input.GetString_())
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

func GetOrder(id string) (*proto.Order, error) {
	row := dbconn.QueryRow("SELECT * FROM public.order WHERE id = $1", id)
	var order proto.Order
	var created string
	var updated string
	var byteSlice []byte
	err := row.Scan(&order.Id, &order.UserId, &byteSlice, &order.PaymentMethod, &order.PaymentUrl, &order.OrderValue, &order.PaymentStatus, &order.MidtransStatus, &created, &updated)
	if err != nil {
		return nil, err 
	}

	created_at, err := helper.ParseTimeToWIB(created, helper.TIME_LAYOUT_ALL)
	if err != nil {
		return nil, err 
	}
	order.Created = &created_at
	updated_at, err := helper.ParseTimeToWIB(updated, helper.TIME_LAYOUT_ALL)
	if err != nil {
		return nil, err 
	}
	order.Updated = &updated_at
	
	order.OrderValue /= 100
	newOrder, err := helper.ParseOrderProducts(&order, byteSlice)
	if err != nil {
		return nil, err 
	}

	return newOrder, nil
}