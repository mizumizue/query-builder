package query_builder

import "testing"

func Test_Hoge(t *testing.T) {
	q := NewInsertQueryBuilder().
		Table("users").
		Column("name", "age", "sex").
		Build()
	t.Log(q)

	expected := "INSERT INTO users(name, age, sex) VALUES(?, ?, ?)"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
}
