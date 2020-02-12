# query-builder

`query-builder`'s purpose is  build sql queries.

execute `Query` or execute `Exec` is not purpose this library.

Please use this with `database/sql` or `jmoiron/sqlx` or etc.

## Install

```
go get -u github.com/trewanek/query_builder
```

## How to use

Query Builder is Immutable.

And You can use method chain.

```
# example model
//type User struct {
//	UserID string `db:"user_id" table:"users"`
//	Name   string `db:"name" table:"users"`
//	Age    int    `db:"age" table:"users"`
//	Sex    string `db:"sex" table:"users"`
//}

# SELECT users.* FROM users;
NewQueryBuilder().Table("users").Build()

# SELECT users.user_id, users.name, users.age, users.sex FROM users;
NewQueryBuilder().Table("users").Model(User{}).Build()

# SELECT users.* FROM users LIMIT ?;
NewQueryBuilder().Table("users").Limit().Build()

# SELECT users.* FROM users WHERE name = ? AND age >= ? AND age <= ? AND sex != ? AND age < ? AND age > ?;
NewQueryBuilder().Table("users").
    Where("name", query_operator.Equal).
    Where("age", query_operator.GraterEqual).
    Where("age", query_operator.LessEqual).
    Where("sex", query_operator.Not).
    Where("age", query_operator.LessThan).
    Where("age", query_operator.GraterThan).
    Build()

# if you want to use named placeholder, `UseNamedPlaceholder()`
# SELECT users.* FROM users WHERE name = :name AND age >= :age AND age <= :age AND sex != :sex AND age < :age AND age > :age;
NewQueryBuilder().Table("users").
    UseNamedPlaceholder().
    Where("name", query_operator.Equal).
    Where("age", query_operator.GraterEqual).
    Where("age", query_operator.LessEqual).
    Where("sex", query_operator.Not).
    Where("age", query_operator.LessThan).
    Where("age", query_operator.GraterThan).
    Build()

# if you want to change named placeholder binding name
# SELECT users.* FROM users WHERE name = :name AND age >= :age1 AND age <= :age2 AND sex != :sex1 AND age < :age3 AND age > :age4;
NewQueryBuilder().Table("users").
    UseNamedPlaceholder().
    Where("name", query_operator.Equal).
    Where("age", query_operator.GraterEqual, "age1").
    Where("age", query_operator.LessEqual, "age2").
    Where("sex", query_operator.Not, "sex1").
    Where("age", query_operator.LessThan, "age3").
    Where("age", query_operator.GraterThan, "age4").
    Build()

# SELECT users.name, users.age, users.sex FROM users;
NewQueryBuilder().Table("users").Select("name", "age", "sex").Build()

# SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.user_id;
joinFields := []string{"user_id"}
NewQueryBuilder().Table("users").UseNamedPlaceholder().
    Join(LeftJoin, "tasks", joinFields, joinFields).Build()

# SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.task_user_id AND users.user_task_id = tasks.task_id;
originFields := []string{"user_id", "user_task_id"}
targetFields := []string{"task_user_id", "task_id"}
NewQueryBuilder().Table("users").
    Join(LeftJoin, "tasks", originFields, targetFields).
    Build()

# SELECT users.* FROM users WHERE user_name = ? AND user_id IN (?, ?, ?);
NewQueryBuilder().Table("users").
    Where("user_name", query_operator.Equal).
    WhereIn("user_id", 3).
    Build()

# SELECT users.* FROM users WHERE user_name = :user_name AND user_id IN (:user_id1, :user_id2, :user_id3);
NewQueryBuilder().Table("users").
    UseNamedPlaceholder().
    Where("user_name", query_operator.Equal).
    WhereIn("user_id", 3).
    Build()

...TODO write other how to use examples
```
