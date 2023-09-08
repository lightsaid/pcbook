package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/lightsaid/pcbook/pb"
	"github.com/lightsaid/pcbook/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Println("--> unary interceptor: ", info.FullMethod)
	return handler(ctx, req)
}

func streamServerInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	log.Println("--> stream interceptor: ", info.FullMethod)
	return handler(srv, ss)
}

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()

	log.Println("start server on port: ", *port)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("assets/uploads")
	ratingStore := service.NewInMemoryRatingStore()

	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.StreamInterceptor(streamServerInterceptor),
	)

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	// 注册gRPC反射服务
	reflection.Register(grpcServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err.Error())
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}
