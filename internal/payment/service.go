package payment

import (
	"context"

	"github.com/jmoiron/sqlx"
	gpayment "github.com/morzhanov/kuber-tools/api/payment"
	uuid "github.com/satori/go.uuid"
)

type service struct {
	db *sqlx.DB
}

type Service interface {
	GetPaymentInfo(ctx context.Context, in *gpayment.GetPaymentInfoRequest) (*gpayment.PaymentMessage, error)
	ProcessPayment(ctx context.Context, in *gpayment.ProcessPaymentRequest) error
}

func (p *service) GetPaymentInfo(ctx context.Context, in *gpayment.GetPaymentInfoRequest) (*gpayment.PaymentMessage, error) {
	var (
		id, orderID, name, status string
		amount                    int32
	)
	row := p.db.QueryRowContext(ctx, `SELECT * FROM payments WHERE order_id = $1`, in.OrderId)
	if err := row.Scan(&id, &orderID, &name, &amount, &status); err != nil {
		return nil, err
	}
	return &gpayment.PaymentMessage{Id: id, OrderId: orderID, Name: name, Status: status, Amount: amount}, nil
}

func (p *service) ProcessPayment(ctx context.Context, in *gpayment.ProcessPaymentRequest) error {
	id := uuid.NewV4().String()
	if _, err := p.db.QueryContext(
		ctx,
		`INSERT INTO payments (id, order_id, name, amount, status) VALUES ($1, $2, $3, $4, $5)`,
		id, in.OrderId, in.Name, in.Amount, in.Status,
	); err != nil {
		return err
	}
	return nil
}

func NewService(db *sqlx.DB) Service {
	return &service{db}
}
