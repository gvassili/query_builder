# query_builder [![GoDoc](https://godoc.org/github.com/gvassili/query_builder?status.svg)](https://godoc.org/github.com/gvassili/query_builder)

A golang MariaDB query builder


> :warning: **work in progress**

Todo list
- [ ] Go documentation
- [X] Select query
- [ ] Insert query
- [ ] Update query
- [ ] Delete query
- [ ] Write missing unit tests

## Installation
``` sh
go get github.com/gvassili/query_builder
```

## Example
``` go
package main

import (
    "fmt"
    qb "github.com/gvassili/query_builder"
)

func main() {
	query := qb.NewQuery()
	query.From(qb.Table("table").As("t")).
		Select(qb.Field("t.field1")).
		InnerJoin(qb.Table("other_table").As("ot"), qb.Field("ot.id").Eq(Field("t.other_id"))).
		Where(qb.Field("ts.Field2").Eq(qb.ParamInt(123)))
	qs, vs := query.Build()
	fmt.Printf("query string: '%s'\nvalues: %+v\n", qs, vs)
}
```
output:
```
query string: 'SELECT t.field1 FROM table AS t INNER JOIN other_table AS ot ON ot.id = t.other_id WHERE ts.Field2 = ?'
values: [123]
```