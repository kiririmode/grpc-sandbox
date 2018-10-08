package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

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
		log.Printf("meatadata postscript: %s", ps)
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
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, newServer())

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
