package main

import (
	"context"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/jmoiron/sqlx"
	"github.com/leshachaplin/config"
	conf "github.com/leshachaplin/grpc-auth-server/internal/config"
	"github.com/leshachaplin/grpc-auth-server/internal/email"
	"github.com/leshachaplin/grpc-auth-server/internal/repository"
	"github.com/leshachaplin/grpc-auth-server/internal/server"
	"github.com/leshachaplin/grpc-auth-server/internal/service"
	"github.com/leshachaplin/grpc-auth-server/protocol"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
)

func main() {
	cfg := conf.NewConfig()

	awsConf, err := config.NewForAws("us-west-2")
	if err != nil {
		log.Fatal("Can't connect to aws", err)
	}

	mongoConn, err := awsConf.GetMongo(cfg.MongoConStr)
	if err != nil {
		log.Fatal("Can't get mongo connection string", err)
	}

	session, err := mgo.Dial(mongoConn.ConnectionString)

	postgresConn, err := awsConf.GetSQL(cfg.PostgresConStr)
	if err != nil {
		log.Fatal("Can't get mongo connection string", err)
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		postgresConn.Username, postgresConn.Password, postgresConn.DB)
	db, err := sqlx.Open(postgresConn.Schema, connStr)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	smtp, err := awsConf.GetSMTP(cfg.SMTPConStr)
	if err != nil {
		log.Fatal("Can't get smtp connection string", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(cfg.GrpcPort))
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	fmt.Printf("Server is listening on %v ...", cfg.Port)

	_, cnsl := context.WithCancel(context.Background())
	userRepo := repository.NewUserRepository(*db)
	claimsRepo := repository.NewRepositoryOfClaims(session.DB("claims"))
	refreshRepo := repository.NewRefreshTokenRepository(db)
	restoreRepo := repository.NewRestoreRepository(*db)
	confirmRepo := repository.NewConfirmationRepository(*db)
	smtpSender := email.NewSMTPEmail(smtp.Username, smtp.Email, smtp.Password, smtp.Host)

	r := service.New(*userRepo, *claimsRepo, *awsConf, *cfg,
		*refreshRepo, *restoreRepo, *confirmRepo, *smtpSender)
	s := grpc.NewServer()
	srv := &server.Server{Rpc: *r}
	authServiceService := protocol.AuthServiceService{
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
	protocol.RegisterAuthServiceService(s, &authServiceService)
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
			if err := db.Close(); err != nil {
				log.Errorf("database not closed %s", err)
			}
			log.Info("Cansel is succesful")
			close(c)
			return
		}
	}
}
