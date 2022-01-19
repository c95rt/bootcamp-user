package config

import (
	"github.com/joeshaw/envdecode"
)

type Config struct {
	HTTPPort  string `env:"HTTP_PORT,default=3001"`
	GRPCConn  gRPCConn
	JWTSecret string `env:"JWT_SECRET,default=5jqo59fAMvi1fj1oi1KDkmwcire9jpp"`
}

type AppConfig struct {
	Config Config
}

type gRPCConn struct {
	URL  string `env:"GRPC_URL,default=docker.for.mac.localhost"`
	Port string `env:"GRPC_PORT,default=50051"`
}

func NewAppConfig() (*AppConfig, error) {
	var conf Config
	if err := envdecode.Decode(&conf); err != nil {
		return nil, err
	}
	return &AppConfig{
		Config: conf,
	}, nil
}
