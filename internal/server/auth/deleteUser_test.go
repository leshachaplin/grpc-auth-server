package auth

import (
	"context"
	"fmt"
	"github.com/leshachaplin/grpc-auth-server/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"testing"
)

func TestDelete(t *testing.T) {
	opts := grpc.WithInsecure()
	clientConnInterface, err := grpc.Dial("0.0.0.0:50051", opts)
	if err != nil {
		log.Error(err)
	}
	defer clientConnInterface.Close()
	client := protocol.NewAuthServiceClient(clientConnInterface)

	login := fmt.Sprintf("les")

	requestDelete := &protocol.DeleteRequest{Login: login}

	responseDelete, err := client.Delete(context.Background(), requestDelete)
	if err == nil {
		fmt.Printf("delete users:%s  %s\n", requestDelete.Login, responseDelete)
	} else {
		t.Errorf("deleting users is failed, got:%s  , want:%s ", err, responseDelete )
	}
}
