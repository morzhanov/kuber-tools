package payment

import (
	"context"

	"github.com/jmoiron/sqlx"
	gpayment "github.com/morzhanov/kuber-tools/api/payment"
	uuid "github.com/satori/go.uuid"
)

type pay struct {
	db *sqlx.DB
}

type Payment interface {
	GetPaymentInfo(ctx context.Context, in *gpayment.GetPaymentInfoRequest) (*gpayment.PaymentMessage, error)
	ProcessPayment(ctx context.Context, in *gpayment.ProcessPaymentMessage) error
}

func (p *pay) GetPaymentInfo(ctx context.Context, in *gpayment.GetPaymentInfoRequest) (*gpayment.PaymentMessage, error) {
	pt := p.tel.Tracer()("postgres")
	dbctx, dbspan := pt.Start(ctx, "process-payment")
	defer dbspan.End()

	var (
		id, orderID, name, status string
		amount                    int32
	)
	rows, err := p.db.QueryContext(dbctx, `SELECT * FROM payments WHERE order_id = $id`, in.OrderId)
	if err != nil {
		return nil, err
	}
	if err := rows.Scan(&id, &orderID, &name, &amount, &status); err != nil {
		return nil, err
	}
	return &gpayment.PaymentMessage{Id: id, OrderId: orderID, Name: name, Status: status, Amount: amount}, nil
}

func (p *pay) ProcessPayment(ctx context.Context, in *gpayment.ProcessPaymentMessage) error {
	pt := p.tel.Tracer()("postgres")
	dbctx, dbspan := pt.Start(ctx, "process-payment")
	defer dbspan.End()

	id := uuid.NewV4().String()
	if _, err := p.db.QueryContext(
		dbctx,
		`INSERT INTO payments (id, order_id, name, amount, status) VALUES ($id, $orderId, $name, $amount, $status)`,
		id, in.OrderId, in.Name, in.Amount, in.Status,
	); err != nil {
		return err
	}
	return nil
}

func NewPayment(db *sqlx.DB) Payment {
	return &pay{db}
}
