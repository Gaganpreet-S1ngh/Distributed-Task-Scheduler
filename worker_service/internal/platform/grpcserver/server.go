package grpcserver

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	Server   *grpc.Server
	Listener net.Listener
}

func NewGrpcServer(port string) (*GRPCServer, error) {

	// if port not given then bind to any port
	if port == "" {
		port = "0"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	srv := grpc.NewServer()

	return &GRPCServer{
		Server:   srv,
		Listener: lis,
	}, nil

}
