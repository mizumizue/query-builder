package query_builder

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
)

func checkQuery(expected, actual string) error {
	if expected != actual {
		return fmt.Errorf("\nexpected: %s \nactual  : %s", expected, actual)
	}
	return nil
}

func checkSqlSyntax(sql string) error {
	_, err := sqlparser.Parse(sql)
	if err != nil {
		return fmt.Errorf("sql syntax err. detail: %v", err)
	}
	return nil
}
