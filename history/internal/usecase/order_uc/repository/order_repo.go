package repository

import (
	"common.local/pkg/postgre"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"history/internal/entity"
)

// Order implements order_uc.OrderRepo
type Order struct {
	db postgre.Db
}

// New is a constructor for Order
func New(db postgre.Db) *Order {
	return &Order{
		db: db,
	}
}

// GetPrices return slice of prices of goods in order for sum calculation
func (r *Order) GetPrices(ctx context.Context, order entity.Order) (entity.Order, error) {
	for i, g := range order.Goods {
		err := r.db.Pool.GetContext(ctx, &order.Goods[i],
			`SELECT price FROM goods WHERE good_id = :id AND status_id = 1`, g)
		if err != nil {
			return entity.Order{}, entity.ErrNoID
		}
	}
	return order, nil
}

// Create make new order in db
func (r *Order) Create(ctx context.Context, order entity.Order) error {
	tx, err := r.db.Pool.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("order_repo - create: %w", err)
	}
	defer tx.Rollback()
	err = orderCreate(ctx, order, tx)
	if err != nil {
		return fmt.Errorf("order_repo - create: %w", err)
	}
	return tx.Commit()
}

// Update archive old state of order and create new
func (r *Order) Update(ctx context.Context, order entity.Order) error {
	tx, err := r.db.Pool.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("order_repo - update: %w", err)
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx,
		`UPDATE orders SET status_id = 2, modified = :modified 
              WHERE order_id = :order_id AND status_id = 1`, order)
	if err != nil {
		return fmt.Errorf("order_repo - update: %w", err)
	}
	err = orderCreate(ctx, order, tx)
	if err != nil {
		return fmt.Errorf("order_repo - update: %w", err)
	}
	return tx.Commit()
}

// Archive delete order from active
func (r *Order) Archive(ctx context.Context, id int) error {
	_, err := r.db.Pool.ExecContext(ctx,
		`UPDATE orders SET status_id = 3, modified = now() 
              WHERE order_id = $1 AND status_id = 1`, id)
	if err != nil {
		return fmt.Errorf("order_repo - archive: %w", err)
	}
	return nil
}

func orderCreate(ctx context.Context, order entity.Order, tx *sqlx.Tx) error {
	res, err := tx.NamedExecContext(ctx,
		`INSERT INTO orders (order_id, user_id, sum, created, modified, status_id) 
VALUES (:order_id, :user_id, :sum, :created, :modified, 1) RETURNING id`, order)
	if err != nil {
		return err
	}
	orderID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	for _, g := range order.Goods {
		res, err = tx.NamedExecContext(ctx,
			`INSERT INTO positions (good_id, amount) VALUES (:good_id, :amount) RETURNING id`, g)
		if err != nil {
			return err
		}
		posID, err := res.LastInsertId()
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx,
			`INSERT INTO order_has_position (order_id, position_id) VALUES ($1, $2)`, orderID, posID)
		if err != nil {
			return err
		}
	}
	return nil
}
