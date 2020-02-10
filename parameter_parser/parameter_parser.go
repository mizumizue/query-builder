package parameter_parser

import (
	"github.com/iancoleman/strcase"
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
	t := reflect.TypeOf(pp.param)
	v := reflect.ValueOf(pp.param)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		fieldPtr := v.Field(i)
		if fieldPtr.Type().Kind() != reflect.Ptr {
			panic("parameter value field is expected ptr value.")
		}

		if fieldPtr.IsNil() {
			continue
		}

		fieldValue := v.Field(i).Elem()
		searchTag := field.Tag.Get("search")
		if searchTag != "" {
			fieldNameSnake := strcase.ToSnake(field.Name)
			namedParam[fieldNameSnake] = fieldValue.Interface()
		}
	}
	return namedParam
}

// ex
//type SearchMachinesParameter struct {
//	MachineNumber *int       `search:"machine_number" operator:"eq"`
//	MachineName   *string    `search:"machine_name" operator:"eq"`
//	BuyDateFrom   *time.Time `search:"buy_date" operator:"ge"`
//	BuyDateTo     *time.Time `search:"buy_date" operator:"lt"`
//	PriceFrom     *int `search:"price" operator:"gt"`
//	PriceTo       *int `search:"price" operator:"le"`
//}
func (pp *ParameterParser) ParseBindMap() []map[string]string {
	t := reflect.TypeOf(pp.param)
	v := reflect.ValueOf(pp.param)
	bindMap := make(map[string]map[string]string)
	dic := make([]string, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// nil は飛ばす
		fieldValue := v.Field(i)
		if fieldValue.Type().Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}

		searchTag := field.Tag.Get("search")
		operatorTag := field.Tag.Get("operator")
		if searchTag != "" && operatorTag != "" {
			fieldNameSnake := strcase.ToSnake(field.Name)
			dic = append(dic, fieldNameSnake)
			bindMap[fieldNameSnake] = map[string]string{
				"target":   searchTag,
				"operator": operatorTag,
			}
		}
	}

	sortedByFieldNumber := make([]map[string]string, 0, len(dic))
	for _, key := range dic {
		bindMap[key]["bind"] = key
		sortedByFieldNumber = append(sortedByFieldNumber, bindMap[key])
	}
	return sortedByFieldNumber
}
