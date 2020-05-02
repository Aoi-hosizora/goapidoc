package yamldoc

type Param struct {
	Name        string
	In          string
	Type        string
	Required    bool
	Description string

	Default interface{}
	Schema  string
	Enums   []interface{}
}

func NewParam(name string, in string, paramType string, required bool, description string) *Param {
	return &Param{Name: name, In: in, Type: paramType, Required: required, Description: description}
}

func (p *Param) SetDefault(defaultValue interface{}) *Param {
	p.Default = defaultValue
	return p
}

func (p *Param) SetSchema(schema string) *Param {
	p.Schema = schema
	return p
}

func (p *Param) SetEnums(enums ...interface{}) *Param {
	p.Enums = enums
	return p
}
