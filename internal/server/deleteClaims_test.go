package server

import (
	"context"
	"fmt"
	"github.com/leshachaplin/grpc-auth-server/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"testing"
)

func TestDeleteClaims(t *testing.T) {
	opts := grpc.WithInsecure()
	clientConnInterface, err := grpc.Dial("0.0.0.0:50051", opts)
	if err != nil {
		log.Error(err)
	}
	defer clientConnInterface.Close()
	client := protocol.NewAuthServiceClient(clientConnInterface)

	login := fmt.Sprintf("les")

	reqDeleteClaims := &protocol.DeleteClaimsRequest{
		Login:  login,
		Claims: map[string]string{"asdasd": "sdsd"},
	}

	responceDeleteClaims, err := client.DeleteClaims(context.Background(), reqDeleteClaims)
	if err == nil {
		fmt.Printf("delete claims%s\n", responceDeleteClaims)
	} else {
 		t.Errorf("deleting claims is failed, got:%s  , want:%s ", err, responceDeleteClaims )
	}
}
