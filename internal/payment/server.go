package payment

import (
	"context"
	"fmt"

	gpayment "github.com/morzhanov/kuber-tools/api/payment"
	gserver "github.com/morzhanov/kuber-tools/internal/grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	gpayment.UnimplementedPaymentServer
	gserver.BaseServer
	server *grpc.Server
	pay    Service
	tracer opentracing.Tracer
}

type Server interface {
	Listen(ctx context.Context, cancel context.CancelFunc)
}

func (s *server) GetPaymentInfo(ctx context.Context, in *gpayment.GetPaymentInfoRequest) (*gpayment.PaymentMessage, error) {
	ctx, span := s.PrepareContext(ctx)
	defer span.Finish()
	dbSpan := s.tracer.StartSpan("postgresql", ext.RPCServerOption(span.Context()))
	dbCtx := context.WithValue(ctx, "span-context", dbSpan.Context())
	defer dbSpan.Finish()
	return s.pay.GetPaymentInfo(dbCtx, in)
}

func (s *server) ProcessPayment(ctx context.Context, in *gpayment.ProcessPaymentRequest) (*gpayment.PaymentMessage, error) {
	ctx, span := s.PrepareContext(ctx)
	defer span.Finish()
	dbSpan := s.tracer.StartSpan("postgresql", ext.RPCServerOption(span.Context()))
	dbCtx := context.WithValue(ctx, "span-context", dbSpan.Context())
	defer dbSpan.Finish()
	dbFindSpan := s.tracer.StartSpan("postgresql", ext.RPCServerOption(span.Context()))
	dbFindCtx := context.WithValue(ctx, "span-context", dbFindSpan.Context())
	defer dbFindSpan.Finish()

	if err := s.pay.ProcessPayment(dbCtx, in); err != nil {
		return nil, err
	}
	fMsg := &gpayment.GetPaymentInfoRequest{OrderId: in.OrderId}
	return s.pay.GetPaymentInfo(dbFindCtx, fMsg)
}

func (s *server) Listen(ctx context.Context, cancel context.CancelFunc) {
	s.BaseServer.Listen(ctx, cancel, s.server)
}

func NewServer(
	grpcAddr string,
	grpcPort string,
	pay Service,
	logger *zap.Logger,
	tracer opentracing.Tracer,
) Server {
	url := fmt.Sprintf("%s:%s", grpcAddr, grpcPort)
	bs := gserver.NewServer(tracer, logger, url)
	s := server{BaseServer: bs, server: grpc.NewServer(), tracer: tracer, pay: pay}
	gpayment.RegisterPaymentServer(s.server, &s)
	reflection.Register(s.server)
	return &s
}
