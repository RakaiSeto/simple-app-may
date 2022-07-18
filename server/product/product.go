package product

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
var varError error

func init() {
	dbconn = db.Db
}

var funcCtx = context.TODO()

func AllProduct(input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	rows, err := dbconn.Query("SELECT id, name, description, price, updated FROM public.product")
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			var errString string = err.Error()
			return &proto.ResponseWrapper{
				Code: 404,
				Message: "not found",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, err
		}
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown error",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, err
	}

	products := make([]*proto.Product, 0)
	for rows.Next() {
		var product proto.Product
		var updated int64
		err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &updated)
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

		updatedString := time.Unix(updated, 0).In(proto.WIB_TIME).Format(proto.TIME_LAYOUT_ALL)

		product.Updated = &updatedString

		products = append(products, &product)
	}

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{Products: &proto.Products{Product: products}}}, nil
}

func OneProduct(input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	row := dbconn.QueryRow("SELECT id, name, description, price, updated FROM public.product where id=$1", input.Id)

	product := proto.Product{}
	var updated int64
	err := row.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &updated)
	updatedString := time.Unix(updated, 0).In(proto.WIB_TIME).Format(proto.TIME_LAYOUT_ALL)

	product.Updated = &updatedString
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			var errString string = err.Error()
			return &proto.ResponseWrapper{
				Code: 404,
				Message: "not found",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, err
		}
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown error",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, err
	}

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{Product: &product}}, nil
}

func AddProduct(input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	var j int
	row := dbconn.QueryRow("SELECT id FROM public.product WHERE name=$1", input.Product.GetName())
	_ = row.Scan(&j)
	if j != 0 {
		var errString string = "product already exists"
		return &proto.ResponseWrapper{
			Code: 409,
			Message: "conflict",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, nil
	}
	
	_, err := dbconn.Exec("INSERT INTO public.product (name, description, price, created, updated) VALUES ($1, $2, $3, $4, $5)", input.Product.GetName(), input.Product.GetDescription(), input.Product.GetPrice(), time.Now().Unix(), time.Now().Unix())
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

	row = dbconn.QueryRow("SELECT id, created, updated FROM public.product WHERE name=$1", input.Product.GetName())
	var created int64
	var updated int64
	err = row.Scan(&input.Product.Id, &created, &updated)

	createdString := time.Unix(created, 0).In(proto.WIB_TIME).Format(proto.TIME_LAYOUT_ALL)
	updatedString := time.Unix(updated, 0).In(proto.WIB_TIME).Format(proto.TIME_LAYOUT_ALL)

	input.Product.Created = &createdString
	input.Product.Updated = &updatedString

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

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{Product: input.Product}}, nil
}

func UpdateProduct(input *proto.RequestBody) (*proto.ResponseWrapper, error){	
	QueryProduct := proto.Product{}
	row := dbconn.QueryRow("SELECT * from public.product where id = $1", input.Product.Id)
	var created int64
	var updated int64
	err := row.Scan(&QueryProduct.Id, &QueryProduct.Name, &QueryProduct.Description, &QueryProduct.Price, &created, &updated)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			var errString string = err.Error()
			return &proto.ResponseWrapper{
				Code: 404,
				Message: "not found",
				ResponseBody: &proto.ResponseBody{
					Error: &errString,
				},
			}, err
		}
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown error",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, err
	}

	if input.Product.Name != "" {QueryProduct.Name = input.Product.Name}
	if input.Product.Description != "" {QueryProduct.Description = input.Product.Description}
	if input.Product.Price != 0 {QueryProduct.Price = input.Product.Price}

	_, err = dbconn.Exec("UPDATE public.product SET name=$2, description=$3, price=$4, updated=$5 WHERE id=$1", QueryProduct.Id, QueryProduct.Name, QueryProduct.Description, QueryProduct.Price, time.Now().Unix())
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

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{
		Product: &QueryProduct,
	}}, nil
}

func DeleteProduct(input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	row := dbconn.QueryRow("SELECT name FROM public.product where id=$1", input.Product.GetId())

	var name string
	
	err := row.Scan(&name)
	if err != nil {
		var errString string = err.Error()
		return &proto.ResponseWrapper{
			Code: 500,
			Message: "unknown error",
			ResponseBody: &proto.ResponseBody{
				Error: &errString,
			},
		}, err
	} else if name == "" {
		varError = fmt.Errorf("product not found")
		return nil, varError
	}

	_, err = dbconn.Exec("DELETE FROM public.product WHERE id=$1", input.Product.GetId())
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

	return &proto.ResponseWrapper{Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{
		ResponseStatus: &proto.ResponseStatus{Response: "success deleting product"},
	}}, nil
}

func GetProduct(id int64) (*proto.Product, error) {
	row := dbconn.QueryRow("SELECT * FROM public.product WHERE id = $1", id)
	var created int
	var updated int
	var product proto.Product
	err := row.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &created, &updated)
	if err != nil {
		return nil, err
	}
	fmt.Println("pe")
	return &product, nil
}