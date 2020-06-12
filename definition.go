package goapidoc

// Model definitions
type Definition struct {
	Name        string
	Description string
	Generics    []string
	Properties  []*Property
}

func NewDefinition(title string, desc string) *Definition {
	return &Definition{Name: title, Description: desc}
}

func (d *Definition) WithGenerics(generics ...string) *Definition {
	d.Generics = generics
	return d
}

func (d *Definition) WithProperties(properties ...*Property) *Definition {
	d.Properties = append(d.Properties, properties...)
	return d
}

// Model property
type Property struct {
	// Deprecated
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
// Deprecated
func NewObjectProperty(name string, ref string, req bool) *Property {
	return &Property{Title: name, Schema: &Schema{Type: OBJECT, Ref: ref, Required: req}}
}

// sugar: array property
// Deprecated
func NewArrayProperty(name string, items *Items, req bool) *Property {
	return &Property{Title: name, Schema: &Schema{Type: ARRAY, Items: items, Required: req}}
}

func (p *Property) WithFormat(format string) *Property {
	p.Format = format
	return p
}

func (p *Property) WithAllowEmptyValue(allow bool) *Property {
	p.AllowEmptyValue = allow
	return p
}

func (p *Property) WithDefault(def interface{}) *Property {
	p.Default = def
	return p
}

func (p *Property) WithEnum(enum ...interface{}) *Property {
	p.Enum = enum
	return p
}

// Deprecated
func (p *Property) WithItems(items *Items) *Property {
	p.Items = items
	return p
}
