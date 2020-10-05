package auth

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
	uuid := fmt.Sprintf("8520e916-f8d4-4341-a3de-8585d42c8c50")

	requestConfirm := &protocol.ConfirmRequest{
		Login:    login,
		UuidConfirm: uuid,
	}

	res, err := client.Confirm(context.Background(), requestConfirm)
	if res != nil && err == nil {
		fmt.Println("create new users adn create token")
	} else {
		t.Errorf("Confirmation is failed, got:%s   ", err)
	}

}