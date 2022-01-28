package gcurd

type MyInterface interface {
	SetID(id int)
	GetID() int
	Table() string
	Columns() []string
	Fields() []interface{}
}
