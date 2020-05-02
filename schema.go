package yamldoc

type Schema struct {
	Title       string
	Type        string
	Required    bool
	Description string

	Format          string
	AllowEmptyValue bool
	Default         interface{}
	Enum            []interface{}

	Ref   string // `type` == object
	Items *Items // `type` == array
}

func NewSchema(schemaType string, required bool) *Schema {
	return &Schema{
		Required: required, Type: schemaType, Format: defaultFormat(schemaType),
	}
}

func NewSchemaRef(ref string) *Schema {
	return &Schema{Ref: ref}
}

func (s *Schema) SetTitle(title string) *Schema {
	s.Title = title
	return s
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
	s.Items = items
	return s
}
