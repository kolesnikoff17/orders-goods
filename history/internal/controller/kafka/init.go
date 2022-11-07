package kafka

import (
	"common.local/pkg/logger"
	"context"
	"github.com/Shopify/sarama"
	"os"
)

func RunConsumers(ctx context.Context, handlers map[string]sarama.ConsumerGroupHandler, l logger.Interface, addr string) {
	kafkaConsumerGroups := initAllConsumerGroups(l, addr)

	for topic, group := range kafkaConsumerGroups {
		go func(topic string, group *sarama.ConsumerGroup) {
			defer func() {
				if r := recover(); r != nil {
					l.Errorf("panic: %s", r)
				}
			}()

			for {
				err := (*group).Consume(ctx, []string{topic}, handlers[topic])
				if err != nil {
					l.Warnf("consumer group error: %w", err)
				}
			}
		}(topic, group)
	}
}

func initAllConsumerGroups(l logger.Interface, addr string) map[string]*sarama.ConsumerGroup {
	return map[string]*sarama.ConsumerGroup{
		os.Getenv("KAFKA_ORDER_TOPIC"): initGroup(os.Getenv("KAFKA_ORDER_TOPIC"), l, addr),
		os.Getenv("KAFKA_GOODS_TOPIC"): initGroup(os.Getenv("KAFKA_GOODS_TOPIC"), l, addr),
	}
}

func initGroup(topic string, l logger.Interface, addr string) *sarama.ConsumerGroup {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_3_0_0
	cfg.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup([]string{addr}, topic, cfg)
	if err != nil {
		l.Warnf("kafka - initgroup: %w", err)
		return nil
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				l.Errorf("panic: %s", r)
			}
		}()

		for err := range group.Errors() {
			l.Warnf("consumer group error: %w", err)
		}
	}()

	return &group
}
