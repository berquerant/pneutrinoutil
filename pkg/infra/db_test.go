package infra_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/berquerant/pneutrinoutil/pkg/infra"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

const createTestTableSql = `CREATE TABLE IF NOT EXISTS test (
  id INT PRIMARY KEY,
  name VARCHAR(255)
);`

const dropTestTableSql = `DROP TABLE IF EXISTS test;`

func setUp(t *testing.T) *sql.DB {
	t.Helper()
	dsn := fmt.Sprintf(
		"test:test@tcp(%s:%s)/test",
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
	)
	t.Logf("DSN=%s", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(dropTestTableSql); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(createTestTableSql); err != nil {
		t.Fatal(err)
	}
	return db
}

func tearDown(t *testing.T, db *sql.DB) {
	t.Helper()
	_, _ = db.Exec(dropTestTableSql)
	_ = db.Close()
}

type Table struct {
	ID   int
	Name string
}

func TestConn(t *testing.T) {
	db := setUp(t)
	defer tearDown(t, db)

	ctx := context.TODO()
	conn := infra.NewConn[Table](db)
	scan := func(f func(...any) error) (*Table, error) {
		var (
			id   int
			name string
		)
		if err := f(&id, &name); err != nil {
			return nil, err
		}
		return &Table{
			ID:   id,
			Name: name,
		}, nil
	}

	t.Run("query empty", func(t *testing.T) {
		got, err := conn.Query(ctx, &infra.QueryRequest[Table]{
			Query: "select id, name from test;",
			Scan:  scan,
		})
		if !assert.Nil(t, err) {
			return
		}
		assert.Equal(t, 0, len(got.Items))
	})

	t.Run("insert", func(t *testing.T) {
		got, err := conn.Exec(ctx, &infra.ExecRequest{
			Query: "insert into test (id, name) values (?, ?);",
			Args: []any{
				1,
				"first",
			},
			AssertResponse: func(r *infra.ExecResponse) error {
				if r.RowsAffected != 1 {
					return errors.New("Error")
				}
				return nil
			},
		})
		if !assert.Nil(t, err) {
			return
		}
		assert.Equal(t, int64(1), got.RowsAffected)
	})

	t.Run("query one", func(t *testing.T) {
		got, err := conn.Query(ctx, &infra.QueryRequest[Table]{
			Query: "select id, name from test;",
			Scan:  scan,
		})
		if !assert.Nil(t, err) {
			return
		}
		if !assert.Equal(t, 1, len(got.Items)) {
			return
		}
		item := got.Items[0]
		assert.Equal(t, &Table{
			ID:   1,
			Name: "first",
		}, item)
	})
}
