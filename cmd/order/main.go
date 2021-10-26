package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/morzhanov/kuber-tools/internal/config"
	"github.com/morzhanov/kuber-tools/internal/logger"
	"github.com/morzhanov/kuber-tools/internal/mongodb"
	"github.com/morzhanov/kuber-tools/internal/order"
	"go.uber.org/zap"
)

func failOnError(l *zap.Logger, step string, err error) {
	if err != nil {
		l.Fatal("initialization error", zap.Error(err), zap.String("step", step))
	}
}

func main() {
	l, err := logger.NewLogger()
	if err != nil {
		log.Fatal("initialization error during logger setup")
	}
	c, err := config.NewConfig()
	failOnError(l, "config", err)
	t, err := telemetry.NewTelemetry(c.JaegerURL, "order", l)
	failOnError(l, "telemetry", err)
	m, err := mongodb.NewMongoDB(c.MongoURL)
	failOnError(l, "mongodb", err)
	msgq, err := mq.NewMq(c.KafkaURL, c.KafkaTopic)
	failOnError(l, "message_queue", err)

	srv := order.NewService(l, t, m, msgq)
	go srv.Listen()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	log.Println("App successfully started!")
	<-quit
	log.Println("received os.Interrupt, exiting...")
}
