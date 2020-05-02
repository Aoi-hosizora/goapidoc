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

// normal property
func NewProperty(name string, schemaType string, required bool, description string) *Schema {
	return &Schema{Title: name, Type: schemaType, Required: required, Description: description}
}

// object property
func NewPropertyObject(name string, ref string, required bool, description string) *Schema {
	return &Schema{Title: name, Type: "object", Ref: ref, Required: required, Description: description}
}

// array property
func NewPropertyArray(name string, items *Items, required bool, description string) *Schema {
	return &Schema{Title: name, Type: "array", Items: items, Required: required, Description: description}
}

// get default format from type
func defaultFormat(t string) string {
	if t == INTEGER {
		return INT32
	} else if t == NUMBER {
		return DOUBLE
	}
	return ""
}
