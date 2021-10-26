package apigw

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morzhanov/kuber-tools/api/order"
	"github.com/morzhanov/kuber-tools/api/payment"
	"github.com/morzhanov/kuber-tools/internal/metrics"
	"github.com/morzhanov/kuber-tools/internal/rest"
	"github.com/morzhanov/kuber-tools/internal/tracing"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type controller struct {
	rest.BaseController
	orderClient   order.OrderClient
	paymentClient payment.PaymentClient
}

type Controller interface {
	Listen(ctx context.Context, cancel context.CancelFunc, port string)
}

func (c *controller) handleHttpErr(ctx *gin.Context, err error) {
	ctx.String(http.StatusInternalServerError, err.Error())
	c.BaseController.Logger().Info("error in the REST handler", zap.Error(err))
}

func (c *controller) handleCreateOrder(ctx *gin.Context) {
	span := c.GetSpan(ctx)
	msg := order.CreateOrderRequest{}
	if err := c.ParseRestBody(ctx, &msg); err != nil {
		c.handleHttpErr(ctx, err)
		return
	}

	reqCtx, err := tracing.InjectGrpcSpan(ctx, *span)
	if err != nil {
		c.handleHttpErr(ctx, err)
	}

	res, err := c.orderClient.CreateOrder(reqCtx, &msg)
	if err != nil {
		c.handleHttpErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *controller) handleProcessOrder(ctx *gin.Context) {
	span := c.GetSpan(ctx)
	id := ctx.Param("id")

	reqCtx, err := tracing.InjectGrpcSpan(ctx, *span)
	if err != nil {
		c.handleHttpErr(ctx, err)
	}

	msg := order.ProcessOrderRequest{Id: id}
	res, err := c.orderClient.ProcessOrder(reqCtx, &msg)
	if err != nil {
		c.handleHttpErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *controller) handleGetPaymentInfo(ctx *gin.Context) {
	span := c.GetSpan(ctx)
	orderID := ctx.Param("orderID")

	reqCtx, err := tracing.InjectGrpcSpan(ctx, *span)
	if err != nil {
		c.handleHttpErr(ctx, err)
	}

	msg := payment.GetPaymentInfoRequest{OrderId: orderID}
	res, err := c.paymentClient.GetPaymentInfo(reqCtx, &msg)
	if err != nil {
		c.handleHttpErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *controller) Listen(ctx context.Context, cancel context.CancelFunc, port string) {
	c.BaseController.Listen(ctx, cancel, port)
}

func NewController(
	oc order.OrderClient,
	pc payment.PaymentClient,
	t opentracing.Tracer,
	l *zap.Logger,
	mc metrics.Collector,
) Controller {
	bc := rest.NewController(t, l, mc)
	c := controller{BaseController: bc, orderClient: oc, paymentClient: pc}
	r := bc.Router()
	r.POST("/order", bc.Handler(c.handleCreateOrder))
	r.PUT("/order/:id", bc.Handler(c.handleProcessOrder))
	r.GET("/payment/:orderID", bc.Handler(c.handleGetPaymentInfo))
	return &c
}
