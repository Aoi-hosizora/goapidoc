package apidoc

// The Schema Object allows the definition of input and output data types.
// These types can be objects, but also primitives and arrays.
//
// Primitive example:
//     NewSchema("integer", true).SetDefault(0)
// Object example: where user is defined in #/Definitions
//     RefSchema("User")
// Array example:
//     ArrSchema(NewItems("integer")           // -> array of integer
//     ArrSchema(RefItems("User")              // -> array of object
//     ArrSchema(ArrItems(NewItems("integer")) // -> array of array
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
func NewSchema(t string, req bool) *Schema {
	return &Schema{Required: req, Type: t, Format: defaultFormat(t)}
}

// Create a schema that is a reference type, can have options
// $ref, options must be (string, *Schema|*Items) pairs
func RefSchema(ref string, options ...interface{}) *Schema {
	return &Schema{Ref: ref, Options: handleWithOptions(options)}
}

// Create a schema that is a array, could not have options, please use it in `Ref`
// items
func ArrSchema(items *Items) *Schema {
	return &Schema{Items: items, Type: ARRAY}
}

func (s *Schema) SetDescription(desc string) *Schema {
	s.Description = desc
	return s
}

func (s *Schema) SetFormat(format string) *Schema {
	s.Format = format
	return s
}

func (s *Schema) SetAllowEmptyValue(allow bool) *Schema {
	s.AllowEmptyValue = allow
	return s
}

func (s *Schema) SetDefault(def interface{}) *Schema {
	s.Default = def
	return s
}

func (s *Schema) SetEnum(enum ...interface{}) *Schema {
	s.Enum = enum
	return s
}

// A limited subset of JSON-Schema's items object.
// It is used by parameter definitions that are not located in "body" -> should use SetSchema()
// example:
//     NewItems("integer")           // -> array of integer
//     RefItems("User")              // -> array of object
//     ArrItems(NewItems("integer")) // -> array of array
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
func NewItems(t string) *Items {
	return &Items{Type: t, Format: defaultFormat(t)}
}

// Create a items that is a reference type, can have options
// $ref, options must be (string, *Schema|*Items) pairs
func RefItems(ref string, options ...interface{}) *Items {
	return &Items{Ref: ref, Options: handleWithOptions(options)}
}

// Create a items that is an array type, could not have options, please use it in `Ref`
// items, options must be (string, *Schema|*Items) pairs
func ArrItems(items *Items) *Items {
	return &Items{Items: items, Type: ARRAY}
}

func (i *Items) SetFormat(format string) *Items {
	i.Format = format
	return i
}

func (i *Items) SetDefault(def interface{}) *Items {
	i.Default = def
	return i
}

func (i *Items) SetEnum(enum ...interface{}) *Items {
	i.Enum = enum
	return i
}

// Addition schema option
type AdditionOption struct {
	Field  string
	Schema *Schema
	Items  *Items
}

func handleWithOptions(options ...interface{}) []*AdditionOption {
	if len(options)&1 == 1 {
		options = options[:len(options)-1]
	}
	out := make([]*AdditionOption, 0)
	idx := 0
	for idx < len(options) {
		field, ok := options[idx].(string)
		if !ok {
			return out
		}
		if schema, ok := options[idx+1].(*Schema); ok {
			out = append(out, &AdditionOption{Field: field, Schema: schema})
		} else if items, ok := options[idx+1].(*Items); ok {
			out = append(out, &AdditionOption{Field: field, Items: items})
		} else {
			return out
		}
		idx += 2
	}
	return out
}
