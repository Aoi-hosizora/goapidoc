package apidoc

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

	Ref   string // `type` == object
	Items *Items // `type` == array
}

func NewSchema(schemaType string, required bool) *Schema {
	return &Schema{Required: required, Type: schemaType, Format: defaultFormat(schemaType)}
}

func NewSchemaRef(ref string) *Schema {
	return &Schema{Ref: ref}
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

// Set object
func (s *Schema) SetRef(ref string) *Schema {
	s.Ref = ref
	return s
}

// Set array
func (s *Schema) SetItems(items *Items) *Schema {
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

	Ref   string
	Items *Items // `type` == array
}

func NewItems(itemType string) *Items {
	return &Items{Type: itemType, Format: defaultFormat(itemType)}
}

func NewItemsRef(ref string) *Items {
	return &Items{Ref: ref}
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

// Set object
func (i *Items) SetRef(ref string) *Items {
	i.Ref = ref
	return i
}

// Set array
func (i *Items) SetItems(items *Items) *Items {
	i.Items = items
	return i
}
