package gcurd

import (
	"context"
	"database/sql"
	entsql "entgo.io/ent/dialect/sql"
	"errors"
	"fmt"
)

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

func DeleteWhere[T MyInterface](c context.Context, db *sql.DB, obj T, wvs []WhereValue) error {
	deleteBuild := sqlBuilder().Delete(obj.Table())

	for _, wv := range wvs {
		if CheckIn(obj.Columns(), wv.Name) {
			p, err := predicate(wv)
			if err != nil {
				return errors.New(fmt.Sprintf("column name [%s] is no found", wv.Name))
			}
			deleteBuild.Where(p)
		}
	}

	query, args := deleteBuild.Query()
	_, err := db.ExecContext(c, query, args...)
	return err
}

func Update[T MyInterface](c context.Context, db *sql.DB, obj T, key string, value interface{}) error {
	if b := CheckIn(obj.Columns(), key); b != true {
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

func GetWhere[T MyInterface](c context.Context, db *sql.DB, obj T, wv []WhereValue) (T, error) {
	builder := sqlBuilder()
	selector := &entsql.Selector{}
	selector = builder.Select(obj.Columns()...).From(entsql.Table(obj.Table()))
	selector = whereBuild(obj, selector, wv)
	selector.Limit(1)
	query, args := selector.Query()
	err := db.QueryRowContext(c, query, args...).Scan(obj.Fields()...)
	return obj, err
}

// First the first record ordered by primary key
func First[T MyInterface](c context.Context, db *sql.DB, obj T) (T, error) {
	query, args := sqlBuilder().Select(obj.Columns()...).From(entsql.Table(obj.Table())).OrderBy(entsql.Asc(AutoIncrementId)).Limit(1).Query()
	err := db.QueryRowContext(c, query, args...).Scan(obj.Fields()...)
	return obj, err
}

// Last  record, ordered by primary key desc
func Last[T MyInterface](c context.Context, db *sql.DB, obj T) (T, error) {
	query, args := sqlBuilder().Select(obj.Columns()...).From(entsql.Table(obj.Table())).OrderBy(entsql.Desc(AutoIncrementId)).Limit(1).Query()
	err := db.QueryRowContext(c, query, args...).Scan(obj.Fields()...)
	return obj, err
}

func Find[T MyInterface](c context.Context, db *sql.DB, obj T, req *Request, f func() T) (objs []T, err error) {
	objs = []T{}
	query, args := buildFind(obj, req, SqlFind)
	rows, err := db.QueryContext(c, query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var o T
		o = f()
		err = rows.Scan(o.Fields()...)
		if err != nil {
			return
		}
		objs = append(objs, o)
	}
	return
}

func Page[T MyInterface](c context.Context, db *sql.DB, obj T, req *Request, f func() T) (objs []T, err error) {
	objs = []T{}
	query, args := buildFind(obj, req, SqlPageList)
	rows, err := db.QueryContext(c, query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var o T
		o = f()
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
