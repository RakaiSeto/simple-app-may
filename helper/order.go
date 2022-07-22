package helper

import (
	"github.com/RakaiSeto/simple-app-may/service"
	"google.golang.org/protobuf/encoding/protojson"
)

func ParseOrderProducts(order *service.Order, byteSlice []byte) (*service.Order, error) {
	var	reqBody *service.OrderProduct
	err = protojson.Unmarshal(byteSlice, reqBody)
	if err != nil {
		panic(err)
	}

	orderProductsDB := make([]*service.OrderProduct, 0)
	orderProductsDB = append(orderProductsDB, reqBody)
	order.OrderProducts.OrderProduct = orderProductsDB

	return order, nil
}