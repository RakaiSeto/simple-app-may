package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"

	// "errors"
	"context"
	"fmt"
	"strings"

	"github.com/RakaiSeto/simple-app-may/db"
	"google.golang.org/grpc"

	proto "github.com/RakaiSeto/simple-app-may/service"
	"github.com/jasonlvhit/gocron"
	_ "github.com/lib/pq"
	"google.golang.org/protobuf/encoding/protojson"
)

var Client proto.ServiceClient
var dbconn *sql.DB
var justcontext context.Context

type ResponseBody proto.ResponseBody

type Request struct {
	Id 		 	string		`json:"Id,omitempty"`
	Method      string      `json:"Method,omitempty"`
	Url         string      `json:"Url,omitempty"`
	RequestBody *proto.RequestBody `json:"RequestBody,omitempty"`
}

func init() {
	justcontext = context.TODO()
	dbconn = db.Db

	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	
	Client = proto.NewServiceClient(conn)
}

func main() {
	s := gocron.NewScheduler()
    s.Every(5).Second().Do(checkDb)
    <- s.Start()
}

func checkDb() {
	rows, err := dbconn.Query("SELECT id, method, url, reqbody FROM public.queue WHERE status = 'pending'")
	if err != nil {
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set"){
				return
			}
		}
	}

	var count int
	var failed int
	
	for rows.Next() {
		count++
		req := new(Request)
		var byteSlice []byte
		err := rows.Scan(&req.Id, &req.Method, &req.Url, &byteSlice)
		if err != nil{
			fmt.Println(err)
			failed++
			continue
		}

		reqBody := new(proto.RequestBody)
		err = protojson.Unmarshal(byteSlice, reqBody)
		if err != nil {
			fmt.Println(err)
			failed++
			continue
		}
		req.RequestBody = reqBody

		fmt.Printf("%v\n", req.Id)

		fmt.Println(reqBody.Order.GetProductid())

		response := switchFunc(req.Method, req.Url, req.RequestBody) 
		input := new(ResponseBody)
		input = (*ResponseBody)(response.ResponseBody)
		if response.ResponseBody.GetError() != "" {
			_, err = dbconn.Exec("UPDATE public.queue SET status='failed', reqbody=$1, updated=$2 WHERE id=$3", input, time.Now().Unix(), req.Id)
			if err != nil{
				panic(err)
			}
			failed++
			return
		}

		_, err = dbconn.Exec("UPDATE public.queue SET status='success', reqbody=$1, updated=$2 WHERE id=$3", input, time.Now().Unix(), req.Id)
		if err != nil{
			panic(err)
		}
	}
	fmt.Printf("Req: %d, failed:%d", count, failed)
}

func (r *Request) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok{
		fmt.Println("failed assertion")
	}

	return json.Unmarshal(b, &r)
}

func switchFunc(method string, url string, body *proto.RequestBody) (*proto.ResponseWrapper) {
	fmt.Printf("%v", body)
	switch method{
		case "POST":
			switch url{
			case "/user":
				resp, _ := Client.AddUser(justcontext, body)
				fmt.Println(resp)
				return resp
			case "/product":
				resp, _ := Client.AddProduct(justcontext, body)
				return resp
			case "/order":
				resp, _ := Client.AddOrder(justcontext, body)
				return resp
			}
		case "PATCH":
			switch url{

			}
		case "DELETE":
			switch url{

			}
	}
	return nil
}

func (r ResponseBody) Value() (driver.Value, error) {
	return json.Marshal(r)
}