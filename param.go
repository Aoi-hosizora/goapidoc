package yamldoc

type Param struct {
	Name        string
	In          string
	Type        string
	Required    bool
	Description string

	Format  string
	Default interface{}
	Schema  string
	Enum    []interface{}
}

func NewParam(name string, in string, paramType string, required bool, description string) *Param {
	return &Param{Name: name, In: in, Type: paramType, Required: required, Description: description}
}

func (p *Param) SetFormat(format string) *Param {
	p.Format = format
	return p
}

func (p *Param) SetDefault(defaultValue interface{}) *Param {
	p.Default = defaultValue
	return p
}

func (p *Param) SetSchema(schema string) *Param {
	p.Schema = schema
	return p
}

func (p *Param) SetEnum(enum ...interface{}) *Param {
	p.Enum = enum
	return p
}
