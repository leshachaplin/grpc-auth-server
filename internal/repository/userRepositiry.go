package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int       `db:"id"`
	Email     string    `db:"email"`
	Username  string    `db:"username"`
	Confirmed bool      `db:"confirmed"`
	Password  []byte    `db:"password"`
	CreatedAt time.Time `db:"createdat"`
	UpdatedAt time.Time `db:"updatedat"`
}

type UsersRepository struct {
	db *sqlx.DB
}

func NewUserRepository(database sqlx.DB) *UsersRepository {
	return &UsersRepository{
		db: &database,
	}
}

func (r *UsersRepository) FindUser(ctx context.Context, username string) (*User, error) {
	var err error
	rows, err := r.db.QueryxContext(ctx, `SELECT
	username, email, password, confirmed, createdat, updatedat, id
	FROM "user" WHERE username = $1`, username)
	if err != nil {
		return nil, err
	}
	user := User{}
	for rows.Next() {
		err := rows.StructScan(&user)
		//return &user, err
		_ = err
	}
	return &user, err
}

func (r *UsersRepository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var err error
	rows, err := r.db.QueryxContext(ctx, `SELECT
	username, email, password, confirmed, createdat, updatedat, id
	FROM "user" WHERE email = $1`, email)
	if err != nil {
		return nil, err
	}
	user := User{}
	for rows.Next() {
		err := rows.StructScan(&user)
		//return &user, err
		_ = err
	}
	return &user, err
}

func (r *UsersRepository) IfExistUserByUsername(ctx context.Context, login string) bool {
	var err error
	rows, err := r.db.QueryxContext(ctx, `SELECT username, confirmed, password FROM "user" WHERE (username) = $1`, login)
	if err != nil {
		log.Errorf("Queue error", err)
		return false
	}
	for rows.Next() {
		return true
	}
	return false
}

func (r *UsersRepository) IfExistUserByEmail(ctx context.Context, email string) bool {
	var err error
	rows, err := r.db.QueryxContext(ctx, `SELECT email, confirmed, password FROM "user" WHERE (email) = $1`, email)
	if err != nil {
		log.Errorf("Queue error", err)
		return false
	}
	for rows.Next() {
		return true
	}
	return false
}

func (r *UsersRepository) Delete(ctx context.Context, login string) error {
	_, err := r.db.QueryContext(ctx, `delete from "user" where username = $1`, login)
	return err
}

func (r *UsersRepository) Create(ctx context.Context, user *User) error {
	_, err := r.db.QueryContext(ctx, `INSERT into "user" 
    (username, email, confirmed, Password, createdat, updatedat)
     values ($1, $2, $3, $4, $5, $6)`, user.Username, user.Email,
		false, user.Password, time.Now().UTC(), time.Now().UTC())
	return err
}

func (r *UsersRepository) Update(ctx context.Context, user *User, state bool) error {
	_, err := r.db.QueryContext(ctx, `UPDATE "user" set confirmed = $1 where username = $2`, state, user.Username)
	return err
}

func (r *UsersRepository) UpdatePassword(ctx context.Context, user *User, newPassword string) error {
	_, err := r.db.QueryContext(ctx, `UPDATE "user" set password = $1 where username = $2`, newPassword, user.Username)
	return err
}

func Hash(password string) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
	}

	return hashedPassword
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func MergeMaps(a map[string]string, b map[string]string) {
	for k, v := range b {
		a[k] = v
	}
}

func DeleteMap(a map[string]string, b map[string]string) {
	for k, _ := range b {
		delete(a, k)
	}
}

func IntefaceToString(claims map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range claims {
		result[k] = v.(string)
	}
	return result
}
