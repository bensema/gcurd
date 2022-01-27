package gcurd

import (
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

type Op int

const (
	// Predicate operators.
	OpEQ      Op = iota // =
	OpNEQ               // <>
	OpGT                // >
	OpGTE               // >=
	OpLT                // <
	OpLTE               // <=
	OpIn                // IN
	OpNotIn             // NOT IN
	OpLike              // LIKE
	OpIsNull            // IS NULL
	OpNotNull           // IS NOT NULL

	// Arithmetic operators.
	OpAdd // +
	OpSub // -
	OpMul // *
	OpDiv // / (Quotient)
	OpMod // % (Reminder)
)

var ops = [...]string{
	OpEQ:      "=",
	OpNEQ:     "<>",
	OpGT:      ">",
	OpGTE:     ">=",
	OpLT:      "<",
	OpLTE:     "<=",
	OpIn:      "IN",
	OpNotIn:   "NOT IN",
	OpLike:    "LIKE",
	OpIsNull:  "IS NULL",
	OpNotNull: "IS NOT NULL",
	OpAdd:     "+",
	OpSub:     "-",
	OpMul:     "*",
	OpDiv:     "/",
	OpMod:     "%",
}

type FindType int

const (
	SqlPageList FindType = iota
	SqlPageCount
	SqlFind
)

func sqlBuilder() *entsql.DialectBuilder {
	return entsql.Dialect(dialect.MySQL)
}
