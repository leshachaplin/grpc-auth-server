package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Claim struct {
	Description string `db:"key"`
	Value       string `db:"value"`
}

type ClaimsRepository struct {
	db *sqlx.DB
}

func NewRepositoryOfClaims(database *sqlx.DB) *ClaimsRepository {
	return &ClaimsRepository{
		db: database,
	}
}

func (r *ClaimsRepository) GetClaims(ctx context.Context, login string) (map[string]string, error) {
	var err error
	result := make(map[string]string, 0)
	rows, err := r.db.QueryxContext(ctx, `SELECT (key, value) FROM "claim" WHERE userid = $1`, login)
	if err != nil {
		return nil, err
	}
	claim := Claim{}
	for rows.Next() {
		err := rows.StructScan(&claim)
		//return &user, err
		_ =err

		result[claim.Description] = claim.Value
	}
	return result, err
}

func (r *ClaimsRepository) IfExistClaim(ctx context.Context, key, login string) (bool, error) {
	rows, err := r.db.QueryxContext(ctx,`SELECT (key, value) FROM "claim" WHERE username = $1 and key = $2 `, login, key)
	if err != nil {
		return false, nil
	}
	for rows.Next() {
		return true, err
	}
	return false, nil
}

func (r *ClaimsRepository) AddClaims(ctx context.Context, claims map[string]string, login string) error {
	var err error
	for k, v := range claims {
		_, err = r.db.QueryContext(ctx, `INSERT into "claim" (key, value, username) values ($1, $2, $3) `, k, v, login)
	}
	return err
}

func (r *ClaimsRepository) DeleteClaims(ctx context.Context, claims map[string]string,  login string) error {
	var err error
	for k, v := range claims {
		_, err = r.db.QueryContext(ctx, `DELETE FROM "claim" WHERE username = $1 and key = $2 and value = $3`,login , k, v)
	}
	return err
}