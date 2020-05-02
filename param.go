package yamldoc

type Param struct {
	Name        string
	In          string
	Required    bool
	Description string

	// when `in` is any value other tha body
	Type            string // string number integer boolean array file
	Format          string
	AllowEmptyValue bool
	Default         interface{}
	Enum            []interface{}

	Schema string // when `in` is body, simplify $ref
	Items  *Items // when `type` is array, simplify items
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

func (p *Param) SetSchema(ref string) *Param {
	p.Schema = ref
	return p
}

func (p *Param) SetItems(items *Items) *Param {
	p.Items = items
	return p
}

type Items struct {
	Type    string // string number integer boolean array
	Format  string
	Default interface{}
	Schema  string
	Items   *Items
}

func NewItems(itemType string) *Items {
	return &Items{Type: itemType}
}

func (i *Items) SetFormat(format string) *Items {
	i.Format = format
	return i
}

func (i *Items) SetDefault(defaultValue string) *Items {
	i.Default = defaultValue
	return i
}

func (i *Items) SetSchema(ref string) *Items {
	i.Schema = ref
	return i
}

func (i *Items) SetItems(items *Items) *Items {
	i.Items = items
	return i
}
