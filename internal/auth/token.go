package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/leshachaplin/grpc-auth-server/internal/repository"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func CreatTokenAuth(login string, claim map[string]string, secretKey string) (string,  error) {
	var err error
	claims := jwt.MapClaims{}
	claims["login"] = login
	repository.MergeMaps(repository.IntefaceToString(claims), claim)
	claims["exp"] = strconv.Itoa(int(time.Now().Add(time.Minute * 30).Unix()))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Errorf("can't create a token", err)
		return "", err
	}
	return t, nil
}

func CreatTokenRefresh(uuuid string, secretKey string) (string,  error) {
	var err error
	rtClaims := jwt.MapClaims{}
	rtClaims["sub"] = "1"
	rtClaims["uuid"] = uuuid
	rtClaims["exp"] = strconv.Itoa(int(time.Now().Add(time.Hour * 24).Unix()))
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rt, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		log.Errorf("can't create a token")
		return "", err
	}
	return rt, nil
}

func GetTokenAuth(tokensString, secretKeyAuth string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokensString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secretKeyAuth), nil
	})

	return token, err
}

func GetTokenRefresh(tokensString, secretKeyRefresh string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokensString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secretKeyRefresh), nil
	})

	return token, err
}

func GetExpirationTimeToRefreshToken(token, secretKeyRefresh string) (time.Time, error) {
	rt, err := GetTokenRefresh(token, secretKeyRefresh)
	if err != nil {
		return time.Unix(0,0), err
	}
	rtClaims := repository.IntefaceToString(rt.Claims.(jwt.MapClaims))
	t, err := strconv.ParseInt(rtClaims["exp"], 10, 64)
	if err != nil {
		return time.Unix(0,0), err
	}
	tm := time.Unix(t, 0)
	return tm, err
}


