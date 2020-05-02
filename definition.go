package yamldoc

// Model definitions
type Definition struct {
	Name        string
	Description string
	Properties  []*Property
}

func NewDefinition(title string, description string) *Definition {
	return &Definition{Name: title, Description: description}
}

func (m *Definition) SetProperties(properties ...*Property) *Definition {
	m.Properties = append(m.Properties, properties...)
	return m
}

// Model property
type Property struct {
	*Schema
	Title string
}

// normal property
func NewProperty(name string, schemaType string, required bool, description string) *Property {
	return &Property{Title: name, Schema: &Schema{Type: schemaType, Required: required, Description: description}}
}

// sugar: object property
func NewObjectProperty(name string, ref string, required bool, description string) *Property {
	return &Property{Title: name, Schema: &Schema{Type: OBJECT, Ref: ref, Required: required, Description: description}}
}

// sugar: array property
func NewArrayProperty(name string, items *Items, required bool, description string) *Property {
	return &Property{Title: name, Schema: &Schema{Type: ARRAY, Items: items, Required: required, Description: description}}
}

func (p *Property) SetFormat(format string) *Property {
	p.Format = format
	return p
}

func (p *Property) SetAllowEmptyValue(allowEmptyValue bool) *Property {
	p.AllowEmptyValue = allowEmptyValue
	return p
}

func (p *Property) SetDefault(defaultValue interface{}) *Property {
	p.Default = defaultValue
	return p
}

func (p *Property) SetEnum(enum ...interface{}) *Property {
	p.Enum = enum
	return p
}

func (p *Property) SetRef(ref string) *Property {
	p.Ref = ref
	return p
}

func (p *Property) SetItems(items *Items) *Property {
	p.Items = items
	return p
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
