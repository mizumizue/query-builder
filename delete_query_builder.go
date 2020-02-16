package query_builder

import "strings"

type DeleteQueryBuilder struct {
	*queryBuilder
}

func NewDeleteQueryBuilder() *DeleteQueryBuilder {
	builder := &DeleteQueryBuilder{}
	builder.queryBuilder = &queryBuilder{}
	builder.placeholderType = Question
	return builder
}

func (builder *DeleteQueryBuilder) copy() *DeleteQueryBuilder {
	return &DeleteQueryBuilder{
		builder.queryBuilder.copy(),
	}
}

func (builder *DeleteQueryBuilder) Placeholder(placeholderType int) *DeleteQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.placeholder(placeholderType)
	return copied
}

func (builder *DeleteQueryBuilder) Table(tableName string) *DeleteQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.table(tableName)
	return copied
}

func (builder *DeleteQueryBuilder) Where(column, operator string, bind ...string) *DeleteQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.where(column, operator, bind...)
	return copied
}

func (builder *DeleteQueryBuilder) WhereIn(column string, listLength int, bind ...string) *DeleteQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.whereIn(column, listLength, bind...)
	return copied
}

func (builder *DeleteQueryBuilder) WhereNotIn(column string, listLength int, bind ...string) *DeleteQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.whereNotIn(column, listLength, bind...)
	return copied
}

func (builder *DeleteQueryBuilder) Build() string {
	if builder.tableName == "" {
		panic("target table is empty!!!")
	}

	copied := builder.copy()
	copied.query = append(copied.query, "DELETE", "FROM", builder.tableName)

	if len(builder.whereConditions) > 0 {
		copied.query = append(copied.query, builder.getWhereParagraphs()...)
	}

	return strings.TrimRight(strings.Join(copied.query, " "), "") + ";"
}
