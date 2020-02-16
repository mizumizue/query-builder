package query_builder

import "testing"

func Test_InsertQueryBuilder_Column(t *testing.T) {
	q := NewInsertQueryBuilder().
		Table("users").
		Column("name", "age", "sex").
		Build()

	expected := "INSERT INTO users(name, age, sex) VALUES(?, ?, ?);"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewInsertQueryBuilder().
		Placeholder(Named).
		Table("users").
		Column("name", "age", "sex").
		Build()

	expected2 := "INSERT INTO users(name, age, sex) VALUES(:name, :age, :sex);"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_InsertQueryBuilder_Model(t *testing.T) {
	q := NewInsertQueryBuilder().
		Table("users").
		Model(User{}).
		Build()

	expected := "INSERT INTO users(user_id, name, age, sex) VALUES(?, ?, ?, ?);"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewInsertQueryBuilder().
		Placeholder(Named).
		Table("users").
		Model(User{}).
		Build()

	expected2 := "INSERT INTO users(user_id, name, age, sex) VALUES(:user_id, :name, :age, :sex);"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}
