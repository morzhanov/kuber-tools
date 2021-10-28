package grpcserver

import (
	"context"
	"fmt"
	"net"

	"github.com/morzhanov/kuber-tools/internal/errors"
	"github.com/morzhanov/kuber-tools/internal/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/common/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type baseServer struct {
	Tracer opentracing.Tracer
	Logger *zap.Logger
	Uri    string
}

type BaseServer interface {
	PrepareContext(ctx context.Context) (context.Context, opentracing.Span)
	Listen(ctx context.Context, cancel context.CancelFunc, server *grpc.Server)
}

func (s *baseServer) PrepareContext(ctx context.Context) (context.Context, opentracing.Span) {
	span := tracing.StartSpanFromGrpcRequest(ctx, s.Tracer)
	return ctx, span
}

func (s *baseServer) Listen(ctx context.Context, cancel context.CancelFunc, server *grpc.Server) {
	log.Info(fmt.Sprintf("Configuring GRPC server %s", s.Uri))
	lis, err := net.Listen("tcp", s.Uri)
	if err != nil {
		cancel()
		errors.LogInitializationError(err, "grpc server", s.Logger)
		return
	}

	go func() {
		if err := server.Serve(lis); err != nil {
			cancel()
			errors.LogInitializationError(err, "grpc server", s.Logger)
			return
		}
	}()
	log.Info("Grpc server started", zap.String("port", s.Uri))
	<-ctx.Done()
	if err := lis.Close(); err != nil {
		cancel()
		errors.LogInitializationError(err, "grpc server", s.Logger)
		return
	}
}

func NewServer(
	tracer opentracing.Tracer,
	logger *zap.Logger,
	uri string,
) BaseServer {
	return &baseServer{tracer, logger, uri}
}
