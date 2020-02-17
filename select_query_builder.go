package query_builder

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type SelectQueryBuilder struct {
	joins         []map[string]interface{}
	groupByColumn string
	order         map[string]string
	limit         map[string]interface{}
	offset        map[string]interface{}
	*queryBuilder
	subQueryBuilder *queryBuilder
}

func NewSelectQueryBuilder() *SelectQueryBuilder {
	builder := &SelectQueryBuilder{}
	builder.queryBuilder = &queryBuilder{}
	builder.placeholderType = Question
	return builder
}

func (builder *SelectQueryBuilder) copy() *SelectQueryBuilder {
	return &SelectQueryBuilder{
		builder.joins,
		builder.groupByColumn,
		builder.order,
		builder.limit,
		builder.offset,
		builder.queryBuilder.copy(),
		nil,
	}
}

// Default placeholder is ?
func (builder *SelectQueryBuilder) Placeholder(placeholderType int) *SelectQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.placeholder(placeholderType)
	return copied
}

func (builder *SelectQueryBuilder) Table(tableName string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.table(tableName)
	return copied
}

func (builder *SelectQueryBuilder) Model(src interface{}) *SelectQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.model(src)
	return copied
}

func (builder *SelectQueryBuilder) Column(columns ...string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.column(columns...)
	return copied
}

func (builder *SelectQueryBuilder) Join(joinType, joinTable string, onOriginFields, onTargetFields []string, otherTable ...string) *SelectQueryBuilder {
	copied := builder.copy()

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

func (builder *SelectQueryBuilder) Where(column, operator string, bind ...string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.where(column, operator, bind...)
	return copied
}

func (builder *SelectQueryBuilder) WhereIn(column string, listLength int, bind ...string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.whereIn(column, listLength, bind...)
	return copied
}

func (builder *SelectQueryBuilder) WhereNotIn(column string, listLength int, bind ...string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.whereNotIn(column, listLength, bind...)
	return copied
}

func (builder *SelectQueryBuilder) WhereSubQuery(column, operator string, subQueryBuilder *SelectQueryBuilder) *SelectQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.whereSubQuery(column, operator, subQueryBuilder)
	return copied
}

func (builder *SelectQueryBuilder) WhereMultiByStruct(src interface{}) *SelectQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.whereMultiByStruct(src)
	return copied
}

func (builder *SelectQueryBuilder) GroupBy(column string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.groupByColumn = column
	return copied
}

// ex. OrderBy("created, user_id", Asc)
func (builder *SelectQueryBuilder) OrderBy(columns, order string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.order = map[string]string{
		"columns": columns,
		"order":   order,
	}
	return copied
}

func (builder *SelectQueryBuilder) Limit(bind ...string) *SelectQueryBuilder {
	bd := "limit"
	if len(bind) != 0 {
		bd = bind[0]
	}
	copied := builder.copy()
	copied.limit = map[string]interface{}{
		"use":  true,
		"bind": bd,
	}
	return copied
}

func (builder *SelectQueryBuilder) Offset(bind ...string) *SelectQueryBuilder {
	bd := "offset"
	if len(bind) != 0 {
		bd = bind[0]
	}
	copied := builder.copy()
	copied.offset = map[string]interface{}{
		"use":  true,
		"bind": bd,
	}
	return copied
}

func (builder *SelectQueryBuilder) Build() string {
	if builder.tableName == "" {
		panic("target table is empty!!!")
	}

	copied := builder.copy()
	columns := builder.columns
	copied.query = append(copied.query, builder.getSelectParagraphs(builder.tableName, columns)...)

	if len(builder.joins) > 0 {
		copied.query = append(copied.query, builder.getJoinParagraphs(builder.tableName)...)
	}

	if len(builder.whereConditions) > 0 {
		copied.query = append(copied.query, builder.getWhereParagraphs()...)
	}

	if builder.groupByColumn != "" {
		copied.query = append(copied.query, builder.getGroupByParagraph())
	}

	if len(builder.order) > 0 {
		copied.query = append(copied.query, builder.getOrderParagraph())
	}

	if builder.limit["use"] != nil && builder.limit["use"].(bool) {
		copied.query = append(copied.query, builder.getLimitParagraph(builder.placeholderType))
	}

	if builder.offset["use"] != nil && builder.offset["use"].(bool) {
		copied.query = append(copied.query, builder.getOffsetParagraph(builder.placeholderType))
	}

	return strings.TrimRight(strings.Join(copied.query, " "), "") + ";"
}

func (builder *SelectQueryBuilder) getSelectParagraphs(tableName string, columns []string) []string {
	paragraph := make([]string, 0, 0)
	paragraph = append(paragraph, "SELECT")

	if len(columns) == 0 {
		paragraph = append(paragraph, tableName+".*")
		paragraph = append(paragraph, "FROM", tableName)
		return paragraph
	}

	format := "%s.%s,"
	for index, column := range columns {
		if index == len(columns)-1 {
			format = strings.TrimRight(format, ",")
		}

		table := tableName
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
	return append(paragraph, "FROM", tableName)
}

func (builder *SelectQueryBuilder) getJoinParagraphs(tableName string) []string {
	paragraph := make([]string, 0, 0)
	for _, join := range builder.joins {
		joinOrginTableBase := tableName
		if join["otherTable"] != nil {
			joinOrginTableBase = join["otherTable"].(string)
		}

		paragraphFormer := fmt.Sprintf("%s %s ON ", join["type"], join["table"])
		paragraphLastHalf := builder.buildOnParagraph(
			joinOrginTableBase,
			join["table"].(string),
			join["onOriginFields"].([]string),
			join["onTargetFields"].([]string),
		)
		paragraph = append(paragraph, paragraphFormer+paragraphLastHalf)
	}
	return paragraph
}

func (builder *SelectQueryBuilder) buildOnParagraph(
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

func (builder *SelectQueryBuilder) getGroupByParagraph() string {
	return fmt.Sprintf("GROUP BY %s", builder.groupByColumn)
}

func (builder *SelectQueryBuilder) getOrderParagraph() string {
	return fmt.Sprintf("ORDER BY %s %s", builder.order["columns"], builder.order["order"])
}

func (builder *SelectQueryBuilder) getLimitParagraph(placeholder int) string {
	bind := "?"
	if placeholder == Named {
		bind = ":" + builder.limit["bind"].(string)
	}
	if placeholder == DollarNumber {
		bind = "$" + strconv.Itoa(builder.argNum+1)
		builder.argNum += 1
	}
	return fmt.Sprintf("LIMIT %s", bind)
}

func (builder *SelectQueryBuilder) getOffsetParagraph(placeholder int) string {
	if builder.limit == nil {
		panic("offset is limit required")
	}

	bind := "?"
	if placeholder == Named {
		bind = ":" + builder.offset["bind"].(string)
	}
	if placeholder == DollarNumber {
		bind = "$" + strconv.Itoa(builder.argNum+1)
		builder.argNum += 1
	}
	return fmt.Sprintf("OFFSET %s", bind)
}
