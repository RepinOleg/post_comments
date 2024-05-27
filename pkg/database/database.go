package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"os"
)

type DBConfig struct {
	Host     string
	Port     uint16
	User     string
	Password string
	DB       string
	SSLMode  string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetUint16("db.port"),
		User:     viper.GetString("db.user"),
		Password: os.Getenv("DB_PASSWORD"),
		DB:       viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}
}

func NewDB(cfg DBConfig) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB, cfg.SSLMode)

	connect, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	return connect, nil
}
