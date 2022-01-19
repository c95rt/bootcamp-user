package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/c95rt/bootcamp-user/grpc/config"
	"github.com/c95rt/bootcamp-user/grpc/endpoints"
	"github.com/c95rt/bootcamp-user/grpc/pb"
	"github.com/c95rt/bootcamp-user/grpc/repository"
	"github.com/c95rt/bootcamp-user/grpc/service"
	"github.com/c95rt/bootcamp-user/grpc/transport"
	"google.golang.org/grpc"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "grpcUserService",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	appConfig, err := config.NewAppConfig()
	if err != nil {
		level.Error(logger).Log("exit", err)
		panic(err)
	}

	mariaDBConn, err := config.CreateConnectionSQL(appConfig.Config.MariaDBConn)
	if err != nil {
		level.Error(logger).Log("exit", err)
		panic(err)
	}
	defer mariaDBConn.Close()

	mongoDBConn, err := config.CreateConnectionMongoDB(appConfig.Config.MongoDBConn)
	if err != nil {
		level.Error(logger).Log("exit", err)
		panic(err)
	}

	repository, err := repository.NewRepository(mariaDBConn, mongoDBConn)
	if err != nil {
		level.Error(logger).Log("exit", err)
		panic(err)
	}

	srv := service.NewService(repository, logger)

	endpoints := endpoints.MakeEndpoints(srv)
	grpcServer := transport.NewGRPCServer(endpoints)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%s", appConfig.Config.GRPCPort))
	if err != nil {
		level.Error(logger).Log("exit", err)
		panic(err)
	}
	defer grpcListener.Close()

	go func() {
		baseServer := grpc.NewServer()
		pb.RegisterUserServiceServer(baseServer, grpcServer)
		level.Info(logger).Log("msg", "Server started")
		baseServer.Serve(grpcListener)
	}()

	level.Error(logger).Log("exit", <-errs)
}
