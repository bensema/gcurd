package gcurd

import (
	"context"
	"database/sql"
	entsql "entgo.io/ent/dialect/sql"
	"errors"
	"fmt"
)

func Create[T Model](c context.Context, db *sql.DB, obj T) (sql.Result, error) {
	query, args := sqlBuilder().Insert(obj.Table()).Columns(obj.Columns()[1:]...).Values(obj.Fields()[1:]...).Query()
	return db.ExecContext(c, query, args...)
}

// Delete delete where id = ?
func Delete[T Model](c context.Context, db *sql.DB, obj T, id int) error {
	obj.SetID(id)
	query, args := sqlBuilder().Delete(obj.Table()).Where(entsql.EQ(AutoIncrementId, obj.GetID())).Query()
	_, err := db.ExecContext(c, query, args...)
	return err
}

// DeleteWhere delete where [something] = ?
func DeleteWhere[T Model](c context.Context, db *sql.DB, obj T, wvs []*WhereValue) error {
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

// Update update set [key] = [value] where id = ?
func Update[T Model](c context.Context, db *sql.DB, obj T, id int, key string, value any) error {
	obj.SetID(id)
	if b := CheckIn(obj.Columns(), key); b != true {
		return errors.New("update column error")
	}
	query, args := sqlBuilder().Update(obj.Table()).Set(key, value).Where(entsql.EQ(AutoIncrementId, obj.GetID())).Query()
	_, err := db.ExecContext(c, query, args...)
	return err
}

// UpdateWhere update [key] = [value] where [something] = ?
func UpdateWhere[T Model](c context.Context, db *sql.DB, obj T, key string, value any, wvs []*WhereValue) error {
	if b := CheckIn(obj.Columns(), key); b != true {
		return errors.New("update column error")
	}
	updateBuild := sqlBuilder().Update(obj.Table()).Set(key, value)

	for _, wv := range wvs {
		if CheckIn(obj.Columns(), wv.Name) {
			p, err := predicate(wv)
			if err != nil {
				return errors.New(fmt.Sprintf("column name [%s] is no found", wv.Name))
			}
			updateBuild.Where(p)
		}
	}

	query, args := updateBuild.Query()
	_, err := db.ExecContext(c, query, args...)
	return err
}

// UpdateWhereKV update [key] = [value], [key] = [value],... where [something] = ?
func UpdateWhereKV[T Model](c context.Context, db *sql.DB, obj T, kvs []KeyValue, wvs []*WhereValue) error {

	updateBuild := sqlBuilder().Update(obj.Table())
	for _, kv := range kvs {
		if b := CheckIn(obj.Columns(), kv.Key); b != true {
			return errors.New("update column error")
		}
		updateBuild = updateBuild.Set(kv.Key, kv.Value)
	}
	for _, wv := range wvs {
		if CheckIn(obj.Columns(), wv.Name) {
			p, err := predicate(wv)
			if err != nil {
				return errors.New(fmt.Sprintf("column name [%s] is no found", wv.Name))
			}
			updateBuild.Where(p)
		}
	}

	query, args := updateBuild.Query()
	_, err := db.ExecContext(c, query, args...)
	return err
}

// Get select where id = ?
func Get[T Model](c context.Context, db *sql.DB, obj T, id int) (T, error) {
	obj.SetID(id)
	query, args := sqlBuilder().Select(obj.Columns()...).From(entsql.Table(obj.Table())).Where(entsql.EQ(AutoIncrementId, obj.GetID())).Query()
	err := db.QueryRowContext(c, query, args...).Scan(obj.Fields()...)
	return obj, err
}

// GetWhere select where[something] = ?
func GetWhere[T Model](c context.Context, db *sql.DB, obj T, wvs []*WhereValue) (T, error) {
	builder := sqlBuilder()
	selector := &entsql.Selector{}
	selector = builder.Select(obj.Columns()...).From(entsql.Table(obj.Table()))
	selector = whereBuild(obj, selector, wvs)
	selector.Limit(1)
	query, args := selector.Query()
	err := db.QueryRowContext(c, query, args...).Scan(obj.Fields()...)
	return obj, err
}

// First the first record ordered by id asc
func First[T Model](c context.Context, db *sql.DB, obj T) (T, error) {
	query, args := sqlBuilder().Select(obj.Columns()...).From(entsql.Table(obj.Table())).OrderBy(entsql.Asc(AutoIncrementId)).Limit(1).Query()
	err := db.QueryRowContext(c, query, args...).Scan(obj.Fields()...)
	return obj, err
}

// Last  record, ordered by id desc
func Last[T Model](c context.Context, db *sql.DB, obj T) (T, error) {
	query, args := sqlBuilder().Select(obj.Columns()...).From(entsql.Table(obj.Table())).OrderBy(entsql.Desc(AutoIncrementId)).Limit(1).Query()
	err := db.QueryRowContext(c, query, args...).Scan(obj.Fields()...)
	return obj, err
}

func Find[T Model](c context.Context, db *sql.DB, obj T, wvs []*WhereValue, f func() T) (objs []T, err error) {
	objs = []T{}
	req := &Request{
		Where: wvs,
	}
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

func PageFind[T Model](c context.Context, db *sql.DB, obj T, req *Request, f func() T) (objs []T, err error) {
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

func PageTotal[T Model](c context.Context, db *sql.DB, obj T, req *Request) (total int, err error) {
	total = 0
	query, args := buildFind(obj, req, SqlPageCount)
	err = db.QueryRowContext(c, query, args...).Scan(&total)
	return
}
