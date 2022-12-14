package main

import (
	"common.local/pkg/httpserver"
	"common.local/pkg/kafkaproducer"
	"common.local/pkg/logger"
	mongodb "common.local/pkg/mongo"
	"github.com/gin-gonic/gin"
	"good/config"
	v1 "good/internal/controller/http/v1"
	"good/internal/usecase"
	"good/internal/usecase/kafka"
	"good/internal/usecase/repository"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title           Good
// @version         1.0
// @description     CRUD API for goods db
// @contact.name   Developer
// @contact.email  vladkolesnikofff@gmail.com
// @license.name  MIT
// @license.url   https://github.com/kolesnikoff17/orders-goods/blob/main/LICENSE
// @host      localhost:8080
// @BasePath  /v1
func main() {
	cfg := config.NewConfig()
	uri := config.DbParams(cfg)

	l, err := logger.New(cfg.Logger.Level)
	if err != nil {
		log.Fatalf("failed to build logger: %s", err)
	}

	db, err := mongodb.New(uri)
	if err != nil {
		l.Fatalf("failed to connect to db: %s", err)
	}
	defer db.Close()

	p, err := kafkaproducer.New(cfg.Kafka.Host, cfg.Kafka.Port)
	if err != nil {
		l.Fatalf("failed to connect to kafka broker: %s", err)
	}
	defer p.Close()

	useCase := usecase.New(repository.New(db), kafka.New(p))

	handler := gin.New()
	v1.NewRouter(handler, useCase, l)
	server := httpserver.New(handler)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		l.Infof("shutting down with signal: %s", sig)
	case err = <-server.Notify():
		l.Infof("server err: %s", err)
	}

	err = server.Shutdown()
	if err != nil {
		l.Infof("server shutdown err: %s", err)
	}
}
