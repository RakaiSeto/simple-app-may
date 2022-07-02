package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var snapClient snap.Client

func main() {
	midtrans.ServerKey = "SB-Mid-server-41bvMytGnmXZ737sy5XFXmIM"
	midtrans.Environment = midtrans.Sandbox

	snapClient.New(midtrans.ServerKey, midtrans.Sandbox)

	req := & snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  "TES-4",
			GrossAmt: 90000,
		}, 
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
	}

	snapClient.Options.SetPaymentOverrideNotification("http://localhost:8080/notification/handling")
	snapResp, _ := snapClient.CreateTransaction(req)
	fmt.Println("Response :", snapResp)

	g := gin.Default()
	g.POST("/notification/handling", Tes)

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func Tes(ctx *gin.Context) {
	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	fmt.Println(string(jsonData))
}

