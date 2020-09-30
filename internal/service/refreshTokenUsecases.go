package service

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	guuid "github.com/google/uuid"
	"github.com/leshachaplin/config"
	"github.com/leshachaplin/grpc-auth-server/internal/auth"
	conf "github.com/leshachaplin/grpc-auth-server/internal/config"
	"github.com/leshachaplin/grpc-auth-server/internal/repository"
	log "github.com/sirupsen/logrus"
)

func getUuids(tokenAuth, tokenReqRefresh,
	secretKeyAuth, secretKeyRefresh string) (string, string, error) {
	token, err := auth.GetTokenAuth(tokenAuth, secretKeyAuth)
	if err != nil {
		log.Errorf("error in get token claims", err)
		return "", "", err
	}

	claims := repository.IntefaceToString(token.Claims.(jwt.MapClaims))

	tokenRefresh, err := auth.GetTokenRefresh(tokenReqRefresh, secretKeyRefresh)
	if err != nil {
		log.Errorf("error in get token claims refresh", err)
		return "", "", err
	}

	claimsR := repository.IntefaceToString(tokenRefresh.Claims.(jwt.MapClaims))

	return claims["login"], claimsR["uuid"], nil
}

func refreshToken(ctx context.Context, users repository.UserRepository,
	refresh repository.Refresh, claims repository.Claims, cfg *conf.Config, cfgAws *config.ConfigAws,
	tokenAuth, tokenReqRefresh string) (newToken string, newRefToken string, err error) {

	secretKeyAuth, err := cfgAws.GetSecret(cfg.SecretKeyAuth)
	if err != nil {
		return "", "", err
	}
	secretKeyRefresh, err := cfgAws.GetSecret(cfg.SecretKeyRefresh)
	if err != nil {
		return "", "", err
	}

	authUuid, refreshUuid, err := getUuids(tokenAuth, tokenReqRefresh, secretKeyAuth.ApiKey, secretKeyRefresh.ApiKey)
	if err != nil {
		return "", "", err
	}

	if authUuid == refreshUuid {

		user, err := users.FindUser(ctx, authUuid)
		if err != nil {
			log.Errorf("User not found when refresh token", err)
			return "", "", err
		}
		_, err = refresh.GetRefreshToken(ctx, authUuid)
		if err != nil {
			log.Errorf("Invalid session")
		}

		if ok := users.Delete(ctx, authUuid); ok != nil {
			log.Errorf("error in deleting users", err)
			return "", "", ok
		}

		err = users.Create(ctx, user)
		if err != nil {
			log.Errorf("error in create users", err)
			return "", "", err
		}

		err = users.Update(ctx, user, true)
		if err != nil {
			log.Errorf("error in update user", err)
			return "", "", err
		}

		newClaims, err := claims.GetClaims(ctx, authUuid)
		if err != nil {
			log.Errorf("error in get Claims in refresh token", err)
			return "", "", nil
		}
		uuid := guuid.New().String()
		newToken, err = auth.CreatTokenAuth(authUuid, newClaims, secretKeyAuth.ApiKey)
		newRefToken, err = auth.CreatTokenRefresh(uuid, secretKeyRefresh.ApiKey)
		if err != nil {
			log.Errorf("error in create new tokens", err)
			return "", "", err
		}

		t, err := auth.GetExpirationTimeToRefreshToken(newRefToken, secretKeyRefresh.ApiKey)
		if err != nil {
			log.Errorf("error in get expiration time", err)
			return "", "", err
		}
		err = refresh.AddRefreshTokens(ctx, authUuid, newRefToken, t)
	} else {
		return "", "", errors.New("not validate in refresh")
	}
	return newToken, newRefToken, nil
}
