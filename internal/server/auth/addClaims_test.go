package auth

import (
	"context"
	"fmt"
	"github.com/leshachaplin/grpc-auth-server/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"testing"
)

func TestAddClaims(t *testing.T) {
	opts := grpc.WithInsecure()
	clientConnInterface, err := grpc.Dial("0.0.0.0:50051", opts)
	if err != nil {
		log.Error(err)
	}
	defer clientConnInterface.Close()
	client := protocol.NewAuthServiceClient(clientConnInterface)

	login := fmt.Sprintf("les")

	requestAddClaims := &protocol.AddClaimsRequest{
		Login: login,
		Claims: map[string]string{
			"role":"admin",
			"asdasd":"sdsd",
		},
	}

	responceAddClaims, err := client.AddClaims(context.Background(), requestAddClaims)
	if err == nil {
		fmt.Printf("add claims to user%s\n", login)
	} else {
		t.Errorf("add claims is failed, got:%s  , want:%s ", err, responceAddClaims )
	}
}

