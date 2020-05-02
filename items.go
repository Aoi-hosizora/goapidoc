package yamldoc

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

func (i *Items) SetItems(items *Items) *Items {
	i.Items = items
	return i
}
