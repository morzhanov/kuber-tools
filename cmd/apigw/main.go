package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/morzhanov/kuber-tools/api/order"
	"github.com/morzhanov/kuber-tools/api/payment"
	"github.com/morzhanov/kuber-tools/internal/apigw"
	"github.com/morzhanov/kuber-tools/internal/apigw/config"
	"github.com/morzhanov/kuber-tools/internal/logger"
	"github.com/morzhanov/kuber-tools/internal/metrics"
	"github.com/morzhanov/kuber-tools/internal/tracing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func failOnError(l *zap.Logger, cancel context.CancelFunc, step string, err error) {
	if err != nil {
		cancel()
		l.Fatal("initialization error", zap.Error(err), zap.String("step", step))
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l, err := logger.NewLogger("apigw")
	if err != nil {
		log.Fatal("initialization error during logger setup")
	}
	c, err := config.NewConfig()
	failOnError(l, cancel, "config", err)
	t, err := tracing.NewTracer(ctx, l, "apigw")
	failOnError(l, cancel, "tracer", err)
	m := metrics.NewMetricsCollector("apigw")

	payUrl := fmt.Sprintf("%s:%s", c.PaymentGRPCurl, c.PaymentGRPCport)
	payConn, err := grpc.Dial(payUrl, grpc.WithInsecure(), grpc.WithBlock())
	failOnError(l, cancel, "payment_grpc", err)
	payClient := payment.NewPaymentClient(payConn)

	orderUrl := fmt.Sprintf("%s:%s", c.OrderGRPCurl, c.OrderGRPCport)
	orderConn, err := grpc.Dial(orderUrl, grpc.WithInsecure(), grpc.WithBlock())
	failOnError(l, cancel, "order_grpc", err)
	orderClient := order.NewOrderClient(orderConn)

	s := apigw.NewController(orderClient, payClient, t, l, m)
	s.Listen(ctx, cancel, c.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	l.Info("API Gateway service successfully started!")
	<-quit
	l.Info("received os.Interrupt, exiting...")
}
