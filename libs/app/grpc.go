package app

import (
	"fmt"
	"net"

	"github.com/LSDXXX/libs/api"
	"github.com/LSDXXX/libs/api/grpc"
	stdgrpc "google.golang.org/grpc"
)

// GrpcServer server
type GrpcServer interface {
	Start(port int) error
}

type grpcServer struct {
	server *stdgrpc.Server
}

// NewGrpcServer new
//  @param services
//  @return GrpcServer
func NewGrpcServer(services []api.GrpcService) GrpcServer {
	s := stdgrpc.NewServer(stdgrpc.UnaryInterceptor(grpc.UnaryServerInterceptor))
	for _, service := range services {
		service.Use(s)
	}
	return &grpcServer{server: s}
}

func (s *grpcServer) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	return s.server.Serve(lis)
}
