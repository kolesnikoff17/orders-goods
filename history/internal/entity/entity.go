package entity

import "time"

// Order -.
type Order struct {
	ID       int           `json:"-" bd:"order_id"`
	UserID   int           `json:"user_id" bd:"user_id"`
	Sum      string        `json:"sum" bd:"sum"`
	Status   string        `json:"status" bd:"status_name"`
	StatusID int           `json:"-" bd:"status_id"`
	Created  time.Time     `json:"created" bd:"created"`
	Modified time.Time     `json:"modified" bd:"modified"`
	Goods    []GoodInOrder `json:"goods"`
}

// GoodInOrder -.
type GoodInOrder struct {
	GoodID   string    `json:"-" bd:"good_id"`
	Name     string    `json:"name" bd:"name"`
	Category string    `json:"category" bd:"category"`
	Price    string    `json:"price" bd:"price"`
	Amount   int       `json:"amount" bd:"amount"`
	Created  time.Time `json:"created" bd:"created"`
	Modified time.Time `json:"modified" bd:"modified"`
}
