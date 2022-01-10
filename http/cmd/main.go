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
	"github.com/joeshaw/envdecode"
	"google.golang.org/grpc"

	"github.com/c95rt/bootcamp-user/http/config"
	userEndpoints "github.com/c95rt/bootcamp-user/http/endpoints"
	userRepository "github.com/c95rt/bootcamp-user/http/repository"
	userService "github.com/c95rt/bootcamp-user/http/service"
	userTransport "github.com/c95rt/bootcamp-user/http/transport"
)

type logLogger struct {
	log.Logger
}

func (l logLogger) exit(err error) {
	level.Error(l).Log("exit", err)
	os.Exit(-1)
}

func main() {
	var (
		grpcUserServiceAddr = flag.String("addr", "localhost:50051", "The gprcUserServer address in the format of host:port")
		httpAddr            = flag.String("http", ":8080", "http listen address")
	)
	var logger logLogger
	{
		logger.Logger = log.NewLogfmtLogger(os.Stderr)
		logger.Logger = log.NewSyncLogger(logger)
		logger.Logger = log.With(logger,
			"service", "httpService",
			"time", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("message", "http service started")
	defer level.Info(logger).Log("message", "http service ended")

	var conf config.Configuration
	if err := envdecode.Decode(&conf); err != nil {
		logger.exit(err)
	}

	flag.Parse()

	context := &config.AppContext{
		Config: conf,
	}

	var err error
	var grpcUserServiceConn *grpc.ClientConn
	opts := []grpc.DialOption{grpc.WithInsecure()}
	grpcUserServiceConn, err = grpc.Dial(*grpcUserServiceAddr, opts...)
	if err != nil {
		logger.exit(err)
	}

	repository := userRepository.NewRepository(grpcUserServiceConn)

	var srv userService.Service
	srv = userService.NewService(repository)

	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	endpoints := userEndpoints.MakeEndpoints(srv)
	go func() {
		httpHandler := userTransport.NewHTTPServer(endpoints)
		errChan <- http.ListenAndServe(*httpAddr, httpHandler)
	}()

	level.Error(logger).Log("exit", <-errChan)
}
