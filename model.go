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
	Title           string
	Description     string
	Type            string
	Required        bool
	AllowEmptyValue bool
	Format          string

	Schema string
	Enum   []string
}

func NewProperty(title string, description string, propType string, required bool) *Property {
	return &Property{Title: title, Description: description, Type: propType, Required: required}
}

func (p *Property) SetAllowEmptyValue(allowEmptyValue bool) *Property {
	p.AllowEmptyValue = allowEmptyValue
	return p
}

func (p *Property) SetFormat(format string) *Property {
	p.Format = format
	return p
}

func (p *Property) SetSchema(schema string) *Property {
	p.Schema = schema
	return p
}

func (p *Property) SetEnum(enum ...string) *Property {
	p.Enum = enum
	return p
}
