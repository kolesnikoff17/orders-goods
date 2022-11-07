package cmd

import (
	"common.local/pkg/httpserver"
	"common.local/pkg/logger"
	"common.local/pkg/postgre"
	"context"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"history/config"
	v1 "history/internal/controller/http/v1"
	"history/internal/controller/kafka"
	"history/internal/usecase/good_uc"
	goodrepo "history/internal/usecase/good_uc/repository"
	"history/internal/usecase/history_uc"
	historyrepo "history/internal/usecase/history_uc/repository"
	"history/internal/usecase/order_uc"
	orderrepo "history/internal/usecase/order_uc/repository"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title           History
// @version         1.0
// @description     API for read operations on history database
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

	db, err := postgre.New(uri)
	if err != nil {
		l.Fatalf("failed to connect to db: %s", err)
	}
	defer db.Close()

	gUc := good_uc.New(goodrepo.New(db))
	oUc := order_uc.New(orderrepo.New(db))

	handlers := map[string]sarama.ConsumerGroupHandler{
		os.Getenv("GOODS_CREATED_TOPIC"): kafka.NewGoodConsumer(gUc, l),
		os.Getenv("ORDER_CREATED_TOPIC"): kafka.NewOrderConsumer(oUc, l),
	}
	ctx, cancel := context.WithCancel(context.Background())
	kafka.RunConsumers(ctx, handlers, l, cfg.Kafka.Host+":"+cfg.Kafka.Port)

	hUc := history_uc.New(historyrepo.New(db))

	handler := gin.New()
	v1.NewRouter(handler, hUc, l)
	server := httpserver.New(handler)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		l.Infof("shutting down with signal: %s", sig)
	case err = <-server.Notify():
		l.Infof("server err: %s", err)
	}

	cancel()
	err = server.Shutdown()
	if err != nil {
		l.Infof("server shutdown err: %s", err)
	}
}
