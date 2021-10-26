package apigw

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morzhanov/kuber-tools/api/order"
	"github.com/morzhanov/kuber-tools/internal/rest"
	"go.uber.org/zap"
)

type controller struct {
	rest.BaseController
	client Client
}

type Controller interface {
	Listen(port string)
}

func (c *controller) handleHttpErr(ctx *gin.Context, err error) {
	ctx.String(http.StatusInternalServerError, err.Error())
	c.BaseController.Logger().Info("error in the REST handler", zap.Error(err))
}

func (c *controller) handleCreateOrder(ctx *gin.Context) {
	c.Meter().IncReqCount()
	t := c.Tracer()("rest")
	sctx, span := t.Start(ctx, "create-order")
	defer span.End()

	d := order.CreateOrderMessage{}
	if err := c.BaseController.ParseRestBody(ctx, &d); err != nil {
		c.handleHttpErr(ctx, err)
		return
	}
	res, err := c.client.CreateOrder(sctx, &d)
	if err != nil {
		c.handleHttpErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *controller) handleProcessOrder(ctx *gin.Context) {
	c.Meter().IncReqCount()
	t := c.Tracer()("rest")
	sctx, span := t.Start(ctx, "process-order")
	ctx.Set("span-context", span.SpanContext())
	defer span.End()

	id := ctx.Param("id")
	res, err := c.client.ProcessOrder(sctx, id)
	if err != nil {
		c.handleHttpErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *controller) handleGetPaymentInfo(ctx *gin.Context) {
	c.Meter().IncReqCount()
	t := c.Tracer()("rest")
	sctx, span := t.Start(ctx, "get-payment-info")
	ctx.Set("span-context", span.SpanContext())
	defer span.End()

	orderID := ctx.Param("orderID")
	res, err := c.client.GetPaymentInfo(sctx, orderID)
	if err != nil {
		c.handleHttpErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *controller) Listen(port string) {
	c.BaseController.Listen(port)
}

func NewController(
	client Client,
	log *zap.Logger,
) Controller {
	bc := rest.NewBaseController(log, tel)
	c := controller{BaseController: bc, client: client}
	r := bc.Router()
	r.POST("/order", bc.Handler(c.handleCreateOrder))
	r.PUT("/order/:id", bc.Handler(c.handleProcessOrder))
	r.GET("/payment/:orderID", bc.Handler(c.handleGetPaymentInfo))
	return &c
}
