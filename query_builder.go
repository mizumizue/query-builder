package query_builder

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	SubQueryEmptyErr           = fmt.Errorf("subQuery is required. this should be not empty")
	SubQueryReturnMultiRowsErr = fmt.Errorf("subQuery returns multi rows. set limit and specify single row")
	UnspecifiedColumnErr       = fmt.Errorf("subQuery column should be specified")
)

type queryBuilder struct {
	query           []string
	tableName       string
	columns         []string
	whereConditions []map[string]string
	placeholderType int
	argNum          int
}

func (builder *queryBuilder) placeholder(placeholderType int) *queryBuilder {
	copied := builder.copy()
	copied.placeholderType = placeholderType
	return copied
}

func (builder *queryBuilder) table(tableName string) *queryBuilder {
	copied := builder.copy()
	copied.tableName = tableName
	return copied
}

func (builder *queryBuilder) column(columns ...string) *queryBuilder {
	copied := builder.copy()
	for _, column := range columns {
		copied.columns = append(copied.columns, column)
	}
	return copied
}

// db・tableタグを見て、FieldをSelect対象としてSet
func (builder *queryBuilder) model(model interface{}) *queryBuilder {
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

func (builder *queryBuilder) where(column, operator string, bind ...string) *queryBuilder {
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

// use in Operator and Placeholder, if bind is empty, IN(:{column}1, :{column}2, :{column}3...})
// use in Operator and Placeholder, if bind passed, IN(:{bind}1, :{bind}2, :{bind}3...})
func (builder *queryBuilder) whereIn(column string, listLength int, bind ...string) *queryBuilder {
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

func (builder *queryBuilder) whereNotIn(column string, listLength int, bind ...string) *queryBuilder {
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

func (builder *queryBuilder) whereMultiByStruct(targetTag string, src interface{}) *queryBuilder {
	copied := builder.copy()
	searchMap := builder.buildBindMap(targetTag, src)
	for _, info := range searchMap {
		op := getOperatorFromTag(info["operator"])
		if op == "" {
			continue
		}
		copied = copied.where(info["target"], op, info["bind"])
	}
	return copied
}

func (builder *queryBuilder) whereSubQuery(column, operator string, subQueryBuilder *SelectQueryBuilder) *queryBuilder {
	copied := builder.copy()

	if subQueryBuilder == nil {
		panic(SubQueryEmptyErr)
	}

	if len(subQueryBuilder.columns) == 0 || len(subQueryBuilder.columns) > 1 {
		panic(UnspecifiedColumnErr)
	}

	if subQueryBuilder.limit["use"] != nil && !subQueryBuilder.limit["use"].(bool) {
		panic(SubQueryReturnMultiRowsErr)
	}

	copied.whereConditions = append(copied.whereConditions, map[string]string{
		"column":   column,
		"operator": operator,
		"subQuery": strings.TrimRight(subQueryBuilder.Build(), ";"),
	})
	return copied
}

func (builder *queryBuilder) copy() *queryBuilder {
	return &queryBuilder{
		query:           builder.query,
		tableName:       builder.tableName,
		columns:         builder.columns,
		whereConditions: builder.whereConditions,
		placeholderType: builder.placeholderType,
	}
}

func (builder *queryBuilder) getWhereParagraphs() []string {
	paragraph := make([]string, 0, 0)
	paragraph = append(paragraph, "WHERE")

	format := "%s %s %s AND"
	for index, condition := range builder.whereConditions {
		if index == len(builder.whereConditions)-1 {
			format = strings.TrimRight(format, " AND")
		}

		if condition["subQuery"] != "" {
			paragraph = append(paragraph, fmt.Sprintf(
				"%s %s (%s)",
				condition["column"],
				condition["operator"],
				condition["subQuery"],
			))
			continue
		}

		bind := "?"
		if builder.placeholderType == Named {
			bind = ":" + condition["bind"]
		}

		if condition["operator"] == In || condition["operator"] == NotIn {
			listLength, _ := strconv.Atoi(condition["listLength"])
			listBind := builder.buildListBind(bind, listLength)
			paragraph = append(paragraph, fmt.Sprintf(format, condition["column"], condition["operator"], listBind))
			continue
		}

		if builder.placeholderType == DollarNumber {
			paragraph = append(paragraph, fmt.Sprintf(format, condition["column"], condition["operator"], "$"+strconv.Itoa(builder.argNum+1)))
			builder.argNum += 1
		} else {
			paragraph = append(paragraph, fmt.Sprintf(format, condition["column"], condition["operator"], bind))
		}
	}
	return paragraph
}

func (builder *queryBuilder) buildListBind(bind string, listLength int) string {
	format := "(%s)"
	list := make([]string, 0, listLength)
	for i := 0; i < listLength; i++ {
		if builder.placeholderType == Named {
			list = append(list, bind+strconv.Itoa(i+1))
			continue
		}
		if builder.placeholderType == DollarNumber {
			list = append(list, "$"+strconv.Itoa(builder.argNum+1))
			builder.argNum += 1
			continue
		}
		list = append(list, bind)
	}
	return fmt.Sprintf(format, strings.Join(list, ", "))
}

func getOperatorFromTag(tag string) string {
	switch tag {
	case "eq":
		return Equal
	case "gt":
		return GraterThan
	case "gte":
		return GraterThanEqual
	case "lt":
		return LessThan
	case "lte":
		return LessThanEqual
	case "ne":
		return NotEqual
	case "like":
		return Like
	case "not-like":
		return NotLike
	case "is-null":
		return IsNull
	case "not-null":
		return IsNotNull
	default:
		return ""
	}
}

func (builder *queryBuilder) buildBindMap(targetTag string, src interface{}) []map[string]string {
	t, v := builder.getReflectTypeAndValue(src)
	bindMap := make(map[string]map[string]string)
	dic := make([]string, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if fieldValue.Type().Kind() == reflect.Ptr && v.Field(i).IsNil() {
			continue
		}

		dbTag, bindTag, operatorTag := field.Tag.Get(DBTag), field.Tag.Get(targetTag), field.Tag.Get(OperatorTag)
		if dbTag == "" || bindTag == "" || operatorTag == "" {
			continue
		}

		dic = append(dic, bindTag)
		bindMap[bindTag] = map[string]string{
			"target":   dbTag,
			"operator": operatorTag,
		}
	}

	sortedByFieldNumber := make([]map[string]string, 0, len(dic))
	for _, key := range dic {
		bindMap[key]["bind"] = key
		sortedByFieldNumber = append(sortedByFieldNumber, bindMap[key])
	}
	return sortedByFieldNumber
}

func (builder *queryBuilder) getReflectTypeAndValue(src interface{}) (reflect.Type, reflect.Value) {
	t := reflect.TypeOf(src)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	v := reflect.ValueOf(src)
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return t, v
}
