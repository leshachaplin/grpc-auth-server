package users

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/leshachaplin/grpc-auth-server/protocol"
	"google.golang.org/grpc"
	"testing"
)

func TestAddUser(t *testing.T) {
	opts := grpc.WithInsecure()
	clientConnInterface, err := grpc.Dial("0.0.0.0:50053", opts)
	if err != nil {
		log.Error(err)
	}
	defer clientConnInterface.Close()
	client := protocol.NewUserServiceClient(clientConnInterface)

	request := &protocol.CreateRequest{User: &protocol.User{
		Username:  "lesha24",
		Confirmed: false,
		Email:     "lesha@chaplin.com",
		Password:  "1234",
	}}

	responceCreate, err := client.CreateUser(context.Background(), request)
	if err == nil {
		fmt.Printf("add claims to users%s\n", responceCreate)
	} else {
		t.Errorf("add claims is failed, got:%s  , want:%s ", err, responceCreate )
	}
}
