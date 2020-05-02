package yamldoc

type Items struct {
	Ref  string
	Type string

	Format  string
	Default interface{}
	Enum    []interface{}
	Items   *Items // `type` == array
}

func NewRefItems(ref string) *Items {
	return &Items{Ref: ref}
}

func NewItems(itemType string) *Items {
	return &Items{Type: itemType}
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
	i.Items = items
	return i
}
