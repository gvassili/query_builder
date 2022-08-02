package query_builder

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type partByte byte

func (b partByte) Build() (string, []interface{}) {
	return string(b), nil
}

type partString string

func (t partString) Build() (string, []interface{}) {
	return string(t), nil
}

type partParam struct {
	v interface{}
}

func (a partParam) Build() (string, []interface{}) {
	return "?", []interface{}{a.v}
}

type parts []builder

func (ps parts) Build() (string, []interface{}) {
	var sb strings.Builder
	var vs []interface{}
	for _, p := range ps {
		s, v := p.Build()
		sb.WriteString(s)
		vs = append(vs, v...)
	}
	return sb.String(), vs
}

type Part struct {
	builder
}

func (p Part) IsZero() bool {
	return p.builder == nil
}

func (p Part) append(vs ...builder) Part {
	if ps, ok := p.builder.(parts); ok {
		return Part{append(ps, vs...)}
	} else {
		return Part{append(parts{p}, vs...)}
	}
}

func (p Part) As(v string) Part {
	return p.append(partString(" AS " + v))
}

func (p Part) Eq(v Part) Part {
	return p.append(partString(" = "), v)
}

func (p Part) Ne(v Part) Part {
	return p.append(partString(" != "), v)
}

func (p Part) Lt(v Part) Part {
	return p.append(partString(" < "), v)
}

func (p Part) Lte(v Part) Part {
	return p.append(partString(" <= "), v)
}

func (p Part) Gt(v Part) Part {
	return p.append(partString(" > "), v)
}

func (p Part) Gte(v Part) Part {
	return p.append(partString(" >= "), v)
}

func (p Part) In(vs []Part) Part {
	if len(vs) == 0 {
		p = p.append(partString(" IN (NULL)"))
	} else {
		p = p.append(partString(" IN ("))
		for i, v := range vs {
			if i != 0 {
				p = p.append(partString(", "))
			}
			p = p.append(v)
		}
		p = p.append(partByte(')'))
	}
	return p
}

func (p Part) NotIn(vs []Part) Part {
	return p.append(partString(" NOT")).In(vs)
}

func (p Part) Add(v Part) Part {
	return p.append(partString(" + "), v)
}

func (p Part) Sub(v Part) Part {
	return p.append(partString(" - "), v)
}

func (p Part) Is(v Part) Part {
	return p.append(partString(" IS "), v)
}

func (p Part) IsNot(v Part) Part {
	return p.append(partString(" IS NOT "), v)
}

func (p Part) And(v Part) Part {
	return p.append(partString(" AND "), v)
}

func And(l, r Part) Part {
	return Part{parts{l, partString(" AND "), r}}
}

func (p Part) Or(v Part) Part {
	return p.append(partString(" OR "), v)
}

func Or(l, r Part) Part {
	return Part{parts{l, partString(" OR "), r}}
}

func (p Part) Xor(v Part) Part {
	return p.append(partString(" XOR "), v)
}

func Xor(l, r Part) Part {
	return Part{parts{l, partString(" XOR "), r}}
}

func List(vs ...Part) Part {
	p := make(parts, 0, len(vs)<<1)
	for i, v := range vs {
		if i != 0 {
			p = append(p, partString(", "))
		}
		p = append(p, v)
	}
	return Part{p}
}

func Cond(v Part) Part {
	return Part{parts{partByte('('), v, partByte(')')}}
}

func Min(v Part) Part {
	return Part{parts{partString("MIN("), v, partByte(')')}}
}

func Max(v Part) Part {
	return Part{parts{partString("MAX("), v, partByte(')')}}
}

func Count(v Part) Part {
	return Part{parts{partString("COUNT("), v, partByte(')')}}
}

func Average(v Part) Part {
	return Part{parts{partString("CAST(AVG("), v, partString(") AS DECIMAL(7,2))")}}
}

func ToBase64(v Part) Part {
	return Part{parts{partString("TO_BASE64("), v, partByte(')')}}
}

func DateOverlaps(startA Part, endA Part, startB time.Time, endB time.Time) Part {
	if startB.Unix() <= 0 {
		startB = time.Now()
	}
	if endB.Unix() <= 0 {
		endB = time.Now()
	}
	return Cond(Or(
		Or(startA.Is(Null()), endA.Is(Null())),
		And(Param(startB).Lte(endA), Param(endB).Gte(startA))))
}

func Concat(vs ...Part) Part {
	p := make(parts, 0, 2+(len(vs)<<1))
	p = append(p, partString("CONCAT("))
	for i, v := range vs {
		if i != 0 {
			p = append(p, partString(", "))
		}
		p = append(p, v)
	}
	p = append(p, partByte(')'))
	return Part{p}
}

func Distinct(v Part) Part {
	return Part{parts{partString("DISTINCT("), v, partByte(')')}}
}

func JsonExtract(v, cmd Part) Part {
	return Part{parts{partString("JSON_EXTRACT("), v, partString(", "), cmd, partByte(')')}}
}

func If(cond, def, v Part) Part {
	return Part{parts{partString("IF("), cond, partString(", "), def, partString(", "), v, partByte(')')}}
}

func Case(cond, then, els Part) Part {
	return Part{parts{partString("CASE WHEN "), cond, partString(" THEN "), then, partString(" ELSE "), els, partString(" END")}}
}

func Exists(query *Query) Part {
	return Part{parts{partString("EXISTS"), query.Part()}}
}

var cAll = Part{partString("*")}

func All() Part {
	return cAll
}

var cNull = Part{partString("NULL")}

func Null() Part {
	return cNull
}

var cTrue = Part{partString("TRUE")}

func True() Part {
	return cTrue
}

var cFalse = Part{partString("FALSE")}

func False() Part {
	return cFalse
}

func Value(v interface{}) Part {
	switch t := v.(type) {
	case string:
		return Part{ValueString(t)}
	case int:
		return Part{ValueInt(t)}
	case time.Time:
		return Part{ValueTime(t)}
	case bool:
		return Part{ValueBool(t)}
	default:
		panic(fmt.Errorf("invalid value type %T", v))
	}
}

func ValueString(v string) Part {
	return Part{partString("'" + v + "'")}
}

func ValueInt(v int) Part {
	return Part{partString(strconv.Itoa(v))}
}

func ValueTime(v time.Time) Part {
	return Part{partString("'" + v.Format("2006-01-02 15:04:05") + "'")}
}

func ValueBool(v bool) Part {
	if v {
		return True()
	} else {
		return False()
	}
}

func Table(name string) Part {
	return Part{partString(name)}
}

func Field(field string) Part {
	return Part{partString(field)}
}

func FieldNp(np string, field string) Part {
	return Part{parts{partString(np), partByte('.'), partString(field)}}
}

func Alias(alias string) Part {
	return Part{partString(alias)}
}

func ParamBool(v bool) Part {
	return Part{partParam{v}}
}

func ParamInt(v int) Part {
	return Part{partParam{v}}
}

func ParamString(v string) Part {
	return Part{partParam{v}}
}

func Param(v interface{}) Part {
	return Part{partParam{v}}
}

func ParamBools(vs []bool) []Part {
	ps := make([]Part, 0, len(vs))
	for _, v := range vs {
		ps = append(ps, ParamBool(v))
	}
	return ps
}

func ParamInts(vs []int) []Part {
	ps := make([]Part, 0, len(vs))
	for _, v := range vs {
		ps = append(ps, ParamInt(v))
	}
	return ps
}

func ParamStrings(vs []string) []Part {
	ps := make([]Part, 0, len(vs))
	for _, v := range vs {
		ps = append(ps, ParamString(v))
	}
	return ps
}

func Params(vs []interface{}) []Part {
	ps := make([]Part, 0, len(vs))
	for _, v := range vs {
		ps = append(ps, Param(v))
	}
	return ps
}

type ValueBuilder parts

func NewValueBuilder() ValueBuilder {
	return (ValueBuilder)(parts{partString("VALUES ")})
}

func (vb *ValueBuilder) Append(values ...Part) {
	if len(*((*parts)(vb))) > 1 {
		*((*parts)(vb)) = append(*((*parts)(vb)), partString(", "))
	}
	*((*parts)(vb)) = append(*((*parts)(vb)), partByte('('))
	for i, v := range values {
		if i > 0 {
			*((*parts)(vb)) = append(*((*parts)(vb)), partString(", "))
		}
		*((*parts)(vb)) = append(*((*parts)(vb)), v)
	}
	*((*parts)(vb)) = append(*((*parts)(vb)), partByte(')'))
}

func (p Part) Build() (string, []interface{}) {
	return p.builder.Build()
}

func (vb ValueBuilder) Part() Part {
	return Part{vb}
}

func (vb ValueBuilder) Build() (string, []interface{}) {
	return (parts)(vb).Build()
}

func TableFields(name string, fields ...Part) Part {
	parts := parts{partString(name), partByte('(')}
	for i, p := range fields {
		if i != 0 {
			parts = append(parts, partString(", "))
		}
		parts = append(parts, p)
	}
	parts = append(parts, partByte(')'))
	return Part{parts}
}
