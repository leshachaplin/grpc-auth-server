package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"strings"
)

func selectFindField(ctx context.Context, db *sqlx.DB, userField interface{}) (rows *sqlx.Rows, err error) {
	switch userField {
	case userField.(string):
		{
			stringField := userField.(string)
			if strings.ContainsAny(stringField, "@") {
				rows, err = db.QueryxContext(ctx, `SELECT
				username, email, password, confirmed, createdat, updatedat, id
				FROM "users" WHERE email = $1`, stringField)
			} else {
				rows, err = db.QueryxContext(ctx, `SELECT
				username, email, password, confirmed, createdat, updatedat, id
				FROM "users" WHERE username = $1`, stringField)
			}

		}
		//TODO: write any cases that you need
	}
	if err != nil {
		return nil, err
	}
	return rows, nil
}
