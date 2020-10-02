module github.com/leshachaplin/grpc-auth-server

go 1.15

require (
	github.com/caarlos0/env/v6 v6.3.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2
	github.com/jmoiron/sqlx v1.2.0
	github.com/labstack/echo/v4 v4.1.17
	github.com/leshachaplin/config v0.0.0-20200929120454-b9660bed7ef5
	github.com/leshachaplin/emailSender v0.0.0-20201002065921-270c9c5ff48b
	github.com/lib/pq v1.0.0
	github.com/sirupsen/logrus v1.7.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0
)
