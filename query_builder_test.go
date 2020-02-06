package query_builder

import (
	"query-builder/query_operator"
	"reflect"
	"testing"
)

type User struct {
	UserID string `db:"user_id" table:"users"`
	Name   string `db:"name" table:"users"`
	Age    int    `db:"age" table:"users"`
	Sex    string `db:"sex" table:"users"`
}

func Test_QueryBuilder(t *testing.T) {
	q := NewQueryBuilder().Table("users").Build()
	expected := "SELECT users.* FROM users;"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
		t.Fail()
	}
}

func Test_QueryBuilderModel(t *testing.T) {
	q := NewQueryBuilder().Table("users").Model(User{}).Build()
	expected := "SELECT users.user_id, users.name, users.age, users.sex FROM users;"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
		t.Fail()
	}
}

func Test_QueryBuilderWithLimit(t *testing.T) {
	q := NewQueryBuilder().Table("users").Limit().Build()
	expected := "SELECT users.* FROM users LIMIT ?;"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
		t.Fail()
	}

	q2 := NewQueryBuilder().Table("users").UseNamedPlaceholder().Limit().Build()
	expected2 := "SELECT users.* FROM users LIMIT :limit;"
	if q2 != expected2 {
		t.Logf("expected: %s, acctual: %s", expected2, q2)
		t.Fail()
	}
}

func Test_QueryBuilderWithOffset(t *testing.T) {
	q := NewQueryBuilder().Table("users").Offset().Build()
	expected := "SELECT users.* FROM users OFFSET ?;"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
		t.Fail()
	}

	q2 := NewQueryBuilder().Table("users").UseNamedPlaceholder().Offset().Build()
	expected2 := "SELECT users.* FROM users OFFSET :offset;"
	if q2 != expected2 {
		t.Logf("expected: %s, acctual: %s", expected2, q2)
		t.Fail()
	}
}

func Test_QueryBuilderWithLimitAndOffset(t *testing.T) {
	q := NewQueryBuilder().Table("users").Limit().Offset().Build()
	expected := "SELECT users.* FROM users LIMIT ? OFFSET ?;"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
		t.Fail()
	}

	q2 := NewQueryBuilder().Table("users").Limit().Offset().Build()
	expected2 := "SELECT users.* FROM users LIMIT :limit OFFSET :offset;"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected2, q2)
		t.Fail()
	}
}

func Test_QueryBuilderSelect(t *testing.T) {
	q := NewQueryBuilder().Table("users").Select("name", "age", "sex").Build()
	expected := "SELECT users.name, users.age, users.sex FROM users;"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
		t.Fail()
	}
}

func Test_QueryBuilderMultiPattern(t *testing.T) {
	q := NewQueryBuilder().Table("users").
		Where("name", query_operator.Equal).
		Where("age", query_operator.GraterEqual).
		Where("age", query_operator.LessEqual).
		Where("sex", query_operator.Not).
		Where("age", query_operator.LessThan).
		Where("age", query_operator.GraterThan).
		Build()

	expected := "SELECT users.* FROM users WHERE name = ? AND age >= ? AND age <= ? AND sex != ? AND age < ? AND age > ?;"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
		t.Fail()
	}

	q2 := NewQueryBuilder().Table("users").
		UseNamedPlaceholder().
		Where("name", query_operator.Equal).
		Where("age", query_operator.GraterEqual).
		Where("age", query_operator.LessEqual).
		Where("sex", query_operator.Not).
		Where("age", query_operator.LessThan).
		Where("age", query_operator.GraterThan).
		Build()

	expected2 := "SELECT users.* FROM users WHERE name = :name AND age >= :age AND age <= :age AND sex != :sex AND age < :age AND age > :age;"
	if q2 != expected2 {
		t.Logf("expected: %s, acctual: %s", expected2, q2)
		t.Fail()
	}

	q3 := NewQueryBuilder().Table("users").
		UseNamedPlaceholder().
		Where("name", query_operator.Equal).
		Where("age", query_operator.GraterEqual, "age1").
		Where("age", query_operator.LessEqual, "age2").
		Where("sex", query_operator.Not, "sex1").
		Where("age", query_operator.LessThan, "age3").
		Where("age", query_operator.GraterThan, "age4").
		Build()

	expected3 := "SELECT users.* FROM users WHERE name = :name AND age >= :age1 AND age <= :age2 AND sex != :sex1 AND age < :age3 AND age > :age4;"
	if q3 != expected3 {
		t.Logf("expected: %s, acctual: %s", expected3, q3)
		t.Fail()
	}
}

func Test_QueryBuilderJoin(t *testing.T) {
	q := NewQueryBuilder().Table("users").UseNamedPlaceholder().Join(LeftJoin, "tasks", "user_id").Build()
	expected := "SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.user_id;"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
		t.Fail()
	}

	q2 := NewQueryBuilder().Table("users").UseNamedPlaceholder().
		Join(LeftJoin, "tasks", "user_id").
		Join(LeftJoin, "subtasks", "task_id", "tasks").
		Build()

	expected2 := "SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.user_id LEFT JOIN subtasks ON tasks.task_id = subtasks.task_id;"
	if q2 != expected2 {
		t.Logf("expected: %s acctual: %s", expected2, q2)
		t.Fail()
	}
}

func Test_QueryBuilderWhereOperator(t *testing.T) {
	q := NewQueryBuilder().Table("users").
		UseNamedPlaceholder().
		Join(LeftJoin, "tasks", "user_id").
		Build()
	expected := "SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.user_id;"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
		t.Fail()
	}
}

func Test_QueryBuilderIsImmutable(t *testing.T) {
	qb := NewQueryBuilder().Table("users").Offset()
	qb2 := qb.Table("tasks")

	if reflect.DeepEqual(qb, qb2) {
		t.Fail()
		t.Log(qb, qb2, " are deepEqual true. query build is not immutable.")
	}
}
