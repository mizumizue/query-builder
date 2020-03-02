package query_builder

import "testing"

func Test_InsertQueryBuilder_Column(t *testing.T) {
	testCommonFunc(
		t,
		"INSERT INTO users(name, age, sex) VALUES(?, ?, ?);",
		NewInsertQueryBuilder().
			Table("users").
			Column("name", "age", "sex").
			Build(),
		true,
	)

	testCommonFunc(
		t,
		"INSERT INTO users(name, age, sex) VALUES(:name, :age, :sex);",
		NewInsertQueryBuilder().
			Placeholder(Named).
			Table("users").
			Column("name", "age", "sex").
			Build(),
		true,
	)
}

func Test_InsertQueryBuilder_Omit(t *testing.T) {
	testCommonFunc(
		t,
		"INSERT INTO users(age, sex) VALUES(?, ?);",
		NewInsertQueryBuilder().
			Table("users").
			Column("name", "age", "sex").
			Omit("name").
			Build(),
		true,
	)

	testCommonFunc(
		t,
		"INSERT INTO users(name, sex) VALUES(:name, :sex);",
		NewInsertQueryBuilder().
			Placeholder(Named).
			Table("users").
			Column("name", "age", "sex").
			Omit("age").
			Build(),
		true,
	)
}

func Test_InsertQueryBuilder_Model(t *testing.T) {
	testCommonFunc(
		t,
		"INSERT INTO users(user_id, name, age, sex) VALUES(?, ?, ?, ?);",
		NewInsertQueryBuilder().
			Table("users").
			Model(User{}, true).
			Build(),
		true,
	)

	testCommonFunc(
		t,
		"INSERT INTO users(name, age, sex) VALUES(?, ?, ?);",
		NewInsertQueryBuilder().
			Table("users").
			Model(User{}, true).
			Omit("user_id").
			Build(),
		true,
	)

	testCommonFunc(
		t,
		"INSERT INTO users(user_id, name, age) VALUES(?, ?, ?);",
		NewInsertQueryBuilder().
			Table("users").
			Model(User{}, true).
			Omit("sex").
			Build(),
		true,
	)

	testCommonFunc(
		t,
		"INSERT INTO users(name) VALUES(:name);",
		NewInsertQueryBuilder().
			Placeholder(Named).
			Table("users").
			Model(User{Name: "hoge"}).
			Build(),
		true,
	)
}
