package schema

import (
	"fmt"

	"github.com/BurntSushi/toml"
	validate "github.com/go-playground/validator/v10"
)

type validator struct {
	Validator *validate.Validate
	rules     map[string]interface{}
}

func NewValidator(v *validate.Validate) validator {
	return validator{
		Validator: v,
		rules:     map[string]interface{}{},
	}
}

// Loads schema file and passes it through recursive function to generate a rules map
func (v *validator) LoadSchema(name string, tomlSchema string) error {
	var raw map[string]interface{}
	_, err := toml.Decode(tomlSchema, &raw)
	if err != nil {
		return err
	}

	v.rules[name] = map[string]interface{}{}

	makeSchema(v.rules[name].(map[string]interface{}), raw)

	return nil
}

func (v *validator) ValidateSchema(name string, data map[string]interface{}) map[string]interface{} {
	// We need to re-cast some values to work with the validator https://github.com/go-playground/validator/issues/1108
	patchedData := data
	patchData(patchedData, data)

	return v.Validator.ValidateMap(data, v.rules[name].(map[string]interface{}))
}

func makeSchema(root map[string]interface{}, raw map[string]interface{}) {
	for k, v := range raw {
		switch val := v.(type) {
		case string:
			root[k] = val
		case map[string]interface{}:
			root[k] = map[string]interface{}{}
			makeSchema(root[k].(map[string]interface{}), val)
		case []map[string]interface{}:
			root[k] = map[string]interface{}{}
			schema := val[0]
			makeSchema(root[k].(map[string]interface{}), schema)
		default:
			panic(fmt.Sprintf("could not parse schema field %v with value %v", k, val))
		}
	}
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
