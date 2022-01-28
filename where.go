package gcurd

import (
	entsql "entgo.io/ent/dialect/sql"
)

func whereBuild[T MyInterface](obj T, selector *entsql.Selector, wvs []WhereValue) *entsql.Selector {
	for _, wv := range wvs {
		if CheckIn(obj.Columns(), wv.Name) {
			switch wv.Op {
			case OpEQ:
				selector = selector.Where(entsql.EQ(wv.Name, wv.Value))
			case OpNEQ:
				selector = selector.Where(entsql.NEQ(wv.Name, wv.Value))
			case OpGT:
				selector = selector.Where(entsql.GT(wv.Name, wv.Value))
			case OpGTE:
				selector = selector.Where(entsql.GTE(wv.Name, wv.Value))
			case OpLT:
				selector = selector.Where(entsql.LT(wv.Name, wv.Value))
			case OpLTE:
				selector = selector.Where(entsql.LTE(wv.Name, wv.Value))
			case OpIn:
				selector = selector.Where(entsql.In(wv.Name, wv.Value))
			case OpNotIn:
				selector = selector.Where(entsql.NotIn(wv.Name, wv.Value))
			case OpLike:
				// todo
			case OpIsNull:
				selector = selector.Where(entsql.IsNull(wv.Name))
			case OpNotNull:
				selector = selector.Where(entsql.NotNull(wv.Name))

			}
		}
	}
	return selector
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

	if CheckIn(obj.Columns(), req.OrderBy.Filed) {
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
	}

	selector.Offset((req.Pagination.Num - 1) * req.Pagination.Size)
	selector.Limit(req.Pagination.Size)
	return selector.Query()
}
