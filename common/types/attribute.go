package types

type Attribute struct {
	Name       string
	Type       string
	Package    string
	List       bool
	Optional   bool
	Deprecated bool
}

type InheritedAttribute struct {
	Owner string
	Attribute
}
