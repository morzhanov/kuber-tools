package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/morzhanov/kuber-tools/internal/errors"

	"github.com/morzhanov/kuber-tools/internal/logger"
	"github.com/morzhanov/kuber-tools/internal/payment"
	"github.com/morzhanov/kuber-tools/internal/payment/config"
	"github.com/morzhanov/kuber-tools/internal/psql"
	"github.com/morzhanov/kuber-tools/internal/tracing"
	"go.uber.org/zap"
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

	l, err := logger.NewLogger()
	if err != nil {
		log.Fatal("initialization error during logger setup")
	}
	c, err := config.NewConfig()
	failOnError(l, cancel, "config", err)
	t, err := tracing.NewTracer(ctx, l, "payment")
	failOnError(l, cancel, "tracer", err)
	db, err := psql.NewDb(c.PostgresURL)
	failOnError(l, cancel, "postgresql", err)
	srv := payment.NewService(db)

	if err := psql.RunMigrations(db, "payment"); err != nil {
		cancel()
		errors.LogInitializationError(err, "migrations", l)
	}
	l.Info("all database migrations applied...")

	s := payment.NewServer(c.URL, c.Port, srv, l, t)

	go s.Listen(ctx, cancel)
	l.Info("Payment service successfully started!")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	l.Info("received os.Interrupt, exiting...")
}
