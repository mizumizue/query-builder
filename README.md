# query-builder

`query-builder`'s purpose is  build sql queries.

execute `Query` or execute `Exec` is not purpose this library.

Please use this with `database/sql` or `jmoiron/sqlx` or etc.

## Characteristic

- Method chain
- Immutable

## Builder Types

- SelectQueryBuilder
- InsertQueryBuilder
- UpdateQueryBuilder
- DeleteQueryBuilder

## How to use examples

### SelectQueryBuilder

```
//example model
type User struct {
	UserID string `db:"user_id" table:"users"`
	Name   string `db:"name" table:"users"`
	Age    int    `db:"age" table:"users"`
	Sex    string `db:"sex" table:"users"`
}
```

```
# All Columns select
# SELECT users.* FROM users;
NewSelectQueryBuilder().Table("users").Build()

# Columns By struct select
# SELECT users.user_id, users.name, users.age, users.sex FROM users;
NewSelectQueryBuilder().Table("users").Model(User{}).Build()

# Columns select
# SELECT users.name, users.age, users.sex FROM users;
NewSelectQueryBuilder().Table("users").Column("name", "age", "sex").Build()

# Use GroupBy
# SELECT users.* FROM users GROUP BY user_id;
NewSelectQueryBuilder().Table("users").GroupBy("user_id").Build()

# Use OrderBy
# SELECT users.* FROM users ORDER BY created ASC;
NewSelectQueryBuilder().Table("users").OrderBy("created", Asc).Build()

# Use Limit
# SELECT users.* FROM users LIMIT ?;
NewSelectQueryBuilder().Table("users").Limit().Build()

# Use Offset
# SELECT users.* FROM users OFFSET ?;
NewSelectQueryBuilder().Table("users").Offset().Build()

# Use Where
# SELECT users.* FROM users WHERE name = ? AND age >= ? AND age <= ? AND sex != ? AND age < ? AND age > ?;
NewSelectQueryBuilder().
    Table("users").
    Where("name", Equal).
    Where("age", GraterEqual).
    Where("age", LessEqual).
    Where("sex", Not).
    Where("age", LessThan).
    Where("age", GraterThan).
    Build()

# Select Placeholder(default is `?`)
# SELECT users.* FROM users WHERE name = :name AND age >= :age AND age <= :age AND sex != :sex AND age < :age AND age > :age;
NewSelectQueryBuilder().
    Placeholder(Named).
    Table("users").
    Where("name", Equal).
    Where("age", GraterEqual).
    Where("age", LessEqual).
    Where("sex", Not).
    Where("age", LessThan).
    Where("age", GraterThan).
    Build()

# Select custom Placeholder(default `Named` is `column_name`)
# SELECT users.* FROM users WHERE name = :name AND age >= :age1 AND age <= :age2 AND sex != :sex1 AND age < :age3 AND age > :age4;
NewSelectQueryBuilder().
    Placeholder(Named).
    Table("users").
    Where("name", Equal).
    Where("age", GraterEqual, "age1").
    Where("age", LessEqual, "age2").
    Where("sex", Not, "sex1").
    Where("age", LessThan, "age3").
    Where("age", GraterThan, "age4").
    Build()

# Use IN(?)
# SELECT users.* FROM users WHERE user_name = ? AND user_id IN (?, ?, ?);
NewSelectQueryBuilder().
    Table("users").
    Where("user_name", Equal).
    WhereIn("user_id", 3).
    Build()

# Use IN(:named)
# SELECT users.* FROM users WHERE user_name = :user_name AND user_id IN (:user_id1, :user_id2, :user_id3);
NewSelectQueryBuilder().
    Placeholder(Named).
    Table("users").
    Where("user_name", Equal).
    WhereIn("user_id", 3).
    Build()

# Use NOT IN(?)
# SELECT users.* FROM users WHERE user_name = ? AND user_id IN (?, ?, ?);
NewSelectQueryBuilder().
    Table("users").
    Where("user_name", Equal).
    WhereNotIn("user_id", 3).
    Build()

# Use NOT IN(:named)
# SELECT users.* FROM users WHERE user_name = :user_name AND user_id IN (:user_id1, :user_id2, :user_id3);
NewSelectQueryBuilder().
    Table("users").
    Placeholder(Named).
    Where("user_name", Equal).
    WhereNotIn("user_id", 3).
    Build()

# Use Where Bind By Struct
# SELECT machines.* FROM machines WHERE machine_number = :machine_number AND machine_name = :machine_name AND buy_date >= :buy_date_from AND buy_date < :buy_date_to AND price > :price_from AND price <= :price_to AND owner != :owner;
# Ex Struct
type SearchMachinesParameter struct { //ex Tagged struct
    MachineNumber *int       `search:"machine_number" operator:"eq"`
    MachineName   *string    `search:"machine_name" operator:"eq"`
    BuyDateFrom   *time.Time `search:"buy_date" operator:"ge"`
    BuyDateTo     *time.Time `search:"buy_date" operator:"lt"`
    PriceFrom     *int       `search:"price" operator:"gt"`
    PriceTo       *int       `search:"price" operator:"le"`
    Owner         *string    `search:"owner" operator:"not"`
}
NewSelectQueryBuilder().
    Placeholder(Named).
    Table("machines").
    WhereMultiByStruct(searchParam).
    Build()

# Use Join
# SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.user_id;
joinFields := []string{"user_id"}
NewSelectQueryBuilder().
    Placeholder(Named).
    Table("users").
    Join(LeftJoin, "tasks", joinFields, joinFields).Build()

# Use Join with Named Parameter
# SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.task_user_id AND users.user_task_id = tasks.task_id;
originFields := []string{"user_id", "user_task_id"}
targetFields := []string{"task_user_id", "task_id"}
NewSelectQueryBuilder().
    Table("users").
    Join(LeftJoin, "tasks", originFields, targetFields).
    Build()

# Multi Field Join
# SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.user_id AND users.task_id = tasks.task_id;
fields := []string{"user_id", "task_id"}
NewSelectQueryBuilder().Table("users").
    Join(LeftJoin, "tasks", fields, fields).
    Build()
```

### InsertQueryBuilder

```
# Select Columns
# INSERT INTO users(name, age, sex) VALUES(?, ?, ?);
NewInsertQueryBuilder().
    Table("users").
    Column("name", "age", "sex").
    Build()

# INSERT INTO users(name, age, sex) VALUES(:name, :age, :sex);
NewInsertQueryBuilder().
    Placeholder(Named).
    Table("users").
    Column("name", "age", "sex").
    Build()

# Select By Model
# INSERT INTO users(user_id, name, age, sex) VALUES(?, ?, ?, ?);
NewInsertQueryBuilder().
    Table("users").
    Model(User{}).
    Build()

# INSERT INTO users(user_id, name, age, sex) VALUES(:user_id, :name, :age, :sex);
NewInsertQueryBuilder().
    Placeholder(Named).
    Table("users").
    Model(User{}).
    Build()
```

### UpdateQueryBuilder

```
# Select Columns
# UPDATE users SET name = ?, age = ?, sex = ?;
NewUpdateQueryBuilder().
    Table("users").
    Column("name", "age", "sex").
    Build()

# UPDATE users SET name = :name, age = :age, sex = :sex;
NewUpdateQueryBuilder().
    Placeholder(Named).
    Table("users").
    Column("name", "age", "sex").
    Build()

# Use Where
# UPDATE users SET name = ?, age = ?, sex = ? WHERE name = ? AND age >= ? AND age <= ? AND sex != ? AND age < ? AND age > ?;
NewUpdateQueryBuilder().
    Table("users").
    Column("name", "age", "sex").
    Where("name", Equal).
    Where("age", GraterEqual).
    Where("age", LessEqual).
    Where("sex", Not).
    Where("age", LessThan).
    Where("age", GraterThan).
    Build()

# Use IN(?)
# UPDATE users SET name = ?, age = ?, sex = ? WHERE user_name = ? AND user_id IN (?, ?, ?);
NewUpdateQueryBuilder().
    Table("users").
    Column("name", "age", "sex").
    Where("user_name", Equal).
    WhereIn("user_id", 3).
    Build()

# Use NOT IN(?)
# UPDATE users SET name = ?, age = ?, sex = ? WHERE user_name = ? AND user_id NOT IN (?, ?, ?);
NewUpdateQueryBuilder().
    Table("users").
    Column("name", "age", "sex").
    Where("user_name", Equal).
    WhereNotIn("user_id", 3).
    Build()
```

### DeleteQueryBuilder

```
# All Delete
# DELETE FROM users;
NewDeleteQueryBuilder().
    Table("users").
    Build()

# Use Where
# DELETE FROM users WHERE name = ? AND age >= ? AND age <= ? AND sex != ? AND age < ? AND age > ?;
NewDeleteQueryBuilder().
    Table("users").
    Where("name", Equal).
    Where("age", GraterEqual).
    Where("age", LessEqual).
    Where("sex", Not).
    Where("age", LessThan).
    Where("age", GraterThan).
    Build()

# Use IN(?)
# DELETE FROM users WHERE user_name = ? AND user_id IN (?, ?, ?);
NewDeleteQueryBuilder().
    Table("users").
    Where("user_name", Equal).
    WhereIn("user_id", 3).
    Build()

# Use NOT IN(?)
# DELETE FROM users WHERE user_name = ? AND user_id NOT IN (?, ?, ?);
NewDeleteQueryBuilder().Table("users").
    Where("user_name", Equal).
    WhereNotIn("user_id", 3).
    Build()
```

## Install

```
go get -u github.com/trewanek/query-builder
```
