package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lightsaid/pcbook/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LaptopServer 提供 laptop （笔记本电脑） 的服务
type LaptopServer struct {
	Store LaptopStore

	pb.UnimplementedLaptopServiceServer
}

// NewLaptopServer 创建返回一个 LaptopServer
func NewLaptopServer(s LaptopStore) *LaptopServer {
	return &LaptopServer{
		Store: s,
	}
}

// CreateLaptop 创建 一个 Laptop 服务
func (server *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	fmt.Printf("receive a create-laptop request with id: %s\n", laptop.Id)

	// 当有uuid时，检查是 uuid 是否正确
	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not valid UUID: %v", err)
		}
	} else {
		// 随机生成 uuid 提供给 laptop.id
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}
		laptop.Id = id.String()
	}

	// NOTE: 在client端设置超时ctx为2秒，这里休眠3秒，client是报错了，但是这里（服务端）还是继续往下执行
	// 这是不正确的，因此这里需要处理
	time.Sleep(3 * time.Second)

	// NOTE: 在保存数据到store或数据库时，需要检查ctx.Err是否存在错误
	if ctx.Err() != nil {
		// 超过截至日期
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("deadline is ecxceeded")
			return nil, status.Errorf(codes.DeadlineExceeded, "deadline is ecxceeded")
		}

		if ctx.Err() == context.Canceled {
			return nil, status.Errorf(codes.Canceled, "request canceled")
		}

		return nil, status.Errorf(codes.Internal, "ctx.Err: %s", ctx.Err())
	}

	// 保存到 store
	err := server.Store.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to the store: %v", err)
	}

	fmt.Println("save laptop with id: ", laptop.Id)

	// 返回
	res := &pb.CreateLaptopResponse{Id: laptop.Id}

	return res, nil
}
