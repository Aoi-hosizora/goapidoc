package goapidoc

// ==========
// Definition
// ==========

// TODO XML

// Definition represents an api definition information of Document.
type Definition struct {
	name string
	desc string

	generics   []string
	properties []*Property
}

// NewDefinition creates a default Definition with given arguments.
func NewDefinition(name, desc string) *Definition {
	return &Definition{name: name, desc: desc}
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
	enum             []interface{}
	minLength        *int
	maxLength        *int
	minItems         *int
	maxItems         *int
	uniqueItems      bool
	collectionFormat string
	minimum          *float64
	maximum          *float64
	exclusiveMin     bool
	exclusiveMax     bool
	multipleOf       float64
	itemOption       *ItemOption
}

// NewProperty creates a default Property with given arguments.
func NewProperty(name, typ string, required bool, desc string) *Property {
	return &Property{name: name, typ: typ, required: required, desc: desc}
}

func (p *Property) GetName() string             { return p.name }
func (p *Property) GetType() string             { return p.typ }
func (p *Property) GetRequired() bool           { return p.required }
func (p *Property) GetDesc() string             { return p.desc }
func (p *Property) GetAllowEmpty() bool         { return p.allowEmpty }
func (p *Property) GetDefault() interface{}     { return p.defaul }
func (p *Property) GetExample() interface{}     { return p.example }
func (p *Property) GetPattern() string          { return p.pattern }
func (p *Property) GetEnum() []interface{}      { return p.enum }
func (p *Property) GetMinLength() *int          { return p.minLength }
func (p *Property) GetMaxLength() *int          { return p.maxLength }
func (p *Property) GetMinItems() *int           { return p.minItems }
func (p *Property) GetMaxItems() *int           { return p.maxItems }
func (p *Property) GetUniqueItems() bool        { return p.uniqueItems }
func (p *Property) GetCollectionFormat() string { return p.collectionFormat }
func (p *Property) GetMinimum() *float64        { return p.minimum }
func (p *Property) GetMaximum() *float64        { return p.maximum }
func (p *Property) GetExclusiveMin() bool       { return p.exclusiveMin }
func (p *Property) GetExclusiveMax() bool       { return p.exclusiveMax }
func (p *Property) GetMultipleOf() float64      { return p.multipleOf }
func (p *Property) GetItemOption() *ItemOption  { return p.itemOption }

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

// Enum sets the whole enums in Property.
func (p *Property) Enum(enums ...interface{}) *Property {
	p.enum = enums
	return p
}

// MinLength sets the minLength in Property.
func (p *Property) MinLength(min int) *Property {
	p.minLength = &min
	return p
}

// MaxLength sets the maxLength in Property.
func (p *Property) MaxLength(max int) *Property {
	p.maxLength = &max
	return p
}

// LengthRange sets the minLength and maxLength in Property.
// TODO BREAK CHANGES
func (p *Property) LengthRange(min, max int) *Property {
	p.minLength = &min
	p.maxLength = &max
	return p
}

// MinItems sets the minItems in Property.
func (p *Property) MinItems(min int) *Property {
	p.minItems = &min
	return p
}

// MaxItems sets the maxItems in Property.
func (p *Property) MaxItems(max int) *Property {
	p.maxItems = &max
	return p
}

// ItemsRange sets the minItems and maxItems in Property.
func (p *Property) ItemsRange(min, max int) *Property {
	p.minItems = &min
	p.maxItems = &max
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
	p.minimum = &min
	return p
}

// Maximum sets the maximum in Property.
// TODO BREAK CHANGES
func (p *Property) Maximum(max float64) *Property {
	p.maximum = &max
	return p
}

// ValueRange sets the minimum and maximum in Property.
// TODO BREAK CHANGES
func (p *Property) ValueRange(min, max float64) *Property {
	p.minimum = &min
	p.maximum = &max
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

// ItemOption sets the itemOption in Property.
func (p *Property) ItemOption(itemOption *ItemOption) *Property {
	p.itemOption = itemOption
	return p
}

// ItemOption represents an array type's item option for Param and Property.
type ItemOption struct {
	allowEmpty       bool
	defaul           interface{}
	example          interface{}
	pattern          string
	enum             []interface{}
	minLength        *int
	maxLength        *int
	minItems         *int
	maxItems         *int
	uniqueItems      bool
	collectionFormat string
	minimum          *float64
	maximum          *float64
	exclusiveMin     bool
	exclusiveMax     bool
	multipleOf       float64
	itemOption       *ItemOption
}

// NewItemOption creates a default ItemOption with no option.
func NewItemOption() *ItemOption {
	return &ItemOption{}
}

func (o *ItemOption) GetAllowEmpty() bool         { return o.allowEmpty }
func (o *ItemOption) GetDefault() interface{}     { return o.defaul }
func (o *ItemOption) GetExample() interface{}     { return o.example }
func (o *ItemOption) GetPattern() string          { return o.pattern }
func (o *ItemOption) GetEnum() []interface{}      { return o.enum }
func (o *ItemOption) GetMinLength() *int          { return o.minLength }
func (o *ItemOption) GetMaxLength() *int          { return o.maxLength }
func (o *ItemOption) GetMinItems() *int           { return o.minItems }
func (o *ItemOption) GetMaxItems() *int           { return o.maxItems }
func (o *ItemOption) GetUniqueItems() bool        { return o.uniqueItems }
func (o *ItemOption) GetCollectionFormat() string { return o.collectionFormat }
func (o *ItemOption) GetMinimum() *float64        { return o.minimum }
func (o *ItemOption) GetMaximum() *float64        { return o.maximum }
func (o *ItemOption) GetExclusiveMin() bool       { return o.exclusiveMin }
func (o *ItemOption) GetExclusiveMax() bool       { return o.exclusiveMax }
func (o *ItemOption) GetMultipleOf() float64      { return o.multipleOf }
func (o *ItemOption) GetItemOption() *ItemOption  { return o.itemOption }

// AllowEmpty sets the allowEmpty in ItemOption.
func (o *ItemOption) AllowEmpty(allow bool) *ItemOption {
	o.allowEmpty = allow
	return o
}

// Default sets the default in ItemOption.
func (o *ItemOption) Default(defaul interface{}) *ItemOption {
	o.defaul = defaul
	return o
}

// Example sets the example in ItemOption.
func (o *ItemOption) Example(example interface{}) *ItemOption {
	o.example = example
	return o
}

// Pattern sets the pattern in ItemOption.
func (o *ItemOption) Pattern(pattern string) *ItemOption {
	o.pattern = pattern
	return o
}

// Enum sets the whole enums in ItemOption.
func (o *ItemOption) Enum(enums ...interface{}) *ItemOption {
	o.enum = enums
	return o
}

// MinLength sets the minLength in ItemOption.
func (o *ItemOption) MinLength(min int) *ItemOption {
	o.minLength = &min
	return o
}

// MaxLength sets the maxLength in ItemOption.
func (o *ItemOption) MaxLength(max int) *ItemOption {
	o.maxLength = &max
	return o
}

// LengthRange sets the minLength and maxLength in ItemOption.
func (o *ItemOption) LengthRange(min, max int) *ItemOption {
	o.minLength = &min
	o.maxLength = &max
	return o
}

// MinItems sets the minItems in ItemOption.
func (o *ItemOption) MinItems(min int) *ItemOption {
	o.minItems = &min
	return o
}

// MaxItems sets the maxItems in ItemOption.
func (o *ItemOption) MaxItems(max int) *ItemOption {
	o.maxItems = &max
	return o
}

// ItemsRange sets the minItems and maxItems in ItemOption.
func (o *ItemOption) ItemsRange(min, max int) *ItemOption {
	o.minItems = &min
	o.maxItems = &max
	return o
}

// UniqueItems sets the uniqueItems in ItemOption.
func (o *ItemOption) UniqueItems(unique bool) *ItemOption {
	o.uniqueItems = unique
	return o
}

// CollectionFormat sets the collectionFormat in ItemOption.
func (o *ItemOption) CollectionFormat(collectionFormat string) *ItemOption {
	o.collectionFormat = collectionFormat
	return o
}

// Minimum sets the minimum in ItemOption.
func (o *ItemOption) Minimum(min float64) *ItemOption {
	o.minimum = &min
	return o
}

// Maximum sets the maximum in ItemOption.
func (o *ItemOption) Maximum(max float64) *ItemOption {
	o.maximum = &max
	return o
}

// ValueRange sets the minimum and maximum in ItemOption.
func (o *ItemOption) ValueRange(min, max float64) *ItemOption {
	o.minimum = &min
	o.maximum = &max
	return o
}

// ExclusiveMin sets the exclusiveMin in ItemOption.
func (o *ItemOption) ExclusiveMin(exclusiveMin bool) *ItemOption {
	o.exclusiveMin = exclusiveMin
	return o
}

// ExclusiveMax sets the exclusiveMax in ItemOption.
func (o *ItemOption) ExclusiveMax(exclusiveMax bool) *ItemOption {
	o.exclusiveMax = exclusiveMax
	return o
}

// MultipleOf sets the multipleOf in ItemOption.
func (o *ItemOption) MultipleOf(multipleOf float64) *ItemOption {
	o.multipleOf = multipleOf
	return o
}

// ItemOption sets the itemOption in ItemOption.
func (o *ItemOption) ItemOption(itemOption *ItemOption) *ItemOption {
	o.itemOption = itemOption
	return o
}

// cloneItemOption clones the given ItemOption recursively.
func cloneItemOption(o *ItemOption) *ItemOption {
	if o == nil {
		return nil
	}
	return &ItemOption{
		allowEmpty:       o.allowEmpty,
		defaul:           o.defaul,
		example:          o.example,
		pattern:          o.pattern,
		enum:             o.enum,
		minLength:        o.minLength,
		maxLength:        o.maxLength,
		minItems:         o.minItems,
		maxItems:         o.maxItems,
		uniqueItems:      o.uniqueItems,
		collectionFormat: o.collectionFormat,
		minimum:          o.minimum,
		maximum:          o.maximum,
		exclusiveMin:     o.exclusiveMin,
		exclusiveMax:     o.exclusiveMax,
		multipleOf:       o.multipleOf,
		itemOption:       cloneItemOption(o.itemOption),
	}
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
		enum:             p.enum,
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
		itemOption:       cloneItemOption(p.itemOption),
	}
}
