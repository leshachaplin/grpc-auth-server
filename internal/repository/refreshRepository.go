package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type RefreshTokenRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(database *sqlx.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db:database,
	}
}

func (r *RefreshTokenRepository) AddRefreshTokens(ctx context.Context, login string,  token string, time time.Time) error{
	var err error
	_, err = r.db.QueryContext(ctx, `INSERT into "refresh" (username, expiration, token) values ($1, $2, $3) `, login, time, token)
	return err
}

func (r *RefreshTokenRepository) GetRefreshToken(ctx context.Context, login string) (string, error) {
	rows, err := r.db.QueryxContext(ctx, `SELECT username, expiration, token FROM "refresh" WHERE username = $1 and expiration < $2`, login, time.Now())
	var refToken string
	for rows.Next() {
		err := rows.StructScan(&refToken)
		_ =err
	}
	return refToken, err
}