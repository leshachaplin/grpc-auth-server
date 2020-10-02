package user

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/leshachaplin/grpc-auth-server/internal/repository"
	"github.com/leshachaplin/grpc-auth-server/internal/service/user"
	"github.com/leshachaplin/grpc-auth-server/protocol"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Rpc user.UserService
}

func (s *Server) CreateUser(ctx context.Context, req *protocol.CreateRequest) (*protocol.EmptyResponse, error) {
	u := repository.User{
		Email:    req.User.Email,
		Username: req.User.Username,
		Password: []byte(req.User.Password),
	}
	err := s.Rpc.CreateUser(ctx, &u)
	if err != nil {
		log.Error("error in CreateUser")
		return nil, err
	}
	return &protocol.EmptyResponse{}, nil
}

func (s *Server) Delete(ctx context.Context, req *protocol.DeleteRequest) (*protocol.EmptyResponse, error) {
	err := s.Rpc.Delete(ctx, req.Login)
	if err != nil {
		log.Error("error in Delete")
		return nil, err
	}
	return &protocol.EmptyResponse{}, nil
}

func (s *Server) Update(ctx context.Context, req *protocol.UpdateRequest) (*protocol.EmptyResponse, error) {
	u := repository.User{
		Email:     req.User.Email,
		Username:  req.User.Username,
		Password:  []byte(req.User.Password),
		Confirmed: req.User.Confirmed,
	}
	err := s.Rpc.Update(ctx, &u)
	if err != nil {
		log.Error("error in Update")
		return nil, err
	}
	return &protocol.EmptyResponse{}, nil
}

func (s *Server) Find(ctx context.Context, req *protocol.FindRequest) (*protocol.User, error) {
	user, err := s.Rpc.Find(ctx, req.UserField)
	if err != nil {
		log.Error("error in Find")
		return nil, err
	}

	createdAt, err := ptypes.TimestampProto(user.CreatedAt)
	if err != nil {
		log.Error("error in converting types golang Time to protobuf Timestamp")
		return nil, err
	}

	updatedAt, err := ptypes.TimestampProto(user.UpdatedAt)
	if err != nil {
		log.Error("error in converting types golang Time to protobuf Timestamp")
		return nil, err
	}

	return &protocol.User{
		Id:        int32(user.ID),
		Username:  user.Username,
		Confirmed: user.Confirmed,
		Email:     user.Email,
		Password:  string(user.Password),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
