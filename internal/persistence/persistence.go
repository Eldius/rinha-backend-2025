package persistence

import (
	"context"
	"github.com/eldius/rinha-backend-2025/internal/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PaymentRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (p *PaymentRepository) Save(ctx context.Context, pay model.PaymentInfo) error {
	_, err := p.db.NamedExecContext(ctx,
		`INSERT INTO payment_info (
					correlation_id,
					amount,
					requested_at,
					message,
					provider,
					status,
					created_at
				) VALUES (
					:correlation_id,
					:amount,
					:requested_at,
					:message,
					:provider,
					:status,
					:created_at
				)`, pay)
	return err
}
