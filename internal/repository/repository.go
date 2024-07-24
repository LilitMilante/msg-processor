package repository

import (
	"context"

	"msg-processor/internal/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r Repository) CreateMsg(ctx context.Context, msg entity.Msg) error {
	q := "INSERT INTO messages (id, text, state, created_at) VALUES ($1, $2, $3, $4)"

	_, err := r.db.Exec(ctx, q, msg.ID, msg.Text, msg.State, msg.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
