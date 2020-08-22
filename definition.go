package goapidoc

// Model definitions
type Definition struct {
	name string
	desc string

	generics   []string
	properties []*Property
}

func NewDefinition(title string, desc string) *Definition {
	return &Definition{name: title, desc: desc}
}

func (d *Definition) Generics(generics ...string) *Definition {
	d.generics = generics
	return d
}

func (d *Definition) Properties(properties ...*Property) *Definition {
	d.properties = append(d.properties, properties...)
	return d
}

// Model property
type Property struct {
	name     string
	typ      string
	required bool
	desc     string

	allowEmpty bool
	def        interface{}
	example    interface{}
	enum       []interface{}
	minLength  int
	maxLength  int
	minimum    int
	maximum    int
}

func NewProperty(name string, typ string, req bool, desc string) *Property {
	return &Property{name: name, typ: typ, required: req, desc: desc}
}

func (p *Property) AllowEmpty(allow bool) *Property {
	p.allowEmpty = allow
	return p
}

func (p *Property) Default(def interface{}) *Property {
	p.def = def
	return p
}

func (p *Property) Example(ex interface{}) *Property {
	p.example = ex
	return p
}

func (p *Property) Enum(enum ...interface{}) *Property {
	p.enum = enum
	return p
}

func (p *Property) MinLength(min int) *Property {
	p.minLength = min
	return p
}

func (p *Property) MaxLength(max int) *Property {
	p.maxLength = max
	return p
}

func (p *Property) Minimum(min int) *Property {
	p.minimum = min
	return p
}

func (p *Property) Maximum(max int) *Property {
	p.maximum = max
	return p
}
