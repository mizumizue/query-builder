package query_builder

import "testing"

func Test_DeleteQueryBuilder_Normal(t *testing.T) {
	q := NewDeleteQueryBuilder().
		Table("users").
		Build()

	expected := "DELETE FROM users;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_DeleteQueryBuilder_Where(t *testing.T) {
	q := NewDeleteQueryBuilder().
		Table("users").
		Where("name", Equal).
		Where("age", GraterEqual).
		Where("age", LessEqual).
		Where("sex", Not).
		Where("age", LessThan).
		Where("age", GraterThan).
		Build()

	expected := "DELETE FROM users WHERE name = ? AND age >= ? AND age <= ? AND sex != ? AND age < ? AND age > ?;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewDeleteQueryBuilder().
		Placeholder(Named).
		Table("users").
		Where("name", Equal).
		Where("age", GraterEqual).
		Where("age", LessEqual).
		Where("sex", Not).
		Where("age", LessThan).
		Where("age", GraterThan).
		Build()

	expected2 := "DELETE FROM users WHERE name = :name AND age >= :age AND age <= :age AND sex != :sex AND age < :age AND age > :age;"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_DeleteQueryBuilder_WhereIn(t *testing.T) {
	q := NewDeleteQueryBuilder().
		Table("users").
		Where("user_name", Equal).
		WhereIn("user_id", 3).
		Build()

	expected := "DELETE FROM users WHERE user_name = ? AND user_id IN (?, ?, ?);"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewDeleteQueryBuilder().
		Placeholder(Named).
		Table("users").
		Where("user_name", Equal).
		WhereIn("user_id", 3).
		Build()

	expected2 := "DELETE FROM users WHERE user_name = :user_name AND user_id IN (:user_id1, :user_id2, :user_id3);"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_DeleteQueryBuilder_WhereNotIn(t *testing.T) {
	q := NewDeleteQueryBuilder().Table("users").
		Where("user_name", Equal).
		WhereNotIn("user_id", 3).
		Build()
	expected := "DELETE FROM users WHERE user_name = ? AND user_id NOT IN (?, ?, ?);"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewDeleteQueryBuilder().Table("users").
		Placeholder(Named).
		Where("user_name", Equal).
		WhereNotIn("user_id", 3).
		Build()
	expected2 := "DELETE FROM users WHERE user_name = :user_name AND user_id NOT IN (:user_id1, :user_id2, :user_id3);"
	if q2 != expected2 {
		t.Logf("expected: %s, acctual: %s", expected2, q2)
		t.Fail()
	}
}
