package config

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type Configuration struct {
	MariaDBConn mariadbConn
}

type AppContext struct {
	Config      Configuration
	MariaDBConn *sqlx.DB
}

type mariadbConn struct {
	URL            string `env:"DATA_BASE_URL,default=localhost"`
	Name           string `env:"DATA_BASE_NAME,default=bootcamp"`
	User           string `env:"DATA_BASE_USER,default=admin"`
	Port           int    `env:"DATA_BASE_PORT,default=3306"`
	Password       string `env:"DATA_BASE_PASSWORD,default=admin"`
	OpenConnection int    `env:"DATA_BASE_MAX_OPEN_CONNECTION,default=5"`
}

func CreateConnectionSQL(conf mariadbConn) (*sqlx.DB, error) {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", conf.User, conf.Password, conf.URL, strconv.Itoa(conf.Port), conf.Name)
	return sqlx.Connect("mysql", conn)
}
