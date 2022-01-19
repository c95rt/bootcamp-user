package config

import (
	"context"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joeshaw/envdecode"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MariaDBConn mariaDBConn
	MongoDBConn mongoDBConn
	GRPCPort    string `env:"GRPC_PORT,default=50051"`
	JWTSecret   string `env:"JWT_SECRET,default=5jqo59fAMvi1fj1oi1KDkmwcire9jpp"`
}

type AppConfig struct {
	Config Config
}

type mariaDBConn struct {
	URL            string `env:"MARIADB_URL,default=docker.for.mac.localhost"`
	Name           string `env:"MARIADB_DATABASE,default=bootcamp"`
	User           string `env:"MARIADB_USER,default=root"`
	Port           int    `env:"MARIADB_PORT,default=3307"`
	Password       string `env:"MARIADB_PASSWORD,default=admin"`
	OpenConnection int    `env:"MARIADB_MAX_OPEN_CONNECTION,default=5"`
}

type mongoDBConn struct {
	User     string `env:"MONGODB_USER,default=root"`
	Password string `env:"MONGODB_PASSWORD,default=admin"`
	Port     int    `env:"MONGODB_PORT,default=27017"`
	URL      string `env:"MONGODB_URL,default=docker.for.mac.localhost"`
	Database string `env:"MONGODB_DATABASE,default=admin"`
}

func CreateConnectionSQL(conf mariaDBConn) (*sqlx.DB, error) {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", conf.User, conf.Password, conf.URL, strconv.Itoa(conf.Port), conf.Name)
	return sqlx.Connect("mysql", conn)
}

func CreateConnectionMongoDB(conf mongoDBConn) (*mongo.Database, error) {
	conn := fmt.Sprintf("mongodb://%s:%s@%s:%s", conf.User, conf.Password, conf.URL, strconv.Itoa(conf.Port))
	client, err := mongo.NewClient(options.Client().ApplyURI(conn))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client.Database(conf.Database), nil
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
