package auth

import (
	"context"
	"github.com/leshachaplin/grpc-auth-server/internal/service/auth"
	"github.com/leshachaplin/grpc-auth-server/protocol"
	log "github.com/sirupsen/logrus"
)

type Server struct{
	Rpc auth.AuthenticationService
}

func (s *Server) AddClaims(ctx context.Context, req *protocol.AddClaimsRequest) (*protocol.EmptyResponse, error) {
	err := s.Rpc.AddClaims(req.Claims, ctx, req.Login)
	if err != nil {
		log.Errorf("erron in adding claims: %s", err)
		return nil, err
	}
	return &protocol.EmptyResponse{}, nil
}


func (s *Server) Delete(ctx context.Context, req *protocol.DeleteRequest) (*protocol.EmptyResponse, error) {
	err :=  s.Rpc.Delete(ctx, req.Login)
	if err != nil {
		log.Errorf("erron in deleting users: %s ", err)
		return nil, err
	}
	return &protocol.EmptyResponse{}, nil
}

func (s *Server) DeleteClaims(ctx context.Context ,req *protocol.DeleteClaimsRequest) (*protocol.EmptyResponse, error) {
	err := s.Rpc.DeleteClaims(ctx,req.Claims, req.Login)
	if err != nil {
		log.Errorf("erron in deleting claims: %s", err)
		return nil, err
	}
	return &protocol.EmptyResponse{}, nil
}

func (s *Server) SignIn(ctx context.Context, req *protocol.SignInRequest) (*protocol.SignInResponse, error) {
	token, refresh, err := s.Rpc.SignIn(ctx, req.Login, req.Password)
	if err != nil {
		log.Errorf("error in sign in: %s", err)
		return nil, err
	}

	responce := &protocol.SignInResponse{
		Token:        token,
		RefreshToken: refresh,
	}
	return responce, nil
}


func (s *Server) SignUp(ctx context.Context, req *protocol.SignUpRequest) (*protocol.EmptyResponse, error) {
	err := s.Rpc.SignUp(ctx, req.Email, req.Login, req.Password)
	if err != nil {
		log.Errorf("error in sign up: %s", err)
		return nil, err
	}

	return &protocol.EmptyResponse{}, nil
}

func (s *Server) RefreshToken(ctx context.Context, req *protocol.RefreshTokenRequest) (*protocol.RefreshTokenResponse, error) {
	token, refToken, err := s.Rpc.RefreshToken(ctx , req.Token, req.TokenRefresh)
	if err != nil {
		log.Errorf("error in refresh token: %s", err)
		return nil, err
	}

	response := &protocol.RefreshTokenResponse{
		Token:        token,
		RefreshToken: refToken,
	}
	return response, nil
}

func (s *Server) Confirm(ctx context.Context, req *protocol.ConfirmRequest) (*protocol.EmptyResponse, error) {
	err := s.Rpc.ConfirmEmail(ctx, req.UuidConfirm, req.Login)
	if err != nil {
		log.Errorf("Email not confirmed: %s", err)
		return nil, err
	}

	return &protocol.EmptyResponse{}, nil
}


func (s *Server) Restore(ctx context.Context, req *protocol.RestoreRequest) (*protocol.EmptyResponse, error) {
	err := s.Rpc.RestorePassword(ctx, req.Token, req.Login, req.NewPassword)
	if err != nil {
		log.Errorf("password no restore: %s", err)
		return nil, err
	}

	return &protocol.EmptyResponse{}, err
}

func (s *Server) ChangePassword(ctx context.Context, req *protocol.ChangePasswordRequest) (*protocol.EmptyResponse, error) {
	err := s.Rpc.ChangePassword(ctx, req.Email, req.OldPassword, req.NewPassword)
	if err != nil {
		log.Errorf("password no changed: %s", err)
		return nil, err
	}

	return &protocol.EmptyResponse{}, err
}

func (s *Server) ForgotPassword(ctx context.Context, req *protocol.ForgotPasswordRequest) (*protocol.EmptyResponse, error) {
	err := s.Rpc.ForgotPassword(ctx, req.Login)
	if err != nil {
		log.Errorf("error in forgot password method: %s", err)
		return nil, err
	}

	return &protocol.EmptyResponse{}, err
}









