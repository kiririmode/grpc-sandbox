package router

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/kiririmode/grpc-sandbox/common/conf"
	"github.com/kiririmode/grpc-sandbox/common/log"
	"github.com/kiririmode/grpc-sandbox/helloworld"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// GrpcServer は gRPC サーバそのものを表現する
type GrpcServer struct {
	server   *grpc.Server
	config   *conf.Configuration
	logger   *log.Log
	Listener net.Listener
}

func NewGrpcServer(server *grpc.Server, conf *conf.Configuration, logger *log.Log) *GrpcServer {
	return &GrpcServer{
		server: server,
		config: conf,
		logger: logger,
	}
}

func (s *GrpcServer) Name() string {
	return "grpc server"
}

func (s *GrpcServer) Initialize() error {
	port := s.config.GetInt("server.port")

	s.logger.Logger.Infof("listening to tcp port %d", port)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrapf(err, "failed to listen port %d", port)
	}
	s.Listener = listener

	helloworld.RegisterGreeterServer(s.server, s)
	reflection.Register(s.server)

	return nil
}

func (s *GrpcServer) Finalize() error {
	err := s.Listener.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close listener")
	}
	return nil
}

func (s *GrpcServer) Serve() error {
	if err := s.server.Serve(s.Listener); err != nil {
		return errors.Errorf("failed to serve: %v", err)
	}
	return nil
}

func (s *GrpcServer) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
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
	return &helloworld.HelloReply{Message: fmt.Sprintf("Hello %s", req.Name)}, nil
}

func (s *GrpcServer) SayHelloToMany(stream helloworld.Greeter_SayHelloToManyServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		err = stream.Send(&helloworld.HelloReply{Message: fmt.Sprintf("Hello %s", req.Name)})
		if err != nil {
			return err
		}
	}
}
