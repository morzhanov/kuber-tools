package order

import (
	"context"
	"fmt"

	porder "github.com/morzhanov/kuber-tools/api/order"
	"github.com/morzhanov/kuber-tools/api/payment"
	grpcserver "github.com/morzhanov/kuber-tools/internal/grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	porder.UnimplementedOrderServer
	server *grpc.Server
	grpcserver.BaseServer
	coll      *mongo.Collection
	tracer    opentracing.Tracer
	payClient payment.PaymentClient
}

type Server interface {
	Listen(ctx context.Context, cancel context.CancelFunc)
}

func (s *server) CreateOrder(ctx context.Context, msg *porder.CreateOrderRequest) (*porder.OrderMessage, error) {
	ctx, span := s.PrepareContext(ctx)
	defer span.Finish()
	dbSpan := s.tracer.StartSpan("mongodb", ext.RPCServerOption(span.Context()))
	dbCtx := context.WithValue(ctx, "span-context", dbSpan.Context())
	defer dbSpan.Finish()

	id := uuid.NewV4().String()
	o := porder.OrderMessage{Id: id, Name: msg.Name, Amount: msg.Amount, Status: "new"}
	if _, err := s.coll.InsertOne(dbCtx, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *server) ProcessOrder(ctx context.Context, msg *porder.ProcessOrderRequest) (*porder.OrderMessage, error) {
	ctx, span := s.PrepareContext(ctx)
	defer span.Finish()
	dbSpan := s.tracer.StartSpan("mongodb", ext.RPCServerOption(span.Context()))
	dbCtx := context.WithValue(ctx, "span-context", dbSpan.Context())
	defer dbSpan.Finish()
	dbFindSpan := s.tracer.StartSpan("mongodb", ext.RPCServerOption(span.Context()))
	dbFindCtx := context.WithValue(ctx, "span-context", dbFindSpan.Context())
	defer dbFindSpan.Finish()
	paySpan := s.tracer.StartSpan("grpc-request", ext.RPCServerOption(span.Context()))
	payCtx := context.WithValue(ctx, "span-context", paySpan.Context())
	defer paySpan.Finish()

	filter := bson.D{{"_id", msg.Id}}
	update := bson.D{{"$set", bson.D{{"status", "processed"}}}}
	if _, err := s.coll.UpdateOne(dbCtx, filter, update); err != nil {
		return nil, err
	}
	res := s.coll.FindOne(dbFindCtx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	resMsg := porder.OrderMessage{}
	if err := res.Decode(&resMsg); err != nil {
		return nil, err
	}

	pMsg := &payment.ProcessPaymentRequest{OrderId: resMsg.Id, Name: resMsg.Name, Amount: resMsg.Amount, Status: resMsg.Status}
	if _, err := s.payClient.ProcessPayment(payCtx, pMsg); err != nil {
		return nil, err
	}
	return &resMsg, nil
}

func (s *server) Listen(ctx context.Context, cancel context.CancelFunc) {
	s.BaseServer.Listen(ctx, cancel, s.server)
}

func NewServer(
	grpcUrl string,
	grpcPort string,
	payClient payment.PaymentClient,
	coll *mongo.Collection,
	tracer opentracing.Tracer,
	logger *zap.Logger,
) Server {
	uri := fmt.Sprintf("%s:%s", grpcUrl, grpcPort)
	bs := grpcserver.NewServer(tracer, logger, uri)
	s := server{BaseServer: bs, server: grpc.NewServer(), tracer: tracer, payClient: payClient, coll: coll}
	porder.RegisterOrderServer(s.server, &s)
	reflection.Register(s.server)
	return &s
}
