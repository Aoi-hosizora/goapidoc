package yamldoc

type Model struct {
	Title       string
	Type        string // object
	Description string
	Properties  []*Property
}

func NewModel(title string, description string) *Model {
	return &Model{Title: title, Type: "object", Description: description}
}

func (m *Model) SetProperties(properties ...*Property) *Model {
	m.Properties = properties
	return m
}

type Property struct {
	Name            string
	Description     string
	Type            string
	Required        bool
	AllowEmptyValue bool
	Format          string
	Enum            []interface{}

	Schema string // if `type` is object
	Items  *Items // is `type` is array
}

func NewProperty(name string, description string, propType string, required bool) *Property {
	return &Property{Name: name, Description: description, Type: propType, Required: required}
}

func (p *Property) SetAllowEmptyValue(allowEmptyValue bool) *Property {
	p.AllowEmptyValue = allowEmptyValue
	return p
}

func (p *Property) SetFormat(format string) *Property {
	p.Format = format
	return p
}

func (p *Property) SetEnum(enum ...interface{}) *Property {
	p.Enum = enum
	return p
}

func (p *Property) SetSchema(ref string) *Property {
	p.Schema = ref
	return p
}

func (p *Property) SetItems(item *Items) *Property {
	p.Items = item
	return p
}
