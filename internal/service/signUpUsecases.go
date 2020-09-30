package service

import (
	"context"
	"errors"
	guuid "github.com/google/uuid"
	"github.com/leshachaplin/grpc-auth-server/internal/email"
	"github.com/leshachaplin/grpc-auth-server/internal/repository"
	log "github.com/sirupsen/logrus"
	"time"
)

/**
 * Minimum eight characters, maximum twenty characters, at least one uppercase letter,
 * one lowercase letter, one number and one special character
 */
const passwordExp = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]`

func checkIsCorrectData(ctx context.Context, users repository.UserRepository,
	mail, username, password string) error {
	if username == "" {
		return errors.New("username field is empty")
	}

	if password == "" {
		return errors.New("password field is empty")
	}

	if mail == "" {
		return errors.New("mail field is empty")
	}

	if users.IfExistUserByUsername(ctx, username) {
		return errors.New("user with such username is exist, enter other username")
	}

	if users.IfExistUserByEmail(ctx, mail) {
		return errors.New("user with such email is exist, enter other email")
	}

	//r, err := regexp.Compile(passwordExp)
	//if err != nil {
	//	log.Error("invalid regexp in signUp Usecases")
	//	return err
	//}

	//if !r.Match([]byte(password)) {
	//	log.Error("invalid password in signUp Usecases, password should contain" +
	//		"minimum eight characters, maximum twenty characters, at least one uppercase" +
	//		" letter, one lowercase letter, one number and one special character  ")
	//	return errors.New("invalid password type")
	//}

	return nil
}

func checkIsConfirmEmailSend(ctx context.Context, username string,
	confirm repository.Confirm, emailSend email.EmailSender) error {
	uuidConfirm := guuid.New().String()

	confirmTemplate := email.PasswordTemplate{Token: uuidConfirm}

	err := confirm.Create(ctx, username, uuidConfirm,  time.Now().Add(time.Minute*60*10).UTC())

	if err != nil {
		log.Errorf("error in confirm", err)
		return err
	}

	err = emailSend.Send(username, confirmTemplate)
	if err != nil {
		log.Errorf("error in sending email", err)
		return err
	}

	return nil
}

func createUser(ctx context.Context, user *repository.User, users repository.UserRepository) error {
	err := users.Create(ctx, user)
	if err != nil {
		log.Errorf("error in create new user", err)
		return err
	}

	err = users.Update(ctx, user, false)
	if err != nil {
		log.Errorf("error in update user", err)
		return err
	}

	return nil
}
