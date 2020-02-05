package query_builder

import (
	"reflect"
	"strings"
)

type QueryBuilder struct {
	tableName       string
	selects         []string
	whereConditions []map[string]string
	limit           bool
	offset          bool
	placeholder     int
}

const (
	Question = iota
	Named
)

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		placeholder: Question,
	}
}

func (qb *QueryBuilder) UseNamedPlaceholder() *QueryBuilder {
	qb.placeholder = Named
	return qb
}

func (qb *QueryBuilder) Table(tableName string) *QueryBuilder {
	qb.tableName = tableName
	return qb
}

func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	for _, column := range columns {
		qb.selects = append(qb.selects, column)
	}
	return qb
}

// db・tableタグを見て、FieldをSelect対象としてSet
func (qb *QueryBuilder) Model(model interface{}) *QueryBuilder {
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
			qb.selects = append(qb.selects, dbTag)
		}
	}
	return qb
}

func (qb *QueryBuilder) Where(column, operator string, bind ...string) *QueryBuilder {
	var bd string
	if len(bind) == 0 {
		bd = column
	} else {
		bd = bind[0]
	}
	qb.whereConditions = append(qb.whereConditions, map[string]string{
		"column":   column,
		"operator": operator,
		"bind":     bd,
	})
	return qb
}

func (qb *QueryBuilder) Limit() *QueryBuilder {
	qb.limit = true
	return qb
}

func (qb *QueryBuilder) Offset() *QueryBuilder {
	qb.offset = true
	return qb
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

	if len(qb.whereConditions) > 0 {
		q += "WHERE "
		for _, condition := range qb.whereConditions {
			if qb.placeholder == Named {
				q += condition["column"] + " " + condition["operator"] + " :" + condition["bind"] + " AND "
			} else {
				q += condition["column"] + " " + condition["operator"] + " ? AND "
			}

		}
		q = strings.TrimRight(q, "AND ")
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
