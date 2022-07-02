package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/RakaiSeto/simple-app-may/db"
	proto "github.com/RakaiSeto/simple-app-may/service"
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
		err := ordersRows.Scan(&order.Id, &order.Userid, &order.Productid, &order.Quantity, &order.Totalprice, &created, &updated)
		if err != nil {
			return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: "*cricket..."}}}, nil
		}

		createdString := time.Unix(created, 0).In(proto.WIB_TIME).Format(proto.TIME_LAYOUT_ALL)
		updatedString := time.Unix(updated, 0).In(proto.WIB_TIME).Format(proto.TIME_LAYOUT_ALL)

		order.Created = &createdString
		order.Updated = &updatedString

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
	err = orderRow.Scan(&order.Id, &order.Userid, &order.Productid, &order.Quantity, &order.Totalprice, &created, &updated)
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

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{Order: &order}}, nil
}

func AddOrder(input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	fmt.Println(*input.String_)
	creden, err := proto.ParseJWT(funcCtx, *input.String_)
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	// deduct wallet
	totalprice := input.Order.GetQuantity() * input.Product.GetPrice()
	var wallet int64
	row := dbconn.QueryRow("SELECT wallet FROM public.user where id=$1", creden["userid"].(float64))
	err = row.Scan(&wallet)
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	if (wallet - totalprice) < 0 {
		var errString string = "insufficient funds"
		return &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}
	
	// input order
	_, err = dbconn.Exec("INSERT INTO public.order (userid, productid, quantity, totalprice, created, updated) VALUES ($1, $2, $3, $4, $5, $6)", creden["userid"].(float64), input.Order.GetProductid(), input.Order.GetQuantity(), totalprice, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	wallet -= totalprice
	// update user
	_, err = dbconn.Exec("UPDATE public.user SET wallet = $1, updated = $2 WHERE id = $3", wallet, time.Now().Unix(), creden["userid"].(float64)) 
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	// get inputted order
	row = dbconn.QueryRow("SELECT id, created FROM public.order WHERE userid=$1 ORDER BY id desc limit 1", creden["userid"].(float64))

	var created int64
	err = row.Scan(&input.Order.Id, &created)
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	createdString := time.Unix(created, 0).In(proto.WIB_TIME).Format(proto.TIME_LAYOUT_ALL)

	input.Order.Created = &createdString

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{Order: input.Order}}, nil
}

func UpdateOrder(input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	creden, err := proto.ParseJWT(funcCtx, input.GetString_())
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}
	
	// get order
	var dborder proto.Order
	row := dbconn.QueryRow("SELECT * FROM public.order where id=$1 AND userid=$2", input.Order.GetId(), creden["userid"].(float64))
	err = row.Scan(&dborder.Id, &dborder.Userid, &dborder.Productid, &dborder.Quantity, &dborder.Totalprice, &dborder.Created, &dborder.Updated)

	if input.Order.Productid != 0 {
		row = dbconn.QueryRow("SELECT id, name, price FROM public.product where id=$1", input.Order.GetProductid())
	} else {row = dbconn.QueryRow("SELECT id, name, price FROM public.product where id=$1", dborder.Productid)}
	// get the product
	product := proto.Product{}
	err = row.Scan(&product.Id, &product.Name, &product.Price)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			var errString string = "no such product"
			return &proto.ResponseWrapper{Code: 404, Message: "not found", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
		}
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	if input.Order.Quantity != 0 {dborder.Quantity = input.Order.Quantity}

	// deduct wallet
	totalprice := input.Order.GetQuantity() * product.Price
	var wallet int64
	row = dbconn.QueryRow("SELECT wallet FROM public.product where id=$1", creden["userid"].(float64))
	err = row.Scan(&wallet)
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	if (wallet - totalprice) < 0 {
		var errString string = "insufficient funds"
		return &proto.ResponseWrapper{Code: 400, Message: "bad request", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	dborder.Totalprice = totalprice

	_, err = dbconn.Exec("UPDATE public.user SET wallet = $1, updated = $2 WHERE id = $3", wallet, time.Now().Unix(), creden["userid"].(float64))
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	// update order
	_, err = dbconn.Exec("UPDATE public.order SET userid=$1, productid=$2, quantity=$3, totalprice=$4 WHERE id=$5", creden["userid"].(float64), dborder.Productid, dborder.Quantity, totalprice, creden["userid"].(float64))
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	// get inputted order
	row = dbconn.QueryRow("SELECT updated FROM public.order WHERE userid=$1", creden["userid"].(float64))

	var updated int64
	err = row.Scan(&updated)
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	updatedString := time.Unix(updated, 0).In(proto.WIB_TIME).Format(proto.TIME_LAYOUT_ALL)

	dborder.Updated = &updatedString

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{Order: &dborder}}, nil
}

func DeleteOrder(input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	creden, err := proto.ParseJWT(funcCtx, input.GetString_())
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 401, Message: "unauthorized", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	// delete order
	_, err = dbconn.Exec("DELETE FROM public.order WHERE id=$1 AND userid=$2", input.Id, creden["userid"].(float64))
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{Code: 500, Message: "unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
	}

	success := "success deleting order number" + string(input.GetId())
	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{ResponseStatus: &proto.ResponseStatus{Response: success}}}, nil
}