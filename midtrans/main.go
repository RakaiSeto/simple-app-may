package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/RakaiSeto/simple-app-may/db"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var dbconn *sql.DB

func main() {
	dbconn = db.Db

	g := gin.Default()
	g.GET("/hello", Hello)
	g.POST("/notification/handling", Tes)

	if err := g.Run(":5000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func Tes(ctx *gin.Context) {
	rawData, _ := ioutil.ReadAll(ctx.Request.Body)

	var jsonData map[string]interface{}
	json.Unmarshal([]byte(rawData), &jsonData)

	if jsonData["transaction_status"].(string) == "settlement" || jsonData["transaction_status"].(string) == "capture" || jsonData["transaction_status"].(string) == "refund" || jsonData["transaction_status"].(string) == "partial_refund" {
		_, err := dbconn.Exec("UPDATE public.order SET payment_status = $1, midtrans_status = $2, updated_at = $3 WHERE id = $4", "completed", "settlement", time.Now().Unix(), jsonData["transaction_id"].(string))
		if err != nil {
			fmt.Println(err)
		}

		row := dbconn.QueryRow("SELECT user_id, total_value FROM public.order WHERE id = $1", jsonData["transaction_id"].(string))
		var user_id int64
		var total_value float64

		err = row.Scan(&user_id, &total_value)

		_, err = dbconn.Exec("INSERT INTO public.transaction_history (user_id, order_id, total_value, midtrans_status, fraud_status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", user_id, jsonData["transaction_id"].(string), total_value, "completed", "settlement", time.Now().Unix(), time.Now().Unix())
		if err != nil {
			panic(err)
		}

	} else if jsonData["transaction_status"].(string) == "deny" || jsonData["transaction_status"].(string) == "cancel" || jsonData["transaction_status"].(string) == "expire" {
		_, err := dbconn.Exec("UPDATE public.order SET payment_status = $1, midtrans_status = $2, updated_at = $3 WHERE id = $4", "failed_payment", jsonData["transaction_status"].(string), time.Now().Unix(), jsonData["transaction_id"].(string))
		if err != nil {
			fmt.Println(err)
		}

		row := dbconn.QueryRow("SELECT user_id, total_value FROM public.order WHERE id = $1", jsonData["transaction_id"].(string))
		var user_id int64
		var total_value float64

		err = row.Scan(&user_id, &total_value)

		_, err = dbconn.Exec("INSERT INTO public.transaction_history (user_id, order_id, total_value, midtrans_status, fraud_status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", user_id, jsonData["transaction_id"].(string), total_value, "failed_payment", jsonData["transaction_status"].(string), time.Now().Unix(), time.Now().Unix())
		if err != nil {
			panic(err)
		}
	} else if jsonData["transaction_status"].(string) == "pending"{
		_, err := dbconn.Exec("UPDATE public.order SET payment_status = $1, midtrans_status = $2, updated_at = $3 WHERE id = $4", "waiting_payment", jsonData["transaction_status"].(string), time.Now().Unix(), jsonData["transaction_id"].(string))
		if err != nil {
			fmt.Println(err)
		}
	} else if jsonData["transaction_status"].(string) == "authorize"{
		_, err := dbconn.Exec("UPDATE public.order SET payment_status = $1, midtrans_status = $2, updated_at = $3 WHERE id = $4", "completed_payment", jsonData["transaction_status"].(string), time.Now().Unix(), jsonData["transaction_id"].(string))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func Hello(ctx *gin.Context) {
	ctx.IndentedJSON(200, gin.H{"Hello": "hello, Rakai here",})
}

