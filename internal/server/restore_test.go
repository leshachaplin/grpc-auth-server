package server

import (
	"context"
	"fmt"
	"github.com/leshachaplin/grpc-auth-server/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"testing"
)

func TestRestore(t *testing.T) {
	opts := grpc.WithInsecure()
	clientConnInterface, err := grpc.Dial("0.0.0.0:50051", opts)
	if err != nil {
		log.Error(err)
	}
	defer clientConnInterface.Close()
	client := protocol.NewAuthServiceClient(clientConnInterface)

	login := fmt.Sprintf("les")
	newPassword := fmt.Sprintf("penisPalladin")
	uuid := fmt.Sprintf("c4685855-f8c3-410c-879a-772901ace9f9")

	requestRestore := &protocol.RestoreRequest{
		Token:       uuid,
		Login:       login,
		NewPassword: newPassword,
	}

	res, err := client.Restore(context.Background(), requestRestore)
	if res != nil && err == nil {
		fmt.Println("Restore password")
	} else {
		t.Errorf("Restored is failed, got:%s   ", err)
	}

}
