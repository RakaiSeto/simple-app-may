package main

import (
	"database/sql"
	"time"

	// "errors"
	"context"
	"fmt"
	"strings"

	"github.com/RakaiSeto/simple-app-may/db"

	proto "github.com/RakaiSeto/simple-app-may/service"
	"github.com/jasonlvhit/gocron"
	_ "github.com/lib/pq"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var Client proto.ServiceClient
var dbconn *sql.DB
var justcontext context.Context
var snapClient snap.Client


func init() {
	justcontext = context.TODO()
	dbconn = db.Db
}

func main() {
	midtrans.ServerKey = "SB-Mid-server-41bvMytGnmXZ737sy5XFXmIM"
	midtrans.Environment = midtrans.Sandbox

	snapClient.New(midtrans.ServerKey, midtrans.Sandbox)

	s := gocron.NewScheduler()
    s.Every(5).Second().Do(checkDb)
    <- s.Start()
}

func checkDb() {
	rows, err := dbconn.Query("SELECT id, payment_method, order_value FROM public.order WHERE payment_status = 'unprocessed'")
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
		var req *proto.Order
		err := rows.Scan(&req.Id, &req.PaymentMethod, &req.OrderValue)
		if err != nil{
			fmt.Println(err)
			failed++
			continue
		}

		if req.PaymentMethod == "admin_confirmation" {
			_, err = dbconn.Exec("UPDATE public.order SET payment_status='waiting_payment', updated_at=$2 WHERE id=$1", req.Id, time.Now().Unix())
			if err != nil{
				panic(err)
			}
			return
		}

		midtransReq := & snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  req.GetId(),
				GrossAmt: (req.GetOrderValue()*100),
			},
		}
	
		snapResp, _ := snapClient.CreateTransaction(midtransReq)

		if req.PaymentMethod == "bca"{
			snapResp.RedirectURL += "/#/bank-transfer/bca-va"
		} else {
			snapResp.RedirectURL += "/#/gopay-qris"
		}
		
		_, err = dbconn.Exec("UPDATE public.order SET payment_url = $1, payment_status = $2, midtrans_status = $3, updated_at=$4 WHERE id = $5", snapResp.RedirectURL, "waiting_payment", "pending", time.Now().Unix(), req.Id)
		if err != nil{
			panic(err)
		}
	}
	fmt.Printf("Req: %d, failed:%d", count, failed)
}