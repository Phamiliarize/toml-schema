# toml-schema

`toml-schema` utilizes TOML to define validation schemas for input in a simple manner. Under the hood, we generate a `map[string]interface{}` rule definition that can be used by [go-playground/validator](https://github.com/go-playground/validator).

It was created for use in the [sabaresu](https://github.com/Phamiliarize/sabaresu) serverless framework to enable configuration-based validation.


# Getting Started

## Schema

Schema are written in [TOML](https://toml.io/), like so:

```toml
# key = "validators"
name = "required,min=1,max=128"
age = "required,number,min=1,max=150"
ecash = "number,required,min=0,max=150000"
is_cool = "required"

# A nested map can be added with [key]
[location]
address1 = "required"
address2 = "required"

# A slice is added with [[key]]

[[ships]]
id = "required,uuid"
make = "oneof=x-wing y-wing a-wing millenium falcon tie-fighter"
```


## Validate

To validate data against your schema and constraints:

```go
dogShelter := LoadSchema(dogShelterSchema)
err := dogShelter.Validate(myData)
if err != nil {
    panic(err)
}
```


## Supported Types & Constraints
`toml-schema` bases it's typing with [JSON](https://www.w3schools.com/js/js_json_datatypes.asp) and [Lua](https://www.lua.org/pil/2.html) in mind; this covers the majority of cases.


| Type | GoLang Type |
| ---- | ---- |
| `string` | `string` |
| `number` | `float64` |
| `boolean` | `bool` |
| `array` | `[]interface{}` |
| `object` | `map[string]interface{}` |
| `nil` | `nil` |

The supported constraints/validations are basically inline with that types offerings for [go-playground/validator](https://github.com/go-playground/validator?tab=readme-ov-file#baked-in-validations).