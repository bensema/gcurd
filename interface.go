package gcurd

type MyInterface interface {
	GetID() int
	Table() string
	Columns() []string
	Fields() []interface{}
}
