package schema

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	validate "github.com/go-playground/validator/v10"
)

type Validator struct {
	Validator *validate.Validate
	rules     map[string]interface{}
}

func NewValidator(v *validate.Validate) Validator {
	v.RegisterValidation("string", isString)
	v.RegisterValidation("boolean", isBoolean)
	return Validator{
		Validator: v,
		rules:     map[string]interface{}{},
	}
}

// Loads schema file and passes it through recursive function to generate a rules map
func (v *Validator) LoadSchema(name string, tomlSchema string) error {
	var raw map[string]interface{}
	_, err := toml.Decode(tomlSchema, &raw)
	if err != nil {
		return err
	}

	v.rules[name] = map[string]interface{}{}

	err = makeSchema(v.rules[name].(map[string]interface{}), raw)
	if err != nil {
		return err
	}

	return nil
}

func (v *Validator) ValidateSchema(name string, data map[string]interface{}) map[string]interface{} {
	// We need to re-cast some values to work with the validator https://github.com/go-playground/validator/issues/1108
	patchedData := data
	patchData(patchedData, data)

	return v.Validator.ValidateMap(data, v.rules[name].(map[string]interface{}))
}

func makeSchema(root map[string]interface{}, raw map[string]interface{}) error {
	for k, v := range raw {
		switch val := v.(type) {
		case string:
			if !slices.Contains([]string{"string", "number", "boolean"}, strings.Split(val, ",")[0]) {
				return fmt.Errorf("field %s is missing basic typing as first validator", k)
			}
			root[k] = val
		case map[string]interface{}:
			root[k] = map[string]interface{}{}
			makeSchema(root[k].(map[string]interface{}), val)
		case []map[string]interface{}:
			root[k] = map[string]interface{}{}
			schema := val[0]
			makeSchema(root[k].(map[string]interface{}), schema)
		default:
			return fmt.Errorf("could not parse schema field %v with value %v", k, val)
		}
	}
	return nil
}

// patchData is a temporary function until go-validator can patch #1108
func patchData(root map[string]interface{}, raw map[string]interface{}) {
	// We need to re-cast some values to work with the validator https://github.com/go-playground/validator/issues/1108
	for k, v := range raw {
		switch val := v.(type) {
		case map[string]interface{}:
			patchData(root[k].(map[string]interface{}), val)
		case []interface{}:
			new := []map[string]interface{}{}

			for _, value := range val {
				process := value.(map[string]interface{})
				patchData(process, process)
				new = append(new, process)
			}

			root[k] = new
		}
	}
}

func isString(fl validate.FieldLevel) bool {
	return fl.Field().Kind() == reflect.String
}

func isBoolean(fl validate.FieldLevel) bool {
	return fl.Field().Kind() == reflect.Bool
}
