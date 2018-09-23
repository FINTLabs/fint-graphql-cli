package types

type Association struct {
	Name          string
	Target        string
	TargetPackage string
	Deprecated    bool
	Multiplicity  string
}
