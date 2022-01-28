package gcurd

import (
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

const (
	MySQL    = dialect.MySQL
	SQLite   = dialect.SQLite
	Postgres = dialect.SQLite
	Gremlin  = dialect.Gremlin
)

var Dialect = MySQL

func sqlBuilder() *entsql.DialectBuilder {
	return entsql.Dialect(Dialect)
}
