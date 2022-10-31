package entity

import "time"

// Order -.
type Order struct {
	ID       int            `json:"-" db:"id"`
	UserID   int            `json:"user_id" db:"user_id"`
	Created  time.Time      `json:"created" db:"created"`
	Modified time.Time      `json:"modified" db:"modified"`
	Goods    []GoodsInOrder `json:"goods"`
}

// GoodsInOrder -.
type GoodsInOrder struct {
	GoodID string `json:"good_id" db:"good_id"`
	Amount int    `json:"amount" db:"amount"`
}
