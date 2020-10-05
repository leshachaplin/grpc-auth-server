package users
//
//import (
//	"context"
//	"fmt"
//	"github.com/labstack/gommon/log"
//	"github.com/leshachaplin/grpc-auth-server/protocol"
//	"google.golang.org/grpc"
//	"google.golang.org/protobuf/types/known/anypb"
//	"testing"
//)
//
//func TestFindUser(t *testing.T) {
//	opts := grpc.WithInsecure()
//	clientConnInterface, err := grpc.Dial("0.0.0.0:50053", opts)
//	if err != nil {
//		log.Error(err)
//	}
//	defer clientConnInterface.Close()
//	client := protocol.NewUserServiceClient(clientConnInterface)
//
//	request := &protocol.FindRequest{UserField: anypb.Any{
//		TypeUrl: "string",
//		Value:   byte("lesha"),
//	} "lesha"}
//
//	responceCreate, err := client.CreateUser(context.Background(), request)
//	if err == nil {
//		fmt.Printf("add claims to users%s\n", responceCreate)
//	} else {
//		t.Errorf("add claims is failed, got:%s  , want:%s ", err, responceCreate )
//	}
//}
