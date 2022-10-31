package kafka

import (
	"common.local/pkg/kafkaproducer"
	"encoding/json"
	"github.com/Shopify/sarama"
	"good/internal/entity"
	"os"
)

// Notifier is a listener interface
type Notifier interface {
	Notify(message Message) error
}

// Producer implements Notifier interface
type Producer struct {
	conn *kafkaproducer.Conn
}

// New is a constructor for Producer
func New(c *kafkaproducer.Conn) *Producer {
	return &Producer{
		conn: c,
	}
}

// Message is a model-listener message adapter
type Message struct {
	Action ActionType
	Data   Good
}

// Good is a kafka message
type Good struct {
	ID   string      `json:"id"`
	Data entity.Good `json:"data,omitempty"`
}

// ActionType is an enum type for possible actions
type ActionType string

// Enum for ActionType
const (
	Create ActionType = "create"
	Update ActionType = "update"
	Delete ActionType = "delete"
)

// Notify publish new message to kafka topic
func (p *Producer) Notify(message Message) error {
	msg := msgPrepare(message)
	_, _, err := p.conn.Producer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

func msgPrepare(m Message) *sarama.ProducerMessage {
	var topic string
	switch m.Action {
	case Create:
		topic = os.Getenv("KAFKA_CREATE_TOPIC")
	case Update:
		topic = os.Getenv("KAFKA_UPDATE_TOPIC")
	case Delete:
		topic = os.Getenv("KAFKA_DELETE_TOPIC")
	}
	msg, _ := json.Marshal(m.Data)
	return &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(msg),
	}
}
