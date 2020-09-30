package server

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
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODE5NDIzNTMsImxvZ2luIjoibGVzaGEyIn0.t_EfPOQ4ap5F79FM_L2jsyrVtdu_1_ikEC0VbPY4wDo",
	}

	respRefresh, err := client.RefreshToken(context.Background(), requestRefresh)
	if respRefresh != nil  && err == nil{
		fmt.Printf("user refresh auth token:%s\n , user refresh token:%s\n ", respRefresh.GetToken(), respRefresh.GetRefreshToken())
	} else {
		t.Errorf("refresh is failed, got:%s  , want:%s\n %s ", err, respRefresh.GetToken() , respRefresh.GetRefreshToken() )
	}

}

