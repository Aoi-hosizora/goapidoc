package yamldoc

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

func NewParam(name string, in string, paramType string, required bool, description string) *Param {
	return &Param{Name: name, In: in, Type: paramType, Required: required, Description: description}
}

func (p *Param) SetFormat(format string) *Param {
	p.Format = format
	return p
}

func (p *Param) SetAllowEmptyValue(allowEmptyValue bool) *Param {
	p.AllowEmptyValue = allowEmptyValue
	return p
}

func (p *Param) SetDefault(defaultValue interface{}) *Param {
	p.Default = defaultValue
	return p
}

func (p *Param) SetEnum(enum ...interface{}) *Param {
	p.Enum = enum
	return p
}

// !
func (p *Param) SetSchema(schema *Schema) *Param {
	p.Schema = schema
	return p
}

// !
func (p *Param) SetItems(items *Items) *Param {
	p.Items = items
	return p
}
