package query_builder

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/trewanek/query-builder/object_parser"
)

var (
	SubQueryEmptyErr     = fmt.Errorf("subQuery is required. this should be not empty")
	UnspecifiedColumnErr = fmt.Errorf("subQuery column should be specified")
)

type queryBuilder struct {
	query           []string
	tableName       string
	columns         []string
	whereConditions []map[string]string
	placeholderType int
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

func (builder *queryBuilder) whereMultiByStruct(src interface{}) *queryBuilder {
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

func (builder *queryBuilder) whereSubQuery(column, operator string, subQueryBuilder *SelectQueryBuilder) *queryBuilder {
	copied := builder.copy()

	if subQueryBuilder == nil {
		panic(SubQueryEmptyErr)
	}

	if len(subQueryBuilder.columns) == 0 || len(subQueryBuilder.columns) > 1 {
		panic(UnspecifiedColumnErr)
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

		paragraph = append(paragraph, fmt.Sprintf(format, condition["column"], condition["operator"], bind))
	}
	return paragraph
}

func (builder *queryBuilder) buildListBind(bind string, listLength int) string {
	format := "(%s)"
	list := make([]string, 0, listLength)
	for i := 0; i < listLength; i++ {
		if builder.placeholderType == Named {
			list = append(list, bind+strconv.Itoa(i+1))
		} else {
			list = append(list, bind)
		}
	}
	return fmt.Sprintf(format, strings.Join(list, ", "))
}
