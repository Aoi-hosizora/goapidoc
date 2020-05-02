package yamldoc

type Model struct {
	Name        string
	Description string
	Properties  []*Schema
}

func NewModel(title string, description string) *Model {
	return &Model{Name: title, Description: description}
}

func (m *Model) SetProperties(schema ...*Schema) *Model {
	m.Properties = append(m.Properties, schema...)
	return m
}
