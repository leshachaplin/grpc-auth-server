package auth

import (
	"context"
	"fmt"
	"github.com/leshachaplin/grpc-auth-server/protocol"
	"google.golang.org/grpc"
	"testing"
)

func TestSignUp(t *testing.T) {
	opts := grpc.WithInsecure()
	clientConnInterface, err := grpc.Dial("0.0.0.0:50051", opts)
	if err != nil {
		t.Error(err)
	}
	defer clientConnInterface.Close()
	client := protocol.NewAuthServiceClient(clientConnInterface)

	requestSignUp := &protocol.SignUpRequest{
		Email:    "lesha.chaplin66@gmail.com",
		Login:    "les",
		Password: "chaplin",
	}

	res, err := client.SignUp(context.Background(), requestSignUp)
	if res != nil && err == nil {
		fmt.Println("create new users adn create token")
	} else {
		t.Errorf("Sign up is failed, got:%s   ", err)
	}
}
