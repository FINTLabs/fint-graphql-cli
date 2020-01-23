package types

import "strings"

var JSON_TYPE_MAP = map[string]string{
	"string":    `"string"`,
	"boolean":   `"boolean"`,
	"date":      `"string", "format": "date-time"`,
	"dateTime":  `"string", "format": "date-time"`,
	"float":     `"number"`,
	"double":    `"number"`,
	"long":      `"integer"`,
	"int":       `"integer"`,
	"referanse": `"string", "format": "uri"`,
}

func GetJsonType(att *Attribute) string {

	result := att.Type

	value, ok := JSON_TYPE_MAP[att.Type]
	if ok {
		result = value
	} else {
		result = `"object", "$ref": "https://schema.fintlabs.no/` + GetComponentName(att.Package) + "/" + strings.ToLower(result) + `.json"`
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
