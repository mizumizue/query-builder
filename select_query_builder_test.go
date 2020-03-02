package query_builder

import (
	"reflect"
	"testing"
	"time"
)

type User struct {
	UserID string `db:"user_id" table:"users"`
	Name   string `db:"name" table:"users"`
	Age    int    `db:"age" table:"users"`
	Sex    string `db:"sex" table:"users"`
}

func Test_SelectQueryBuilder_OnlyTable(t *testing.T) {
	testCommonFunc(
		t,
		"SELECT users.* FROM users;",
		NewSelectQueryBuilder().Table("users").Build(),
		true,
	)
}

func Test_SelectQueryBuilder_Model(t *testing.T) {
	testCommonFunc(
		t,
		"SELECT users.user_id, users.name, users.age, users.sex FROM users;",
		NewSelectQueryBuilder().
			Table("users").
			Model(User{}).
			Build(),
		true,
	)

	testCommonFunc(
		t,
		"SELECT users.user_id, users.name, users.age, users.sex FROM users;",
		NewSelectQueryBuilder().
			Table("users").
			Model(&User{}).
			Build(),
		true,
	)
}

func Test_SelectQueryBuilder_Column(t *testing.T) {
	testCommonFunc(
		t,
		"SELECT users.name, users.age, users.sex FROM users;",
		NewSelectQueryBuilder().
			Table("users").
			Column("name", "age", "sex").
			Build(),
		true,
	)
}

func Test_SelectQueryBuilder_Column_DBMethod(t *testing.T) {
	testCommonFunc(
		t,
		"SELECT users.user_id, COALESCE(name, '') as name, users.age, users.sex FROM users;",
		NewSelectQueryBuilder().
			Table("users").
			Column("user_id", "COALESCE(name, '') as name", "age", "sex").
			Build(),
		true,
	)

	testCommonFunc(
		t,
		"SELECT users.user_id, COALESCE(name, '') as name FROM users;",
		NewSelectQueryBuilder().
			Table("users").
			Column("user_id", "COALESCE(name, '') as name").
			Build(),
		true,
	)
}

func Test_SelectQueryBuilder_OrderBy(t *testing.T) {
	testCommonFunc(
		t,
		"SELECT users.* FROM users ORDER BY created ASC;",
		NewSelectQueryBuilder().
			Table("users").
			OrderBy("created", Asc).
			Build(),
		true,
	)

	testCommonFunc(
		t,
		"SELECT users.* FROM users ORDER BY created, user_id DESC;",
		NewSelectQueryBuilder().
			Table("users").
			OrderBy("created, user_id", Desc).
			Build(),
		true,
	)
}

func Test_SelectQueryBuilder_GroupBy(t *testing.T) {
	testCommonFunc(
		t,
		"SELECT users.* FROM users GROUP BY user_id;",
		NewSelectQueryBuilder().
			Table("users").
			GroupBy("user_id").
			Build(),
		true,
	)
}

func Test_SelectQueryBuilder_Limit(t *testing.T) {
	q := NewSelectQueryBuilder().
		Table("users").
		Limit().
		Build()
	expected := "SELECT users.* FROM users LIMIT ?;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewSelectQueryBuilder().
		Table("users").
		Placeholder(Named).
		Limit().
		Build()
	expected2 := "SELECT users.* FROM users LIMIT :limit;"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q2); err != nil {
		t.Log(err)
		t.Fail()
	}

	q3 := NewSelectQueryBuilder().
		Table("users").
		Placeholder(DollarNumber).
		Limit().
		Build()
	expected3 := "SELECT users.* FROM users LIMIT $1;"
	if err := checkQuery(expected3, q3); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_Offset(t *testing.T) {
	q := NewSelectQueryBuilder().
		Table("users").
		Limit().
		Offset().
		Build()
	expected := "SELECT users.* FROM users LIMIT ? OFFSET ?;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewSelectQueryBuilder().
		Table("users").
		Placeholder(Named).
		Limit().
		Offset().
		Build()
	expected2 := "SELECT users.* FROM users LIMIT :limit OFFSET :offset;"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q2); err != nil {
		t.Log(err)
		t.Fail()
	}

	q3 := NewSelectQueryBuilder().
		Table("users").
		Placeholder(DollarNumber).
		Limit().
		Offset().
		Build()
	expected3 := "SELECT users.* FROM users LIMIT $1 OFFSET $2;"
	if err := checkQuery(expected3, q3); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_OrderBy_GroupBy_Limit_Offset(t *testing.T) {
	q := NewSelectQueryBuilder().
		Table("users").
		OrderBy("created", Asc).
		GroupBy("user_id").
		Limit().
		Offset().
		Build()
	expected := "SELECT users.* FROM users GROUP BY user_id ORDER BY created ASC LIMIT ? OFFSET ?;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewSelectQueryBuilder().
		Placeholder(Named).
		Table("users").
		OrderBy("created", Desc).
		GroupBy("user_id").
		Limit().
		Offset().
		Build()
	expected2 := "SELECT users.* FROM users GROUP BY user_id ORDER BY created DESC LIMIT :limit OFFSET :offset;"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_Where(t *testing.T) {
	// ? bind
	testCommonFunc(
		t,
		"SELECT users.* FROM users "+
			"WHERE name = ? "+
			"AND name LIKE ? "+
			"AND name NOT LIKE ? "+
			"AND age >= ? "+
			"AND age <= ? "+
			"AND sex != ? "+
			"AND age < ? "+
			"AND age > ? "+
			"AND age IS NULL "+
			"AND age IS NOT NULL;",
		NewSelectQueryBuilder().
			Table("users").
			Where("name", Equal).
			Where("name", Like).
			Where("name", NotLike).
			Where("age", GraterThanEqual).
			Where("age", LessThanEqual).
			Where("sex", NotEqual).
			Where("age", LessThan).
			Where("age", GraterThan).
			Where("age", IsNull).
			Where("age", IsNotNull).
			Build(),
		true,
	)

	// column name bind
	testCommonFunc(
		t,
		"SELECT users.* FROM users "+
			"WHERE name = :name "+
			"AND name LIKE :name "+
			"AND name NOT LIKE :name "+
			"AND age >= :age "+
			"AND age <= :age "+
			"AND sex != :sex "+
			"AND age < :age "+
			"AND age > :age "+
			"AND age IS NULL "+
			"AND age IS NOT NULL;",
		NewSelectQueryBuilder().Table("users").
			Placeholder(Named).
			Where("name", Equal).
			Where("name", Like).
			Where("name", NotLike).
			Where("age", GraterThanEqual).
			Where("age", LessThanEqual).
			Where("sex", NotEqual).
			Where("age", LessThan).
			Where("age", GraterThan).
			Where("age", IsNull).
			Where("age", IsNotNull).
			Build(),
		true,
	)

	// custom name bind
	testCommonFunc(
		t,
		"SELECT users.* FROM users "+
			"WHERE name = :name1 "+
			"AND name LIKE :name2 "+
			"AND name NOT LIKE :name3 "+
			"AND age >= :age1 "+
			"AND age <= :age2 "+
			"AND sex != :sex1 "+
			"AND age < :age3 "+
			"AND age > :age4 "+
			"AND age IS NULL "+
			"AND age IS NOT NULL;",
		NewSelectQueryBuilder().Table("users").
			Placeholder(Named).
			Where("name", Equal, "name1").
			Where("name", Like, "name2").
			Where("name", NotLike, "name3").
			Where("age", GraterThanEqual, "age1").
			Where("age", LessThanEqual, "age2").
			Where("sex", NotEqual, "sex1").
			Where("age", LessThan, "age3").
			Where("age", GraterThan, "age4").
			Where("age", IsNull, "age5").
			Where("age", IsNotNull, "age6").
			Build(),
		true,
	)

	// column name bind
	testCommonFunc(
		t,
		"SELECT users.* FROM users "+
			"WHERE name = $1 "+
			"AND name LIKE $2 "+
			"AND name NOT LIKE $3 "+
			"AND age >= $4 "+
			"AND age <= $5 "+
			"AND sex != $6 "+
			"AND age < $7 "+
			"AND age > $8 "+
			"AND age IS NULL "+
			"AND age IS NOT NULL;",
		NewSelectQueryBuilder().Table("users").
			Placeholder(DollarNumber).
			Where("name", Equal).
			Where("name", Like).
			Where("name", NotLike).
			Where("age", GraterThanEqual).
			Where("age", LessThanEqual).
			Where("sex", NotEqual).
			Where("age", LessThan).
			Where("age", GraterThan).
			Where("age", IsNull).
			Where("age", IsNotNull).
			Build(),
		false,
	)
}

func Test_SelectQueryBuilder_OR(t *testing.T) {
	// ? bind
	testCommonFunc(
		t,
		"SELECT users.* FROM users WHERE name = ? OR name = ?;",
		NewSelectQueryBuilder().
			Table("users").
			Where("name", Equal).
			Or("name", Equal).
			Build(),
		true,
	)
}

func Test_SelectQueryBuilder_WhereIn(t *testing.T) {
	q := NewSelectQueryBuilder().
		Table("users").
		Where("user_name", Equal).
		WhereIn("user_id", 3).
		Build()

	expected := "SELECT users.* FROM users WHERE user_name = ? AND user_id IN (?, ?, ?);"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewSelectQueryBuilder().
		Placeholder(Named).
		Table("users").
		Where("user_name", Equal).
		WhereIn("user_id", 3).
		Build()

	expected2 := "SELECT users.* FROM users WHERE user_name = :user_name AND user_id IN (:user_id1, :user_id2, :user_id3);"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q2); err != nil {
		t.Log(err)
		t.Fail()
	}

	q3 := NewSelectQueryBuilder().
		Placeholder(DollarNumber).
		Table("users").
		Where("user_name", Equal).
		WhereIn("user_id", 3).
		Build()

	expected3 := "SELECT users.* FROM users WHERE user_name = $1 AND user_id IN ($2, $3, $4);"
	if err := checkQuery(expected3, q3); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_WhereNotIn(t *testing.T) {
	q := NewSelectQueryBuilder().Table("users").
		Where("user_name", Equal).
		WhereNotIn("user_id", 3).
		Build()
	expected := "SELECT users.* FROM users WHERE user_name = ? AND user_id NOT IN (?, ?, ?);"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewSelectQueryBuilder().Table("users").
		Placeholder(Named).
		Where("user_name", Equal).
		WhereNotIn("user_id", 3).
		Build()
	expected2 := "SELECT users.* FROM users WHERE user_name = :user_name AND user_id NOT IN (:user_id1, :user_id2, :user_id3);"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q2); err != nil {
		t.Log(err)
		t.Fail()
	}

	q3 := NewSelectQueryBuilder().Table("users").
		Placeholder(DollarNumber).
		Where("user_name", Equal).
		WhereNotIn("user_id", 3).
		Build()
	expected3 := "SELECT users.* FROM users WHERE user_name = $1 AND user_id NOT IN ($2, $3, $4);"
	if err := checkQuery(expected3, q3); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_WhereMultiByStruct(t *testing.T) {
	type SearchMachinesParameter struct { //ex Tagged struct
		MachineNumber int       `db:"machine_number" search:"machine_number" operator:"eq"`
		MachineName   string    `db:"machine_name" search:"machine_name" operator:"eq"`
		BuyDateFrom   time.Time `db:"buy_date" search:"buy_date_from" operator:"gte"`
		BuyDateTo     time.Time `db:"buy_date" search:"buy_date_to" operator:"lt"`
		PriceFrom     int       `db:"price" search:"price_from" operator:"gt"`
		PriceTo       int       `db:"price" search:"price_to" operator:"lte"`
		Owner         string    `db:"owner" search:"owner" operator:"ne"`
	}

	machineNumber := 150
	machineName := "machine1"
	price := 1000
	now := time.Now()
	owner := "owner1"

	searchParam := SearchMachinesParameter{
		MachineNumber: machineNumber,
		MachineName:   machineName,
		BuyDateFrom:   now,
		BuyDateTo:     now,
		PriceFrom:     price,
		PriceTo:       price,
		Owner:         owner,
	}

	q := NewSelectQueryBuilder().
		Placeholder(Named).
		Table("machines").
		WhereMultiByStruct(searchParam).
		Build()

	expected := "SELECT machines.* FROM " +
		"machines " +
		"WHERE machine_number = :machine_number " +
		"AND machine_name = :machine_name " +
		"AND buy_date >= :buy_date_from " +
		"AND buy_date < :buy_date_to " +
		"AND price > :price_from " +
		"AND price <= :price_to " +
		"AND owner != :owner;"

	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewSelectQueryBuilder().
		Placeholder(DollarNumber).
		Table("machines").
		WhereMultiByStruct(searchParam).
		Build()

	expected2 := "SELECT machines.* FROM " +
		"machines " +
		"WHERE machine_number = $1 " +
		"AND machine_name = $2 " +
		"AND buy_date >= $3 " +
		"AND buy_date < $4 " +
		"AND price > $5 " +
		"AND price <= $6 " +
		"AND owner != $7;"

	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_WhereMultiByStructPtr(t *testing.T) {
	type SearchMachinesParameter struct { //ex Tagged struct
		MachineNumber *int       `db:"machine_number" search:"machine_number" operator:"eq"`
		MachineName   *string    `db:"machine_name" search:"machine_name" operator:"eq"`
		BuyDateFrom   *time.Time `db:"buy_date" search:"buy_date_from" operator:"gte"`
		BuyDateTo     *time.Time `db:"buy_date" search:"buy_date_to" operator:"lt"`
		PriceFrom     *int       `db:"price" search:"price_from" operator:"gt"`
		PriceTo       *int       `db:"price" search:"price_to" operator:"lte"`
		Owner         *string    `db:"owner" search:"owner" operator:"ne"`
	}

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

	q := NewSelectQueryBuilder().
		Placeholder(Named).
		Table("machines").
		WhereMultiByStruct(searchParam).
		Build()

	expected := "SELECT machines.* FROM " +
		"machines " +
		"WHERE machine_number = :machine_number " +
		"AND machine_name = :machine_name " +
		"AND buy_date >= :buy_date_from " +
		"AND buy_date < :buy_date_to " +
		"AND price > :price_from " +
		"AND price <= :price_to " +
		"AND owner != :owner;"

	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewSelectQueryBuilder().
		Placeholder(DollarNumber).
		Table("machines").
		WhereMultiByStruct(searchParam).
		Build()

	expected2 := "SELECT machines.* FROM " +
		"machines " +
		"WHERE machine_number = $1 " +
		"AND machine_name = $2 " +
		"AND buy_date >= $3 " +
		"AND buy_date < $4 " +
		"AND price > $5 " +
		"AND price <= $6 " +
		"AND owner != $7;"

	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_Join(t *testing.T) {
	joinFields := []string{"user_id"}
	q := NewSelectQueryBuilder().
		Placeholder(Named).
		Table("users").
		Join(LeftJoin, "tasks", joinFields, joinFields).
		Join(RightJoin, "tasks", joinFields, joinFields).
		Join(InnerJoin, "tasks", joinFields, joinFields).
		Build()
	expected := "SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.user_id RIGHT JOIN tasks ON users.user_id = tasks.user_id INNER JOIN tasks ON users.user_id = tasks.user_id;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q); err != nil {
		t.Log(err)
		t.Fail()
	}

	joinFields2 := []string{"user_id"}
	joinFields3 := []string{"task_id"}
	q2 := NewSelectQueryBuilder().
		Placeholder(Named).
		Table("users").
		Join(LeftJoin, "tasks", joinFields2, joinFields2).
		Join(RightJoin, "tasks", joinFields2, joinFields2).
		Join(InnerJoin, "tasks", joinFields2, joinFields2).
		Join(LeftJoin, "subtasks", joinFields3, joinFields3, "tasks").
		Build()
	expected2 := "SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.user_id RIGHT JOIN tasks ON users.user_id = tasks.user_id INNER JOIN tasks ON users.user_id = tasks.user_id LEFT JOIN subtasks ON tasks.task_id = subtasks.task_id;"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_JoinMultipleFields(t *testing.T) {
	fields := []string{"user_id", "task_id"}
	q := NewSelectQueryBuilder().Table("users").
		Join(LeftJoin, "tasks", fields, fields).
		Build()
	expected := "SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.user_id AND users.task_id = tasks.task_id;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q); err != nil {
		t.Log(err)
		t.Fail()
	}

	// JOIN 先と元でField名が異なる場合のJOIN
	originFields := []string{"user_id", "user_task_id"}
	targetFields := []string{"task_user_id", "task_id"}
	q2 := NewSelectQueryBuilder().Table("users").
		Join(LeftJoin, "tasks", originFields, targetFields).
		Build()
	expected2 := "SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.task_user_id AND users.user_task_id = tasks.task_id;"

	if q2 != expected2 {
		t.Logf("expected: %s\n acctual: %s", expected2, q2)
		t.Fail()
	}
	if err := checkSqlSyntax(q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_SubQuery(t *testing.T) {
	q := NewSelectQueryBuilder().
		Table("users").
		WhereSubQuery(
			"user_id",
			Equal,
			NewSelectQueryBuilder().Table("users").Column("user_id"),
		).
		Build()

	expected := "SELECT users.* FROM users WHERE user_id = (SELECT users.user_id FROM users);"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_DoubleSubQuery(t *testing.T) {
	sub2 := NewSelectQueryBuilder().
		Table("users").
		Column("user_id").
		Limit()

	sub1 := NewSelectQueryBuilder().
		Table("users").
		Column("user_id").
		WhereSubQuery("user_id", Equal, sub2).
		Limit()

	q := NewSelectQueryBuilder().
		Table("users").
		WhereSubQuery(
			"user_id",
			Equal,
			sub1,
		).
		Build()

	expected := "SELECT users.* FROM users WHERE user_id = (SELECT users.user_id FROM users WHERE user_id = (SELECT users.user_id FROM users LIMIT ?) LIMIT ?);"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := checkSqlSyntax(q); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_UnSpecifiedColumnSubQueryNotSelected(t *testing.T) {
	defer func() {
		err := recover()
		if err.(error) != UnspecifiedColumnErr {
			t.Log(err)
			t.Fail()
		}
	}()
	_ = NewSelectQueryBuilder().
		Table("users").
		WhereSubQuery(
			"user_id",
			Equal,
			NewSelectQueryBuilder().
				Table("users"),
		).
		Build()
}

func Test_SelectQueryBuilder_UnSpecifiedColumnSubQueryToMany(t *testing.T) {
	defer func() {
		err := recover()
		if err.(error) != UnspecifiedColumnErr {
			t.Log(err)
			t.Fail()
		}
	}()

	_ = NewSelectQueryBuilder().
		Table("users").
		Column("user_id", "name").
		WhereSubQuery(
			"user_id",
			Equal,
			NewSelectQueryBuilder().
				Table("users"),
		).
		Build()
}

func Test_SelectQueryBuilder_IsImmutable(t *testing.T) {
	qb := NewSelectQueryBuilder().
		Table("users").
		Offset()

	copied := qb.
		Table("tasks")

	if reflect.DeepEqual(qb, copied) {
		t.Fail()
		t.Log(qb, copied, " are deepEqual true. object is not immutable.")
	}
}
