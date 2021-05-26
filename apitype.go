package goapidoc

import (
	"regexp"
	"strings"
)

type apiTypeKind int

const (
	apiPrimeKind apiTypeKind = iota
	apiObjectKind
	apiArrayKind
)

// Example: Obj<integer#int64, string[]>[]
type apiType struct {
	name string
	kind apiTypeKind

	prime  *apiPrime
	object *apiObject
	array  *apiArray
}

// Example: integer#int64
type apiPrime struct {
	typ    string
	format string
}

// Example: xxx<yyy,...>
type apiObject struct {
	generics []*apiType // yyy,...
	typ      string     // xxx
}

// Example: xxx[][]
type apiArray struct {
	item *apiType // xxx[]
}

// parseApiType parses type string to apiType.
// TODO WAITING FOR REFACTORING, NEED TO BE ROBUST, WITHOUT TEST COMPLETELY
func parseApiType(typ string) *apiType {
	typ = strings.TrimSpace(typ)
	l := len(typ)
	if !regexp.MustCompile(`[a-zA-Z0-9\-_#, <>«»\[\]]+`).MatchString(typ) {
		panic("Invalid type `" + typ + "`")
	}

	// array: X[] | X[][]
	if l >= 3 && typ[l-2:] == "[]" {
		item := parseApiType(typ[:l-2]) // X | X[]
		return &apiType{
			name:  typ,
			kind:  apiArrayKind,
			array: &apiArray{item: item},
		}
	}

	// prime with and without format: X | X# | X#Y
	for _, prime := range []string{INTEGER, NUMBER, STRING, BOOLEAN, FILE, ARRAY, OBJECT} {
		primeLen := len(prime)
		if typ == prime || (strings.HasPrefix(typ, prime) && typ[primeLen] == '#') {
			format := ""
			if l >= primeLen+2 {
				format = strings.TrimSpace(typ[primeLen+1:]) // Y
			}
			if format == "" {
				format = defaultFormat(prime)
			}
			return &apiType{
				name:  typ,
				kind:  apiPrimeKind,
				prime: &apiPrime{typ: prime, format: format},
			}
		}
	}

	// object with generic: X<Y> | X<Y,Z<A>> | X<Y,Z<A,B<C>>>
	start := strings.Index(typ, "<")
	if l > 3 && start > 0 && typ[l-1] == '>' && typ[l-2:] != "<>" {
		paramStr := typ[start+1 : l-1]         // Y | Y,Z<A> | Y,Z<A,B<C>>
		params := strings.Split(paramStr, ",") // Y || Y|Z<A> || Y|Z<A|B<C>>
		paramsLen := len(params)
		generics := make([]string, 0, 1)
		for idx, param := range params {
			if strings.Count(param, "<") == strings.Count(param, ">") { // Y | Z<A>
				generics = append(generics, param)
			} else { // Z<A | B<C>>
				if idx+1 < paramsLen {
					params[idx+1] = params[idx] + "," + params[idx+1] // Z<A,B<C>>
				} else {
					panic("Invalid type `" + typ + "`")
				}
			}
		}
		objectName := typ[:start]
		genericsType := make([]*apiType, 0, 1)
		for _, generic := range generics {
			genericsType = append(genericsType, parseApiType(generic))
		}
		return &apiType{
			name:   typ,
			kind:   apiObjectKind,
			object: &apiObject{typ: objectName, generics: genericsType},
		}
	}

	// object without generic: X | X<>
	objectName := typ
	if l >= 2 && objectName[l-2:] == "<>" {
		objectName = typ[:l-2]
	}
	return &apiType{
		name:   typ,
		kind:   apiObjectKind,
		object: &apiObject{typ: objectName, generics: []*apiType{}},
	}
}

// defaultFormat returns the default format for given type.
func defaultFormat(typ string) string {
	if typ == INTEGER {
		return INT32
	}
	if typ == NUMBER {
		return DOUBLE
	}
	return ""
}

// prehandleGenericName parses and returns the prehandled generic name and properties types for Definition.
// TODO WAITING FOR REFACTORING, NEED TO BE CONST FUNCTION
func prehandleGenericName(definition *Definition) *Definition {
	// clone definition
	out := &Definition{
		name:       definition.name,
		desc:       definition.desc,
		generics:   append(make([]string, 0, len(definition.generics)), definition.generics...),
		properties: make([]*Property, 0, len(definition.properties)),
	}
	for _, prop := range definition.properties {
		out.properties = append(out.properties, &Property{
			name:       prop.name,
			typ:        prop.typ,
			required:   prop.required,
			desc:       prop.desc,
			allowEmpty: prop.allowEmpty,
			defaul:     prop.defaul,
			example:    prop.example,
			enums:      prop.enums,
			minLength:  prop.minLength,
			maxLength:  prop.maxLength,
			minimum:    prop.minimum,
			maximum:    prop.maximum,
		})
	}

	// update generic type in properties
	for _, prop := range out.properties {
		for _, gen := range out.generics { // T -> «T»
			if strings.HasPrefix(gen, "«") && strings.HasSuffix(gen, "»") {
				continue
			}

			newGen := "«" + gen + "»"
			re, err := regexp.Compile(`(^|[, <])` + gen + `([, <>\[]|$)`) // {, <} {, <>[]}
			if err != nil {
				panic("Failed to parse type of: `" + gen + "`")
			}

			prop.typ = strings.ReplaceAll(prop.typ, " ", "")
			prop.typ = re.ReplaceAllString(prop.typ, "$1"+newGen+"$2")
			prop.typ = strings.ReplaceAll(prop.typ, ",", ", ")
		}
	}

	// change generic type in generic list
	for idx := range out.generics {
		gen := out.generics[idx]
		if !strings.HasPrefix(gen, "«") || !strings.HasSuffix(gen, "»") {
			out.generics[idx] = "«" + gen + "»"
		}
	}
	return out
}

// prehandleGenericList parses and prehandles related generic list for Definition.
// TODO WAITING FOR REFACTORING, NEED TO BE CONST FUNCTION
func prehandleGenericList(definitions []*Definition, allTypes []string) []*Definition {
	genericDefs := make(map[string]*Definition)
	normalDefs := newOrderedMap(0) // old normal definitions
	addedDefs := newOrderedMap(0)  // new definitions to add
	for _, definition := range definitions {
		if len(definition.generics) == 0 {
			normalDefs.Set(definition.name, definition)
		} else {
			genericDefs[definition.name] = definition
		}
	}

	// preHandle
	var preHandle func(string)
	preHandle = func(typ string) {
		at := parseApiType(typ) // *apiType
		for at.kind == apiArrayKind {
			at = at.array.item
		}
		if at.kind != apiObjectKind || len(at.object.generics) == 0 {
			return
		}
		genDef, ok := genericDefs[at.object.typ]
		if !ok {
			return
		}

		properties := make([]*Property, 0, len(genDef.properties))
		for _, prop := range genDef.properties {
			properties = append(properties, &Property{
				name:       prop.name,
				typ:        prop.typ, // << need to be parsed
				required:   prop.required,
				desc:       prop.desc,
				allowEmpty: prop.allowEmpty,
				defaul:     prop.defaul,
				example:    prop.example,
				enums:      prop.enums,
				minLength:  prop.minLength,
				maxLength:  prop.maxLength,
				minimum:    prop.minimum,
				maximum:    prop.maximum,
			})
		}
		addedDef := &Definition{
			name:       genDef.name, // << need to be parsed
			desc:       genDef.desc,
			generics:   []string{},
			properties: properties,
		}

		specNames := make([]string, 0, len(genDef.generics))
		for idx, genName := range genDef.generics {
			if len(at.object.generics) < idx {
				break
			}
			specName := at.object.generics[idx].name
			specNames = append(specNames, specName)
			for _, prop := range addedDef.properties {
				prop.typ = strings.ReplaceAll(prop.typ, genName, specName) // replace directly
			}
		}
		addedDef.name += "<" + strings.Join(specNames, ", ") + ">"
		addedDefs.Set(addedDef.name, addedDef)

		for _, prop := range addedDef.properties {
			if !addedDefs.Has(prop.typ) {
				preHandle(prop.typ) // << preHandle properties types
			}
		}
	}

	for _, typ := range allTypes {
		preHandle(typ)
	}

	out := make([]*Definition, 0, len(normalDefs.Keys())+len(addedDefs.Keys()))
	for _, key := range normalDefs.Keys() {
		val, _ := normalDefs.Get(key)
		out = append(out, val.(*Definition))
	}
	for _, key := range addedDefs.Keys() {
		val, _ := addedDefs.Get(key)
		out = append(out, val.(*Definition))
	}

	return out
}
