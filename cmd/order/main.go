package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/morzhanov/kuber-tools/api/payment"
	"github.com/morzhanov/kuber-tools/internal/logger"
	"github.com/morzhanov/kuber-tools/internal/mongodb"
	"github.com/morzhanov/kuber-tools/internal/order"
	"github.com/morzhanov/kuber-tools/internal/order/config"
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

	l, err := logger.NewLogger("order")
	if err != nil {
		log.Fatal("initialization error during logger setup")
	}
	c, err := config.NewConfig()
	failOnError(l, cancel, "config", err)
	t, err := tracing.NewTracer(ctx, l, "order")
	failOnError(l, cancel, "tracer", err)
	db, err := mongodb.NewMongoDB(c.MongoURL)
	failOnError(l, cancel, "mongodb", err)

	payUrl := fmt.Sprintf("%s:%s", c.PaymentGRPCurl, c.PaymentGRPCport)
	payConn, err := grpc.Dial(payUrl, grpc.WithInsecure(), grpc.WithBlock())
	failOnError(l, cancel, "payment_grpc", err)
	payClient := payment.NewPaymentClient(payConn)

	s := order.NewServer(c.URL, c.Port, payClient, db, t, l)

	l.Info("Order service successfully started!")
	s.Listen(ctx, cancel)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	l.Info("received os.Interrupt, exiting...")
}
