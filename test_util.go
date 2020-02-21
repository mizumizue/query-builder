package query_builder

import (
	"fmt"
	"testing"

	"github.com/xwb1989/sqlparser"
)

func testCommonFunc(t *testing.T, expected, actual string, sqlSyntaxCheck bool) {
	t.Run("expected query test", func(t *testing.T) {
		if err := checkQuery(expected, actual); err != nil {
			t.Log(err)
			t.Fail()
		}
	})

	t.Run("sql syntax check test", func(t *testing.T) {
		if sqlSyntaxCheck {
			if err := checkSqlSyntax(actual); err != nil {
				t.Log(err)
				t.Fail()
			}
		}
	})
}

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
