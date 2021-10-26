package apigw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/morzhanov/kuber-tools/api/order"
	"github.com/morzhanov/kuber-tools/api/payment"
)

type client struct {
	orderUrl      string
	paymentClient payment.PaymentClient
}

type Client interface {
	CreateOrder(ctx context.Context, msg *order.CreateOrderMessage) (*order.OrderMessage, error)
	ProcessOrder(ctx context.Context, orderID string) (*order.OrderMessage, error)
	GetPaymentInfo(ctx context.Context, orderID string) (*payment.PaymentMessage, error)
}

func (c *client) CreateOrder(ctx context.Context, msg *order.CreateOrderMessage) (*order.OrderMessage, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.orderUrl, bytes.NewReader(b))
	body, err := rest.PerformRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	o := order.OrderMessage{}
	if err := json.Unmarshal(body, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (c *client) ProcessOrder(ctx context.Context, orderID string) (*order.OrderMessage, error) {
	url := fmt.Sprintf("%s/%s", c.orderUrl, orderID)
	req, err := http.NewRequest("PUT", url, nil)
	body, err := rest.PerformRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	o := order.OrderMessage{}
	if err := json.Unmarshal(body, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (c *client) GetPaymentInfo(ctx context.Context, orderID string) (*payment.PaymentMessage, error) {
	msg := payment.GetPaymentInfoRequest{OrderId: orderID}
	return c.paymentClient.GetPaymentInfo(ctx, &msg)
}

func NewClient(orderUrl string, paymentClient payment.PaymentClient) Client {
	return &client{orderUrl, paymentClient}
}
