package goapidoc

// ==========
// Definition
// ==========

// Definition represents an api definition information of Document.
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

// Generics sets the whole generics in Definition.
func (d *Definition) Generics(generics ...string) *Definition {
	d.generics = generics
	return d
}

// Properties sets the whole properties in Definition.
func (d *Definition) Properties(properties ...*Property) *Definition {
	d.properties = properties
	return d
}

// AddProperties add some properties into Definition.
func (d *Definition) AddProperties(properties ...*Property) *Definition {
	d.properties = append(d.properties, properties...)
	return d
}

// ========
// Property
// ========

// Property represents a definition property information of Definition.
type Property struct {
	name     string
	typ      string // string number integer boolean array
	required bool
	desc     string

	allowEmpty       bool
	defaul           interface{}
	example          interface{}
	pattern          string
	enums            []interface{}
	minLength        int
	maxLength        int
	minItems         int
	maxItems         int
	uniqueItems      bool
	collectionFormat string
	minimum          float64
	maximum          float64
	exclusiveMin     bool
	exclusiveMax     bool
	multipleOf       float64
}

// NewProperty creates a default Property with given arguments.
func NewProperty(name, typ string, required bool, desc string) *Property {
	return &Property{name: name, typ: typ, required: required, desc: desc}
}

func (p *Property) GetName() string             { return p.name }
func (p *Property) GetType() string             { return p.typ }
func (p *Property) GetRequired() bool           { return p.required }
func (p *Property) GetAllowEmpty() bool         { return p.allowEmpty }
func (p *Property) GetDefault() interface{}     { return p.defaul }
func (p *Property) GetExample() interface{}     { return p.example }
func (p *Property) GetPattern() string          { return p.pattern }
func (p *Property) GetEnums() []interface{}     { return p.enums }
func (p *Property) GetMinLength() int           { return p.minLength }
func (p *Property) GetMaxLength() int           { return p.maxLength }
func (p *Property) GetMinItems() int            { return p.minItems }
func (p *Property) GetMaxItems() int            { return p.maxItems }
func (p *Property) GetUniqueItems() bool        { return p.uniqueItems }
func (p *Property) GetCollectionFormat() string { return p.collectionFormat }
func (p *Property) GetMinimum() float64         { return p.minimum }
func (p *Property) GetMaximum() float64         { return p.maximum }
func (p *Property) GetExclusiveMin() bool       { return p.exclusiveMin }
func (p *Property) GetExclusiveMax() bool       { return p.exclusiveMax }
func (p *Property) GetMultipleOf() float64      { return p.multipleOf }

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

// Pattern sets the pattern in Property.
func (p *Property) Pattern(pattern string) *Property {
	p.pattern = pattern
	return p
}

// Enums sets the whole enums in Property.
// TODO BREAK CHANGES
func (p *Property) Enums(enums ...interface{}) *Property {
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

// LengthRange sets the minLength and maxLength in Property.
// TODO BREAK CHANGES
func (p *Property) LengthRange(min, max int) *Property {
	p.minLength = min
	p.maxLength = max
	return p
}

// MinItems sets the minItems in Property.
func (p *Property) MinItems(min int) *Property {
	p.minItems = min
	return p
}

// MaxItems sets the maxItems in Property.
func (p *Property) MaxItems(max int) *Property {
	p.maxItems = max
	return p
}

// ItemsRange sets the minItems and maxItems in Property.
func (p *Property) ItemsRange(min, max int) *Property {
	p.minItems = min
	p.maxItems = max
	return p
}

// UniqueItems sets the uniqueItems in Property.
func (p *Property) UniqueItems(unique bool) *Property {
	p.uniqueItems = unique
	return p
}

// CollectionFormat sets the collectionFormat in Property.
func (p *Property) CollectionFormat(collectionFormat string) *Property {
	p.collectionFormat = collectionFormat
	return p
}

// Minimum sets the minimum in Property.
// TODO BREAK CHANGES
func (p *Property) Minimum(min float64) *Property {
	p.minimum = min
	return p
}

// Maximum sets the maximum in Property.
// TODO BREAK CHANGES
func (p *Property) Maximum(max float64) *Property {
	p.maximum = max
	return p
}

// ValueRange sets the minimum and maximum in Property.
// TODO BREAK CHANGES
func (p *Property) ValueRange(min, max float64) *Property {
	p.minimum = min
	p.maximum = max
	return p
}

// ExclusiveMin sets the exclusiveMin in Property.
func (p *Property) ExclusiveMin(exclusiveMin bool) *Property {
	p.exclusiveMin = exclusiveMin
	return p
}

// ExclusiveMax sets the exclusiveMax in Property.
func (p *Property) ExclusiveMax(exclusiveMax bool) *Property {
	p.exclusiveMax = exclusiveMax
	return p
}

// MultipleOf sets the multipleOf in Property.
func (p *Property) MultipleOf(multipleOf float64) *Property {
	p.multipleOf = multipleOf
	return p
}

// cloneProperty clones the given Property.
func cloneProperty(p *Property) *Property {
	return &Property{
		name:             p.name,
		typ:              p.typ,
		required:         p.required,
		desc:             p.desc,
		allowEmpty:       p.allowEmpty,
		defaul:           p.defaul,
		example:          p.example,
		pattern:          p.pattern,
		enums:            p.enums,
		minLength:        p.minLength,
		maxLength:        p.maxLength,
		minItems:         p.minItems,
		maxItems:         p.maxItems,
		uniqueItems:      p.uniqueItems,
		collectionFormat: p.collectionFormat,
		minimum:          p.minimum,
		maximum:          p.maximum,
		exclusiveMin:     p.exclusiveMin,
		exclusiveMax:     p.exclusiveMax,
		multipleOf:       p.multipleOf,
	}
}

// createParamFromProperty clones the given Property to Param.
func createParamFromProperty(p *Property) *Param {
	return &Param{
		name: p.name,
		// in: p.in,
		typ:              p.typ,
		required:         p.required,
		desc:             p.desc,
		allowEmpty:       p.allowEmpty,
		defaul:           p.defaul,
		example:          p.example,
		pattern:          p.pattern,
		enums:            p.enums,
		minLength:        p.minLength,
		maxLength:        p.maxLength,
		minItems:         p.minItems,
		maxItems:         p.maxItems,
		uniqueItems:      p.uniqueItems,
		collectionFormat: p.collectionFormat,
		minimum:          p.minimum,
		maximum:          p.maximum,
		exclusiveMin:     p.exclusiveMin,
		exclusiveMax:     p.exclusiveMax,
		multipleOf:       p.multipleOf,
	}
}
