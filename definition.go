package goapidoc

// ==========
// Definition
// ==========

// Definition represents an api definition information of Document.
type Definition struct {
	name string
	desc string

	xmlRepr    *XMLRepr
	generics   []string
	properties []*Property
}

// NewDefinition creates a default Definition with given arguments.
func NewDefinition(name, desc string) *Definition {
	return &Definition{name: name, desc: desc}
}

// GetName returns the name from Definition.
func (d *Definition) GetName() string { return d.name }

// GetDesc returns the desc from Definition.
func (d *Definition) GetDesc() string { return d.desc }

// GetXMLRepr returns the xml repr from Definition.
func (d *Definition) GetXMLRepr() *XMLRepr { return d.xmlRepr }

// GetGenerics returns the whole generics from Definition.
func (d *Definition) GetGenerics() []string { return d.generics }

// GetProperties returns the whole properties from Definition.
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

// XMLRepr sets the xml repr in Definition, this is only supported in Swagger.
func (d *Definition) XMLRepr(repr *XMLRepr) *Definition {
	d.xmlRepr = repr
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

// =======
// XMLRepr
// =======

// XMLRepr represents a xml representation format information of Definition, Property, Param an ItemOption, this is only supported in Swagger.
type XMLRepr struct {
	name      string
	namespace string
	prefix    string
	attribute bool
	wrapped   bool
}

// NewXMLRepr creates a default XMLRepr with given arguments.
func NewXMLRepr(name string) *XMLRepr {
	return &XMLRepr{name: name}
}

// GetName returns the name from XMLRepr.
func (x *XMLRepr) GetName() string { return x.name }

// GetNamespace returns the namespace from XMLRepr.
func (x *XMLRepr) GetNamespace() string { return x.namespace }

// GetPrefix returns the prefix from XMLRepr.
func (x *XMLRepr) GetPrefix() string { return x.prefix }

// GetAttribute returns the attribute from XMLRepr.
func (x *XMLRepr) GetAttribute() bool { return x.attribute }

// GetWrapped returns the wrapped from XMLRepr.
func (x *XMLRepr) GetWrapped() bool { return x.wrapped }

// Name sets the name in XMLRepr.
func (x *XMLRepr) Name(name string) *XMLRepr {
	x.name = name
	return x
}

// Namespace sets the namespace in XMLRepr.
func (x *XMLRepr) Namespace(namespace string) *XMLRepr {
	x.namespace = namespace
	return x
}

// Prefix sets the prefix in XMLRepr.
func (x *XMLRepr) Prefix(prefix string) *XMLRepr {
	x.prefix = prefix
	return x
}

// Attribute sets the attribute in XMLRepr.
func (x *XMLRepr) Attribute(attribute bool) *XMLRepr {
	x.attribute = attribute
	return x
}

// Wrapped sets the wrapped in XMLRepr.
func (x *XMLRepr) Wrapped(wrapped bool) *XMLRepr {
	x.wrapped = wrapped
	return x
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
	xmlRepr          *XMLRepr
}

// NewProperty creates a default Property with given arguments.
func NewProperty(name, typ string, required bool, desc string) *Property {
	return &Property{name: name, typ: typ, required: required, desc: desc}
}

// GetName returns the name from Property.
func (p *Property) GetName() string { return p.name }

// GetType returns the type from Property.
func (p *Property) GetType() string { return p.typ }

// GetRequired returns the required from Property.
func (p *Property) GetRequired() bool { return p.required }

// GetDesc returns the desc from Property.
func (p *Property) GetDesc() string { return p.desc }

// GetAllowEmpty returns the allowEmpty from Property.
func (p *Property) GetAllowEmpty() bool { return p.allowEmpty }

// GetDefault returns the default from Property.
func (p *Property) GetDefault() interface{} { return p.defaul }

// GetExample returns the example from Property.
func (p *Property) GetExample() interface{} { return p.example }

// GetPattern returns the pattern from Property.
func (p *Property) GetPattern() string { return p.pattern }

// GetEnum returns the enum from Property.
func (p *Property) GetEnum() []interface{} { return p.enum }

// GetMinLength returns the minLength from Property.
func (p *Property) GetMinLength() *int { return p.minLength }

// GetMaxLength returns the maxLength from Property.
func (p *Property) GetMaxLength() *int { return p.maxLength }

// GetMinItems returns the minItems from Property.
func (p *Property) GetMinItems() *int { return p.minItems }

// GetMaxItems returns the maxItems from Property.
func (p *Property) GetMaxItems() *int { return p.maxItems }

// GetUniqueItems returns the uniqueItems from Property.
func (p *Property) GetUniqueItems() bool { return p.uniqueItems }

// GetCollectionFormat returns the collectionFormat from Property.
func (p *Property) GetCollectionFormat() string { return p.collectionFormat }

// GetMinimum returns the minimum from Property.
func (p *Property) GetMinimum() *float64 { return p.minimum }

// GetMaximum returns the maximum from Property.
func (p *Property) GetMaximum() *float64 { return p.maximum }

// GetExclusiveMin returns the exclusiveMin from Property.
func (p *Property) GetExclusiveMin() bool { return p.exclusiveMin }

// GetExclusiveMax returns the exclusiveMax from Property.
func (p *Property) GetExclusiveMax() bool { return p.exclusiveMax }

// GetMultipleOf returns the multipleOf from Property.
func (p *Property) GetMultipleOf() float64 { return p.multipleOf }

// GetItemOption returns the item option from Property.
func (p *Property) GetItemOption() *ItemOption { return p.itemOption }

// GetXMLRepr returns the xml repr from Property.
func (p *Property) GetXMLRepr() *XMLRepr { return p.xmlRepr }

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
func (p *Property) Minimum(min float64) *Property {
	p.minimum = &min
	return p
}

// Maximum sets the maximum in Property.
func (p *Property) Maximum(max float64) *Property {
	p.maximum = &max
	return p
}

// ValueRange sets the minimum and maximum in Property.
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

// ItemOption sets the item option in Property, this is only supported in Swagger.
func (p *Property) ItemOption(itemOption *ItemOption) *Property {
	p.itemOption = itemOption
	return p
}

// XMLRepr sets the xml repr in Property, this is only supported in Swagger.
func (p *Property) XMLRepr(repr *XMLRepr) *Property {
	p.xmlRepr = repr
	return p
}

// ==========
// ItemOption
// ==========

// ItemOption represents an array type's item option of Param, Property and ItemOption itself, this is only supported in Swagger.
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
	xmlRepr          *XMLRepr
}

// NewItemOption creates a default ItemOption.
func NewItemOption() *ItemOption {
	return &ItemOption{}
}

// GetAllowEmpty returns the allowEmpty from ItemOption.
func (o *ItemOption) GetAllowEmpty() bool { return o.allowEmpty }

// GetDefault returns the default from ItemOption.
func (o *ItemOption) GetDefault() interface{} { return o.defaul }

// GetExample returns the example from ItemOption.
func (o *ItemOption) GetExample() interface{} { return o.example }

// GetPattern returns the pattern from ItemOption.
func (o *ItemOption) GetPattern() string { return o.pattern }

// GetEnum returns the enum from ItemOption.
func (o *ItemOption) GetEnum() []interface{} { return o.enum }

// GetMinLength returns the minLength from ItemOption.
func (o *ItemOption) GetMinLength() *int { return o.minLength }

// GetMaxLength returns the maxLength from ItemOption.
func (o *ItemOption) GetMaxLength() *int { return o.maxLength }

// GetMinItems returns the minItems from ItemOption.
func (o *ItemOption) GetMinItems() *int { return o.minItems }

// GetMaxItems returns the maxItems from ItemOption.
func (o *ItemOption) GetMaxItems() *int { return o.maxItems }

// GetUniqueItems returns the uniqueItems from ItemOption.
func (o *ItemOption) GetUniqueItems() bool { return o.uniqueItems }

// GetCollectionFormat returns the collectionFormat from ItemOption.
func (o *ItemOption) GetCollectionFormat() string { return o.collectionFormat }

// GetMinimum returns the minimum from ItemOption.
func (o *ItemOption) GetMinimum() *float64 { return o.minimum }

// GetMaximum returns the maximum from ItemOption.
func (o *ItemOption) GetMaximum() *float64 { return o.maximum }

// GetExclusiveMin returns the exclusiveMin from ItemOption.
func (o *ItemOption) GetExclusiveMin() bool { return o.exclusiveMin }

// GetExclusiveMax returns the exclusiveMax from ItemOption.
func (o *ItemOption) GetExclusiveMax() bool { return o.exclusiveMax }

// GetMultipleOf returns the multipleOf from ItemOption.
func (o *ItemOption) GetMultipleOf() float64 { return o.multipleOf }

// GetItemOption returns the itemOption from ItemOption.
func (o *ItemOption) GetItemOption() *ItemOption { return o.itemOption }

// GetXMLRepr returns the xmlRepr from ItemOption.
func (o *ItemOption) GetXMLRepr() *XMLRepr { return o.xmlRepr }

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

// ItemOption sets the item option in ItemOption, this is only supported in Swagger.
func (o *ItemOption) ItemOption(itemOption *ItemOption) *ItemOption {
	o.itemOption = itemOption
	return o
}

// XMLRepr sets the xml repr in ItemOption, this is only supported in Swagger.
func (o *ItemOption) XMLRepr(repr *XMLRepr) *ItemOption {
	o.xmlRepr = repr
	return o
}

// ========
// cloneXXX
// ========

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
		xmlRepr:          o.xmlRepr,
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
		xmlRepr:          p.xmlRepr,
	}
}
