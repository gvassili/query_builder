package query_builder

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func partToString(p Part) string {
	s, _ := p.Build()
	return s
}

func partsToStrings(ps []Part) []string {
	strings := make([]string, 0, len(ps))
	for _, p := range ps {
		strings = append(strings, partToString(p))
	}
	return strings
}

func TestPart_As(t *testing.T) {
	assert.Equal(t, `field_name AS alias_name`, partToString(Field("field_name").As("alias_name")))
}

func TestPart_Eq(t *testing.T) {
	assert.Equal(t, `field_name = 123`, partToString(Field("field_name").Eq(ValueInt(123))))
}

func TestPart_Neq(t *testing.T) {
	assert.Equal(t, `field_name != 123`, partToString(Field("field_name").Ne(ValueInt(123))))
}

func TestPart_Lt(t *testing.T) {
	assert.Equal(t, `field_name < 123`, partToString(Field("field_name").Lt(ValueInt(123))))
}

func TestPart_Lte(t *testing.T) {
	assert.Equal(t, `field_name <= 123`, partToString(Field("field_name").Lte(ValueInt(123))))
}

func TestPart_Gt(t *testing.T) {
	assert.Equal(t, `field_name > 123`, partToString(Field("field_name").Gt(ValueInt(123))))
}

func TestPart_Gte(t *testing.T) {
	assert.Equal(t, `field_name >= 123`, partToString(Field("field_name").Gte(ValueInt(123))))
}

func TestPart_In(t *testing.T) {
	assert.Equal(t, `field_name IN (1, 2, 3)`, partToString(Field("field_name").In([]Part{ValueInt(1), ValueInt(2), ValueInt(3)})))
}

func TestPart_NotIn(t *testing.T) {
	assert.Equal(t, `field_name NOT IN (1, 2, 3)`, partToString(Field("field_name").NotIn([]Part{ValueInt(1), ValueInt(2), ValueInt(3)})))
}

func TestPart_Add(t *testing.T) {
	assert.Equal(t, `field_name + 123`, partToString(Field("field_name").Add(ValueInt(123))))
}

func TestPart_Sub(t *testing.T) {
	assert.Equal(t, `field_name - 123`, partToString(Field("field_name").Sub(ValueInt(123))))
}

func TestPart_Is(t *testing.T) {
	assert.Equal(t, `field_name IS NULL`, partToString(Field("field_name").Is(Null())))
}

func TestPart_IsNot(t *testing.T) {
	assert.Equal(t, `field_name IS NOT NULL`, partToString(Field("field_name").IsNot(Null())))
}

func TestPart_And(t *testing.T) {
	assert.Equal(t, `field_name_1 = 1 AND field_name_2 = 2`,
		partToString(Field("field_name_1").Eq(ValueInt(1)).
			And(Field("field_name_2").Eq(ValueInt(2)))),
	)
}

func Test_And(t *testing.T) {
	assert.Equal(t, `field_name_1 = 1 AND field_name_2 = 2`,
		partToString(And(
			Field("field_name_1").Eq(ValueInt(1)),
			Field("field_name_2").Eq(ValueInt(2)),
		)),
	)
}

func TestPart_Or(t *testing.T) {
	assert.Equal(t, `field_name_1 = 1 OR field_name_2 = 2`,
		partToString(Field("field_name_1").Eq(ValueInt(1)).
			Or(Field("field_name_2").Eq(ValueInt(2)))),
	)
}

func Test_Or(t *testing.T) {
	assert.Equal(t, `field_name_1 = 1 OR field_name_2 = 2`,
		partToString(Or(
			Field("field_name_1").Eq(ValueInt(1)),
			Field("field_name_2").Eq(ValueInt(2)),
		)),
	)
}

func TestPart_Xor(t *testing.T) {
	assert.Equal(t, `0 XOR 1`,
		partToString(ValueInt(0).Xor(ValueInt(1))),
	)
}

func Test_Xor(t *testing.T) {
	assert.Equal(t, `0 XOR 1`,
		partToString(Xor(ValueInt(0), ValueInt(1))),
	)
}

func TestCond(t *testing.T) {
	assert.Equal(t, `(field_name_1 = 1 AND field_name_2 = 2)`,
		partToString(Cond(Field("field_name_1").Eq(ValueInt(1)).
			And(Field("field_name_2").Eq(ValueInt(2))))),
	)
}

func TestMin(t *testing.T) {
	assert.Equal(t, `MIN(field_name)`, partToString(Min(Field("field_name"))))
}

func TestMax(t *testing.T) {
	assert.Equal(t, `MAX(field_name)`, partToString(Max(Field("field_name"))))
}

func TestCount(t *testing.T) {
	assert.Equal(t, `COUNT(field_name)`, partToString(Count(Field("field_name"))))
}

func TestAverage(t *testing.T) {
	assert.Equal(t, `CAST(AVG(field_name) AS DECIMAL(7,2))`, partToString(Average(Field("field_name"))))
}

func TestToBase64(t *testing.T) {
	assert.Equal(t, `TO_BASE64(field_name)`, partToString(ToBase64(Field("field_name"))))
}

func TestDateOverlaps(t *testing.T) {
	assert.Equal(t, `(a.starts_at IS NULL OR a.ends_at IS NULL OR ? <= a.ends_at AND ? >= a.starts_at)`, partToString(DateOverlaps(Field("a.starts_at"), Field("a.ends_at"), time.Time{}, time.Time{})))
}

func TestConcat(t *testing.T) {
	assert.Equal(t, `CONCAT(field_name)`, partToString(Concat(Field("field_name"))))
	assert.Equal(t, `CONCAT(field_name_1, field_name_2)`, partToString(Concat(Field("field_name_1"), Field("field_name_2"))))
}

func TestDistinct(t *testing.T) {
	assert.Equal(t, "DISTINCT(field_name)", partToString(Distinct(Field("field_name"))))
}

func TestJsonExtract(t *testing.T) {
	assert.Equal(t, `JSON_EXTRACT(field_name, '$cmd')`, partToString(JsonExtract(Field("field_name"), Value("$cmd"))))
}

func TestIf(t *testing.T) {
	assert.Equal(t, `IF(field_name = 'string', 1, 0)`, partToString(If(Field("field_name").Eq(ValueString("string")), ValueInt(1), ValueInt(0))))
}

func TestCase(t *testing.T) {
	assert.Equal(t, `CASE WHEN field_name = 'string' THEN 1 ELSE 0 END`, partToString(Case(Field("field_name").Eq(ValueString("string")), ValueInt(1), ValueInt(0))))
}

func TestExists(t *testing.T) {
	assert.Equal(t, `EXISTS(SELECT * FROM table)`, partToString(Exists(NewQuery().From(Table("table")).Select(All()))))
}

func TestAll(t *testing.T) {
	assert.Equal(t, "*", partToString(All()))
}

func TestNull(t *testing.T) {
	assert.Equal(t, "NULL", partToString(Null()))
}

func TestTrue(t *testing.T) {
	assert.Equal(t, "TRUE", partToString(True()))
}

func TestFalse(t *testing.T) {
	assert.Equal(t, "FALSE", partToString(False()))
}

func TestValue(t *testing.T) {
	assert.Equal(t, "'string'", partToString(Value("string")))
	assert.Equal(t, "123", partToString(Value(123)))
	assert.Equal(t, "-123", partToString(Value(-123)))
	if ts, err := time.Parse(time.RFC3339, "2001-02-03T04:05:06Z"); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "'2001-02-03 04:05:06'", partToString(Value(ts)))
	}
	assert.Equal(t, "TRUE", partToString(Value(true)))
	assert.Equal(t, "FALSE", partToString(Value(false)))
	assert.Panics(t, func() { partToString(Value(nil)) })
	assert.Panics(t, func() { partToString(Value(struct{}{})) })
	assert.Panics(t, func() { partToString(Value([]interface{}{"value", 1, false})) })
}

func TestValueString(t *testing.T) {
	assert.Equal(t, "'string'", partToString(ValueString("string")))
}

func TestValueInt(t *testing.T) {
	assert.Equal(t, "123", partToString(ValueInt(123)))
	assert.Equal(t, "-123", partToString(ValueInt(-123)))
}

func TestValueTime(t *testing.T) {
	if ts, err := time.Parse(time.RFC3339, "2001-02-03T04:05:06Z"); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "'2001-02-03 04:05:06'", partToString(ValueTime(ts)))
	}
}

func TestValueBool(t *testing.T) {
	assert.Equal(t, "TRUE", partToString(ValueBool(true)))
	assert.Equal(t, "FALSE", partToString(ValueBool(false)))
}

func TestTable(t *testing.T) {
	assert.Equal(t, "table_name", partToString(Table("table_name")))
}

func TestField(t *testing.T) {
	assert.Equal(t, "field_name", partToString(Field("field_name")))
}

func TestFieldNp(t *testing.T) {
	assert.Equal(t, "namespace_name.field_name", partToString(FieldNp("namespace_name", "field_name")))
}

func TestAlias(t *testing.T) {
	assert.Equal(t, `alias_name`, partToString(Alias("alias_name")))
}

func TestParamBool(t *testing.T) {
	s, v := ParamBool(true).Build()
	assert.Equal(t, `?`, s)
	assert.Equal(t, true, v[0])
	assert.Len(t, v, 1)
	s, v = ParamBool(false).Build()
	assert.Equal(t, `?`, s)
	assert.Equal(t, false, v[0])
	assert.Len(t, v, 1)
}

func TestParamInt(t *testing.T) {
	s, v := ParamInt(123).Build()
	assert.Equal(t, `?`, s)
	assert.Equal(t, 123, v[0])
	assert.Len(t, v, 1)
}

func TestParamString(t *testing.T) {
	s, v := ParamString("string").Build()
	assert.Equal(t, `?`, s)
	assert.Equal(t, "string", v[0])
	assert.Len(t, v, 1)
}

func TestParam(t *testing.T) {
	s, v := Param(123).Build()
	assert.Equal(t, `?`, s)
	assert.Equal(t, 123, v[0])
	assert.Len(t, v, 1)
	s, v = Param("string").Build()
	assert.Equal(t, `?`, s)
	assert.Equal(t, "string", v[0])
	assert.Len(t, v, 1)
	s, v = Param(true).Build()
	assert.Equal(t, `?`, s)
	assert.Equal(t, true, v[0])
	assert.Len(t, v, 1)
}

func TestParamBools(t *testing.T) {
	ps := ParamBools([]bool{true, false})
	assert.Equal(t, ps, []Part{ParamBool(true), ParamBool(false)})
}

func TestParamInts(t *testing.T) {
	ps := ParamInts([]int{123, -123})
	assert.Equal(t, ps, []Part{ParamInt(123), ParamInt(-123)})
}

func TestParamStrings(t *testing.T) {
	ps := ParamStrings([]string{"string_1", "string_2"})
	assert.Equal(t, ps, []Part{ParamString("string_1"), ParamString("string_2")})
}

func TestParamParams(t *testing.T) {
	ps := Params([]interface{}{true, 123, "string_1"})
	assert.Equal(t, ps, []Part{Param(true), Param(123), Param("string_1")})
	assert.Equal(t, ps, []Part{ParamBool(true), ParamInt(123), ParamString("string_1")})
}
