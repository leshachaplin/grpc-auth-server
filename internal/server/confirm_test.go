package server

import (
	"context"
	"fmt"
	"github.com/leshachaplin/grpc-auth-server/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"testing"
)

func TestConfirm(t *testing.T) {
	opts := grpc.WithInsecure()
	clientConnInterface, err := grpc.Dial("0.0.0.0:50051", opts)
	if err != nil {
		log.Error(err)
	}
	defer clientConnInterface.Close()
	client := protocol.NewAuthServiceClient(clientConnInterface)

	login := fmt.Sprintf("les")
	uuid := fmt.Sprintf("a16e276e-3421-4d8a-9db7-cf1a3096f291")

	requestConfirm := &protocol.ConfirmRequest{
		Login:    login,
		UuidConfirm: uuid,
	}

	res, err := client.Confirm(context.Background(), requestConfirm)
	if res != nil && err == nil {
		fmt.Println("create new user adn create token")
	} else {
		t.Errorf("Confirmation is failed, got:%s   ", err)
	}

}