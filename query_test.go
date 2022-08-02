package query_builder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewQueryFrom(t *testing.T) {
	q := NewQueryFrom(Table("table_name").As("tb"))
	assert.Equal(t, `table_name AS tb`, partToString(q.from))
}

func TestQuery_Select(t *testing.T) {
	q := NewQueryFrom(Table("table_name"))
	q.Select(Field("field_name_1").As("field_alias_1"))
	q.Select(Count(All()), Field("field_name_3"))
	assert.Equal(t, []string{`field_name_1 AS field_alias_1`, `COUNT(*)`, "field_name_3"}, partsToStrings(q.selectParts))
}

func TestQuery_InnerJoin(t *testing.T) {
	q := NewQueryFrom(Table("table_name"))
	q.InnerJoin(Table("table_name_2"), Field("table_name.id").Eq(Field("table_name_2.eid")))
	q.InnerJoin(Table("table_name_3"), Field("table_name_2.id").Eq(Field("table_name_3.eid")))
	assert.Equal(t, []string{"INNER JOIN table_name_2 ON table_name.id = table_name_2.eid", `INNER JOIN table_name_3 ON table_name_2.id = table_name_3.eid`}, partsToStrings(q.joinParts))
}

func TestQuery_Where(t *testing.T) {
	q := NewQueryFrom(Table("table_name"))
	q.Where(Field("field_name_1").Eq(Value(123)))
	q.Where(Field("field_name_2").Eq(Value("string")), Field("field_name_3").Is(True()))
	assert.Equal(t, []string{"field_name_1 = 123", `field_name_2 = 'string'`, "field_name_3 IS TRUE"}, partsToStrings(q.whereParts))
}

func TestQuery_GroupBy(t *testing.T) {
	q := NewQueryFrom(Table("table_name"))
	q.GroupBy(Field("field_name_1"))
	q.GroupBy(Field("field_name_2"), Field("field_name_3"))
	assert.Equal(t, []string{"field_name_1", "field_name_2", "field_name_3"}, partsToStrings(q.groupByParts))
}

func TestQuery_OrderBy(t *testing.T) {
	q := NewQueryFrom(Table("table_name"))
	q.OrderBy(Field("field_name_1"), OrderDirectionAsc)
	q.OrderBy(Field("field_name_2"), OrderDirectionDesc)
	assert.Equal(t, []string{"field_name_1 ASC", "field_name_2 DESC"}, partsToStrings(q.orderByParts))
}

func TestQuery_Limit(t *testing.T) {
	q := NewQueryFrom(Table("table_name"))
	q.Limit(1)
	assert.Equal(t, q.limit, 1)
}

func TestQuery_BuildQuery(t *testing.T) {
	q := NewQueryFrom(Table("table_name").As("tb"))
	q.Select(Field("field_name_1").As("field_alias_1"))
	q.Select(Count(All()), Field("field_name_3"))
	q.InnerJoin(Table("table_name_2"), Field("table_name.id").Eq(Field("table_name_2.eid")))
	q.InnerJoin(Table("table_name_3"), Field("table_name_2.id").Eq(Field("table_name_3.eid")))
	q.Where(Field("field_name_1").Eq(Value(123)))
	q.Where(Field("field_name_2").Eq(Value("string")), Field("field_name_3").Is(True()))
	q.Where(Alias("field_alias_1").In(Params([]interface{}{123, 456})), Alias("field_alias_1").Lt(ParamInt(789)))
	q.GroupBy(Field("field_name_1"))
	q.GroupBy(Field("field_name_2"), Field("field_name_3"))
	q.OrderBy(Field("field_name_1"), OrderDirectionAsc)
	q.OrderBy(Field("field_name_2"), OrderDirectionDesc)
	q.Limit(1)
	s, i := q.Build()
	assert.Equal(t, `SELECT field_name_1 AS field_alias_1, COUNT(*), field_name_3 FROM table_name AS tb INNER JOIN table_name_2 ON table_name.id = table_name_2.eid INNER JOIN table_name_3 ON table_name_2.id = table_name_3.eid WHERE field_name_1 = 123 AND field_name_2 = 'string' AND field_name_3 IS TRUE AND field_alias_1 IN (?, ?) AND field_alias_1 < ? GROUP BY field_name_1, field_name_2, field_name_3 ORDER BY field_name_1 ASC, field_name_2 DESC LIMIT 1`, s)
	assert.Equal(t, i, []interface{}{123, 456, 789})
}
