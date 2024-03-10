package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	validate "github.com/go-playground/validator/v10"
)

func TestValidator_LoadSchema(t *testing.T) {
	tomlSchema := NewValidator(validate.New())
	loadTestSchema(&tomlSchema, "character")

	control := map[string]interface{}{
		"age":             "required,number,min=1,max=1500",
		"credits":         "number,required,min=0,max=150000",
		"force_sensitive": "required",
		"location": map[string]interface{}{
			"address1": "required",
			"address2": "required",
		},
		"name": "required,min=1,max=128",
		"ships": map[string]interface{}{
			"id":   "required,uuid",
			"make": "oneOf=x-wing y-wing a-wing millenium falcon tie-fighter",
		},
	}
	if !reflect.DeepEqual(control, tomlSchema.rules["character"]) {
		t.Errorf("\ngot %v\nwant %v", tomlSchema.rules["character"], control)
	}
}

func TestValidator_LoadSchema_BadSchema(t *testing.T) {
	tomlSchema := NewValidator(validate.New())
	err := tomlSchema.LoadSchema("test", `
	test =
	`)
	if err == nil {
		t.Errorf("expected bad schema load to raise an error")
	}
}

func TestValidator_ValidateSchema(t *testing.T) {
	tomlSchema := NewValidator(validate.New())
	loadTestSchema(&tomlSchema, "basic")
	loadTestSchema(&tomlSchema, "character")
	data := loadTestData()

	cases := []struct {
		schemaName  string
		data        map[string]interface{}
		expectedErr int
	}{
		{
			schemaName:  "character",
			data:        data["character"].(map[string]interface{})["1"].(map[string]interface{}),
			expectedErr: 1,
		},
		{
			schemaName:  "basic",
			data:        data["basic"].(map[string]interface{})["1"].(map[string]interface{}),
			expectedErr: 0,
		},
		{
			schemaName:  "basic",
			data:        data["basic"].(map[string]interface{})["2"].(map[string]interface{}),
			expectedErr: 4,
		},
		{
			schemaName:  "basic",
			data:        data["basic"].(map[string]interface{})["3"].(map[string]interface{}),
			expectedErr: 4,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.schemaName, func(t *testing.T) {
			t.Parallel()
			errs := tomlSchema.ValidateSchema(tt.schemaName, tt.data)

			got := len(errs)
			if got != tt.expectedErr {
				t.Errorf("\ngot %v validation errors\nwant %v\n%v", got, tt.expectedErr, errs)
			}
		})
	}
}

func loadTestSchema(v *validator, name string) string {
	toml, err := os.ReadFile(fmt.Sprintf("./testdata/%s.toml", name))
	if err != nil {
		panic(err)
	}

	err = v.LoadSchema(name, string(toml))
	if err != nil {
		panic(err)
	}

	return name
}

func loadTestData() map[string]interface{} {
	jsonBytes, err := os.ReadFile("./testdata/data.json")
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{}
	json.Unmarshal(jsonBytes, &data)
	return data
}