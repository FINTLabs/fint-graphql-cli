package types

var GRAPHQL_TYPE_MAP = map[string]string{
	"string":    "String",
	"boolean":   "Boolean",
	"date":      "String",
	"dateTime":  "String",
	"float":     "Float",
	"double":    "Double",
	"long":      "Long",
	"int":       "Int",
	"referanse": "String",
}

func GetGraphQlType(att *Attribute) string {

	result := att.Type

	value, ok := GRAPHQL_TYPE_MAP[att.Type]
	if ok {
		result = value
	}

	if att.List {
		result = "[" + result + "]"
	}

	if !att.Optional {
		result = result + "!"
	}

	return result
}
