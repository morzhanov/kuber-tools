package grpc

import (
	"context"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type baseServer struct {
	url string
	log *zap.Logger
}

type BaseServer interface {
	Listen(ctx context.Context, cancel context.CancelFunc, server *grpc.Server)
	Logger() *zap.Logger
	Tracer() tracer.TraceFn
	Meter() meter.Meter
}

func (s *baseServer) Listen(ctx context.Context, cancel context.CancelFunc, server *grpc.Server) {
	lis, err := net.Listen("tcp", s.url)
	if err != nil {
		cancel()
		s.log.Fatal("error during grpc server setup", zap.Error(err))
		return
	}

	if err := server.Serve(lis); err != nil {
		cancel()
		s.log.Fatal("error during grpc server setup", zap.Error(err))
		return
	}
	s.log.Info("Grpc server started", zap.String("port", s.url))
	<-ctx.Done()
	if err := lis.Close(); err != nil {
		cancel()
		s.log.Fatal("error during grpc server setup", zap.Error(err))
		return
	}
}

func (s *baseServer) Logger() *zap.Logger       { return s.log }
func (s *baseServer) Tracer() telemetry.TraceFn { return s.tel.Tracer() }
func (s *baseServer) Meter() meter.Meter        { return s.tel.Meter() }

func NewServer(url string, log *zap.Logger) BaseServer {
	return &baseServer{log: log, url: url}
}
