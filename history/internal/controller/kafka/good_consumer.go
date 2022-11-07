package kafka

import (
	"common.local/pkg/logger"
	"encoding/json"
	"github.com/Shopify/sarama"
	"history/internal/entity"
	"history/internal/usecase/good_uc"
	"time"
)

// GoodConsumer is kafka consumer for good topic
type GoodConsumer struct {
	uc good_uc.Good
	l  logger.Interface
}

// NewGoodConsumer is a constructor for GoodConsumer
func NewGoodConsumer(uc good_uc.Good, l logger.Interface) *GoodConsumer {
	return &GoodConsumer{
		uc: uc,
		l:  l,
	}
}

// GoodMessage is a kafka message value
type GoodMessage struct {
	Action string `json:"action"`
	Good   struct {
		ID   string `json:"id"`
		Data struct {
			Name     string `json:"name"`
			Category string `json:"category"`
			Price    string `json:"price"`
		} `json:"data,omitempty"`
	} `json:"good"`
}

// ConsumeClaim is a message handler
func (c GoodConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg := <-claim.Messages():
			var m GoodMessage
			err := json.Unmarshal(msg.Value, &m)
			if err != nil {
				c.l.Warnf("can't parse %s from topic %s with err: %w", string(msg.Value), msg.Topic, err)
				session.MarkMessage(msg, "")
				continue
			}
			switch m.Action {
			case "create":
				err = c.uc.CreateGood(session.Context(), entity.GoodInOrder{
					GoodID:   m.Good.ID,
					Name:     m.Good.Data.Name,
					Category: m.Good.Data.Category,
					Price:    m.Good.Data.Price,
				})
			case "update":
				err = c.uc.UpdateGood(session.Context(), entity.GoodInOrder{
					GoodID:   m.Good.ID,
					Name:     m.Good.Data.Name,
					Category: m.Good.Data.Category,
					Price:    m.Good.Data.Price,
				})
			case "delete":
				err = c.uc.DeleteGood(session.Context(), m.Good.ID)
			}
			if err != nil {
				c.l.Warnf("err %w with msg value %s", err, string(msg.Value))
				time.Sleep(time.Second)
			} else {
				session.MarkMessage(msg, "")
			}
		case <-session.Context().Done():
			return nil
		}
	}
}

// Setup -.
func (GoodConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup -.
func (GoodConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
