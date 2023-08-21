package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/lightsaid/pcbook/pb"
	"github.com/lightsaid/pcbook/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	laptopClient := pb.NewLaptopServiceClient(conn)

	laptop := sample.NewLaptop()
	// laptop.Id = ""
	// laptop.Id = "invalid id"
	// laptop.Id = "a828d1ae-07db-4fef-8779-4f1cedcec686"

	req := pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := laptopClient.CreateLaptop(ctx, &req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Println("laptop already exists")
		} else {
			log.Fatal("cannot create laptop: ", err)
		}
		return
	}
	log.Println("created laptop with id: ", res.Id)
}
