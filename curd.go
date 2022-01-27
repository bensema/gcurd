package gcurd

import (
	"context"
	"database/sql"
	entsql "entgo.io/ent/dialect/sql"
	"errors"
	"fmt"
)

const AutoIncrementId = "id"

func Create[T MyInterface](c context.Context, db *sql.DB, obj T) error {
	query, args := sqlBuilder().Insert(obj.Table()).Columns(obj.Columns()...).Values(obj.Fields()...).Query()
	_, err := db.ExecContext(c, query, args...)
	return err
}

func Delete[T MyInterface](c context.Context, db *sql.DB, obj T) error {
	query, args := sqlBuilder().Delete(obj.Table()).Where(entsql.EQ(AutoIncrementId, obj.GetID())).Query()
	_, err := db.ExecContext(c, query, args...)
	return err
}

func Update[T MyInterface](c context.Context, db *sql.DB, obj T, key string, value interface{}) error {
	if b := checkIn(obj.Columns(), key); b != true {
		return errors.New("update field error")
	}
	query, args := sqlBuilder().Update(obj.Table()).Set(key, value).Where(entsql.EQ(AutoIncrementId, obj.GetID())).Query()
	_, err := db.ExecContext(c, query, args...)
	return err
}

func Get[T MyInterface](c context.Context, db *sql.DB, obj T) (T, error) {
	query, args := sqlBuilder().Select(obj.Columns()...).From(entsql.Table(obj.Table())).Where(entsql.EQ(AutoIncrementId, obj.GetID())).Query()
	err := db.QueryRowContext(c, query, args...).Scan(obj.Fields()...)
	return obj, err
}

func Find[T MyInterface](c context.Context, db *sql.DB, obj T, req *Request) (objs []T, err error) {
	objs = []T{}
	query, args := buildFind(obj, req, SqlFind)
	rows, err := db.QueryContext(c, query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		//var o T
		o := *new(T)
		err = rows.Scan(o.Fields()...)
		if err != nil {
			return
		}
		fmt.Println(o)
		objs = append(objs, o)
	}
	return
}

func Page[T MyInterface](c context.Context, db *sql.DB, obj T, req *Request) (objs []T, err error) {
	objs = []T{}
	query, args := buildFind(obj, req, SqlPageList)
	rows, err := db.QueryContext(c, query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		//var o T
		o := *new(T)
		err = rows.Scan(o.Fields()...)
		if err != nil {
			return
		}
		objs = append(objs, o)
	}
	return
}

func PageTotal[T MyInterface](c context.Context, db *sql.DB, obj T, req *Request) (total int, err error) {
	total = 0
	query, args := buildFind(obj, req, SqlPageCount)
	err = db.QueryRowContext(c, query, args...).Scan(&total)
	return
}

func buildFind[T MyInterface](obj T, req *Request, findType FindType) (string, []interface{}) {
	builder := sqlBuilder()
	selector := &entsql.Selector{}
	switch findType {
	case SqlPageList, SqlFind:
		selector = builder.Select(obj.Columns()...)
	case SqlPageCount:
		selector = builder.Select("Count(*)")
	}

	selector = selector.From(entsql.Table(obj.Table()))

	selector = whereBuild(obj, selector, req.Where)

	// count 返回
	if findType == SqlPageCount {
		return selector.Query()
	}

	if findType == SqlFind {
		return selector.Query()
	}

	orderBy := ""
	switch req.OrderBy.Direction {
	case "desc", "DESC":
		orderBy = entsql.Desc(req.OrderBy.Filed)
	case "asc", "ASC":
		orderBy = entsql.Asc(req.OrderBy.Filed)
	default:
		orderBy = entsql.Asc(req.OrderBy.Filed)
	}

	selector = selector.OrderBy(orderBy)

	selector.Offset((req.Pagination.Num - 1) * req.Pagination.Size)
	selector.Limit(req.Pagination.Size)
	return selector.Query()
}
