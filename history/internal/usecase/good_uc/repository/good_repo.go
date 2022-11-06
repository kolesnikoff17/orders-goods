package repository

import (
	"common.local/pkg/postgre"
	"context"
	"database/sql"
	"fmt"
	"history/internal/entity"
)

// Good implement good_uc.GoodRepo
type Good struct {
	db postgre.Db
}

// New is a constructor for Good
func New(db postgre.Db) *Good {
	return &Good{
		db: db,
	}
}

// Create put new good in repository
func (r *Good) Create(ctx context.Context, good entity.GoodInOrder) error {
	_, err := r.db.Pool.NamedExecContext(ctx,
		`INSERT INTO goods (good_id, name, price, category, status_id) 
VALUES (:good_id, :name, :price, :category, 1)`, good)
	if err != nil {
		return fmt.Errorf("good_repo - create: %w", err)
	}
	return nil
}

// Update archive old state of good and create new
func (r *Good) Update(ctx context.Context, good entity.GoodInOrder) error {
	tx, err := r.db.Pool.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("good_repo - update: %w", err)
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx,
		`UPDATE goods SET status_id = 2 WHERE good_id = $1 AND status_id = 1`, good.GoodID)
	if err != nil {
		return fmt.Errorf("good_repo - update: %w", err)
	}
	_, err = tx.NamedExecContext(ctx,
		`INSERT INTO goods (good_id, name, price, category, status_id) 
VALUES (:good_id, :name, :price, :category, 1)`, good)
	if err != nil {
		return fmt.Errorf("good_repo - update: %w", err)
	}
	return tx.Commit()
}

// Archive mark good as archived
func (r *Good) Archive(ctx context.Context, id string) error {
	_, err := r.db.Pool.ExecContext(ctx, `UPDATE goods SET status_id = 2 WHERE good_id = $1 AND status_id = 1`, id)
	if err != nil {
		return fmt.Errorf("good_repo - archive: %w", err)
	}
	return nil
}
