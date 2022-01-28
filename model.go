package gcurd

type Model interface {
	SetID(id int)
	GetID() int
	Table() string
	Columns() []string
	Fields() []interface{}
}
