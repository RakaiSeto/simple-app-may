package main

import (
	"context"
	"fmt"
	"net"

	// order "github.com/RakaiSeto/simple-app-may/server/order"
	// product "github.com/RakaiSeto/simple-app-may/server/product"
	auth "github.com/RakaiSeto/simple-app-may/server/auth"
	user "github.com/RakaiSeto/simple-app-may/server/user"
	proto "github.com/RakaiSeto/simple-app-may/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct{
	proto.ServiceServer
}

func main() {
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	s := Server{}

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

func (s *Server) Login(ctx context.Context, Request *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	user := Request.RequestBody.GetUser()
	response, err := auth.Login(user)
	if err != nil {
		return response, nil
	}
	return response, nil
}

func (s *Server) Logout(ctx context.Context, Request *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	string := Request.RequestBody.GetString_()
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
	response, err := user.OneUser(int(id.GetId()))
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
	userInput := Request.GetUser()
	response, err := user.UpdateUser(userInput)
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

// func (s *Server) AllProduct(ctx context.Context, empty *proto.EmptyStruct) (*proto.Products, error) {
// 	response, err := product.AllProduct()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var returned proto.Products
// 	returned.Product = response
// 	return &returned, nil 
// }

// func (s *Server) OneProduct(ctx context.Context, id *proto.Id) (*proto.Product, error) {
// 	response, err := product.OneProduct(int(id.GetId()))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return response, nil 
// }

// func (s *Server) AddProduct(ctx context.Context, productInput *proto.AdminProduct) (*proto.AddProductStatus, error) {
// 	response, err := product.AddProduct(productInput)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return response, nil 
// }

// func (s *Server) UpdateProduct(ctx context.Context, productInput *proto.AdminProduct) (*proto.ResponseStatus, error) {
// 	response, err := product.UpdateProduct(productInput)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return response, nil 
// }

// func (s *Server) DeleteProduct(ctx context.Context, productInput *proto.AdminProduct) (*proto.ResponseStatus, error) {
// 	response, err := product.DeleteProduct(productInput)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return response, nil 
// }

// func (s *Server) AllOrder(ctx context.Context, userInput *proto.User) (*proto.Orders, error) {
// 	response, err := order.AllOrder(userInput)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var returned proto.Orders
// 	returned.Order = response
// 	return &returned, nil 
// }

// func (s *Server) OneOrder(ctx context.Context, orderInput *proto.Order) (*proto.Order, error) {
// 	response, err := order.OneOrder(orderInput)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return response, nil 
// }

// func (s *Server) AddOrder(ctx context.Context, orderInput *proto.Order) (*proto.AddOrderStatus, error) {
// 	response, err := order.AddOrder(orderInput)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return response, nil 
// }

// func (s *Server) UpdateOrder(ctx context.Context, orderInput *proto.Order) (*proto.ResponseStatus, error) {
// 	response, err := order.UpdateOrder(orderInput)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return response, nil 
// }

// func (s *Server) DeleteOrder(ctx context.Context, orderInput *proto.Order) (*proto.ResponseStatus, error) {
// 	response, err := order.DeleteOrder(orderInput)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return response, nil 
// }