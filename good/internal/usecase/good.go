package usecase

import (
	"context"
	"errors"
	"fmt"
	"good/internal/entity"
	"good/internal/usecase/kafka"
)

// Good is model layer, it keeps kafka and db connection and implements business-logic
type Good struct {
	r GoodRepo
	p kafka.Notifier
}

// New is a constructor for Good
func New(r GoodRepo, p kafka.Notifier) *Good {
	return &Good{
		r: r,
		p: p,
	}
}

// GetGood return good by its id
func (uc *Good) GetGood(ctx context.Context, id string) (entity.Good, error) {
	g, err := uc.r.GetByID(ctx, id)
	switch {
	case errors.Is(err, entity.ErrNoID):
		return entity.Good{}, err
	case err != nil:
		return entity.Good{}, fmt.Errorf("usecase - getgood: %w", err)
	}
	return g, nil
}

// NewGood creates new good in db
func (uc *Good) NewGood(ctx context.Context, good entity.Good) (string, error) {
	id, err := uc.r.Create(ctx, good)
	if err != nil {
		return "", fmt.Errorf("usecase - newgood: %w", err)
	}
	err = uc.p.Notify(kafka.Message{Action: kafka.Create, Data: kafka.Good{ID: id, Data: good}})
	if err != nil {
		return "", fmt.Errorf("usecase - newgood: %w", err)
	}
	return id, nil
}

// UpdateGood updates good by its id
func (uc *Good) UpdateGood(ctx context.Context, good entity.Good) error {
	_, err := uc.r.GetByID(ctx, good.ID)
	switch {
	case errors.Is(err, entity.ErrNoID):
		return err
	case err != nil:
		return fmt.Errorf("usecase - updategood: %w", err)
	}
	err = uc.r.Update(ctx, good)
	if err != nil {
		return fmt.Errorf("usecase - updategood: %w", err)
	}
	err = uc.p.Notify(kafka.Message{Action: kafka.Update, Data: kafka.Good{ID: good.ID, Data: good}})
	if err != nil {
		return fmt.Errorf("usecase - updategood: %w", err)
	}
	return nil
}

// DeleteGood deletes good by its id
func (uc *Good) DeleteGood(ctx context.Context, id string) error {
	_, err := uc.r.GetByID(ctx, id)
	switch {
	case errors.Is(err, entity.ErrNoID):
		return err
	case err != nil:
		return fmt.Errorf("usecase - deletegood: %w", err)
	}
	err = uc.r.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("usecase - deletegood: %w", err)
	}
	err = uc.p.Notify(kafka.Message{Action: kafka.Delete, Data: kafka.Good{ID: id}})
	if err != nil {
		return fmt.Errorf("usecase - deletegood: %w", err)
	}
	return nil
}
