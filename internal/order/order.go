package order

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	porder "github.com/morzhanov/kuber-tools/api/order"
	"github.com/morzhanov/kuber-tools/api/payment"
	"github.com/morzhanov/kuber-tools/internal/rest"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type service struct {
	rest.BaseController
	coll *mongo.Collection
}

type Service interface {
	Listen()
}

func (s *service) handleHttpErr(ctx *gin.Context, err error) {
	ctx.String(http.StatusInternalServerError, err.Error())
	s.BaseController.Logger().Info("error in the REST handler", zap.Error(err))
}

func (s *service) handleCreateOrder(ctx *gin.Context) {
	s.Meter().IncReqCount()
	t := s.Tracer()("rest")
	dbt := s.Tracer()("mongodb")
	parentCtx, err := rest.GetSpanContext(ctx)
	if err != nil {
		s.handleHttpErr(ctx, err)
		return
	}
	_, span := t.Start(*parentCtx, "create-order")
	defer span.End()
	dbctx, dbspan := dbt.Start(*parentCtx, "create-order")
	defer dbspan.End()

	jsonData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return
	}
	d := porder.CreateOrderMessage{}
	if err = json.Unmarshal(jsonData, &d); err != nil {
		s.handleHttpErr(ctx, err)
		return
	}

	id := uuid.NewV4().String()
	msg := porder.OrderMessage{Id: id, Name: d.Name, Amount: d.Amount, Status: "new"}
	_, err = s.coll.InsertOne(dbctx, &msg)
	if err != nil {
		s.handleHttpErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, &msg)
}

func (s *service) handleProcessOrder(ctx *gin.Context) {
	s.Meter().IncReqCount()
	t := s.Tracer()("rest")
	dbt := s.Tracer()("mongodb")
	et := s.Tracer()("kafka")
	parentCtx, err := rest.GetSpanContext(ctx)
	if err != nil {
		s.handleHttpErr(ctx, err)
		return
	}
	_, span := t.Start(*parentCtx, "process-order")
	defer span.End()
	dbctx, dbspan := dbt.Start(*parentCtx, "process-order")
	defer dbspan.End()
	dbctxInsert, dbspanInsert := dbt.Start(*parentCtx, "get-order")
	defer dbspanInsert.End()
	ectx, espan := et.Start(*parentCtx, "process-order")
	defer espan.End()

	id := ctx.Param("id")
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"status", "processed"}}}}
	_, err = s.coll.UpdateOne(dbctx, filter, update)
	if err != nil {
		s.handleHttpErr(ctx, err)
		return
	}
	res := s.coll.FindOne(dbctxInsert, filter)
	if res.Err() != nil {
		s.handleHttpErr(ctx, res.Err())
		return
	}
	msg := porder.OrderMessage{}
	if err := res.Decode(&msg); err != nil {
		s.handleHttpErr(ctx, res.Err())
		return
	}

	if err := s.mq.WriteMessage(ectx, &payment.ProcessPaymentMessage{OrderId: msg.Id, Name: msg.Name, Amount: msg.Amount, Status: msg.Status}); err != nil {
		s.handleHttpErr(ctx, res.Err())
		return
	}
	ctx.JSON(http.StatusOK, &msg)
}

func (s *service) Listen() {
	r := s.BaseController.Router()
	r.POST("/", s.handleCreateOrder)
	r.POST("/:id", s.handleProcessOrder)
	r.Run()
}

func NewService(log *zap.Logger, coll *mongo.Collection) Service {
	bc := rest.NewBaseController(log, tel)
	return &service{BaseController: bc, coll: coll}
}
