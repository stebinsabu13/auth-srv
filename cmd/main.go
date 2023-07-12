package main

import (
	"fmt"
	"log"
	"net"

	"github.com/stebin13/auth-srv/pkg/config"
	"github.com/stebin13/auth-srv/pkg/db"
	"github.com/stebin13/auth-srv/pkg/pb"
	"github.com/stebin13/auth-srv/pkg/services"
	"github.com/stebin13/auth-srv/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed to load config", err)
	}
	h := db.Initdb(c)
	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "go-grpc-auth-svc",
		ExpirationHours: 24,
	}
	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Auth Svc on", c.Port)
	s := services.Server{
		H:   h,
		JWT: jwt,
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
