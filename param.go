package apidoc

type Param struct {
	Name        string
	In          string
	Type        string // string number integer boolean array (file)
	Required    bool
	Description string

	// `in` != body
	Format          string
	AllowEmptyValue bool
	Default         interface{}
	Enum            []interface{}

	Schema *Schema // `in` == body
	Items  *Items  // `in` != body && `type` == array
}

func NewParam(name string, in string, t string, req bool, desc string) *Param {
	return &Param{Name: name, In: in, Required: req, Description: desc, Type: t, Format: defaultFormat(t)}
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

func (p *Param) WithFormat(format string) *Param {
	p.Format = format
	return p
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

// Set object (when `in` == body)
// Deprecated
func (p *Param) WithSchema(schema *Schema) *Param {
	p.Type = ""
	p.Schema = schema
	return p
}

// Set array (when `type` == array)
// Deprecated
func (p *Param) WithItems(items *Items) *Param {
	p.Type = ARRAY
	p.Items = items
	return p
}
