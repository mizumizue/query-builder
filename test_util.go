package query_builder

import (
	"fmt"
)

func checkQuery(expected, actual string) error {
	if expected != actual {
		return fmt.Errorf("\nexpected: %s \nactual  : %s", expected, actual)
	}
	return nil
}
