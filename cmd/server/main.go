package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/lightsaid/pcbook/pb"
	"github.com/lightsaid/pcbook/service"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()

	log.Println("start server on port: ", *port)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("assets/uploads")

	laptopServer := service.NewLaptopServer(laptopStore, imageStore)
	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

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
