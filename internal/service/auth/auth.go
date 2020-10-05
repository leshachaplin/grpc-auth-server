package auth

import (
	"context"
	"errors"
	guuid "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	emailProto "github.com/leshachaplin/emailSender/protocol"
	conf "github.com/leshachaplin/grpc-auth-server/internal/config"
	"github.com/leshachaplin/grpc-auth-server/internal/repository"
	log "github.com/sirupsen/logrus"
	"time"
)

type AuthenticationService struct {
	users     repository.UserRepository
	refresh   repository.Refresh
	claims    repository.Claims
	email     repository.Email
	cfg       *conf.Config
	restore   repository.RestorePassword
	confirm   repository.Confirm
}

//TODO: change databases
//TODO: automatizate createdAt & updatedAt
//TODO: Restore - delete after usage
//TODO: create uuids to AT & RT and save them

func New(usersRepository repository.UsersRepository, claims repository.ClaimsRepository,
	mail repository.EmailRepository, config conf.Config, tokenRepository repository.RefreshTokenRepository,
	restoreR repository.RestoreRepository, confirmR repository.ConfirmationRepository) *AuthenticationService {
	return &AuthenticationService{
		users:     &usersRepository,
		refresh:   &tokenRepository,
		claims:    &claims,
		email:     &mail,
		cfg:       &config,
		restore:   &restoreR,
		confirm:   &confirmR,
	}
}

func (s *AuthenticationService) SignIn(ctx context.Context, login string, password string) (string, string, error) {
	user, err := s.users.FindUser(ctx, login)
	if err != nil {
		log.Errorf("error in find users", err)
		return "", "", err
	}

	token, refreshToken, err := createTokens(ctx, user, s.cfg, s.claims,
		s.refresh, password, login)

	return token, refreshToken, err
}

func (s *AuthenticationService) SignUp(ctx context.Context, mail string, username string, password string) error {
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

	err = checkIsConfirmEmailSend(ctx, username, s.cfg, s.email, s.confirm)
	if err != nil {
		return err
	}

	return err
}

func (s *AuthenticationService) ResendEmail(ctx context.Context, login, mail string) error {
	err := s.confirm.Delete(ctx, login)
	if err != nil {
		log.Errorf("error in delete uuid confirmation", err)
		return err
	}

	err = checkIsConfirmEmailSend(ctx, login, s.cfg, s.email, s.confirm)

	return err
}

func (s *AuthenticationService) ConfirmEmail(ctx context.Context, uuidConf, login string) error {
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

func (s *AuthenticationService) ForgotPassword(ctx context.Context, login string) error {

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

	c := ctx.(echo.Context)
	req := &emailProto.SendRequest{}

	if err := c.Bind(&req); err != nil {
		log.Errorf("echo.Context binding Error ForgotPassword %s", err)
		return err
	}

	err = s.email.Create(ctx, req.Email)
	if err != nil {
		log.Errorf("Error in create new email: EmailRepository %s", err)
		return err
	}

	_, err = s.cfg.SMTPClient.Send(ctx, req)

	return err
}

func (s *AuthenticationService) RestorePassword(ctx context.Context, uuidRestore, login, newPassword string) error {
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

func (s *AuthenticationService) RefreshToken(ctx context.Context, tokenReqAuth,
	tokenReqRefresh string) (string, string, error) {
	newToken, newRefToken, err := refreshToken(ctx, s.users, s.refresh, s.claims,
		s.cfg, tokenReqAuth, tokenReqRefresh)
	return newToken, newRefToken, err
}

func (s *AuthenticationService) ChangePassword(ctx context.Context, mail, oldPassword, newPassword string) error {
	user, err := s.users.FindUser(ctx, mail)
	if err != nil {
		log.Errorf("error in find users by email", err)
		return err
	}

	err = repository.VerifyPassword(string(user.Password), oldPassword)
	if err != nil {
		return errors.New("password invalid")
	}

	user.Password = repository.Hash(newPassword)

	err = s.users.Update(ctx, user)

	return err
}

func (s *AuthenticationService) DeleteClaims(ctx context.Context, claims map[string]string, login string) error {
	return s.claims.DeleteClaims(ctx, claims, login)
}

func (s *AuthenticationService) AddClaims(claims map[string]string, ctx context.Context, login string) error {
	return s.claims.AddClaims(ctx, claims, login)
}
