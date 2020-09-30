package service

import (
	"context"
	"github.com/leshachaplin/config"
	"github.com/leshachaplin/grpc-auth-server/internal/auth"
	conf "github.com/leshachaplin/grpc-auth-server/internal/config"
	"github.com/leshachaplin/grpc-auth-server/internal/repository"
	log "github.com/sirupsen/logrus"
	"time"
)

func createTokens(ctx context.Context, user *repository.User, cfg *conf.Config,
	claims repository.Claims, refresh repository.Refresh, cfgAws *config.ConfigAws,
	password, username string) (token string, refreshToken string, err error) {
	secretKeyAuth, err := cfgAws.GetSecret(cfg.SecretKeyAuth)
	if err != nil {
		return "", "", err
	}
	secretKeyRefresh, err := cfgAws.GetSecret(cfg.SecretKeyRefresh)
	if err != nil {
		return "", "", err
	}
	if user != nil {
		err = repository.VerifyPassword(string(user.Password), password)
		if user.Username == username && err == nil {

			claims, err := claims.GetClaims(ctx, username)
			if err != nil {
				log.Errorf("error in get claims", err)
				return "", "", err
			}

			token, err = auth.CreatTokenAuth(username, claims, secretKeyAuth.ApiKey)
			if err != nil {
				log.Errorf("error in create new token", err)
				return "", "", err
			}

			refreshToken, err = auth.CreatTokenRefresh(username, secretKeyRefresh.ApiKey)
			if err != nil {
				log.Errorf("error in create new refresh token", err)
				return "", "", err
			}

			err = refresh.AddRefreshTokens(ctx, user.Username, refreshToken, time.Now().Add(time.Hour*24))
			if err != nil {
				log.Errorf("error in adding in table new refresh token", err)
				return "", "", err
			}
		}
	} else {
		log.Errorf("user not found", err)
		return "", "", err
	}
	return token, refreshToken, nil
}
