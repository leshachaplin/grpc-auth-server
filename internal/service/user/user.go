package user

import (
	"context"
	"github.com/leshachaplin/grpc-auth-server/internal/repository"
	log "github.com/sirupsen/logrus"
)

type UserService struct {
	users     repository.UserRepository
}

func New(usersRepository repository.UsersRepository) *UserService {
	return &UserService{
		users: &usersRepository,
	}
}

func (u *UserService) CreateUser(ctx context.Context, user *repository.User) error {
	err := u.users.Create(ctx, user)
	if err != nil {
		log.Error("error in CreateUser")
		return err
	}
	return nil
}

func (u *UserService) Delete(ctx context.Context, login string) error {
	err := u.users.Delete(ctx, login)
	if err != nil {
		log.Error("error in Delete")
		return err
	}
	return nil
}

func (u *UserService) Update(ctx context.Context, user *repository.User) error {
	err := u.users.Update(ctx, user)
	if err != nil {
		log.Error("error in Update")
		return err
	}
	return nil
}

func (u *UserService) Find(ctx context.Context, userField interface{}) (*repository.User, error) {
	user, err := u.users.FindUser(ctx, userField)
	if err != nil {
		log.Error("error in Find")
		return nil, err
	}
	return user, nil
}

