package query_builder

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewTestDB() *sqlx.DB {
	return sqlx.MustConnect("mysql", "mysql:password@/example")
}
