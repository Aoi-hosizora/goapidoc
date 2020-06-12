package apidoc

import (
	"strings"
)

// Addition schema option
// Deprecated
type AdditionOption struct {
	Field  string
	Ref    string
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
		if ref, ok := options[idx+1].(string); ok {
			out = append(out, &AdditionOption{Field: field, Ref: ref})
		} else if schema, ok := options[idx+1].(*Schema); ok {
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

func mapRefOptions(doc *innerDocument, ref string, options []*AdditionOption) (newRef string) {
	oldDef, ok := doc.Definitions[ref]
	if !ok {
		return ref
	}

	def := &innerDefinition{
		Type:        oldDef.Type,
		Required:    oldDef.Required,
		Description: oldDef.Description,
		Properties:  make(map[string]*innerSchema),
	}
	for k, v := range oldDef.Properties {
		def.Properties[k] = v
	}
	types := make([]string, len(options))

	for i, o := range options {
		if oldProp, ok := def.Properties[o.Field]; ok {
			if o.Schema != nil { // object schema (ref)
				def.Properties[o.Field], types[i] = mapObjectOption(doc, o.Schema)
			} else if o.Items != nil { // array schema (items)
				def.Properties[o.Field], types[i] = mapArrayOption(doc, o.Items, oldProp)
			} else if o.Ref != "" {
				if ok {
					if oldProp.Type != ARRAY { // array
						def.Properties[o.Field], types[i] = mapObjectOption(doc, RefSchema(o.Ref))
					} else { // object
						def.Properties[o.Field], types[i] = mapArrayOption(doc, RefItems(o.Ref), oldProp)
					}
				}
			}
		} else {
			types[i] = "NULL"
		}
	}

	newRef = ref + "<" + strings.Join(types, ",") + ">"
	// for {
	// 	if _, ok := doc.Definitions[newRef]; ok {
	// 		newRef = newRef + "'"
	// 	} else {
	// 		break
	// 	}
	// }
	doc.Definitions[newRef] = def
	return newRef
}

func mapObjectOption(doc *innerDocument, s *Schema) (*innerSchema, string) {
	schema := mapSchema(doc, s)
	if schema.OriginRef != "" {
		return schema, schema.OriginRef
	} else {
		return schema, schema.Type
	}
}

func mapArrayOption(doc *innerDocument, i *Items, prop *innerSchema) (*innerSchema, string) {
	items := mapItems(doc, i)
	newSchema := &innerSchema{}
	*newSchema = *prop
	newSchema.Type = ARRAY
	newSchema.Items = items
	if items.OriginRef != "" {
		return newSchema, items.OriginRef
	} else {
		return newSchema, items.Type
	}
}
