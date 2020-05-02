package yamldoc

type Schema struct {
	Ref string

	Title       string
	Type        string
	Required    bool
	Description string

	Format          string
	AllowEmptyValue bool
	Default         interface{}
	Enum            []interface{}
	Items           *Items // `type` == array
}

func NewProperty(name string, schemaType string, required bool, description string) *Schema {
	return &Schema{Title: name, Type: schemaType, Required: required, Description: description}
}

func NewRefProperty(name string, ref string, required bool, description string) *Schema {
	return &Schema{Title: name, Type: "object", Ref: ref, Required: required, Description: description}
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
	s.Items = items
	return s
}
