package goapidoc

// Model definitions
type Definition struct {
	Name        string
	Description string
	Properties  []*Property
}

func NewDefinition(title string, desc string) *Definition {
	return &Definition{Name: title, Description: desc}
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
func NewProperty(name string, t string, req bool, desc string) *Property {
	return &Property{Title: name, Schema: &Schema{
		Type: t, Format: defaultFormat(t), Required: req, Description: desc,
	}}
}

// sugar: object property
func NewObjectProperty(name string, ref string, req bool) *Property {
	return &Property{Title: name, Schema: &Schema{Type: OBJECT, Ref: ref, Required: req}}
}

// sugar: array property
func NewArrayProperty(name string, items *Items, req bool) *Property {
	return &Property{Title: name, Schema: &Schema{Type: ARRAY, Items: items, Required: req}}
}

func (p *Property) SetFormat(format string) *Property {
	p.Format = format
	return p
}

func (p *Property) SetAllowEmptyValue(allow bool) *Property {
	p.AllowEmptyValue = allow
	return p
}

func (p *Property) SetDefault(def interface{}) *Property {
	p.Default = def
	return p
}

func (p *Property) SetEnum(enum ...interface{}) *Property {
	p.Enum = enum
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
