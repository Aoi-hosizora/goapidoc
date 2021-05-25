package goapidoc

// Definition represents an api definition information.
type Definition struct {
	name string
	desc string

	generics   []string
	properties []*Property
}

// NewDefinition creates a default Definition with given arguments.
func NewDefinition(title, desc string) *Definition {
	return &Definition{name: title, desc: desc}
}

func (d *Definition) GetName() string            { return d.name }
func (d *Definition) GetDesc() string            { return d.desc }
func (d *Definition) GetGenerics() []string      { return d.generics }
func (d *Definition) GetProperties() []*Property { return d.properties }

// Name sets the name in Definition.
func (d *Definition) Name(name string) *Definition {
	d.name = name
	return d
}

// Desc sets the desc in Definition.
func (d *Definition) Desc(desc string) *Definition {
	d.desc = desc
	return d
}

// Generics sets the generics in Definition.
func (d *Definition) Generics(generics ...string) *Definition {
	d.generics = generics
	return d
}

// Properties sets the whole properties in Definition.
func (d *Definition) Properties(properties ...*Property) *Definition {
	d.properties = properties
	return d
}

// Properties add some properties into Definition.
func (d *Definition) AddProperties(properties ...*Property) *Definition {
	d.properties = append(d.properties, properties...)
	return d
}

// Property represents a definition property information.
type Property struct {
	name     string
	typ      string
	required bool
	desc     string

	allowEmpty bool
	defaul     interface{}
	example    interface{}
	enums      []interface{}
	minLength  int
	maxLength  int
	minimum    int
	maximum    int
}

// NewProperty creates a default Property with given arguments.
func NewProperty(name, typ string, required bool, desc string) *Property {
	return &Property{name: name, typ: typ, required: required, desc: desc}
}

func (p *Property) GetName() string         { return p.name }
func (p *Property) GetType() string         { return p.typ }
func (p *Property) GetRequired() bool       { return p.required }
func (p *Property) GetDesc() string         { return p.desc }
func (p *Property) GetAllowEmpty() bool     { return p.allowEmpty }
func (p *Property) GetDefault() interface{} { return p.defaul }
func (p *Property) GetExample() interface{} { return p.example }
func (p *Property) GetEnums() []interface{} { return p.enums }
func (p *Property) GetMinLength() int       { return p.minLength }
func (p *Property) GetMaxLength() int       { return p.maxLength }
func (p *Property) GetMinimum() int         { return p.minimum }
func (p *Property) GetMaximum() int         { return p.maximum }

// Name sets the name in Property.
func (p *Property) Name(name string) *Property {
	p.name = name
	return p
}

// Type sets the type in Property.
func (p *Property) Type(typ string) *Property {
	p.typ = typ
	return p
}

// Required sets the required in Property.
func (p *Property) Required(required bool) *Property {
	p.required = required
	return p
}

// Desc sets the desc in Property.
func (p *Property) Desc(desc string) *Property {
	p.desc = desc
	return p
}

// AllowEmpty sets the allowEmpty in Property.
func (p *Property) AllowEmpty(allow bool) *Property {
	p.allowEmpty = allow
	return p
}

// Default sets the default in Property.
func (p *Property) Default(defaul interface{}) *Property {
	p.defaul = defaul
	return p
}

// Example sets the example in Property.
func (p *Property) Example(example interface{}) *Property {
	p.example = example
	return p
}

// Enum sets the enums in Property.
func (p *Property) Enum(enums ...interface{}) *Property {
	p.enums = enums
	return p
}

// MinLength sets the minLength in Property.
func (p *Property) MinLength(min int) *Property {
	p.minLength = min
	return p
}

// MaxLength sets the maxLength in Property.
func (p *Property) MaxLength(max int) *Property {
	p.maxLength = max
	return p
}

// Length sets the minLength and maxLength in Property.
func (p *Property) Length(min, max int) *Property {
	p.minLength = min
	p.maxLength = max
	return p
}

// Minimum sets the minimum in Property.
func (p *Property) Minimum(min int) *Property {
	p.minimum = min
	return p
}

// Maximum sets the maximum in Property.
func (p *Property) Maximum(max int) *Property {
	p.maximum = max
	return p
}

// Minimum sets the minimum and maximum in Property.
func (p *Property) MinMaximum(min, max int) *Property {
	p.minimum = min
	p.maximum = max
	return p
}
