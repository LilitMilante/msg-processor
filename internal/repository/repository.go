package repository

import (
	"context"

	"msg-processor/internal/entity"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateMsg(ctx context.Context, msg entity.Msg) error {
	q := "INSERT INTO messages (id, text, state, created_at) VALUES ($1, $2, $3, $4)"

	_, err := r.db.Exec(ctx, q, msg.ID, msg.Text, msg.State, msg.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UnprocessedMsgs(ctx context.Context) (msgs []entity.Msg, err error) {
	q := "SELECT id, text, state, created_at FROM messages WHERE state = $1"

	rows, err := r.db.Query(ctx, q, entity.StateUnprocessed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg entity.Msg

		err = rows.Scan(&msg.ID, &msg.Text, &msg.State, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (r *Repository) UpdateMsgsState(ctx context.Context, state entity.State, msgs ...entity.Msg) error {
	q := "UPDATE messages SET state = $1 WHERE id = any ($2)"

	ids := make([]uuid.UUID, 0, len(msgs))
	for _, msg := range msgs {
		ids = append(ids, msg.ID)
	}

	_, err := r.db.Exec(ctx, q, state, ids)
	if err != nil {
		return err
	}

	return nil
}
