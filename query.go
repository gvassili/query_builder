package query_builder

import (
	"strconv"
)

type builder interface {
	Build() (string, []interface{})
}

type partQuery struct {
	s      string
	params []interface{}
}

func (pq partQuery) Build() (string, []interface{}) {
	return pq.s, pq.params
}

type Query struct {
	from         Part
	selectParts  []Part
	joinParts    []Part
	whereParts   []Part
	groupByParts []Part
	orderByParts []Part
	withParts    []Part
	limit        int
}

func NewQueryFrom(table Part) *Query {
	return &Query{
		from: table,
	}
}

func NewQuery() *Query {
	return &Query{}
}

func (q *Query) From(table Part) *Query {
	q.from = table
	return q
}

func Union(lhs, rhs *Query) Part {
	return Part{parts{partByte('('), lhs, partByte(')'), partString(" UNION "), partByte('('), rhs, partByte(')')}}
}

func (q *Query) Select(v ...Part) *Query {
	q.selectParts = append(q.selectParts, v...)
	return q
}

func (q *Query) LeftJoin(table Part, cond Part) *Query {
	q.joinParts = append(q.joinParts, Part{parts{partString("LEFT JOIN "), table, partString(" ON "), cond}})
	return q
}

func (q *Query) InnerJoin(table Part, cond Part) *Query {
	q.joinParts = append(q.joinParts, Part{parts{partString("INNER JOIN "), table, partString(" ON "), cond}})
	return q
}

func (q *Query) Where(v ...Part) *Query {
	q.whereParts = append(q.whereParts, v...)
	return q
}

func (q *Query) GroupBy(v ...Part) *Query {
	q.groupByParts = append(q.groupByParts, v...)
	return q
}

type OrderDirection int

const (
	OrderDirectionAsc = iota
	OrderDirectionDesc
)

func (q *Query) OrderBy(v Part, dir OrderDirection) *Query {
	switch dir {
	case OrderDirectionAsc:
		q.orderByParts = append(q.orderByParts, Part{parts{v, partString(" ASC")}})
	case OrderDirectionDesc:
		q.orderByParts = append(q.orderByParts, Part{parts{v, partString(" DESC")}})
	}
	return q
}

func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

func (q *Query) With(part Part, name Part) *Query {
	q.withParts = append(q.withParts, Part{parts{name, partString(" AS ("), part.builder, partByte(')')}})
	return q
}

func (q *Query) WithRecursive(part Part, name string) *Query {
	q.withParts = append(q.withParts, Part{parts{partString("RECURSIVE " + name + " AS ("), part.builder, partByte(')')}})
	return q
}

func (q *Query) Part() Part {
	s, ps := q.Build()
	return Part{parts{partByte('('), partQuery{s, ps}, partByte(')')}}
}

func (q *Query) Build() (string, []interface{}) {
	parts := parts{}
	if len(q.withParts) > 0 {
		parts = append(parts, partString("WITH "))
		for i, v := range q.withParts {
			if i != 0 {
				parts = append(parts, partString(", "))
			}
			parts = append(parts, v)
		}
		parts = append(parts, partByte(' '))
	}

	parts = append(parts, partString("SELECT "))
	for i, v := range q.selectParts {
		if i != 0 {
			parts = append(parts, partString(", "))
		}
		parts = append(parts, v)
	}
	parts = append(parts, partString(" FROM "))
	parts = append(parts, q.from)
	if len(q.joinParts) != 0 {
		for _, v := range q.joinParts {
			parts = append(parts, partByte(' '), v)
		}
	}
	if len(q.whereParts) != 0 {
		parts = append(parts, partString(" WHERE "))
		for i, v := range q.whereParts {
			if i != 0 {
				parts = append(parts, partString(" AND "))
			}
			parts = append(parts, v)
		}
	}
	if len(q.groupByParts) != 0 {
		parts = append(parts, partString(" GROUP BY "))
		for i, v := range q.groupByParts {
			if i != 0 {
				parts = append(parts, partString(", "))
			}
			parts = append(parts, v)
		}
	}
	if len(q.orderByParts) != 0 {
		parts = append(parts, partString(" ORDER BY "))
		for i, v := range q.orderByParts {
			if i != 0 {
				parts = append(parts, partString(", "))
			}
			parts = append(parts, v)
		}
	}
	if q.limit != 0 {
		parts = append(parts, partString(" LIMIT "+strconv.Itoa(q.limit)))
	}
	return parts.Build()
}
