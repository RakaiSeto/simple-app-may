package main

import (
	"context"
	"fmt"
	"net"

	// product "github.com/RakaiSeto/simple-app-may/server/product"
	auth "github.com/RakaiSeto/simple-app-may/server/auth"
	"github.com/RakaiSeto/simple-app-may/server/product"
	user "github.com/RakaiSeto/simple-app-may/server/user"
	proto "github.com/RakaiSeto/simple-app-may/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct{
	proto.ServiceServer
}
var Client proto.ServiceClient
var Context context.Context

func main() {
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	s := Server{}

	Context = context.TODO()

	conn, err := grpc.Dial("localhost:0404", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	
	Client = proto.NewServiceClient(conn)

	srv := grpc.NewServer()
	proto.RegisterServiceServer(srv, &s)
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (s *Server) Tes(ctx context.Context, empty *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	return &proto.ResponseWrapper{
		Code: 200, Message: "success", ResponseBody: &proto.ResponseBody{
			ResponseStatus: &proto.ResponseStatus{Response: "Hello"},
			},
		}, nil
}

func (s *Server) LoginGithub(ctx context.Context, Request *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	response, err := auth.LoginGithub()
	fmt.Println(response.ResponseBody.ResponseStatus.GetResponse())
	if err != nil {
		return response, nil
	}
	return response, nil
}

func (s *Server) LoginGithubCallback(ctx context.Context, Request *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	fmt.Println(Request.RequestBody.GetString_())
	response, err := auth.LoginGithubCallback(Request.RequestBody.OauthCallback.GetState(), Request.RequestBody.OauthCallback.GetCode())
	if err != nil {
		return response, nil
	}
	return response, nil
}

func (s *Server) LoginGoogle(ctx context.Context, Request *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	response, err := auth.LoginGoogle()
	fmt.Println(response.ResponseBody.ResponseStatus.GetResponse())
	if err != nil {
		return response, nil
	}
	return response, nil
}

func (s *Server) LoginGoogleCallback(ctx context.Context, Request *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	fmt.Println(Request.RequestBody.GetString_())
	response, err := auth.LoginGoogleCallback(Request.RequestBody.OauthCallback.GetState(), Request.RequestBody.OauthCallback.GetCode())
	if err != nil {
		return response, nil
	}
	return response, nil
}

func (s *Server) LoginFacebook(ctx context.Context, Request *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	response, err := auth.LoginFacebook()
	fmt.Println(response.ResponseBody.ResponseStatus.GetResponse())
	if err != nil {
		return response, nil
	}
	return response, nil
}

func (s *Server) LoginFacebookCallback(ctx context.Context, Request *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	fmt.Println(Request.RequestBody.GetString_())
	response, err := auth.LoginFacebookCallback(Request.RequestBody.OauthCallback.GetState(), Request.RequestBody.OauthCallback.GetCode())
	if err != nil {
		return response, nil
	}
	return response, nil
}

func (s *Server) Login(ctx context.Context, Request *proto.RequestBody) (*proto.ResponseWrapper, error) {
	user := Request.GetUser()
	response, err := auth.Login(user)
	if err != nil {
		return response, nil
	}
	return response, nil
}

func (s *Server) Logout(ctx context.Context, Request *proto.RequestBody) (*proto.ResponseWrapper, error) {
	string := Request.GetString_()
	response, err := auth.Logout(string)
	if err != nil {
		return response, nil
	}
	return response, nil
}

func (s *Server) ReqError(ctx context.Context, Request *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	string_ := Request.RequestBody.GetString_()
	response, err := user.CheckRequest(string_)
	if err != nil {
		return response, nil
	}
	return response, nil
}

func (s *Server) AllUser(ctx context.Context, Request *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	response, err := user.AllUser()
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) OneUser(ctx context.Context, Request *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	id := Request.RequestBody.GetId()
	response, err := user.OneUser(int(id))
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) AddUser(ctx context.Context, Request *proto.RequestBody) (*proto.ResponseWrapper, error) {
	userInput := Request.GetUser()
	response, err := user.AddUser(userInput)
	fmt.Println(response)
	fmt.Println(err)
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) UpdateUser(ctx context.Context, Request *proto.RequestBody) (*proto.ResponseWrapper, error) {
	response, err := user.UpdateUser(Request)
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) DeleteUser(ctx context.Context, Request *proto.RequestBody) (*proto.ResponseWrapper, error) {
	userInput := Request.GetUser()
	response, err := user.DeleteUser(userInput)
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) AllProduct(ctx context.Context, empty *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	response, err := product.AllProduct(empty.RequestBody)
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) OneProduct(ctx context.Context, id *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	response, err := product.OneProduct(id.RequestBody)
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) AddProduct(ctx context.Context, productInput *proto.RequestBody) (*proto.ResponseWrapper, error) {
	response, err := product.AddProduct(productInput)
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) UpdateProduct(ctx context.Context, productInput *proto.RequestBody) (*proto.ResponseWrapper, error) {
	response, err := product.UpdateProduct(productInput)
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) DeleteProduct(ctx context.Context, productInput *proto.RequestBody) (*proto.ResponseWrapper, error) {
	response, err := product.DeleteProduct(productInput)
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) AllOrder(ctx context.Context, input *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	response, err := Client.AllOrder(Context, input)
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) OneOrder(ctx context.Context, input *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	response, err := Client.OneOrder(Context, input)
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) MyUser(ctx context.Context, input *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	response, err := user.MyUser(input.RequestBody.GetString_())
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) AddOrder(ctx context.Context, input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	userGet, _ := user.MyUser(input.GetString_())
	if userGet.ResponseBody.Error != nil {
		return userGet, nil
	}

	var products []*proto.Product
	for _, v := range input.Order.OrderProducts.OrderProduct {
		fmt.Println(int64(v.GetId()))
		productInput, err := product.GetProduct(int64(v.GetId()))
		if err != nil {
			errString := err.Error()
			return &proto.ResponseWrapper{Code: 500, Message:"unknown", ResponseBody: &proto.ResponseBody{Error: &errString}}, nil
		}
		products = append(products, productInput)
	}
	
	output := &proto.RequestBody{Order: input.GetOrder(), Products: &proto.Products{Product: products}, User: userGet.ResponseBody.GetUser()}
	response, err := Client.AddOrder(Context, output)
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) DeleteOrder(ctx context.Context, input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	response, err := Client.DeleteOrder(Context, input)
	if err != nil {
		return response, nil
	}
	return response, nil 
}