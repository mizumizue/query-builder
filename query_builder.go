package query_builder

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/trewanek/query-builder/object_parser"
)

type QueryBuilder struct {
	query           []string
	tableName       string
	selects         []string
	joins           []map[string]interface{}
	whereConditions []map[string]string
	groupBy         string
	order           map[string]string
	limit           map[string]interface{}
	offset          map[string]interface{}
	placeholder     int
}

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

func (qb *QueryBuilder) Join(joinType, joinTable string, onOriginFields, onTargetFields []string, otherTable ...string) *QueryBuilder {
	copied := qb.copy()

	if len(onOriginFields) != len(onTargetFields) {
		panic("origin fields and target fields need to be same length")
	}

	m := make(map[string]interface{})
	m["type"] = joinType
	m["table"] = joinTable
	m["onOriginFields"] = onOriginFields
	m["onTargetFields"] = onTargetFields

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

// use in Operator and UseNamedPlaceholder, if bind is empty, IN(:{column}1, :{column}2, :{column}3...})
// use in Operator and UseNamedPlaceholder, if bind passed, IN(:{bind}1, :{bind}2, :{bind}3...})
func (qb *QueryBuilder) WhereIn(column string, listLength int, bind ...string) *QueryBuilder {
	copied := qb.copy()
	bd := column
	if len(bind) != 0 {
		bd = bind[0]
	}
	copied.whereConditions = append(copied.whereConditions, map[string]string{
		"column":     column,
		"listLength": strconv.Itoa(listLength),
		"operator":   In,
		"bind":       bd,
	})
	return copied
}

func (qb *QueryBuilder) WhereNotIn(column string, listLength int, bind ...string) *QueryBuilder {
	copied := qb.copy()
	bd := column
	if len(bind) != 0 {
		bd = bind[0]
	}
	copied.whereConditions = append(copied.whereConditions, map[string]string{
		"column":     column,
		"listLength": strconv.Itoa(listLength),
		"operator":   NotIn,
		"bind":       bd,
	})
	return copied
}

func (qb *QueryBuilder) WhereMultiByStruct(src interface{}) *QueryBuilder {
	copied := qb.copy()

	searchMap := object_parser.NewObjectParser(src).SearchBindMap()
	for _, info := range searchMap {
		var op string
		switch info["operator"] {
		case "eq":
			op = Equal
		case "lt":
			op = LessThan
		case "le":
			op = LessEqual
		case "gt":
			op = GraterThan
		case "ge":
			op = GraterEqual
		case "not":
			op = Not
		}
		copied = copied.Where(info["target"], op, info["bind"])
	}
	return copied
}

func (qb *QueryBuilder) GroupBy(column string) *QueryBuilder {
	copied := qb.copy()
	copied.groupBy = column
	return copied
}

// ex. OrderBy("created, user_id", Asc)
func (qb *QueryBuilder) OrderBy(columns, order string) *QueryBuilder {
	copied := qb.copy()
	copied.order = map[string]string{
		"columns": columns,
		"order":   order,
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

	if qb.groupBy != "" {
		copied.query = append(copied.query, qb.getGroupByParagraph())
	}

	if len(qb.order) > 0 {
		copied.query = append(copied.query, qb.getOrderParagraph())
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

		table := qb.tableName
		selectColumn := column
		split := strings.Split(column, ".")
		if len(split) > 1 {
			table = split[0]
			selectColumn = split[1]
		}

		if regexp.MustCompile(`^.*\(.*\)`).Match([]byte(column)) {
			paragraph = append(paragraph, fmt.Sprintf("%s,", selectColumn))
			continue
		}

		paragraph = append(paragraph, fmt.Sprintf(
			format,
			table,
			selectColumn,
		))
	}
	return append(paragraph, "FROM", qb.tableName)
}

func (qb *QueryBuilder) getJoinParagraphs() []string {
	paragraph := make([]string, 0, 0)
	for _, join := range qb.joins {
		joinOrginTableBase := qb.tableName
		if join["otherTable"] != nil {
			joinOrginTableBase = join["otherTable"].(string)
		}

		paragraphFormer := fmt.Sprintf("%s %s ON ", join["type"], join["table"])
		paragraphLastHalf := qb.buildOnParagraph(
			joinOrginTableBase,
			join["table"].(string),
			join["onOriginFields"].([]string),
			join["onTargetFields"].([]string),
		)
		paragraph = append(paragraph, paragraphFormer+paragraphLastHalf)
	}
	return paragraph
}

func (qb *QueryBuilder) buildOnParagraph(
	joinOriginTable,
	joinTargetTable string,
	originFields,
	targetFields []string,
) string {
	onParagraph := make([]string, 0, 0)
	for index, originField := range originFields {
		onParagraph = append(onParagraph, fmt.Sprintf(
			"%s.%s = %s.%s",
			joinOriginTable,
			originField,
			joinTargetTable,
			targetFields[index],
		))
	}
	return strings.Join(onParagraph, " AND ")
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

		if condition["operator"] == In || condition["operator"] == NotIn {
			listLength, _ := strconv.Atoi(condition["listLength"])
			listBind := qb.buildListBind(bind, listLength)
			paragraph = append(paragraph, fmt.Sprintf(format, condition["column"], condition["operator"], listBind))
			continue
		}

		paragraph = append(paragraph, fmt.Sprintf(format, condition["column"], condition["operator"], bind))
	}
	return paragraph
}

func (qb *QueryBuilder) buildListBind(bind string, listLength int) string {
	format := "(%s)"
	list := make([]string, 0, listLength)
	for i := 0; i < listLength; i++ {
		if qb.placeholder == Named {
			list = append(list, bind+strconv.Itoa(i+1))
		} else {
			list = append(list, bind)
		}
	}
	return fmt.Sprintf(format, strings.Join(list, ", "))
}

func (qb *QueryBuilder) getGroupByParagraph() string {
	return fmt.Sprintf("GROUP BY %s", qb.groupBy)
}

func (qb *QueryBuilder) getOrderParagraph() string {
	return fmt.Sprintf("ORDER BY %s %s", qb.order["columns"], qb.order["order"])
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
		groupBy:         qb.groupBy,
		order:           qb.order,
		limit:           qb.limit,
		offset:          qb.offset,
		placeholder:     qb.placeholder,
	}
}
