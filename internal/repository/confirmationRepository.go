package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type Confirmation struct {
	UuidConfirmation string `db:"uuidconfirmation"`
}

type ConfirmationRepository struct {
	db *sqlx.DB
}

func NewConfirmationRepository(database sqlx.DB) *ConfirmationRepository {
	return &ConfirmationRepository{
		db: &database,
	}
}

func (r *ConfirmationRepository) Delete(ctx context.Context, login string) error {
	_, err := r.db.QueryContext(ctx, `delete from "confirmation" where username = $1`, login)
	return err
}

func (r *ConfirmationRepository) Create(ctx context.Context, login string, uuid string, exp time.Time) error {
	_, err := r.db.QueryContext(ctx, `INSERT into "confirmation" (username, uuidconfirmation, expiration) values ($1, $2, $3)`, login, uuid, exp)
	return err
}

func (r *ConfirmationRepository) Get(ctx context.Context, login string) (*Confirmation, error) {
	rows, err := r.db.QueryxContext(ctx, `SELECT uuidconfirmation FROM "confirmation" WHERE username = $1 and  expiration > $2`, login, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	confirm := Confirmation{}
	for rows.Next() {
		err := rows.StructScan(&confirm)
		_ = err
	}
	return &confirm, err
}
