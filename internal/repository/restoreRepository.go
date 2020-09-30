package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type Restore struct {
	Login       string    `db:"username"`
	UuidRestore string    `db:"uuidrestore"`
	Expiration  time.Time `db:"exp"`
	Id          string    `db:"id"`
}

type RestoreRepository struct {
	db *sqlx.DB
}

func NewRestoreRepository(database sqlx.DB) *RestoreRepository {
	return &RestoreRepository{
		db: &database,
	}
}

func (r *RestoreRepository) Delete(ctx context.Context, login string) error {
	_, err := r.db.QueryContext(ctx, `delete from "restore" where username = $1`, login)
	return err
}

func (r *RestoreRepository) Create(ctx context.Context, login string, uuid string, exp time.Time) error {
	_, err := r.db.QueryContext(ctx, `INSERT into "restore" (username, uuidrestore, exp) values ($1, $2, $3)`, login, uuid, exp)
	return err
}

func (r *RestoreRepository) Get(ctx context.Context, login string) (*Restore, error) {
	rows, err := r.db.QueryxContext(ctx, `SELECT username, uuidrestore, exp FROM "restore" WHERE username = $1`, login)
	if err != nil {
		return nil, err
	}
	restore := Restore{}
	for rows.Next() {
		err := rows.StructScan(&restore)
		//return &user, err
		_ = err
	}
	return &restore, err
}
