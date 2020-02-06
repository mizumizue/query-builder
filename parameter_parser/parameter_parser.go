package parameter_parser

import (
	"reflect"
)

type ParameterParser struct {
	param interface{}
}

func NewParameterParser(param interface{}) *ParameterParser {
	return &ParameterParser{
		param: param,
	}
}

func (pp *ParameterParser) ParseNamedParam() map[string]interface{} {
	namedParam := make(map[string]interface{})
	// TODO
	return namedParam
}

// ex
//type SearchMachinesParameter struct {
//	MachineNumber *int       `search:"machine_number" operator:"eq"`
//	MachineName   *string    `search:"machine_name" operator:"eq"`
//	BuyDateFrom   *time.Time `search:"buy_date" operator:"ge"`
//	BuyDateTo     *time.Time `search:"buy_date" operator:"lt"`
//	PriceFrom     *time.Time `search:"price" operator:"gt"`
//	PriceTo       *time.Time `search:"price" operator:"le"`
//}
func (pp *ParameterParser) ParseBindMap() map[string]map[string]string {
	bindMap := make(map[string]map[string]string)

	t := reflect.TypeOf(pp.param)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		searchTag := field.Tag.Get("search")
		operatorTag := field.Tag.Get("operator")
		if searchTag != "" && operatorTag != "" {
			// TODO 好みの問題だけど field.Name を snake_case にしたい
			bindMap[field.Name] = map[string]string{
				"target":   searchTag,
				"operator": operatorTag,
			}
		}
	}
	return bindMap
}
