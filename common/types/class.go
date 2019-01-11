package types

type Class struct {
	Name                      string
	Abstract                  bool
	Deprecated                bool
	Extends                   string
	Package                   string
	Imports                   []string
	Namespace                 string
	Using                     []string
	Documentation             string
	Attributes                []Attribute
	AttributesWithInheritance []Attribute
	Relations                 []Association
	Resources                 []Attribute
	Resource                  bool
	ExtendsResource           bool
	Identifiable              bool
	GitTag                    string
	Stereotype                string
	Identifiers               []Identifier
}
