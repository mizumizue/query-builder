package query_builder

import (
	"reflect"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

func Test_QueryBuilderWithOrder(t *testing.T) {
	q := NewQueryBuilder().
		Table("users").
		OrderBy("created, user_id", Asc).
		Build()
	expected := "SELECT users.* FROM users ORDER BY created, user_id ASC;"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
		t.Fail()
	}

	q2 := NewQueryBuilder().
		Table("users").
		OrderBy("created, user_id", Asc).
		Limit().
		Offset().
		Build()
	expected2 := "SELECT users.* FROM users ORDER BY created, user_id ASC LIMIT ? OFFSET ?;"
	if expected2 != q2 {
		t.Logf("expected: %s\n acctual: %s", expected2, q2)
		t.Fail()
	}
}

func Test_QueryBuilderWithGroupBy(t *testing.T) {
	q := NewQueryBuilder().
		Table("users").
		GroupBy("user_id").
		Build()

	expected := "SELECT users.* FROM users GROUP BY user_id;"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
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

func Test_QueryBuilderSelectCOALESCE(t *testing.T) {
	q := NewQueryBuilder().Table("users").Select("COALESCE(name, 0) as name", "age", "sex").Build()
	expected := "SELECT COALESCE(name, 0) as name, users.age, users.sex FROM users;"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
		t.Fail()
	}
}

func Test_QueryBuilderMultiPattern(t *testing.T) {
	q := NewQueryBuilder().Table("users").
		Where("name", Equal).
		Where("age", GraterEqual).
		Where("age", LessEqual).
		Where("sex", Not).
		Where("age", LessThan).
		Where("age", GraterThan).
		Build()

	expected := "SELECT users.* FROM users WHERE name = ? AND age >= ? AND age <= ? AND sex != ? AND age < ? AND age > ?;"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
		t.Fail()
	}

	q2 := NewQueryBuilder().Table("users").
		UseNamedPlaceholder().
		Where("name", Equal).
		Where("age", GraterEqual).
		Where("age", LessEqual).
		Where("sex", Not).
		Where("age", LessThan).
		Where("age", GraterThan).
		Build()

	expected2 := "SELECT users.* FROM users WHERE name = :name AND age >= :age AND age <= :age AND sex != :sex AND age < :age AND age > :age;"
	if q2 != expected2 {
		t.Logf("expected: %s, acctual: %s", expected2, q2)
		t.Fail()
	}

	q3 := NewQueryBuilder().Table("users").
		UseNamedPlaceholder().
		Where("name", Equal).
		Where("age", GraterEqual, "age1").
		Where("age", LessEqual, "age2").
		Where("sex", Not, "sex1").
		Where("age", LessThan, "age3").
		Where("age", GraterThan, "age4").
		Build()

	expected3 := "SELECT users.* FROM users WHERE name = :name AND age >= :age1 AND age <= :age2 AND sex != :sex1 AND age < :age3 AND age > :age4;"
	if q3 != expected3 {
		t.Logf("expected: %s, acctual: %s", expected3, q3)
		t.Fail()
	}
}

func Test_QueryBuilderJoin(t *testing.T) {
	joinFields := []string{"user_id"}
	q := NewQueryBuilder().Table("users").UseNamedPlaceholder().
		Join(LeftJoin, "tasks", joinFields, joinFields).Build()
	expected := "SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.user_id;"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
		t.Fail()
	}

	joinFields2 := []string{"user_id"}
	joinFields3 := []string{"task_id"}
	q2 := NewQueryBuilder().Table("users").UseNamedPlaceholder().
		Join(LeftJoin, "tasks", joinFields2, joinFields2).
		Join(LeftJoin, "subtasks", joinFields3, joinFields3, "tasks").
		Build()

	expected2 := "SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.user_id LEFT JOIN subtasks ON tasks.task_id = subtasks.task_id;"
	if q2 != expected2 {
		t.Logf("expected: %s acctual: %s", expected2, q2)
		t.Fail()
	}
}

func Test_QueryBuilderWhereOperator(t *testing.T) {
	joinFields := []string{"user_id"}
	q := NewQueryBuilder().Table("users").
		UseNamedPlaceholder().
		Join(LeftJoin, "tasks", joinFields, joinFields).
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

//ex Tag
type SearchMachinesParameter struct {
	MachineNumber *int       `search:"machine_number" operator:"eq"`
	MachineName   *string    `search:"machine_name" operator:"eq"`
	BuyDateFrom   *time.Time `search:"buy_date" operator:"ge"`
	BuyDateTo     *time.Time `search:"buy_date" operator:"lt"`
	PriceFrom     *int       `search:"price" operator:"gt"`
	PriceTo       *int       `search:"price" operator:"le"`
	Owner         *string    `search:"owner" operator:"not"`
}

func Test_WhereMulti(t *testing.T) {
	machineNumber := 150
	machineName := "machine1"
	price := 1000
	now := time.Now()
	owner := "owner1"

	searchParam := SearchMachinesParameter{
		MachineNumber: &machineNumber,
		MachineName:   &machineName,
		BuyDateFrom:   &now,
		BuyDateTo:     &now,
		PriceFrom:     &price,
		PriceTo:       &price,
		Owner:         &owner,
	}

	qb := NewQueryBuilder().Table("machines").
		UseNamedPlaceholder().
		WhereMultiByStruct(searchParam)

	expected := "SELECT machines.* FROM " +
		"machines " +
		"WHERE machine_number = :machine_number " +
		"AND machine_name = :machine_name " +
		"AND buy_date >= :buy_date_from " +
		"AND buy_date < :buy_date_to " +
		"AND price > :price_from " +
		"AND price <= :price_to " +
		"AND owner != :owner;"

	q := qb.Build()
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
		t.Fail()
	}
}

func Test_WhereIn(t *testing.T) {
	q := NewQueryBuilder().Table("users").
		Where("user_name", Equal).
		WhereIn("user_id", 3).
		Build()
	expected := "SELECT users.* FROM users WHERE user_name = ? AND user_id IN (?, ?, ?);"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
		t.Fail()
	}

	q2 := NewQueryBuilder().Table("users").
		UseNamedPlaceholder().
		Where("user_name", Equal).
		WhereIn("user_id", 3).
		Build()
	expected2 := "SELECT users.* FROM users WHERE user_name = :user_name AND user_id IN (:user_id1, :user_id2, :user_id3);"
	if q2 != expected2 {
		t.Logf("expected: %s, acctual: %s", expected2, q2)
		t.Fail()
	}
}

func Test_WhereNotIn(t *testing.T) {
	q := NewQueryBuilder().Table("users").
		Where("user_name", Equal).
		WhereNotIn("user_id", 3).
		Build()
	expected := "SELECT users.* FROM users WHERE user_name = ? AND user_id NOT IN (?, ?, ?);"
	if q != expected {
		t.Logf("expected: %s, acctual: %s", expected, q)
		t.Fail()
	}

	q2 := NewQueryBuilder().Table("users").
		UseNamedPlaceholder().
		Where("user_name", Equal).
		WhereNotIn("user_id", 3).
		Build()
	expected2 := "SELECT users.* FROM users WHERE user_name = :user_name AND user_id NOT IN (:user_id1, :user_id2, :user_id3);"
	if q2 != expected2 {
		t.Logf("expected: %s, acctual: %s", expected2, q2)
		t.Fail()
	}
}

func Test_JoinMultipleFields(t *testing.T) {
	fields := []string{"user_id", "task_id"}
	q := NewQueryBuilder().Table("users").
		Join(LeftJoin, "tasks", fields, fields).
		Build()
	expected := "SELECT users.* FROM users " +
		"LEFT JOIN tasks " +
		"ON users.user_id = tasks.user_id AND users.task_id = tasks.task_id;"
	if q != expected {
		t.Logf("expected: %s\n acctual: %s", expected, q)
		t.Fail()
	}

	// JOIN 先と元でField名が異なる場合のJOIN
	originFields := []string{"user_id", "user_task_id"}
	targetFields := []string{"task_user_id", "task_id"}
	q2 := NewQueryBuilder().Table("users").
		Join(LeftJoin, "tasks", originFields, targetFields).
		Build()
	expected2 := "SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.task_user_id AND users.user_task_id = tasks.task_id;"

	if q2 != expected2 {
		t.Logf("expected: %s\n acctual: %s", expected2, q2)
		t.Fail()
	}
}
