package types

var GRAPHQL_TYPE_MAP = map[string]string{
	"string":      "String",
	"boolean":     "Boolean",
	"date":        "String",
	"dateTime":    "String",
	"float":       "Float",
	"double":      "Double",
	"long":        "Long",
	"int":         "Int",
}

func GetGraphQlType(t string) string {

	value, ok := GRAPHQL_TYPE_MAP[t]
	if ok {
		return value
	} else {
		return t
	}
}
