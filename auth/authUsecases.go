package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/leshachaplin/grpc-auth-server/internal/repository"
	"strconv"
	"time"
)

type JWT struct {
	authSecretKey    []byte
	refreshSecretKey []byte
}

func New(authSecretKey, refreshSecretKey string) *JWT {
	return &JWT{
		authSecretKey:    []byte(authSecretKey),
		refreshSecretKey: []byte(refreshSecretKey),
	}
}

func (j *JWT) ParseAuth(token string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return j.authSecretKey, nil
	})

	if err != nil {
		return nil, err
	}
	claims := t.Claims.(jwt.MapClaims)
	return claims, nil
}

func (j *JWT) ParseRefresh(token string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return j.refreshSecretKey, nil
	})

	if err != nil {
		return nil, err
	}
	claims := t.Claims.(jwt.MapClaims)
	return claims, nil
}

func (j *JWT) GetExpTimeToRefreshToken(token string) (time.Time, error) {
	rt, err := j.ParseRefresh(token)
	if err != nil {
		return time.Unix(0, 0), err
	}
	rtClaims := repository.IntefaceToString(rt)
	t, err := strconv.ParseInt(rtClaims["exp"], 10, 64)
	if err != nil {
		return time.Unix(0, 0), err
	}
	tm := time.Unix(t, 0)
	return tm, err
}
