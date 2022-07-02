package main

import (
	"log"

	proto "github.com/RakaiSeto/simple-app-may/service"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var Client proto.ServiceClient
func init() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	
	Client = proto.NewServiceClient(conn)
}

func main() {
	g := gin.Default()
	g.GET("/hello", Tes)
	g.POST("/login", Login)
	g.GET("/logingithub", LoginGithub)
	g.GET("/login/github/callback", LoginGithubCallback)
	g.GET("/login/google", LoginGoogle)
	g.GET("/login/google/callback", LoginGoogleCallback)
	g.GET("/login/facebook", LoginFacebook)
	g.GET("/login/facebook/callback", LoginFacebookCallback)
	g.POST("/logout", Logout)
	g.GET("/checkreq/:id", CheckError)
	g.GET("/user", AllUser)
	g.GET("/user/:id", OneUser)
	g.GET("/profile", MyUser)
	g.POST("/user", PostUser)
	g.PATCH("/user", PatchUser)
	g.DELETE("/user/:id", DeleteUser)
	g.GET("/product", AllProduct)
	g.GET("/product/:id", OneProduct)
	g.POST("/product", PostUser)
	g.PATCH("/product/:id", PatchUser)
	g.DELETE("/product/:id", DeleteUser)
	g.GET("/order", AllOrder)
	g.GET("/order/:id", OneOrder)
	g.POST("/order", PostOrder)
	g.PATCH("/order/:id", PatchOrder)
	g.DELETE("/order/:id", DeleteOrder)
	g.POST("/admin/topup", AdminTopup)

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
	defer Rabconn.Close()
	defer Rabchan.Close()
}