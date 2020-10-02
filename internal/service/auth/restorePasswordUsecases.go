package auth

import (
	"context"
	"errors"
	"github.com/leshachaplin/grpc-auth-server/internal/repository"
	log "github.com/sirupsen/logrus"
	"time"
)

func checkIsValidUuid(uuid *repository.Restore, uuidRestore string) error {
	//str := strconv.FormatInt(uuid.Expiration, 10)
	//t, err := strconv.ParseInt(str, 10, 64)
	//if err != nil {
	//	return err
	//}

	//tm := time.Unix(t, 0)
	if time.Now().After(uuid.Expiration) {
		return errors.New("Time out")
	}

	if uuid.UuidRestore != uuidRestore {
		return errors.New("restore uuid not matched")
	}
	return nil
}

func updatePassword(ctx context.Context, users repository.UserRepository,
	username, newPassword string) error {
	user, err := users.FindUser(ctx, username)
	if err != nil {
		log.Errorf("User not found", err)
		return err
	}

	user.Password = repository.Hash(newPassword)

	err = users.Update(ctx, user)
	if err != nil {
		log.Errorf("error in updating", err)
		return err
	}
	return nil
}
