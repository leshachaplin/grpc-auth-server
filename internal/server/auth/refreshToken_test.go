package auth

import (
	"context"
	"fmt"
	"github.com/leshachaplin/grpc-auth-server/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"testing"
)

func TestRefreshToken(t *testing.T) {
	opts := grpc.WithInsecure()
	clientConnInterface, err := grpc.Dial("0.0.0.0:50051", opts)
	if err != nil {
		log.Error(err)
	}
	defer clientConnInterface.Close()
	client := protocol.NewAuthServiceClient(clientConnInterface)

	requestRefresh := &protocol.RefreshTokenRequest{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIxNjAxNDg3OTYyIiwibG9naW4iOiJsZXMifQ.6frwn3JiTZ94gjJ5214trKiKNTHc8UT7cUZjfkyLve0",
		TokenRefresh: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIxNjAxNTcyNTYyIiwic3ViIjoiMSIsInV1aWQiOiJsZXMifQ.54soODQk54fz-8aC8Bf2-fzZNx6NGT7o0KZZ6aV9-5k",
	}

	respRefresh, err := client.RefreshToken(context.Background(), requestRefresh)
	if respRefresh != nil  && err == nil{
		fmt.Printf("user refresh auth token:%s\n , user refresh token:%s\n ", respRefresh.GetToken(), respRefresh.GetRefreshToken())
	} else {
		t.Errorf("refresh is failed, got:%s  , want:%s\n %s ", err, respRefresh.GetToken() , respRefresh.GetRefreshToken() )
	}

}

