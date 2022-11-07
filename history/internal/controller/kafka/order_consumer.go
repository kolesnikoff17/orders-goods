package kafka

import (
	"common.local/pkg/logger"
	"encoding/json"
	"github.com/Shopify/sarama"
	"history/internal/entity"
	"history/internal/usecase/order_uc"
	"time"
)

// OrderConsumer is kafka consumer for good topic
type OrderConsumer struct {
	uc order_uc.Order
	l  logger.Interface
}

// NewOrderConsumer is a constructor for OrderConsumer
func NewOrderConsumer(uc order_uc.Order, l logger.Interface) *OrderConsumer {
	return &OrderConsumer{
		uc: uc,
		l:  l,
	}
}

// OrderMessage is a kafka message value
type OrderMessage struct {
	Action string `json:"action"`
	Order  struct {
		ID   int `json:"id"`
		Data struct {
			UserID   int       `json:"user_id" db:"user_id"`
			Created  time.Time `json:"created" db:"created"`
			Modified time.Time `json:"modified" db:"modified"`
			Goods    []struct {
				GoodID string `json:"good_id" db:"good_id"`
				Amount int    `json:"amount" db:"amount"`
			} `json:"goods"`
		} `json:"data,omitempty"`
	} `json:"order"`
}

// ConsumeClaim is a message handler
func (c OrderConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg := <-claim.Messages():
			var m OrderMessage
			err := json.Unmarshal(msg.Value, &m)
			if err != nil {
				c.l.Warnf("can't parse %s from topic %s with err: %w", string(msg.Value), msg.Topic, err)
				session.MarkMessage(msg, "")
				continue
			}
			switch m.Action {
			case "create":
				or := entity.Order{
					ID:       m.Order.ID,
					UserID:   m.Order.Data.UserID,
					Created:  m.Order.Data.Created,
					Modified: m.Order.Data.Modified,
					Goods:    make([]entity.GoodInOrder, 0, len(m.Order.Data.Goods)),
				}
				for _, v := range m.Order.Data.Goods {
					or.Goods = append(or.Goods, entity.GoodInOrder{GoodID: v.GoodID, Amount: v.Amount})
				}
				err = c.uc.CreateOrder(session.Context(), or)
			case "update":
				or := entity.Order{
					ID:       m.Order.ID,
					UserID:   m.Order.Data.UserID,
					Created:  m.Order.Data.Created,
					Modified: m.Order.Data.Modified,
					Goods:    make([]entity.GoodInOrder, 0, len(m.Order.Data.Goods)),
				}
				for _, v := range m.Order.Data.Goods {
					or.Goods = append(or.Goods, entity.GoodInOrder{GoodID: v.GoodID, Amount: v.Amount})
				}
				err = c.uc.UpdateOrder(session.Context(), or)
			case "delete":
				err = c.uc.DeleteOrder(session.Context(), m.Order.ID)
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
func (OrderConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup -.
func (OrderConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
