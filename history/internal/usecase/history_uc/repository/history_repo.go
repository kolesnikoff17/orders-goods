package historyrepo

import (
	"common.local/pkg/postgre"
	"context"
	"fmt"
	"history/internal/entity"
)

// History implements history_uc.HistoryRepo
type History struct {
	db *postgre.Db
}

// New is a constructor for History
func New(db *postgre.Db) *History {
	return &History{
		db: db,
	}
}

// GetByID return order in its last state
func (r *History) GetByID(ctx context.Context, id int) (entity.Order, error) {
	var order entity.Order
	err := r.db.Pool.GetContext(ctx, &order,
		`SELECT o.id AS order_id, o.user_id, o.sum, o.created, o.modified, s.status_name 
FROM orders AS o JOIN status AS s ON o.status_id = s.status_id WHERE o.order_id = $1 
                                                               ORDER BY o.modified DESC LIMIT 1`, id)
	if err != nil {
		return entity.Order{}, entity.ErrNoID
	}
	err = r.db.Pool.SelectContext(ctx, &order.Goods,
		`SELECT g.good_id, g.name, g.category, g.price, p.amount, g.created, g.modified
FROM goods AS g JOIN positions AS p ON p.good_id = g.good_id
JOIN order_has_position AS ohp ON ohp.position_id = p.id
JOIN orders AS o ON o.id = ohp.order_id
WHERE o.id = $1`, order.ID)
	if err != nil {
		return entity.Order{}, fmt.Errorf("history_repo - getbyid: %w", err)
	}
	return order, nil
}

func (r *History) GetHistory(ctx context.Context, id int) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.Pool.SelectContext(ctx, &orders,
		`SELECT o.id AS order_id, o.user_id, o.sum, o.created, o.modified, s.status_name 
FROM orders AS o JOIN status AS s ON o.status_id = s.status_id WHERE o.order_id = $1 
                                                               ORDER BY o.modified DESC`, id)
	if err != nil {
		return nil, fmt.Errorf("history_repo - gethistory: %w", err)
	}
	if orders == nil {
		return nil, entity.ErrNoID
	}
	for i, o := range orders {
		err = r.db.Pool.SelectContext(ctx, &orders[i].Goods,
			`SELECT g.good_id, g.name, g.category, g.price, p.amount, g.created, g.modified
FROM goods AS g JOIN positions AS p ON p.good_id = g.good_id
JOIN order_has_position AS ohp ON ohp.position_id = p.id
JOIN orders AS o ON o.id = ohp.order_id
WHERE o.id = $1`, o.ID)
		if err != nil {
			return nil, fmt.Errorf("history_repo - gethistory: %w", err)
		}
	}
	return orders, nil
}
