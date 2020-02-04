package server

import (
	"fmt"
	"net"

	"github.com/dapperlabs/flow-go/protobuf/services/observation"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	logger     *logrus.Logger
	port       int
	grpcServer *grpc.Server
}

func NewGRPCServer(logger *logrus.Logger, backend *Backend, port int, debug bool) *GRPCServer {
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcprometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpcprometheus.UnaryServerInterceptor),
	)

	observation.RegisterObserveServiceServer(grpcServer, backend)
	grpcprometheus.Register(grpcServer)

	if debug {
		reflection.Register(grpcServer)
	}

	return &GRPCServer{
		logger:     logger,
		port:       port,
		grpcServer: grpcServer,
	}
}

func (g *GRPCServer) Server() *grpc.Server {
	return g.grpcServer
}

func (g *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		return err
	}

	return g.grpcServer.Serve(lis)
}

func (g *GRPCServer) Stop() {
	g.grpcServer.GracefulStop()
}
