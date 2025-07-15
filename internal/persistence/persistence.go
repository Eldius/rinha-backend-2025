package persistence

import (
	"github.com/eldius/rinha-backend-2025/internal/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strings"
)

type PaymentRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (p *PaymentRepository) Save(pay model.PaymentInfo) error {
	_, err := p.db.NamedExec(
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
				)`, &pay)
	return err
}

func (p *PaymentRepository) Summary() (*model.PaymentsSummary, error) {
	rows, err := p.db.Queryx(`select
    pi.provider
    , sum(pi.amount) AS "total_amount"
    , count(1) AS "total_count"
from
    payment_info pi
group by pi.provider`)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	var summary model.PaymentsSummary
	for rows.Next() {
		var r summaryDBResults
		if err := rows.StructScan(&r); err != nil {
			return nil, err
		}
		if strings.EqualFold(r.Provider, "default") {
			summary.Default = model.ProviderSummary{
				TotalRequests: r.Count,
				TotalAmount:   r.Amount,
			}
		} else {
			summary.Fallback = model.ProviderSummary{
				TotalRequests: r.Count,
				TotalAmount:   r.Amount,
			}
		}
	}

	return &summary, nil
}

type summaryDBResults struct {
	Provider string  `db:"provider"`
	Amount   float64 `db:"total_amount"`
	Count    int     `db:"total_count"`
}
