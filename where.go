package gcurd

import (
	entsql "entgo.io/ent/dialect/sql"
)

func whereBuild[T MyInterface](obj T, selector *entsql.Selector, wvs []WhereValue) *entsql.Selector {
	for _, wv := range wvs {
		if checkIn(obj.Columns(), wv.Name) {
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
				//
			case OpIsNull:
				selector = selector.Where(entsql.IsNull(wv.Name))
			case OpNotNull:
				selector = selector.Where(entsql.NotNull(wv.Name))

			}
		}
	}
	return selector
}
