package gcurd

type level int

const (
	Debug Op = iota
	Info
)

var Level = Info

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

const AutoIncrementId = "id"
