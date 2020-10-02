package auth

import (
	"context"
	"errors"
	"github.com/leshachaplin/grpc-auth-server/internal/repository"
	log "github.com/sirupsen/logrus"
)

func confirmUserEmail(ctx context.Context, users repository.UserRepository,
	login, uuidConfirmation, uuidConf string) error {
	user, err := users.FindUser(ctx, login)
	if err != nil {
		log.Errorf("User not found", err)
		return err
	}

	if uuidConfirmation != uuidConf {
		return errors.New("confirm uuid not matched")
	}

	user.Confirmed = true

	err = users.Update(ctx, user)
	if err != nil {
		log.Errorf("error in update user", err)
		return err
	}
	return nil
}
