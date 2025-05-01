package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Conn[T any] struct {
	db *sql.DB
}

func NewConn[T any](db *sql.DB) *Conn[T] {
	return &Conn[T]{
		db: db,
	}
}

type QueryRequest[T any] struct {
	Query string
	Args  []any
	Scan  func(f func(...any) error) (*T, error)
}

type QueryResponse[T any] struct {
	Items []*T
}

func (r *QueryResponse[T]) AssertRows(rows int) error {
	n := len(r.Items)
	if n != rows {
		return fmt.Errorf("%w: want %d got %d", ErrAssertRows, rows, n)
	}
	return nil
}

type Queryer[T any] interface {
	Query(ctx context.Context, req *QueryRequest[T]) (*QueryResponse[T], error)
}

var _ Queryer[int] = &Conn[int]{}

var (
	ErrAssertRows = errors.New("AssertRows")
)

func (c *Conn[T]) Query(ctx context.Context, req *QueryRequest[T]) (*QueryResponse[T], error) {
	conn, err := c.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, req.Query, req.Args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*T
	for rows.Next() {
		item, err := req.Scan(rows.Scan)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &QueryResponse[T]{
		Items: items,
	}, nil
}

type ExecRequest struct {
	Query          string
	Args           []any
	AssertResponse func(*ExecResponse) error
}

var (
	ErrRowsAffected    = errors.New("RowsAffected")
	ErrNotLastInserted = errors.New("NotLastInserted")
)

func AssertRowsAffected(rows int64) func(*ExecResponse) error {
	return func(r *ExecResponse) error {
		if r.RowsAffected != rows {
			return fmt.Errorf("%w: want %d got %d", ErrRowsAffected, rows, r.RowsAffected)
		}
		return nil
	}
}

func AssertLastInserted() func(*ExecResponse) error {
	return func(r *ExecResponse) error {
		if _, ok := r.LastInsertId(); ok {
			return nil
		}
		return ErrNotLastInserted
	}
}

func AssertJoin(f ...func(*ExecResponse) error) func(*ExecResponse) error {
	return func(r *ExecResponse) error {
		for _, a := range f {
			if err := a(r); err != nil {
				return err
			}
		}
		return nil
	}
}

func (r *ExecRequest) assert(res *ExecResponse) error {
	if f := r.AssertResponse; f != nil {
		return f(res)
	}
	return nil
}

type ExecResponse struct {
	RowsAffected int64
	lastInsertId *int64
}

func (r *ExecResponse) LastInsertId() (int64, bool) {
	if x := r.lastInsertId; x != nil {
		return *x, true
	}
	return 0, false
}

type Execer interface {
	Exec(ctx context.Context, req *ExecRequest) (*ExecResponse, error)
}

var _ Execer = &Conn[int]{}

func (c *Conn[T]) Exec(ctx context.Context, req *ExecRequest) (*ExecResponse, error) {
	conn, err := c.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		return nil, err
	}

	result, err := tx.ExecContext(ctx, req.Query, req.Args...)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	response := &ExecResponse{
		RowsAffected: rows,
	}
	if x, err := result.LastInsertId(); err == nil {
		response.lastInsertId = &x
	}

	if err := req.assert(response); err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return response, nil
}
