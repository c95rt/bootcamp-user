package main

import (
	"net"

	"github.com/joeshaw/envdecode"
	joonix "github.com/joonix/log"
	log "github.com/sirupsen/logrus"

	"github.com/c95rt/bootcamp-user/grpc/config"
	"github.com/c95rt/bootcamp-user/grpc/endpoints"
	"github.com/c95rt/bootcamp-user/grpc/pb"
	"github.com/c95rt/bootcamp-user/grpc/repository"
	"github.com/c95rt/bootcamp-user/grpc/service"
	"github.com/c95rt/bootcamp-user/grpc/transport"
	"google.golang.org/grpc"
)

func main() {
	log.SetFormatter(joonix.NewFormatter())

	var conf config.Configuration
	if err := envdecode.Decode(&conf); err != nil {
		log.Fatalf("could not load the app configuration: %v", err)
	}
	context := &config.AppContext{
		Config: conf,
	}

	conn, err := config.CreateConnectionSQL(context.Config.MariaDBConn)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	repository, err := repository.NewRepository(conn)
	if err != nil {
		log.Fatal(err)
	}

	srv := service.NewService(repository)

	endpoints := endpoints.MakeEndpoints(srv)
	grpcServer := transport.NewGRPCServer(endpoints)

	grpcListener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		baseServer := grpc.NewServer()
		pb.RegisterUserServiceServer(baseServer, grpcServer)
		log.Info("server started")
		baseServer.Serve(grpcListener)
	}()
}
