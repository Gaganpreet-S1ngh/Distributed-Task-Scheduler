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
