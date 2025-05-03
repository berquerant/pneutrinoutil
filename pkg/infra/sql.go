package infra

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type SQLParam struct {
	DSN             string
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

func NewSQL(ctx context.Context, param *SQLParam) (*sql.DB, error) {
	db, err := sql.Open("mysql", param.DSN)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(param.ConnMaxLifetime)
	db.SetMaxIdleConns(param.MaxIdleConns)
	db.SetMaxOpenConns(param.MaxOpenConns)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
