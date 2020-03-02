package query_builder

import (
	"fmt"
	"strings"
)

type UpdateQueryBuilder struct {
	*queryBuilder
}

func NewUpdateQueryBuilder() *UpdateQueryBuilder {
	builder := &UpdateQueryBuilder{}
	builder.queryBuilder = newQueryBuilder()
	builder.placeholderType = Question
	return builder
}

func (builder *UpdateQueryBuilder) copy() *UpdateQueryBuilder {
	return &UpdateQueryBuilder{
		builder.queryBuilder.copy(),
	}
}

func (builder *UpdateQueryBuilder) Placeholder(placeholderType int) *UpdateQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.placeholder(placeholderType)
	return copied
}

func (builder *UpdateQueryBuilder) Table(tableName string) *UpdateQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.table(tableName)
	return copied
}

func (builder *UpdateQueryBuilder) Model(src interface{}, notIgnoreZeroValue ...bool) *UpdateQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.model(src, notIgnoreZeroValue...)
	return copied
}

func (builder *UpdateQueryBuilder) Column(columns ...string) *UpdateQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.column(columns...)
	return copied
}

func (builder *UpdateQueryBuilder) Where(column, operator string, bind ...string) *UpdateQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.where(column, operator, bind...)
	return copied
}

func (builder *UpdateQueryBuilder) WhereIn(column string, listLength int, bind ...string) *UpdateQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.whereIn(column, listLength, bind...)
	return copied
}

func (builder *UpdateQueryBuilder) WhereNotIn(column string, listLength int, bind ...string) *UpdateQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.whereNotIn(column, listLength, bind...)
	return copied
}

func (builder *UpdateQueryBuilder) Build() string {
	if builder.tableName == "" {
		panic("target table is empty!!!")
	}

	if len(builder.columns) == 0 {
		panic("target columns is empty!!!")
	}

	copied := builder.copy()
	columns := builder.columns

	copied.query = append(copied.query, "UPDATE", builder.tableName)
	copied.query = append(copied.query, builder.getSetParagraphs(columns...))

	if len(builder.whereConditions) > 0 {
		copied.query = append(copied.query, builder.getWhereParagraphs()...)
	}

	return strings.TrimRight(strings.Join(copied.query, " "), "") + ";"
}

func (builder *UpdateQueryBuilder) getSetParagraphs(columns ...string) string {
	setContents := make([]string, 0, len(columns))
	bind := "?"
	format := "%s = %s"
	for _, column := range columns {
		if builder.placeholderType == Named {
			bind = ":" + column
		}
		setContents = append(setContents, fmt.Sprintf(format, column, bind))
	}
	return fmt.Sprintf("SET %s", strings.Join(setContents, ", "))
}
