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
# SELECT users.* FROM users;
NewSelectQueryBuilder().Table("users").Build()

# SELECT users.user_id, users.name, users.age, users.sex FROM users;
NewSelectQueryBuilder().Table("users").Model(User{}).Build()

# SELECT users.name, users.age, users.sex FROM users;
NewSelectQueryBuilder().Table("users").Column("name", "age", "sex").Build()

# SELECT users.* FROM users LIMIT ?;
NewSelectQueryBuilder().Table("users").Limit().Build()

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

# if you want to use named placeholder, `Placeholder(Named)`
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

# if you want to change named placeholder binding name
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

# SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.user_id;
joinFields := []string{"user_id"}
NewSelectQueryBuilder().
    Placeholder(Named).
    Table("users").
    Join(LeftJoin, "tasks", joinFields, joinFields).Build()

# SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.task_user_id AND users.user_task_id = tasks.task_id;
originFields := []string{"user_id", "user_task_id"}
targetFields := []string{"task_user_id", "task_id"}
NewSelectQueryBuilder().
    Table("users").
    Join(LeftJoin, "tasks", originFields, targetFields).
    Build()

# SELECT users.* FROM users WHERE user_name = ? AND user_id IN (?, ?, ?);
NewSelectQueryBuilder().
    Table("users").
    Where("user_name", Equal).
    WhereIn("user_id", 3).
    Build()

# SELECT users.* FROM users WHERE user_name = :user_name AND user_id IN (:user_id1, :user_id2, :user_id3);
NewSelectQueryBuilder().
    Placeholder(Named).
    Table("users").
    Where("user_name", Equal).
    WhereIn("user_id", 3).
    Build()

...TODO write other how to use examples
```

### InsertQueryBuilder

```
...TODO write other how to use examples
```

### UpdateQueryBuilder

```
...TODO write other how to use examples
```

### DeleteQueryBuilder

```
...TODO write other how to use examples
```

## Install

```
go get -u github.com/trewanek/query-builder
```
