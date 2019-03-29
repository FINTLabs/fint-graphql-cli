package types

import "strings"

var JSON_TYPE_MAP = map[string]string{
	"string":    "string",
	"boolean":   "boolean",
	"date":      "string",
	"dateTime":  "string",
	"float":     "number",
	"double":    "number",
	"long":      "integer",
	"int":       "integer",
	"referanse": "string",
}

func GetJsonType(att *Attribute) string {

	result := att.Type

	value, ok := JSON_TYPE_MAP[att.Type]
	if ok {
		result = "\"" + value + "\""
	} else {
		result = `"object", "$ref": "https://fintlabs.no/schema/` + GetComponentName(att.Package) + "/" + strings.ToLower(result) + `.json"`
	}

	if att.List {
		return `"array",
				"items": { "type": ` + result + ` }`
	}

	return result
}

func GetJsonRelationType(rel *Association) string {
	if rel.Stereotype != "hovedklasse" {
		return "[String]"
	}

	result := rel.Target

	if rel.List {
		result = "[" + result + "]"
	}

	if !rel.Optional {
		result = result + "!"
	}

	return result
}
