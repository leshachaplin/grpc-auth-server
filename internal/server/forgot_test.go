package server

import (
	"context"
	"fmt"
	"github.com/leshachaplin/grpc-auth-server/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"testing"
)

func TestForgot (t *testing.T) {
	opts := grpc.WithInsecure()
	clientConnInterface, err := grpc.Dial("0.0.0.0:50051", opts)
	if err != nil {
		log.Error(err)
	}
	defer clientConnInterface.Close()
	client := protocol.NewAuthServiceClient(clientConnInterface)

	login := fmt.Sprintf("les")
	mail := fmt.Sprintf("lesha.chaplin66@gmail.com")

	requestForgot := &protocol.ForgotPasswordRequest{
		Login: login,
		Email: mail,
	}

	res, err := client.ForgotPassword(context.Background(), requestForgot)
	if res != nil && err == nil {
		fmt.Println("create restore token")
	} else {
		t.Errorf("Forgot is failed, got:%s   ", err)
	}
}
