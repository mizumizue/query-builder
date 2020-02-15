package query_builder

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/trewanek/query-builder/object_parser"
)

type SelectQueryBuilder struct {
	*selectQueryBuilder
	*commonQueryBuilder
}

func (builder *SelectQueryBuilder) copy() *SelectQueryBuilder {
	return &SelectQueryBuilder{
		builder.selectQueryBuilder.copy(),
		builder.commonQueryBuilder.copy(),
	}
}

func (builder *SelectQueryBuilder) UseNamedPlaceholder() *SelectQueryBuilder {
	copied := builder.copy()
	copied.commonQueryBuilder = builder.useNamedPlaceholder()
	return copied
}

func (builder *SelectQueryBuilder) Table(tableName string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.commonQueryBuilder = builder.table(tableName)
	return copied
}

func (builder *SelectQueryBuilder) Model(src interface{}) *SelectQueryBuilder {
	copied := builder.copy()
	copied.commonQueryBuilder = builder.model(src)
	return copied
}

func (builder *SelectQueryBuilder) Column(columns ...string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.commonQueryBuilder = builder.column(columns...)
	return copied
}

func (builder *SelectQueryBuilder) Join(joinType, joinTable string, onOriginFields, onTargetFields []string, otherTable ...string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.selectQueryBuilder = builder.join(joinType, joinTable, onOriginFields, onTargetFields, otherTable...)
	return copied
}

func (builder *SelectQueryBuilder) Where(column, operator string, bind ...string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.commonQueryBuilder = builder.where(column, operator, bind...)
	return copied
}

func (builder *SelectQueryBuilder) WhereIn(column string, listLength int, bind ...string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.commonQueryBuilder = builder.whereIn(column, listLength, bind...)
	return copied
}

func (builder *SelectQueryBuilder) WhereNotIn(column string, listLength int, bind ...string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.commonQueryBuilder = builder.whereNotIn(column, listLength, bind...)
	return copied
}

func (builder *SelectQueryBuilder) WhereMultiByStruct(src interface{}) *SelectQueryBuilder {
	copied := builder.copy()
	copied.commonQueryBuilder = builder.whereMultiByStruct(src)
	return copied
}

func (builder *SelectQueryBuilder) GroupBy(column string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.selectQueryBuilder = builder.groupBy(column)
	return copied
}

func (builder *SelectQueryBuilder) OrderBy(columns, order string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.selectQueryBuilder = builder.orderBy(columns, order)
	return copied
}

func (builder *SelectQueryBuilder) Limit(bind ...string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.selectQueryBuilder = builder.limit(bind...)
	return copied
}

func (builder *SelectQueryBuilder) Offset(bind ...string) *SelectQueryBuilder {
	copied := builder.copy()
	copied.selectQueryBuilder = builder.offset(bind...)
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

	if builder.groupByStr != "" {
		copied.query = append(copied.query, builder.getGroupByParagraph())
	}

	if len(builder.order) > 0 {
		copied.query = append(copied.query, builder.getOrderParagraph())
	}

	if builder.limitMap["use"] != nil && builder.limitMap["use"].(bool) {
		copied.query = append(copied.query, builder.getLimitParagraph(builder.placeholder))
	}

	if builder.offsetMap["use"] != nil && builder.offsetMap["use"].(bool) {
		copied.query = append(copied.query, builder.getOffsetParagraph(builder.placeholder))
	}

	return strings.TrimRight(strings.Join(copied.query, " "), "") + ";"
}

type selectQueryBuilder struct {
	joins      []map[string]interface{}
	groupByStr string
	order      map[string]string
	limitMap   map[string]interface{}
	offsetMap  map[string]interface{}
}

func (builder *selectQueryBuilder) copy() *selectQueryBuilder {
	return &selectQueryBuilder{
		joins:      builder.joins,
		groupByStr: builder.groupByStr,
		order:      builder.order,
		limitMap:   builder.limitMap,
		offsetMap:  builder.offsetMap,
	}
}

func (builder *selectQueryBuilder) getSelectParagraphs(tableName string, columns []string) []string {
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

		paragraph = append(paragraph, fmt.Sprintf(
			format,
			table,
			selectColumn,
		))
	}
	return append(paragraph, "FROM", tableName)
}

func (builder *selectQueryBuilder) getJoinParagraphs(tableName string) []string {
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

func (builder *selectQueryBuilder) buildOnParagraph(
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

func (builder *selectQueryBuilder) getGroupByParagraph() string {
	return fmt.Sprintf("GROUP BY %s", builder.groupByStr)
}

func (builder *selectQueryBuilder) getOrderParagraph() string {
	return fmt.Sprintf("ORDER BY %s %s", builder.order["columns"], builder.order["order"])
}

func (builder *selectQueryBuilder) getLimitParagraph(placeholder int) string {
	bind := "?"
	if placeholder == Named {
		bind = ":" + builder.limitMap["bind"].(string)
	}
	return fmt.Sprintf("LIMIT %s", bind)
}

func (builder *selectQueryBuilder) getOffsetParagraph(placeholder int) string {
	bind := "?"
	if placeholder == Named {
		bind = ":" + builder.offsetMap["bind"].(string)
	}
	return fmt.Sprintf("OFFSET %s", bind)
}

func (builder *commonQueryBuilder) copy() *commonQueryBuilder {
	return &commonQueryBuilder{
		query:           builder.query,
		tableName:       builder.tableName,
		columns:         builder.columns,
		whereConditions: builder.whereConditions,
		placeholder:     builder.placeholder,
	}
}

func (builder *commonQueryBuilder) getWhereParagraphs() []string {
	paragraph := make([]string, 0, 0)
	paragraph = append(paragraph, "WHERE")

	format := "%s %s %s AND"
	for index, condition := range builder.whereConditions {
		if index == len(builder.whereConditions)-1 {
			format = strings.TrimRight(format, " AND")
		}

		bind := "?"
		if builder.placeholder == Named {
			bind = ":" + condition["bind"]
		}

		if condition["operator"] == In || condition["operator"] == NotIn {
			listLength, _ := strconv.Atoi(condition["listLength"])
			listBind := builder.buildListBind(bind, listLength)
			paragraph = append(paragraph, fmt.Sprintf(format, condition["column"], condition["operator"], listBind))
			continue
		}

		paragraph = append(paragraph, fmt.Sprintf(format, condition["column"], condition["operator"], bind))
	}
	return paragraph
}

func (builder *commonQueryBuilder) buildListBind(bind string, listLength int) string {
	format := "(%s)"
	list := make([]string, 0, listLength)
	for i := 0; i < listLength; i++ {
		if builder.placeholder == Named {
			list = append(list, bind+strconv.Itoa(i+1))
		} else {
			list = append(list, bind)
		}
	}
	return fmt.Sprintf(format, strings.Join(list, ", "))
}

type InsertQueryBuilder struct {
	*insertQueryBuilder
	*commonQueryBuilder
}

type UpdateQueryBuilder struct {
	*updateQueryBuilder
	*commonQueryBuilder
}

type DeleteQueryBuilder struct {
	*deleteQueryBuilder
	*commonQueryBuilder
}

type ISelectQueryBuilder interface {
	Select() ISelectQueryBuilder
	Join(joinType, joinTable string, onOriginFields, onTargetFields []string, otherTable ...string) ISelectQueryBuilder
	GroupBy(column string) ISelectQueryBuilder
	OrderBy(columns, order string) ISelectQueryBuilder
	Limit(bind ...string) ISelectQueryBuilder
	Offset(bind ...string) ISelectQueryBuilder
	getGroupByParagraph() string
	getOrderParagraph() string
	getLimitParagraph() string
	getOffsetParagraph() string
}

type IQueryBuilder interface {
	UseNamedPlaceholder() IQueryBuilder
	Table() IQueryBuilder
	Model(model interface{}) IQueryBuilder
	Where(column, operator string, bind ...string) IQueryBuilder
	WhereIn(column string, listLength int, bind ...string) IQueryBuilder
	WhereNotIn(column string, listLength int, bind ...string) IQueryBuilder
	WhereMultiByStruct(src interface{}) IQueryBuilder
	copy() IQueryBuilder
}

type insertQueryBuilder struct {
}

type updateQueryBuilder struct {
}

type deleteQueryBuilder struct {
}

type commonQueryBuilder struct {
	query           []string
	tableName       string
	columns         []string
	whereConditions []map[string]string
	placeholder     int
}

func NewSelectQueryBuilder() *SelectQueryBuilder {
	qb := &SelectQueryBuilder{}
	qb.selectQueryBuilder = &selectQueryBuilder{}
	qb.commonQueryBuilder = &commonQueryBuilder{}
	qb.placeholder = Question
	return qb
}

func (builder *commonQueryBuilder) useNamedPlaceholder() *commonQueryBuilder {
	copied := builder.copy()
	copied.placeholder = Named
	return copied
}

func (builder *commonQueryBuilder) table(tableName string) *commonQueryBuilder {
	copied := builder.copy()
	copied.tableName = tableName
	return copied
}

func (builder *commonQueryBuilder) column(columns ...string) *commonQueryBuilder {
	copied := builder.copy()
	for _, column := range columns {
		copied.columns = append(copied.columns, column)
	}
	return copied
}

// db・tableタグを見て、FieldをSelect対象としてSet
func (builder *commonQueryBuilder) model(model interface{}) *commonQueryBuilder {
	copied := builder.copy()
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
		if dbTag != "" && tableTag == builder.tableName {
			copied.columns = append(copied.columns, dbTag)
		}
	}
	return copied
}

func (builder *selectQueryBuilder) join(joinType, joinTable string, onOriginFields, onTargetFields []string, otherTable ...string) *selectQueryBuilder {
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

func (builder *commonQueryBuilder) where(column, operator string, bind ...string) *commonQueryBuilder {
	copied := builder.copy()
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
func (builder *commonQueryBuilder) whereIn(column string, listLength int, bind ...string) *commonQueryBuilder {
	copied := builder.copy()
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

func (builder *commonQueryBuilder) whereNotIn(column string, listLength int, bind ...string) *commonQueryBuilder {
	copied := builder.copy()
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

func (builder *commonQueryBuilder) whereMultiByStruct(src interface{}) *commonQueryBuilder {
	copied := builder.copy()

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
		copied = copied.where(info["target"], op, info["bind"])
	}
	return copied
}

func (builder *selectQueryBuilder) groupBy(column string) *selectQueryBuilder {
	copied := builder.copy()
	copied.groupByStr = column
	return copied
}

// ex. OrderBy("created, user_id", Asc)
func (builder *selectQueryBuilder) orderBy(columns, order string) *selectQueryBuilder {
	copied := builder.copy()
	copied.order = map[string]string{
		"columns": columns,
		"order":   order,
	}
	return copied
}

func (builder *selectQueryBuilder) limit(bind ...string) *selectQueryBuilder {
	bd := "limitMap"
	if len(bind) != 0 {
		bd = bind[0]
	}
	copied := builder.copy()
	copied.limitMap = map[string]interface{}{
		"use":  true,
		"bind": bd,
	}
	return copied
}

func (builder *selectQueryBuilder) offset(bind ...string) *selectQueryBuilder {
	bd := "offsetMap"
	if len(bind) != 0 {
		bd = bind[0]
	}
	copied := builder.copy()
	copied.offsetMap = map[string]interface{}{
		"use":  true,
		"bind": bd,
	}
	return copied
}
