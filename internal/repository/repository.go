package repository

import (
	"context"
	emailProto "github.com/leshachaplin/emailSender/protocol"
	"time"
)

type Claims interface {
	GetClaims(ctx context.Context, login string) (map[string]string, error)
	IfExistClaim(ctx context.Context, key, login string) (bool, error)
	AddClaims(ctx context.Context, claims map[string]string, login string) error
	DeleteClaims(ctx context.Context, claims map[string]string, login string) error
}

type Refresh interface {
	AddRefreshTokens(ctx context.Context, login string, token string, exp time.Time) error
	GetRefreshToken(ctx context.Context, login string) (string, error)
}

type UserRepository interface {
	FindUser(ctx context.Context, userField interface{}) (*User, error)
	IfExistUser(ctx context.Context, userField interface{}) bool
	Delete(ctx context.Context, login string) error
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

type RestorePassword interface {
	Delete(ctx context.Context, login string) error
	Create(ctx context.Context, login string, uuid string, exp time.Time) error
	Get(ctx context.Context, uid string) (*Restore, error)
}

type Email interface {
	Create(ctx context.Context, mail *emailProto.SMTPEmail) error
}

type Confirm interface {
	Delete(ctx context.Context, login string) error
	Create(ctx context.Context, user string, uuid string, exp time.Time) error
	Get(ctx context.Context, login string) (*Confirmation, error)
}
