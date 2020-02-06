package query_builder

import (
	"reflect"
	"strings"
)

type QueryBuilder struct {
	tableName       string
	selects         []string
	joins           []map[string]string
	whereConditions []map[string]string
	limit           bool
	offset          bool
	placeholder     int
}

const (
	Question = iota
	Named
)

const (
	LeftJoin  = "LEFT JOIN"
	RightJoin = "RIGHT JOIN"
)

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		placeholder: Question,
	}
}

func (qb *QueryBuilder) UseNamedPlaceholder() *QueryBuilder {
	copied := qb.copy()
	copied.placeholder = Named
	return copied
}

func (qb *QueryBuilder) Table(tableName string) *QueryBuilder {
	copied := qb.copy()
	copied.tableName = tableName
	return copied
}

func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	copied := qb.copy()
	for _, column := range columns {
		copied.selects = append(copied.selects, column)
	}
	return copied
}

// db・tableタグを見て、FieldをSelect対象としてSet
func (qb *QueryBuilder) Model(model interface{}) *QueryBuilder {
	copied := qb.copy()
	t := reflect.TypeOf(model)

	if t.Kind() == reflect.Ptr {
		panic("model should be not pointer value")
	}
	if t.Kind() != reflect.Struct {
		panic("model should be struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		tableTag := field.Tag.Get("table")
		if dbTag != "" && tableTag == qb.tableName {
			copied.selects = append(copied.selects, dbTag)
		}
	}
	return copied
}

func (qb *QueryBuilder) Join(joinType, joinTable, onField string, otherTable ...string) *QueryBuilder {
	copied := qb.copy()

	m := make(map[string]string)
	m["type"] = joinType
	m["table"] = joinTable
	m["onField"] = onField

	if len(otherTable) > 0 && otherTable[0] != "" {
		m["otherTable"] = otherTable[0]
	}

	copied.joins = append(copied.joins, m)
	return copied
}

func (qb *QueryBuilder) Where(column, operator string, bind ...string) *QueryBuilder {
	copied := qb.copy()

	var bd string
	if len(bind) == 0 {
		bd = column
	} else {
		bd = bind[0]
	}
	copied.whereConditions = append(copied.whereConditions, map[string]string{
		"column":   column,
		"operator": operator,
		"bind":     bd,
	})
	return copied
}

func (qb *QueryBuilder) Limit() *QueryBuilder {
	copied := qb.copy()
	copied.limit = true
	return copied
}

func (qb *QueryBuilder) Offset() *QueryBuilder {
	copied := qb.copy()
	copied.offset = true
	return copied
}

func (qb *QueryBuilder) Build() string {
	q := "SELECT "
	if len(qb.selects) == 0 {
		q += qb.tableName + ".*"
	} else {
		for _, column := range qb.selects {
			q += qb.tableName + "." + column + ", "
		}
		q = strings.TrimRight(q, ", ")
	}
	q += " FROM " + qb.tableName + " "

	if len(qb.joins) > 0 {
		for _, join := range qb.joins {
			if join["otherTable"] != "" {
				q += join["type"] + " " + join["table"] + " ON " + join["otherTable"] + "." + join["onField"] + " = " + join["table"] + "." + join["onField"] + " "
			} else {
				q += join["type"] + " " + join["table"] + " ON " + qb.tableName + "." + join["onField"] + " = " + join["table"] + "." + join["onField"] + " "
			}

		}
	}

	if len(qb.whereConditions) > 0 {
		q += "WHERE "
		for _, condition := range qb.whereConditions {
			if qb.placeholder == Named {
				q += condition["column"] + " " + condition["operator"] + " :" + condition["bind"] + " AND "
			} else {
				q += condition["column"] + " " + condition["operator"] + " ? AND "
			}

		}
		q = strings.TrimRight(q, " AND") + " "
	}

	if qb.limit {
		if qb.placeholder == Named {
			q += "LIMIT :limit "
		} else {
			q += "LIMIT ? "
		}
	}

	if qb.offset {
		if qb.placeholder == Named {
			q += "OFFSET :offset "
		} else {
			q += "OFFSET ? "
		}
	}

	return strings.TrimRight(q, " ") + ";"
}

func (qb *QueryBuilder) copy() *QueryBuilder {
	return &QueryBuilder{
		tableName:       qb.tableName,
		selects:         qb.selects,
		joins:           qb.joins,
		whereConditions: qb.whereConditions,
		limit:           qb.limit,
		offset:          qb.offset,
		placeholder:     qb.placeholder,
	}
}
