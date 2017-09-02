package qtypes_plugin

type Plugin struct {
	Base
	MyID		int
	Typ			string
	Pkg			string
	Version 	string
	Name 		string
}


func NewNamedPlugin(b Base, typ, pkg, name, version string) *Plugin {
	p := &Plugin{
		Base: b,
		Typ:   		typ,
		Pkg:  		pkg,
		Version:	version,
		Name: 		name,
	}
	return p
}
