package query_builder

import (
	"fmt"
	"query-builder/parameter_parser"
	"query-builder/query_operator"
	"reflect"
	"strings"
)

type QueryBuilder struct {
	query           []string
	tableName       string
	selects         []string
	joins           []map[string]string
	whereConditions []map[string]string
	limit           map[string]interface{}
	offset          map[string]interface{}
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
	bd := column
	if len(bind) != 0 {
		bd = bind[0]
	}
	copied.whereConditions = append(copied.whereConditions, map[string]string{
		"column":   column,
		"operator": operator,
		"bind":     bd,
	})
	return copied
}

func (qb *QueryBuilder) WhereMultiByStruct(src interface{}) *QueryBuilder {
	copied := qb.copy()

	paramMap := parameter_parser.NewParameterParser(src).ParseBindMap()
	for _, info := range paramMap {
		var op string
		switch info["operator"] {
		case "eq":
			op = query_operator.Equal
		case "lt":
			op = query_operator.LessThan
		case "le":
			op = query_operator.LessEqual
		case "gt":
			op = query_operator.GraterThan
		case "ge":
			op = query_operator.GraterEqual
		case "not":
			op = query_operator.Not
		}
		copied = copied.Where(info["target"], op, info["bind"])
	}
	return copied
}

func (qb *QueryBuilder) Limit(bind ...string) *QueryBuilder {
	bd := "limit"
	if len(bind) != 0 {
		bd = bind[0]
	}
	copied := qb.copy()
	copied.limit = map[string]interface{}{
		"use":  true,
		"bind": bd,
	}
	return copied
}

func (qb *QueryBuilder) Offset(bind ...string) *QueryBuilder {
	bd := "offset"
	if len(bind) != 0 {
		bd = bind[0]
	}
	copied := qb.copy()
	copied.offset = map[string]interface{}{
		"use":  true,
		"bind": bd,
	}
	return copied
}

func (qb *QueryBuilder) Build() string {
	if qb.tableName == "" {
		panic("target table is empty!!!")
	}

	copied := qb.copy()
	copied.query = append(copied.query, qb.getSelectParagraphs()...)

	if len(qb.joins) > 0 {
		copied.query = append(copied.query, qb.getJoinParagraphs()...)
	}

	if len(qb.whereConditions) > 0 {
		copied.query = append(copied.query, qb.getWhereParagraphs()...)
	}

	if qb.limit["use"] != nil && qb.limit["use"].(bool) {
		copied.query = append(copied.query, qb.getLimitParagraph())
	}

	if qb.offset["use"] != nil && qb.offset["use"].(bool) {
		copied.query = append(copied.query, qb.getOffsetParagraph())
	}

	return strings.TrimRight(strings.Join(copied.query, " "), "") + ";"
}

func (qb *QueryBuilder) getSelectParagraphs() []string {
	paragraph := make([]string, 0, 0)
	paragraph = append(paragraph, "SELECT")

	if len(qb.selects) == 0 {
		paragraph = append(paragraph, qb.tableName+".*")
		paragraph = append(paragraph, "FROM", qb.tableName)
		return paragraph
	}

	format := "%s.%s,"
	for index, column := range qb.selects {
		if index == len(qb.selects)-1 {
			format = strings.TrimRight(format, ",")
		}
		paragraph = append(paragraph, fmt.Sprintf(
			format,
			qb.tableName,
			column,
		))
	}
	return append(paragraph, "FROM", qb.tableName)
}

func (qb *QueryBuilder) getJoinParagraphs() []string {
	paragraph := make([]string, 0, 0)
	for _, join := range qb.joins {
		joinBase := qb.tableName
		if join["otherTable"] != "" {
			joinBase = join["otherTable"]
		}

		paragraph = append(paragraph, fmt.Sprintf(
			"%s %s ON %s.%s = %s.%s",
			join["type"],
			join["table"],
			joinBase,
			join["onField"],
			join["table"],
			join["onField"],
		))
	}
	return paragraph
}

func (qb *QueryBuilder) getWhereParagraphs() []string {
	paragraph := make([]string, 0, 0)
	paragraph = append(paragraph, "WHERE")

	format := "%s %s %s AND"
	for index, condition := range qb.whereConditions {
		if index == len(qb.whereConditions)-1 {
			format = strings.TrimRight(format, " AND")
		}
		bind := "?"
		if qb.placeholder == Named {
			bind = ":" + condition["bind"]
		}
		paragraph = append(paragraph, fmt.Sprintf(format, condition["column"], condition["operator"], bind))
	}
	return paragraph
}

func (qb *QueryBuilder) getLimitParagraph() string {
	bind := "?"
	if qb.placeholder == Named {
		bind = ":" + qb.limit["bind"].(string)
	}
	return fmt.Sprintf("LIMIT %s", bind)
}

func (qb *QueryBuilder) getOffsetParagraph() string {
	bind := "?"
	if qb.placeholder == Named {
		bind = ":" + qb.offset["bind"].(string)
	}
	return fmt.Sprintf("OFFSET %s", bind)
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
