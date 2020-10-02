package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	emailProto "github.com/leshachaplin/emailSender/protocol"
	"time"
)

type EmailRepository struct {
	db *sqlx.DB
}

func NewEmailRepository(database sqlx.DB) *EmailRepository {
	return &EmailRepository{
		db: &database,
	}
}

//TODO: auto update time
func (r *EmailRepository) Create(ctx context.Context, mail *emailProto.SMTPEmail) error {
	_, err := r.db.QueryContext(ctx, `INSERT into "email" (fromemail, toemail, createdat) values ($1, $2, $3)`,
		mail.From, mail.To, time.Now().UTC())
	return err
}
