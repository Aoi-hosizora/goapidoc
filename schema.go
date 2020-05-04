package apidoc

// !! Only `schema.go` has `Ref` definition: `NewSchemaRef` `NewItemsRef`

// The Schema Object allows the definition of input and output data types.
// These types can be objects, but also primitives and arrays.
//
// Primitive example:
//     NewSchema("integer", true).SetDefault(0)
// Object example: where user is defined in #/Definitions
//     NewSchemaRef("User")
// Array example:
//     NewSchema("array", true).SetItems(NewItems("integer"))                             // -> array of integer
//     NewSchema("array", true).SetItems(NewItemsRef("User"))                             // -> array of object
//     NewSchema("array", true).SetItems(NewItems("array").SetItems(NewItems("integer"))) // -> array of array
type Schema struct {
	Type     string
	Required bool

	Description     string
	Format          string
	AllowEmptyValue bool
	Default         interface{}
	Enum            []interface{}

	Options []*AdditionOption // for ref
	Ref     string            // `type` == object
	Items   *Items            // `type` == array
}

// Schema for response and parameter
func NewSchema(schemaType string, required bool) *Schema {
	return &Schema{Required: required, Type: schemaType, Format: defaultFormat(schemaType)}
}

// $ref
func NewSchemaRef(ref string, options ...*AdditionOption) *Schema {
	return &Schema{Ref: ref, Options: options}
}

func (s *Schema) SetDescription(description string) *Schema {
	s.Description = description
	return s
}

func (s *Schema) SetFormat(format string) *Schema {
	s.Format = format
	return s
}

func (s *Schema) SetAllowEmptyValue(allowEmptyValue bool) *Schema {
	s.AllowEmptyValue = allowEmptyValue
	return s
}

func (s *Schema) SetDefault(defaultValue interface{}) *Schema {
	s.Default = defaultValue
	return s
}

func (s *Schema) SetEnum(enum ...interface{}) *Schema {
	s.Enum = enum
	return s
}

func (s *Schema) SetItems(items *Items) *Schema {
	s.Type = ARRAY
	s.Items = items
	return s
}

// A limited subset of JSON-Schema's items object.
// It is used by parameter definitions that are not located in "body" -> should use SetSchema()
type Items struct {
	Type string

	Format  string
	Default interface{}
	Enum    []interface{}

	Options []*AdditionOption // for ref
	Ref     string
	Items   *Items // `type` == array
}

// Items for response and parameter
func NewItems(itemType string) *Items {
	return &Items{Type: itemType, Format: defaultFormat(itemType)}
}

// $ref
func NewItemsRef(ref string, options ...*AdditionOption) *Items {
	return &Items{Ref: ref, Options: options}
}

func (i *Items) SetFormat(format string) *Items {
	i.Format = format
	return i
}

func (i *Items) SetDefault(defaultValue interface{}) *Items {
	i.Default = defaultValue
	return i
}

func (i *Items) SetEnum(enum ...interface{}) *Items {
	i.Enum = enum
	return i
}

func (i *Items) SetItems(items *Items) *Items {
	i.Type = ARRAY
	i.Items = items
	return i
}

// Addition schema option
type AdditionOption struct {
	Field  string
	Schema *Schema
	Items  *Items
}

func NewSchemaOption(field string, schema *Schema) *AdditionOption {
	return &AdditionOption{Field: field, Schema: schema}
}

func NewItemsOption(field string, items *Items) *AdditionOption {
	return &AdditionOption{Field: field, Items: items}
}
