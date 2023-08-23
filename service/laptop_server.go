package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/google/uuid"
	"github.com/lightsaid/pcbook/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxImageSize = 1 << 20 // 1MB

// LaptopServer 提供 laptop （笔记本电脑） 的服务
type LaptopServer struct {
	laptopStore LaptopStore
	imageStore  ImageStore
	ratingStore RatingStore

	pb.UnimplementedLaptopServiceServer
}

// NewLaptopServer 创建返回一个 LaptopServer
func NewLaptopServer(s LaptopStore, i ImageStore, r RatingStore) *LaptopServer {
	return &LaptopServer{
		laptopStore: s,
		imageStore:  i,
		ratingStore: r,
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
	// time.Sleep(3 * time.Second)

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
	err := server.laptopStore.Save(laptop)
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

// SearchLaptop 查找 laptop，  server-streaming gRPC (服务端流gRPC)
func (server *LaptopServer) SearchLaptop(req *pb.SearchLaptopRequest, stream pb.LaptopService_SearchLaptopServer) error {
	filter := req.GetFilter()
	log.Printf("receive a search-laptop request with filter: %v\n", filter)

	err := server.laptopStore.Search(
		stream.Context(),
		filter,
		func(laptop *pb.Laptop) error {
			res := &pb.SearchLaptopResponse{Laptop: laptop}
			if err := stream.Send(res); err != nil {
				return err
			}
			log.Printf("sent laptop with id: %s", laptop.GetId())

			return nil
		},
	)
	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil
}

// UploadImage 上传图片 client-streaming gRPC (客户端流)
func (server *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {
	// 首先接受图片信息
	req, err := stream.Recv()
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot receive image info"))
	}

	laptopID := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()
	log.Printf("receive an upload-image request for laptop %s with image type %s", laptopID, imageType)

	// 查询是否存在laptop
	laptop, err := server.laptopStore.Find(laptopID)
	if err != nil {
		return logError(status.Errorf(codes.NotFound, "cannot find laptop: %v", err))
	}

	if laptop == nil {
		return logError(status.Errorf(codes.InvalidArgument, "laptop id %s doesn't exist", laptopID))
	}

	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}

		// log.Print("waiting to receive more data")

		req, err = stream.Recv()
		if err == io.EOF {
			log.Println("no more data")
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err))
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		// log.Printf("received a chunk with size: %d", size)

		imageSize += size

		if imageSize > maxImageSize {
			return status.Errorf(codes.InvalidArgument, "to long")
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
		}
	}

	imageID, err := server.imageStore.Save(laptopID, imageType, imageData)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot save image to the store: %v", err))
	}

	res := &pb.UploadImageResponse{
		Id:   imageID,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}

	log.Printf("saved image with id: %s, size: %d", imageID, imageSize)
	return nil
}

// RateLaptop 评分 bidirectional-streaming gRPC (双向流)
func (server *LaptopServer) RateLaptop(stream pb.LaptopService_RateLaptopServer) error {
	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot receive stream request: %v", err))
		}

		laptopID := req.GetLaptopId()
		score := req.GetScore()

		log.Printf("received a rate-laptop request: id = %s, score = %.2f", laptopID, score)

		found, err := server.laptopStore.Find(laptopID)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot find laptop: %v", err))
		}
		if found == nil {
			return logError(status.Errorf(codes.NotFound, "laptopID %s is not found", laptopID))
		}

		rating, err := server.ratingStore.Add(laptopID, score)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot add rating to the store: %v", err))
		}

		res := &pb.RateLaptopResponse{
			LaptopId:     laptopID,
			RatedCount:   rating.Count,
			AverageScore: rating.Sum / float64(rating.Count),
		}

		err = stream.Send(res)
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot send stream response: %v", err))
		}
	}

	return nil
}

func logError(err error) error {
	if err != nil {
		log.Println(err)
	}

	return err
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is canceled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded"))
	default:
		return nil
	}
}
