# toml-schema

`toml-schema` utilizes TOML to define validation schemas for input in a simple manner. Under the hood, we generate a `map[string]interface{}` rule definition that can be used by [go-playground/validator](https://github.com/go-playground/validator).

It was created for use in the [sabaresu](https://github.com/Phamiliarize/sabaresu) serverless framework to enable configuration-based validation.


# Getting Started

## Schema

Schema are written in [TOML](https://toml.io/), like so:

```toml
# key = "validators"
name = "string,required,min=1,max=128"
age = "number,required,int,min=1,max=150"
ecash = "number,required,min=0,max=150000"
is_cool = "boolean,required"

# A nested map can be added with [key]
[location]
address1 = "string,required"
address2 = "string,required"

# A slice is added with [[key]]

[[ships]]
id = "string,required,uuid"
make = "string,oneof=x-wing y-wing a-wing millenium falcon tie-fighter"
```


## Validate

To validate data against your schema and constraints:

```go
import (
    "github.com/Phamiliarize/toml-schema"
    "github.com/go-playground/validator/v10"
)

// Use any v10 instance of a go-playground/validator
// This allows you to extend the potential validators
v := schema.NewValidator(validator.New())
err := v.LoadSchema("mySchema", `test = "string,required"`)
if err != nil {
    panic(err)
}

// Returns a map[string]validation.ValidationErrors
err := v.ValidateSchema("mySchema", myData)
if err != nil {
    panic(err)
}
```


## Basic Types & Constraints
The supported constraints/validations are basically inline with [go-playground/validator](https://github.com/go-playground/validator?tab=readme-ov-file#baked-in-validations) *but* since the end usecase is for [JSON](https://www.w3schools.com/js/js_json_datatypes.asp) and [Lua](https://www.lua.org/pil/2.html) we require a "type constraint" on all fields.

Every field **must start with a basic type**. The basic types supported are:

| Type | GoLang Type | Validators | 
| ---- | ---- | ---- |
| string | `string` | `string` |
| number | `float64` or `int` | `number`|
| boolean | `bool` | `boolean` |
