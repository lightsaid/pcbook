package service_test

import (
	"context"
	"net"
	"testing"

	"github.com/lightsaid/pcbook/pb"
	"github.com/lightsaid/pcbook/sample"
	"github.com/lightsaid/pcbook/serializer"
	"github.com/lightsaid/pcbook/service"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	// 启动 grpc 服务
	laptopServer, serverAddress := startTestLaptopServer(t)
	laptopClient := newTestLaptopClient(t, serverAddress)

	// new alaptop
	laptop := sample.NewLaptop()
	expectID := laptop.Id
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	// 调用 grpc CreateLaptop 服务，创建一个新的laptop
	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res)
	require.Equal(t, expectID, res.Id)

	// 检查laptop是否存在和正确
	other, err := laptopServer.Store.Find(res.Id)
	require.NoError(t, err)

	requireSampleLaptop(t, laptop, other)
}

func startTestLaptopServer(t *testing.T) (*service.LaptopServer, string) {
	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore())

	grpcServer := grpc.NewServer()

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	lis, err := net.Listen("tcp", ":0") // 随机端口
	require.NoError(t, err)

	go grpcServer.Serve(lis)

	return laptopServer, lis.Addr().String()
}

func newTestLaptopClient(t *testing.T, serverAddress string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(conn)
}

func requireSampleLaptop(t *testing.T, laptop *pb.Laptop, laptop2 *pb.Laptop) {
	json1, err := serializer.ProtobufToJSON(laptop)
	require.NoError(t, err)

	json2, err := serializer.ProtobufToJSON(laptop2)
	require.NoError(t, err)

	require.Equal(t, json1, json2)
}
