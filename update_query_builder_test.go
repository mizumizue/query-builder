package query_builder

import (
	"testing"
)

func Test_UpdateQueryBuilder_Column(t *testing.T) {
	q := NewUpdateQueryBuilder().
		Table("users").
		Column("name", "age", "sex").
		Build()

	expected := "UPDATE users SET name = ?, age = ?, sex = ?;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewUpdateQueryBuilder().
		Placeholder(Named).
		Table("users").
		Column("name", "age", "sex").
		Build()

	expected2 := "UPDATE users SET name = :name, age = :age, sex = :sex;"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_UpdateQueryBuilder_Where(t *testing.T) {
	q := NewUpdateQueryBuilder().
		Table("users").
		Column("name", "age", "sex").
		Where("name", Equal).
		Where("age", GraterEqual).
		Where("age", LessEqual).
		Where("sex", Not).
		Where("age", LessThan).
		Where("age", GraterThan).
		Build()

	expected := "UPDATE users SET name = ?, age = ?, sex = ? WHERE name = ? AND age >= ? AND age <= ? AND sex != ? AND age < ? AND age > ?;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewUpdateQueryBuilder().
		Placeholder(Named).
		Table("users").
		Column("name", "age", "sex").
		Where("name", Equal).
		Where("age", GraterEqual).
		Where("age", LessEqual).
		Where("sex", Not).
		Where("age", LessThan).
		Where("age", GraterThan).
		Build()

	expected2 := "UPDATE users SET name = :name, age = :age, sex = :sex WHERE name = :name AND age >= :age AND age <= :age AND sex != :sex AND age < :age AND age > :age;"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_UpdateQueryBuilder_WhereIn(t *testing.T) {
	q := NewUpdateQueryBuilder().
		Table("users").
		Column("name", "age", "sex").
		Where("user_name", Equal).
		WhereIn("user_id", 3).
		Build()

	expected := "UPDATE users SET name = ?, age = ?, sex = ? WHERE user_name = ? AND user_id IN (?, ?, ?);"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewUpdateQueryBuilder().
		Placeholder(Named).
		Table("users").
		Column("name", "age", "sex").
		Where("user_name", Equal).
		WhereIn("user_id", 3).
		Build()

	expected2 := "UPDATE users SET name = :name, age = :age, sex = :sex WHERE user_name = :user_name AND user_id IN (:user_id1, :user_id2, :user_id3);"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_UpdateQueryBuilder_WhereNotIn(t *testing.T) {
	q := NewUpdateQueryBuilder().
		Table("users").
		Column("name", "age", "sex").
		Where("user_name", Equal).
		WhereNotIn("user_id", 3).
		Build()
	expected := "UPDATE users SET name = ?, age = ?, sex = ? WHERE user_name = ? AND user_id NOT IN (?, ?, ?);"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewUpdateQueryBuilder().Table("users").
		Placeholder(Named).
		Column("name", "age", "sex").
		Where("user_name", Equal).
		WhereNotIn("user_id", 3).
		Build()
	expected2 := "UPDATE users SET name = :name, age = :age, sex = :sex WHERE user_name = :user_name AND user_id NOT IN (:user_id1, :user_id2, :user_id3);"
	if q2 != expected2 {
		t.Logf("expected: %s, acctual: %s", expected2, q2)
		t.Fail()
	}
}
