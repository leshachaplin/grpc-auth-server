package main

import (
	"context"
	"fmt"
	conf "github.com/leshachaplin/grpc-auth-server/internal/config"
	"github.com/leshachaplin/grpc-auth-server/internal/repository"
	auth2 "github.com/leshachaplin/grpc-auth-server/internal/server/auth"
	"github.com/leshachaplin/grpc-auth-server/internal/service/auth"
	"github.com/leshachaplin/grpc-auth-server/protocol"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
)

func main() {
	cfg, err := conf.NewConfig()
	if err != nil {
		log.Fatal()
	}


	lis, err := net.Listen("tcp", fmt.Sprintf(cfg.GrpcPort))
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	fmt.Printf("Server is listening on %v ...", cfg.GrpcPort)

	_, cnsl := context.WithCancel(context.Background())

	userRepo := repository.NewUserRepository(*cfg.PostgresConnection)
	claimsRepo := repository.NewRepositoryOfClaims(cfg.PostgresConnection)
	refreshRepo := repository.NewRefreshTokenRepository(cfg.PostgresConnection)
	restoreRepo := repository.NewRestoreRepository(*cfg.PostgresConnection)
	confirmRepo := repository.NewConfirmationRepository(*cfg.PostgresConnection)
	mailRepo := repository.NewEmailRepository(*cfg.PostgresConnection)

	r := auth.New(*userRepo, *claimsRepo, *mailRepo, *cfg,
		*refreshRepo, *restoreRepo, *confirmRepo)
	srv := &auth2.Server{Rpc: *r}

	authServiceService := &protocol.AuthServiceService{
		SignIn:         srv.SignIn,
		SignUp:         srv.SignUp,
		DeleteClaims:   srv.DeleteClaims,
		Delete:         srv.Delete,
		AddClaims:      srv.AddClaims,
		RefreshToken:   srv.RefreshToken,
		Confirm:        srv.Confirm,
		Restore:        srv.Restore,
		ForgotPassword: srv.ForgotPassword,
	}

	s := grpc.NewServer()
	protocol.RegisterAuthServiceService(s, authServiceService)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("server is not connect %s", err)
		}
	}()

	for {
		select {
		case <-c:
			cnsl()
			if err := cfg.PostgresConnection.Close(); err != nil {
				log.Errorf("database not closed %s", err)
			}
			log.Info("Cancel is successful")
			close(c)
			return
		}
	}
}
