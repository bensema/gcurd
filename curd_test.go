package gcurd

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

type Admin struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Status   int    `json:"status"`
}

func (m *Admin) Table() string {
	return "admin"
}

func (m *Admin) SetID(id int) {
	m.Id = id
}

func (m *Admin) GetID() int {
	return m.Id
}

func (m *Admin) Columns() []string {
	return []string{"id", "username", "status"}
}

func (m *Admin) Fields() []any {
	return []any{&m.Id, &m.Username, &m.Status}
}

func (Admin) New() *Admin {
	return &Admin{}
}

func TestCreate(t *testing.T) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	sqlTable := `
		DROP TABLE IF EXISTS "admin";
		CREATE TABLE IF NOT EXISTS "admin" (
		   "id" INTEGER PRIMARY KEY AUTOINCREMENT,
		   "username" VARCHAR(64) NULL,
		   "status" INTEGER  
		);`
	db.Exec(sqlTable)

	Dialect = SQLite
	c := context.Background()
	obj := &Admin{
		Username: "foo",
		Status:   1,
	}
	_, err = Create(c, db, obj)
	if err != nil {
		fmt.Println(err)
	}
	obj = &Admin{
		Username: "far",
		Status:   2,
	}
	_, err = Create(c, db, obj)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGet(t *testing.T) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	Dialect = SQLite
	c := context.Background()
	obj := &Admin{}
	obj, err = Get(c, db, obj, 1)
	if err != nil {
		t.Error(err)
	}
	if obj.Username == "foo" {
		t.Log("Get test success")
	} else {
		t.Error("Get test fail")
	}
}
