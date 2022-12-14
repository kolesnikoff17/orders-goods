package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type (
	// Config is a struct with ENV variables
	Config struct {
		HTTP
		DB
		Logger
		Kafka
	}
	// HTTP -.
	HTTP struct {
		Port            string
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		ShutdownTimeout time.Duration
	}
	// DB -.
	DB struct {
		Host    string
		Port    string
		User    string
		Pwd     string
		Name    string
		MaxConn int
	}
	// Logger -.
	Logger struct {
		Level string
	}
	// Kafka -.
	Kafka struct {
		Host string
		Port string
	}
)

// NewConfig gets values from ENV
func NewConfig() *Config {
	cfg := &Config{}
	cfg.HTTP.Port = os.Getenv("PORT")
	cfg.HTTP.ReadTimeout, _ = time.ParseDuration(os.Getenv("SERVER_READ_TIMEOUT") + "s")
	cfg.HTTP.WriteTimeout, _ = time.ParseDuration(os.Getenv("SERVER_WRITE_TIMEOUT") + "s")
	cfg.HTTP.ShutdownTimeout, _ = time.ParseDuration(os.Getenv("SERVER_SHUTDOWN_TIMEOUT") + "s")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Pwd = os.Getenv("DB_PWD")
	cfg.DB.Name = os.Getenv("DB_NAME")
	cfg.DB.MaxConn, _ = strconv.Atoi(os.Getenv("DB_MAXCONNS"))
	cfg.Logger.Level = os.Getenv("LOG_LVL")
	cfg.Kafka.Host = os.Getenv("KAFKA_HOST")
	cfg.Kafka.Port = os.Getenv("KAFKA_PORT")
	return cfg
}

// DbParams formats connection string from config
func DbParams(cfg *Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=5",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Pwd,
		cfg.DB.Name)
}
