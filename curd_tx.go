package gcurd

import (
	"context"
	"database/sql"
	entsql "entgo.io/ent/dialect/sql"
	"errors"
	"fmt"
)

// CreateTx
func CreateTx[T Model](c context.Context, db *sql.Tx, obj T) (sql.Result, error) {
	query, args := sqlBuilder().Insert(obj.Table()).Columns(obj.Columns()[1:]...).Values(obj.Fields()[1:]...).Query()
	return db.ExecContext(c, query, args...)
}

// DeleteTx delete where id = ?
func DeleteTx[T Model](c context.Context, db *sql.Tx, obj T, id int) error {
	obj.SetID(id)
	query, args := sqlBuilder().Delete(obj.Table()).Where(entsql.EQ(AutoIncrementId, obj.GetID())).Query()
	_, err := db.ExecContext(c, query, args...)
	return err
}

// DeleteWhereTx delete where [something] = ?
func DeleteWhereTx[T Model](c context.Context, db *sql.Tx, obj T, wvs []*WhereValue) error {
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

// UpdateTx update set [key] = [value] where id = ?
func UpdateTx[T Model](c context.Context, db *sql.Tx, obj T, id int, key string, value any) error {
	obj.SetID(id)
	if b := CheckIn(obj.Columns(), key); b != true {
		return errors.New("update column error")
	}
	query, args := sqlBuilder().Update(obj.Table()).Set(key, value).Where(entsql.EQ(AutoIncrementId, obj.GetID())).Query()
	_, err := db.ExecContext(c, query, args...)
	return err
}

// UpdateWhereTx update [key] = [value] where [something] = ?
func UpdateWhereTx[T Model](c context.Context, db *sql.Tx, obj T, key string, value any, wvs []*WhereValue) error {
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

// UpdateWhereKVTx update [key] = [value], [key] = [value],... where [something] = ?
func UpdateWhereKVTx[T Model](c context.Context, db *sql.Tx, obj T, kvs []KeyValue, wvs []*WhereValue) error {

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

// GetTx select where id = ?
func GetTx[T Model](c context.Context, db *sql.Tx, obj T, id int) (T, error) {
	obj.SetID(id)
	query, args := sqlBuilder().Select(obj.Columns()...).From(entsql.Table(obj.Table())).Where(entsql.EQ(AutoIncrementId, obj.GetID())).Query()
	err := db.QueryRowContext(c, query, args...).Scan(obj.Fields()...)
	return obj, err
}

// GetWhereTx select where[something] = ?
func GetWhereTx[T Model](c context.Context, db *sql.Tx, obj T, wvs []*WhereValue) (T, error) {
	builder := sqlBuilder()
	selector := &entsql.Selector{}
	selector = builder.Select(obj.Columns()...).From(entsql.Table(obj.Table()))
	selector = whereBuild(obj, selector, wvs)
	selector.Limit(1)
	query, args := selector.Query()
	err := db.QueryRowContext(c, query, args...).Scan(obj.Fields()...)
	return obj, err
}

func FindTx[T Model](c context.Context, db *sql.Tx, obj T, wvs []*WhereValue, f func() T) (objs []T, err error) {
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

func PageFindTx[T Model](c context.Context, db *sql.Tx, obj T, req *Request, f func() T) (objs []T, err error) {
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

func PageTotalTx[T Model](c context.Context, db *sql.Tx, obj T, req *Request) (total int, err error) {
	total = 0
	query, args := buildFind(obj, req, SqlPageCount)
	err = db.QueryRowContext(c, query, args...).Scan(&total)
	return
}
