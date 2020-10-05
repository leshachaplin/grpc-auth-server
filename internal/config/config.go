package config

// if using go modules

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/jmoiron/sqlx"
	"github.com/leshachaplin/config"
	emailProto "github.com/leshachaplin/emailSender/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Config struct {
	Port                     int    `env:"Port" envDefault:"1323"`
	GrpcPort                 string `env:"GrpcPort" envDefault:"0.0.0.0:50051"`
	UserGrpcPort             string `env:"UserGrpcPort" envDefault:"0.0.0.0:50053"`
	EmailGrpcPort            string `env:"EmailGrpcPort" envDefault:"0.0.0.0:50052"`
	SecretKeyAuth            string
	AuthSecret               string `env:"AuthSecret"`
	SecretKeyRefresh         string
	RefreshSecret            string `env:"RefreshSecret"`
	PostgresConnectionSecret string `env:"PostgresConnectionSecret"`
	SMTPConnectionSecret     string `env:"SMTPConStr"`
	PostgresConnection       *sqlx.DB
	SMTPClient               emailProto.EmailServiceClient
}

//TODO initialize secrets
func NewConfig() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	awsConf, err := config.NewForAws("us-west-2")
	if err != nil {
		log.Errorf("Can't connect to aws", err)
		return nil, err
	}

	postgresConn, err := awsConf.GetSQL(cfg.PostgresConnectionSecret)
	if err != nil {
		log.Errorf("Can't get SQL connection string", err)
		return nil, err
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		postgresConn.Username, postgresConn.Password, postgresConn.DB)
	db, err := sqlx.Open(postgresConn.Schema, connStr)
	if err != nil {
		log.Fatal("Can't connect to database ", err)
	}

	secretKeyAuth, err := awsConf.GetSecret(cfg.AuthSecret)
	if err != nil {
		log.Errorf("Can't get auth secret from aws ")
		return nil, err
	}
	secretKeyRefresh, err := awsConf.GetSecret(cfg.RefreshSecret)
	if err != nil {
		log.Errorf("Can't get refresh secret from aws ")
		return nil, err
	}

	opts := grpc.WithInsecure()
	clientConnInterface, err := grpc.Dial(cfg.EmailGrpcPort, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer clientConnInterface.Close()
	client := emailProto.NewEmailServiceClient(clientConnInterface)

	cfg.PostgresConnection = db
	cfg.SMTPClient = client
	cfg.SecretKeyAuth = secretKeyAuth.ApiKey
	cfg.SecretKeyRefresh = secretKeyRefresh.ApiKey

	return &cfg, nil
}
