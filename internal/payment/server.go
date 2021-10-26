package payment

import (
	"context"
	"fmt"
	"net"

	gpayment "github.com/morzhanov/kuber-tools/api/payment"
	gserver "github.com/morzhanov/kuber-tools/internal/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	gpayment.UnimplementedPaymentServer
	gserver.BaseServer
	srv *grpc.Server
	url string
	pay Payment
}

type Server interface {
	Listen(ctx context.Context, cancel context.CancelFunc, server *grpc.Server)
}

func (s *server) GetPaymentInfo(ctx context.Context, in *gpayment.GetPaymentInfoRequest) (*gpayment.PaymentMessage, error) {
	s.Meter().IncReqCount()
	t := s.Tracer()("rest")
	sctx, span := t.Start(ctx, "get-payment-info")
	defer span.End()
	return s.pay.GetPaymentInfo(sctx, in)
}

func (s *server) Listen(ctx context.Context, cancel context.CancelFunc, srv *grpc.Server) {
	lis, err := net.Listen("tcp", s.url)
	if err != nil {
		cancel()
		s.BaseServer.Logger().Fatal("error during grpc service start")
		return
	}
	if err := srv.Serve(lis); err != nil {
		cancel()
		s.BaseServer.Logger().Fatal("error during grpc service start")
		return
	}
	s.BaseServer.Logger().Info("Grpc srv started", zap.String("port", s.url))
	<-ctx.Done()
	if err := lis.Close(); err != nil {
		cancel()
		s.BaseServer.Logger().Fatal("error during grpc service start")
		return
	}
}

func NewServer(
	grpcAddr string,
	grpcPort string,
	logger *zap.Logger,
	pay Payment,
) Server {
	url := fmt.Sprintf("%s:%s", grpcAddr, grpcPort)
	bs := gserver.NewServer(url, logger)
	s := &server{BaseServer: bs, srv: grpc.NewServer(), url: url, pay: pay}
	gpayment.RegisterPaymentServer(s.srv, s)
	reflection.Register(s.srv)
	return s
}
