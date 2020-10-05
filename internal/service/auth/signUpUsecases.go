package auth

import (
	"context"
	"errors"
	"fmt"
	guuid "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	emailProto "github.com/leshachaplin/emailSender/protocol"
	"github.com/leshachaplin/grpc-auth-server/internal/config"
	"github.com/leshachaplin/grpc-auth-server/internal/repository"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
	unicode "unicode"
)

/**
 * Minimum eight characters, maximum twenty characters, at least one uppercase letter,
 * one lowercase letter, one number and one special character
 */
//const passwordExp = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]`

func verifyPassword(password string) error {
	var uppercasePresent bool
	var lowercasePresent bool
	var numberPresent bool
	var specialCharPresent bool
	const minPassLength = 8
	const maxPassLength = 20
	var passLen int
	var errorString string

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			numberPresent = true
			passLen++
		case unicode.IsUpper(ch):
			uppercasePresent = true
			passLen++
		case unicode.IsLower(ch):
			lowercasePresent = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			specialCharPresent = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}
	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	if !lowercasePresent {
		appendError("lowercase letter missing")
	}
	if !uppercasePresent {
		appendError("uppercase letter missing")
	}
	if !numberPresent {
		appendError("atleast one numeric character required")
	}
	if !specialCharPresent {
		appendError("special character missing")
	}
	if !(minPassLength <= passLen && passLen <= maxPassLength) {
		appendError(fmt.Sprintf("password length must be between %d to %d characters long", minPassLength, maxPassLength))
	}

	if len(errorString) != 0 {
		return fmt.Errorf(errorString)
	}
	return nil
}

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

	if users.IfExistUser(ctx, username) {
		return errors.New("users with such username is exist, enter other username")
	}

	if users.IfExistUser(ctx, mail) {
		return errors.New("users with such email is exist, enter other email")
	}

	err := verifyPassword(password)
	if err != nil {
		return err
	}

	return nil
}

func checkIsConfirmEmailSend(ctx context.Context, username string, cfg *config.Config,
	mail repository.Email, confirm repository.Confirm) error {

	uuidConfirm := guuid.New().String()


	err := confirm.Create(ctx, username, uuidConfirm,  time.Now().Add(time.Minute*60*10).UTC())
	if err != nil {
		log.Errorf("error in confirm", err)
		return err
	}

	c := ctx.(echo.Context)
	req := &emailProto.SendRequest{}

	if err := c.Bind(&req); err != nil {
		log.Errorf("echo.Context binding Error ForgotPassword %s", err)
		return err
	}

	err = mail.Create(ctx, req.Email)
	if err != nil {
		log.Errorf("Error in create new email: EmailRepository %s", err)
		return err
	}

	_, err = cfg.SMTPClient.Send(ctx, req)

	return nil
}

func createUser(ctx context.Context, user *repository.User, users repository.UserRepository) error {
	err := users.Create(ctx, user)
	if err != nil {
		log.Errorf("error in create new users", err)
		return err
	}

	user.Confirmed = false

	err = users.Update(ctx, user)
	if err != nil {
		log.Errorf("error in update users", err)
		return err
	}

	return nil
}
