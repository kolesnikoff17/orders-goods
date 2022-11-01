package history_uc

import (
	"context"
	"errors"
	"fmt"
	"history/internal/entity"
)

// HistoryUseCase implements ReadOpInterface and keeps HistoryRepo
type HistoryUseCase struct {
	r HistoryRepo
}

// New is a constructor for HistoryUseCase
func New(r HistoryRepo) *HistoryUseCase {
	return &HistoryUseCase{
		r: r,
	}
}

// GetOrderByID return order with given id
func (uc *HistoryUseCase) GetOrderByID(ctx context.Context, id int) (entity.Order, error) {
	o, err := uc.r.GetByID(ctx, id)
	switch {
	case errors.Is(err, entity.ErrNoID):
		return entity.Order{}, err
	case err != nil:
		return entity.Order{}, fmt.Errorf("history_uc - getorderbyid: %w", err)
	}
	return o, nil
}

// GetOrderHistory return changes history of order with given id
func (uc *HistoryUseCase) GetOrderHistory(ctx context.Context, id int) ([]entity.Order, error) {
	h, err := uc.r.GetHistory(ctx, id)
	switch {
	case errors.Is(err, entity.ErrNoID):
		return nil, err
	case err != nil:
		return nil, fmt.Errorf("history_uc - getorderhistory: %w", err)
	}
	return h, nil
}
