package query_builder

import "strings"

type InsertQueryBuilder struct {
	*queryBuilder
}

func NewInsertQueryBuilder() *InsertQueryBuilder {
	builder := &InsertQueryBuilder{}
	builder.queryBuilder = &queryBuilder{}
	builder.placeholderType = Question
	return builder
}

func (builder *InsertQueryBuilder) copy() *InsertQueryBuilder {
	return &InsertQueryBuilder{
		builder.queryBuilder.copy(),
	}
}

func (builder *InsertQueryBuilder) Placeholder(placeholderType int) *InsertQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.placeholder(placeholderType)
	return copied
}

func (builder *InsertQueryBuilder) Table(tableName string) *InsertQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.table(tableName)
	return copied
}

func (builder *InsertQueryBuilder) Model(src interface{}) *InsertQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.model(src)
	return copied
}

func (builder *InsertQueryBuilder) Column(columns ...string) *InsertQueryBuilder {
	copied := builder.copy()
	copied.queryBuilder = builder.column(columns...)
	return copied
}

func (builder *InsertQueryBuilder) Build() string {
	if builder.tableName == "" {
		panic("target table is empty!!!")
	}

	copied := builder.copy()
	columns := builder.columns

	copied.query = append(copied.query, builder.getInsertIntoParagraphs()...)
	copied.query = append(copied.query, builder.getTableAndColumnsParagraphs(builder.tableName, columns...))
	copied.query = append(copied.query, builder.getValuesParagraphs(columns...))

	if len(builder.whereConditions) > 0 {
		copied.query = append(copied.query, builder.getWhereParagraphs()...)
	}

	return strings.TrimRight(strings.Join(copied.query, " "), "") + ";"
}

func (builder *InsertQueryBuilder) getInsertIntoParagraphs() []string {
	paragraph := make([]string, 0, 0)
	paragraph = append(paragraph, "INSERT")
	paragraph = append(paragraph, "INTO")
	return paragraph
}

func (builder *InsertQueryBuilder) getTableAndColumnsParagraphs(tableName string, columns ...string) string {
	return ""
}

func (builder *InsertQueryBuilder) getValuesParagraphs(columns ...string) string {
	return ""
}
