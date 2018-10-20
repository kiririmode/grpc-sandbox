package main

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/kiririmode/grpc-sandbox/common"
	"github.com/kiririmode/grpc-sandbox/common/conf"
	"github.com/kiririmode/grpc-sandbox/common/log"
	pb "github.com/kiririmode/grpc-sandbox/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	if req.Name == "error" {
		return nil, status.Error(codes.Internal, "Internal Error")
	}

	// retrieve metadata
	md, _ := metadata.FromIncomingContext(ctx)
	postscripts := md.Get("postscript")
	ps := ""
	if len(postscripts) > 0 {
		ps = postscripts[0]
	}

	// send reply with metadata
	grpc.SendHeader(ctx, metadata.Pairs("postscript", ps))
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %s", req.Name)}, nil
}

func (s *server) SayHelloToMany(stream pb.Greeter_SayHelloToManyServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		err = stream.Send(&pb.HelloReply{Message: fmt.Sprintf("Hello %s", req.Name)})
		if err != nil {
			return err
		}
	}
}

func newServer() *server {
	s := &server{}
	return s
}

func main() {

	// リソースの準備
	config := conf.NewConfiguration("stubserver", "development", []string{"conf"})
	logr := log.NewLog(config)

	// リソースの開始・終了処理
	rm := common.NewResourceManager([]common.Resource{config, logr})
	rm.Initialize()
	defer rm.Finalize()

	logger := logr.Logger
	logger.Info("initialization succeeds")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8000))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, newServer())

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	if err = grpcServer.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}
