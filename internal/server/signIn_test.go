package server

import (
	"context"
	"fmt"
	"github.com/leshachaplin/grpc-auth-server/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"testing"
)

func TestSignIn(t *testing.T) {
	opts := grpc.WithInsecure()
	clientConnInterface, err := grpc.Dial("0.0.0.0:50051", opts)
	if err != nil {
		log.Error(err)
	}
	defer clientConnInterface.Close()
	client := protocol.NewAuthServiceClient(clientConnInterface)

	requestSignIn := &protocol.SignInRequest{
		Login:    "les",
		Password: "chaplin",
	}

	respSignIn, err := client.SignIn(context.Background(), requestSignIn)
	log.Println( respSignIn.Token)
	log.Println( respSignIn.RefreshToken)
	if respSignIn != nil  && err == nil{
		fmt.Printf("user autorizate token:%s\n", respSignIn.GetToken())
	} else {
		t.Errorf("sign in is failed, got:%s  , want:%s ", err, respSignIn.GetToken() )
	}
}
