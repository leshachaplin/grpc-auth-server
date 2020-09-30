package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type Confirmation struct {
	UuidConfirmation string   `db:"uuidconfirmation"`
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


func (r *ConfirmationRepository) Get(ctx context.Context, login string) (string , error) {
	rows, err := r.db.QueryxContext(ctx, `SELECT userid, uuidconfirmation, expiration FROM "confirmation" WHERE userid = $1 and  expiration < $2`, login, time.Now())
	if err != nil {
		return "", err
	}
	var confirm string
	for rows.Next() {
		err := rows.StructScan(&confirm)
		_ = err
	}
	return confirm, err
}
