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
	return &Param{
		Name: name, In: in, Required: req, Description: desc,
		Type: t, Format: defaultFormat(t),
	}
}

func (p *Param) SetFormat(format string) *Param {
	p.Format = format
	return p
}

func (p *Param) SetAllowEmptyValue(allow bool) *Param {
	p.AllowEmptyValue = allow
	return p
}

func (p *Param) SetDefault(def interface{}) *Param {
	p.Default = def
	return p
}

func (p *Param) SetEnum(enum ...interface{}) *Param {
	p.Enum = enum
	return p
}

// Set object (when `in` == body)
func (p *Param) SetSchema(schema *Schema) *Param {
	p.Type = ""
	p.Schema = schema
	return p
}

// Set array (when `type` == array)
func (p *Param) SetItems(items *Items) *Param {
	p.Type = ARRAY
	p.Items = items
	return p
}
