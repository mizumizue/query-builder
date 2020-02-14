package object_parser

import (
	"reflect"

	"github.com/iancoleman/strcase"
)

type ObjectParser struct {
	object      interface{}
	objectType  reflect.Type
	objectValue reflect.Value
}

func NewObjectParser(object interface{}) *ObjectParser {
	return &ObjectParser{
		object:      object,
		objectType:  reflect.TypeOf(object),
		objectValue: reflect.ValueOf(object),
	}
}

func (objectParser *ObjectParser) NamedParam() map[string]interface{} {
	namedParam := make(map[string]interface{})
	t, v := objectParser.getTypeAndValue()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldPtr := v.Field(i)

		if fieldPtr.Type().Kind() != reflect.Ptr {
			panic("objecteter value field is expected ptr value.")
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
func (objectParser *ObjectParser) SearchBindMap() []map[string]string {
	t, v := objectParser.getTypeAndValue()
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

func (objectParser *ObjectParser) getTypeAndValue() (reflect.Type, reflect.Value) {
	if objectParser.objectType.Kind() == reflect.Ptr {
		return objectParser.objectType.Elem(), objectParser.objectValue.Elem()
	}
	return objectParser.objectType, objectParser.objectValue
}
