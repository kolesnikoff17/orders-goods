package repository

import (
	"common.local/pkg/mysql"
	"context"
	"database/sql"
	"fmt"
	"order/internal/entity"
)

// OrderRepo implements repo interface
type OrderRepo struct {
	conn *mysql.Db
}

// New is a constructor for OrderRepo
func New(c *mysql.Db) *OrderRepo {
	return &OrderRepo{
		conn: c,
	}
}

// GetByID return order with given id from db, entity.ErrNoID if there is no one
func (r *OrderRepo) GetByID(ctx context.Context, id int) (entity.Order, error) {
	var o entity.Order
	err := r.conn.Pool.GetContext(ctx, &o,
		`SELECT user_id, created, modified FROM orders WHERE id = ?`, id)
	if err != nil {
		return entity.Order{}, entity.ErrNoID
	}
	err = r.conn.Pool.SelectContext(ctx, &o.Goods,
		`SELECT good_id, amount FROM goods WHERE order_id = ?`, id)
	if err != nil {
		return entity.Order{}, fmt.Errorf("orderrepo - getbyid: %w", err)
	}
	return o, nil
}

// Create insert new order in db
func (r *OrderRepo) Create(ctx context.Context, order entity.Order) (int, error) {
	tx, err := r.conn.Pool.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("orderrepo - create: %w", err)
	}
	defer tx.Rollback()
	res, err := tx.NamedExecContext(ctx,
		`INSERT INTO orders (user_id) VALUES (?)`, order)
	if err != nil {
		return 0, fmt.Errorf("orderrepo - create: %w", err)
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("orderrepo - create: %w", err)
	}
	for _, v := range order.Goods {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO goods (order_id, good_id, amount) VALUES (?, ?, ?)`, lastID, v.GoodID, v.Amount)
		if err != nil {
			return 0, fmt.Errorf("orderrepo - create: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("orderrepo - create: %w", err)
	}
	return int(lastID), nil
}

// Update change order data in db
func (r *OrderRepo) Update(ctx context.Context, order entity.Order) error {
	tx, err := r.conn.Pool.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("orderrepo - update: %w", err)
	}
	defer tx.Rollback()
	_, err = tx.NamedExecContext(ctx,
		`UPDATE orders SET user_id = ? WHERE id = ?`, order)
	if err != nil {
		return fmt.Errorf("orderrepo - update: %w", err)
	}
	_, err = tx.ExecContext(ctx,
		`DELETE FROM goods WHERE order_id = ?`, order.ID)
	if err != nil {
		return fmt.Errorf("orderrepo - update: %w", err)
	}
	for _, v := range order.Goods {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO goods (order_id, good_id, amount) VALUES (?, ?, ?)`, order.ID, v.GoodID, v.Amount)
		if err != nil {
			return fmt.Errorf("orderrepo - update: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("orderrepo - update: %w", err)
	}
	return nil
}

// Delete remove order from db
func (r *OrderRepo) Delete(ctx context.Context, id int) error {
	tx, err := r.conn.Pool.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("orderrepo - delete: %w", err)
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx,
		`DELETE FROM goods WHERE order_id = ?`, id)
	if err != nil {
		return fmt.Errorf("orderrepo - delete: %w", err)
	}
	_, err = tx.ExecContext(ctx,
		`DELETE FROM orders WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("orderrepo - delete: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("orderrepo - delete: %w", err)
	}
	return nil
}
