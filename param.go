package goapidoc

type Param struct {
	Name        string
	In          string
	Type        string // string number integer boolean array (file)
	Required    bool
	Description string

	// `in` != body
	AllowEmptyValue bool
	Default         interface{}
	Enum            []interface{}
}

func NewParam(name string, in string, t string, req bool, desc string) *Param {
	return &Param{Name: name, In: in, Required: req, Description: desc, Type: t}
}

func NewPathParam(name string, t string, req bool, desc string) *Param {
	return NewParam(name, PATH, t, req, desc)
}

func NewQueryParam(name string, t string, req bool, desc string) *Param {
	return NewParam(name, QUERY, t, req, desc)
}

func NewFormParam(name string, t string, req bool, desc string) *Param {
	return NewParam(name, FORM, t, req, desc)
}

func NewBodyParam(name string, t string, req bool, desc string) *Param {
	return NewParam(name, BODY, t, req, desc)
}

func NewHeaderParam(name string, t string, req bool, desc string) *Param {
	return NewParam(name, HEADER, t, req, desc)
}

func (p *Param) WithAllowEmptyValue(allow bool) *Param {
	p.AllowEmptyValue = allow
	return p
}

func (p *Param) WithDefault(def interface{}) *Param {
	p.Default = def
	return p
}

func (p *Param) WithEnum(enum ...interface{}) *Param {
	p.Enum = enum
	return p
}
