package service

import (
	"context"
	"errors"
	guuid "github.com/google/uuid"
	"github.com/leshachaplin/config"
	conf "github.com/leshachaplin/grpc-auth-server/internal/config"
	"github.com/leshachaplin/grpc-auth-server/internal/email"
	"github.com/leshachaplin/grpc-auth-server/internal/repository"
	log "github.com/sirupsen/logrus"
	"time"
)

type UserService struct {
	users     repository.UserRepository
	refresh   repository.Refresh
	claims    repository.Claims
	cfgAws    *config.ConfigAws
	cfg       *conf.Config
	restore   repository.RestorePassword
	confirm   repository.Confirm
	emailSend email.EmailSender
}

//TODO: Restore - delete after usage
func New(usersRepository repository.UsersRepository, claims repository.RepositoryOfClaims,
	configAws config.ConfigAws, config conf.Config, tokenRepository repository.RefreshTokenRepository,
	restoreR repository.RestoreRepository, confirmR repository.ConfirmationRepository,
	smtp email.SMTPEmail) *UserService {
	return &UserService{
		users:     &usersRepository,
		refresh:   &tokenRepository,
		claims:    &claims,
		cfgAws:    &configAws,
		cfg:       &config,
		restore:   &restoreR,
		confirm:   &confirmR,
		emailSend: &smtp,
	}
}

func (s *UserService) SignIn(ctx context.Context, login string, password string) (string, string, error) {
	user, err := s.users.FindUser(ctx, login)
	if err != nil {
		log.Errorf("error in find user", err)
		return "", "", err
	}

	token, refreshToken, err := createTokens(ctx, user, s.cfg, s.claims,
		s.refresh, s.cfgAws, password, login)

	return token, refreshToken, err
}

func (s *UserService) SignUp(ctx context.Context, mail string, username string, password string) error {
	var err error

	err = checkIsCorrectData(ctx, s.users, mail, username, password)
	if err != nil {
		return err
	}

	user := &repository.User{
		Email:    mail,
		Username: username,
		Password: repository.Hash(password),
	}

	err = createUser(ctx, user, s.users)
	if err != nil {
		return err
	}

	err = checkIsConfirmEmailSend(ctx, username, mail, s.confirm, s.emailSend)
	if err != nil {
		return err
	}

	return err
}

func (s *UserService) ResendEmail(ctx context.Context, login, mail string) error {
	err := s.confirm.Delete(ctx, login)
	if err != nil {
		log.Errorf("error in delete uuid confirmation", err)
		return err
	}

	err = checkIsConfirmEmailSend(ctx, login, mail, s.confirm, s.emailSend)

	return err
}

func (s *UserService) ConfirmEmail(ctx context.Context, uuidConf, login string) error {
	uuidConfirmation, err := s.confirm.Get(ctx, login)
	if err != nil {
		log.Errorf("error in confirm email repo")
		return err
	}

	err = confirmUserEmail(ctx, s.users, login, uuidConfirmation.UuidConfirmation, uuidConf)
	if err != nil {
		return err
	}

	err = s.confirm.Delete(ctx, login)
	if err != nil {
		log.Errorf("error in delete uuid confirmation", err)
		return err
	}

	return nil
}

func (s *UserService) ForgotPassword(ctx context.Context, login, mail string) error {
	uuidRestore := guuid.New().String()

	err := s.restore.Delete(ctx, login)
	if err != nil {
		return err
	}

	err = s.restore.Create(ctx, login, uuidRestore, time.Now().Add(time.Minute*60))
	if err != nil {
		log.Errorf("error in create restore uuid", err)
		return err
	}

	restoreTemplate := email.PasswordTemplate{Token: uuidRestore}

	err = s.emailSend.Send(mail, restoreTemplate)
	if err != nil {
		log.Errorf("error in send email", err)
		return err
	}

	return err
}

func (s *UserService) RestorePassword(ctx context.Context, uuidRestore, login, newPassword string) error {
	uuidRest, err := s.restore.Get(ctx, login)
	if err != nil {
		log.Errorf("can't find uuid restore", err)
		return err
	}

	err = checkIsValidUuid(uuidRest, uuidRestore)
	if err != nil {
		return err
	}

	err = updatePassword(ctx, s.users, login, newPassword)
	if err != nil {
		return err
	}

	return err
}

func (s *UserService) RefreshToken(ctx context.Context, tokenReqAuth,
	tokenReqRefresh string) (string, string, error) {
	newToken, newRefToken, err := refreshToken(ctx, s.users, s.refresh, s.claims,
	s.cfg, s.cfgAws, tokenReqAuth, tokenReqRefresh)
	return newToken, newRefToken, err
}

func (s *UserService) ChangePassword(ctx context.Context, mail, oldPassword, newPassword string) error {
	user, err := s.users.FindUserByEmail(ctx, mail)
	if err != nil {
		log.Errorf("error in find user by email", err)
		return err
	}

	err = repository.VerifyPassword(string(user.Password), oldPassword)
	if err != nil {
		return errors.New("password invalid")
	}

	err = s.users.UpdatePassword(ctx, user, newPassword)

	return err
}

func (s *UserService) DeleteClaims(ctx context.Context, claims map[string]string, login string) error {
	return s.claims.DeleteClaims(ctx, claims, login)
}

func (s *UserService) AddClaims(claims map[string]string, ctx context.Context, login string) error {
	return s.claims.AddClaims(ctx, claims, login)
}

func (s *UserService) Delete(ctx context.Context, login string) error {
	return s.users.Delete(ctx, login)
}
