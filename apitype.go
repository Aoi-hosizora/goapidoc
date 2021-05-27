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

// prehandleDefinition prehandles and returns a cloned Definition with new generic type name.
func prehandleDefinition(definition *Definition) *Definition {
	// check and clone definition
	out := &Definition{
		name:       definition.name,
		desc:       definition.desc,
		generics:   make([]string, 0, len(definition.generics)),
		properties: make([]*Property, 0, len(definition.properties)),
	}
	re1 := regexp.MustCompile(`[a-zA-Z0-9_]`)
	re2 := regexp.MustCompile(`[a-zA-Z0-9\-_#, <>\[\]]+`)
	for _, gen := range definition.generics {
		if !re1.MatchString(gen) {
			panic("Invalid generic type `" + gen + "`")
		}
		out.generics = append(out.generics, gen)
	}
	for _, prop := range definition.properties {
		if !re2.MatchString(prop.typ) {
			panic("Invalid type `" + prop.typ + "`")
		}
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

	// update generic type name
	for idx, gen := range out.generics {
		newGen := "«" + gen + "»"
		re, err := regexp.Compile(`(^|[, <])` + gen + `([, <>\[]|$)`) // {, <} {, <>[}
		if err != nil {
			panic("Invalid generic type `" + gen + "`")
		}
		for _, prop := range out.properties {
			prop.typ = strings.ReplaceAll(prop.typ, " ", "")
			prop.typ = re.ReplaceAllString(prop.typ, "$1"+newGen+"$2") // T -> «T»
			prop.typ = strings.ReplaceAll(prop.typ, ",", ", ")
		}
		out.generics[idx] = newGen
	}
	return out
}

// prehandleDefinitionList prehandles and returns final Definition list with given definition list and type list.
func prehandleDefinitionList(allDefinitions []*Definition, allTypes []string) []*Definition {
	// split definitions to generic-defs and normal-defs
	genericDefs := make(map[string]*Definition)
	normalDefs := newOrderedMap(0) // old normal definitions
	addedDefs := newOrderedMap(0)  // new definitions to add
	for _, definition := range allDefinitions {
		if len(definition.generics) != 0 {
			genericDefs[definition.name] = definition
		} else {
			normalDefs.Set(definition.name, definition)
		}
	}

	// extract from types
	var extractFn func(typ string)
	extractFn = func(typ string) {
		at := parseApiType(typ)
		for at.kind == apiArrayKind {
			at = at.array.item
		}
		if at.kind != apiObjectKind || len(at.object.generics) == 0 {
			return
		}

		// object with generic
		genDef, ok := genericDefs[at.object.typ]
		if !ok {
			return
		}
		// new definition need to be added
		newDef := &Definition{
			name:       genDef.name, // TypeName<GenericName, ...>
			desc:       genDef.desc,
			generics:   make([]string, 0),
			properties: make([]*Property, 0, len(genDef.properties)),
		}
		for _, prop := range genDef.properties {
			newDef.properties = append(newDef.properties, &Property{
				name:       prop.name,
				typ:        prop.typ, // << need to extract recurrently
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

		// spec name for new definition
		specNames := make([]string, 0, len(genDef.generics))
		for idx, genName := range genDef.generics {
			if idx >= len(at.object.generics) {
				break
			}
			specName := at.object.generics[idx].name
			specNames = append(specNames, specName)
			for _, prop := range newDef.properties {
				prop.typ = strings.ReplaceAll(prop.typ, genName, specName) // «T» -> XXX, replace directly
			}
		}
		newDef.name += "<" + strings.Join(specNames, ", ") + ">"

		// append to addedDefs
		addedDefs.Set(newDef.name, newDef)
		for _, prop := range newDef.properties {
			if !addedDefs.Has(prop.typ) {
				extractFn(prop.typ) // << extract from properties types
			}
		}
	}
	for _, typ := range allTypes {
		extractFn(typ)
	}

	// combine result definition list
	out := make([]*Definition, 0, len(normalDefs.Keys())+len(addedDefs.Keys()))
	for _, key := range normalDefs.Keys() {
		val := normalDefs.MustGet(key).(*Definition)
		out = append(out, val)
	}
	for _, key := range addedDefs.Keys() {
		val := addedDefs.MustGet(key).(*Definition)
		out = append(out, val)
	}
	return out
}
