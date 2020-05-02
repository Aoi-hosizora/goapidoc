package yamldoc

type Model struct {
	Title       string
	Description string
	Properties  []*Property
}

func NewModel(title string, description string) *Model {
	return &Model{Title: title, Description: description}
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

	Ref   string
	Enums []string
}

func NewProperty(title string, description string, propType string, required bool) *Property {
	return &Property{Title: title, Description: description, Type: propType, Required: required}
}

func (p *Property) SetAllowEmptyValue(allowEmptyValue bool) *Property {
	p.AllowEmptyValue = allowEmptyValue
	return p
}

func (p *Property) SetRef(ref string) *Property {
	p.Ref = ref
	return p
}

func (p *Property) SetEnums(enums ...string) *Property {
	p.Enums = enums
	return p
}
