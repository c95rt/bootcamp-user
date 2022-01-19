package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"

	"github.com/c95rt/bootcamp-user/http/config"
	userEndpoints "github.com/c95rt/bootcamp-user/http/endpoints"
	userRepository "github.com/c95rt/bootcamp-user/http/repository"
	userService "github.com/c95rt/bootcamp-user/http/service"
	userTransport "github.com/c95rt/bootcamp-user/http/transport"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "httpUserService",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("message", "http service started")
	defer level.Info(logger).Log("message", "http service ended")

	appConfig, err := config.NewAppConfig()
	if err != nil {
		level.Error(logger).Log("exit", err)
		panic(err)
	}

	var (
		grpcUserServiceAddr = flag.String("addr", fmt.Sprintf("%s:%s", appConfig.Config.GRPCConn.URL, appConfig.Config.GRPCConn.Port), "The gprcUserServer address in the format of host:port")
		httpAddr            = flag.String("http", fmt.Sprintf(":%s", appConfig.Config.HTTPPort), "http listen address")
	)

	flag.Parse()

	var grpcUserServiceConn *grpc.ClientConn
	opts := []grpc.DialOption{grpc.WithInsecure()}
	grpcUserServiceConn, err = grpc.Dial(*grpcUserServiceAddr, opts...)
	if err != nil {
		level.Error(logger).Log("exit", err)
		panic(err)
	}

	repository := userRepository.NewRepository(grpcUserServiceConn)

	var srv userService.Service
	srv = userService.NewService(repository, logger)

	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	endpoints := userEndpoints.MakeEndpoints(srv)
	go func() {
		httpHandler := userTransport.NewHTTPServer(appConfig, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, httpHandler)
	}()

	level.Error(logger).Log("exit", <-errChan)
}
