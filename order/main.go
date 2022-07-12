package main

import (
	"context"
	// "fmt"
	"net"
	proto "github.com/RakaiSeto/simple-app-may/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct{
	proto.ServiceServer
}

func main() {
	listener, err := net.Listen("tcp", ":0404")
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


func (s *Server) AllOrder(ctx context.Context, input *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	response, err := AllOrder(input.RequestBody)
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) OneOrder(ctx context.Context, input *proto.RequestWrapper) (*proto.ResponseWrapper, error) {
	response, err := OneOrder(input.RequestBody)
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) AddOrder(ctx context.Context, input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	response, err := AddOrder(input)
	if err != nil {
		return response, nil
	}
	return response, nil 
}

func (s *Server) DeleteOrder(ctx context.Context, input *proto.RequestBody) (*proto.ResponseWrapper, error) {
	response, err := DeleteOrder(input)
	if err != nil {
		return response, nil
	}
	return response, nil 
}