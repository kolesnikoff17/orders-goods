package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	defaultMaxConns = 5
)

// Db keeps pool of connections to db
type Db struct {
	maxConns int
	Pool     *sqlx.DB
}

// New is constructor for Db
func New(uri string, opts ...Option) (*Db, error) {
	c, err := sqlx.Connect("mysql", uri)
	if err != nil {
		return nil, err
	}
	db := &Db{Pool: c,
		maxConns: defaultMaxConns}
	for _, opt := range opts {
		opt(db)
	}
	db.Pool.SetMaxIdleConns(db.maxConns)
	db.Pool.SetMaxOpenConns(db.maxConns)
	return db, nil
}

// Close closes db's connection pool
func (db *Db) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
