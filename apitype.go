package goapidoc

import (
	"regexp"
	"strings"
)

type apiTypeKind int8

const (
	apiPrimeKind apiTypeKind = iota + 1
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

// Example: string or string#date-time
type apiPrime struct {
	typ    string
	format string
}

// Example: xxx or xxx<yyy, zzz>
type apiObject struct {
	typ      string
	generics []*apiType
}

// Example: xxx[][]
type apiArray struct {
	item *apiType
}

var (
	typeNameRe    = regexp.MustCompile(`^[a-zA-Z0-9_]+(?:(?:<(.+)>)|(?:#[a-zA-Z0-9\-_]*))?(?:\[])*$`) // xxx(?:<(yyy)>|#zzz)?(?:[])*
	genericNameRe = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

// checkTypeName checks given type name from Param, Response and Definition.
func checkTypeName(typ string) {
	// re
	if !typeNameRe.MatchString(typ) {
		panic("Invalid type `" + typ + "`")
	}
	// <(.+)>
	for _, subTyp := range typeNameRe.FindStringSubmatch(typ)[1:] {
		genParts := strings.Split(subTyp, ",")
		genStrings := make([]string, 0, 2) // cap defaults to 2
		for idx, part := range genParts {
			if strings.Count(part, "<") == strings.Count(part, ">") {
				// T1, T2<T3>
				genStrings = append(genStrings, part)
			} else {
				// T1, T2<T3, T4> => T1 | T2<T3 | T4> => T1 | T2<T3,T4>
				if idx+1 < len(genParts) {
					genParts[idx+1] = genParts[idx] + "," + genParts[idx+1]
				} else {
					panic("Invalid type `" + typ + "`")
				}
			}
		}
		for _, genTyp := range genStrings {
			if genTyp = strings.TrimSpace(genTyp); genTyp != "" {
				checkTypeName(genTyp)
			}
		}
	}
}

// parseApiType parses type string to three kinds of apiType.
func parseApiType(typ string) *apiType {
	typ = strings.TrimSpace(typ)

	// 1. array: X[] | X[][][]
	if strings.HasSuffix(typ, "[]") {
		item := parseApiType(typ[:len(typ)-2]) // X | X[][]
		return &apiType{
			name:  typ,
			kind:  apiArrayKind,
			array: &apiArray{item: item},
		}
	}

	// 2. prime with format: X# | X#Y
	if strings.Contains(typ, "#") {
		for _, prime := range []string{INTEGER, NUMBER, STRING, BOOLEAN, FILE, ARRAY, OBJECT} {
			if strings.HasPrefix(typ, prime) {
				if prime == ARRAY || prime == OBJECT {
					panic("Use array or object in type `" + typ + "` invalidly")
				}
				fmtIdx := strings.Index(typ, "#") + 1
				fmt := ""
				if fmtIdx < len(typ) {
					fmt = typ[fmtIdx:]
				}
				return &apiType{
					name:  typ,
					kind:  apiPrimeKind,
					prime: &apiPrime{typ: prime, format: fmt},
				}
			}
		}
	}

	// 3. object with generic: X<Y<Z>> | X<Y, Z<A, B<C>>>
	if strings.Contains(typ, "<") {
		genIdx := strings.Index(typ, "<") + 1
		genParts := strings.Split(typ[genIdx:len(typ)-1], ",") // Y<Z> || {Y, Z<A, B<C>>}
		genStrings := make([]string, 0, 2)                     // cap defaults to 2
		for idx, part := range genParts {
			part = strings.TrimSpace(part)
			if strings.Count(part, "<") == strings.Count(part, ">") { // Y<Z>
				genStrings = append(genStrings, part)
			} else { // Z<A, B<C>>
				if idx+1 < len(genParts) {
					genParts[idx+1] = genParts[idx] + "," + genParts[idx+1] // -> Z<A + , + B<C>>
				} else {
					panic("Invalid type `" + typ + "`")
				}
			}
		}
		genTypes := make([]*apiType, 0, len(genStrings))
		for _, gen := range genStrings {
			genTypes = append(genTypes, parseApiType(gen))
		}
		object := typ[:genIdx-1]
		switch object {
		case INTEGER, NUMBER, STRING, BOOLEAN, FILE, ARRAY, OBJECT:
			panic("Invalid type `" + typ + "`")
		}
		return &apiType{
			name:   typ,
			kind:   apiObjectKind,
			object: &apiObject{typ: object, generics: genTypes},
		}
	}

	// 4. prime without format: X
	switch typ {
	case INTEGER, NUMBER, STRING, BOOLEAN, FILE:
		return &apiType{
			name:  typ,
			kind:  apiPrimeKind,
			prime: &apiPrime{typ: typ, format: defaultFormat(typ)},
		}
	case ARRAY, OBJECT:
		panic("Use array or object as type invalidly")
	}

	// 5. object without generic: X
	if strings.Contains(typ, "#") {
		panic("Invalid type `" + typ + "`")
	}
	return &apiType{
		name:   typ,
		kind:   apiObjectKind,
		object: &apiObject{typ: typ, generics: []*apiType{}},
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

// collectAllSpecTypes checks and collects all specific types.
func collectAllSpecTypes(doc *Document) []string {
	// check all type names (param, resp, prop)
	cnt := 0
	for _, op := range doc.operations {
		cnt += len(op.params) + len(op.responses)
		for _, param := range op.params {
			checkTypeName(param.typ)
		}
		for _, resp := range op.responses {
			if resp.typ != "" {
				checkTypeName(resp.typ)
			}
		}
	}
	for _, def := range doc.definitions {
		for _, prop := range def.properties {
			checkTypeName(prop.typ)
		}
		if len(def.generics) == 0 {
			cnt += len(def.properties)
		}
	}

	// collect all specific types
	out := make([]string, 0, cnt)
	for _, op := range doc.operations {
		for _, param := range op.params {
			out = append(out, param.typ)
		}
		for _, resp := range op.responses {
			if resp.typ != "" {
				out = append(out, resp.typ)
			}
		}
	}
	for _, def := range doc.definitions {
		if len(def.generics) == 0 {
			for _, prop := range def.properties {
				out = append(out, prop.typ)
			}
		}
	}
	return out
}

// prehandleDefinition deduplicates, checks and prehandles generic names, and returns a new cloned Definition.
func prehandleDefinition(definition *Definition) *Definition {
	// deduplicate and check generic names
	generics := make([]string, 0, len(definition.generics))
	for _, gen := range definition.generics {
		contained := false
		for _, g := range generics {
			if g == gen {
				contained = true
				break
			}
		}
		if !contained {
			if !genericNameRe.MatchString(gen) { // a-zA-Z0-9_
				panic("Invalid generic type `" + gen + "`")
			}
			generics = append(generics, gen)
		}
	}

	// clone definition
	out := &Definition{
		name:       definition.name,
		desc:       definition.desc,
		generics:   generics,
		properties: make([]*Property, 0, len(definition.properties)),
	}
	for _, prop := range definition.properties {
		out.properties = append(out.properties, cloneProperty(prop))
	}

	// update generic type name
	for idx, gen := range out.generics {
		newGen := "«" + gen + "»"
		re := regexp.MustCompile(`(^|[,\s<])` + gen + `([,\s<>\[]|$)`) // {,\s<} xxx {,\s<>[}
		for _, prop := range out.properties {
			for { // replace all property type
				curr := prop.typ
				prop.typ = strings.ReplaceAll(prop.typ, " ", "")
				prop.typ = re.ReplaceAllString(prop.typ, "$1"+newGen+"$2") // T -> «T»
				prop.typ = strings.ReplaceAll(prop.typ, ",", ", ")
				if prop.typ == curr {
					break
				}
			}
		}
		out.generics[idx] = newGen
	}
	return out
}

// prehandleDefinitionList prehandles and returns the final Definition list with given and type list.
func prehandleDefinitionList(allDefinitions []*Definition, allTypes []string) []*Definition {
	// extract generic definitions from given definitions
	allDefMap := make(map[string]*Definition, len(allDefinitions))
	out := make([]*Definition, 0, len(allDefinitions))    // out definition slice
	outKeys := make(map[string]bool, len(allDefinitions)) // out definition key map
	for _, def := range allDefinitions {
		if _, ok := allDefMap[def.name]; ok {
			panic("Duplicate definition `" + def.name + "`")
		}
		allDefMap[def.name] = def
		if len(def.generics) == 0 {
			out = append(out, def) // definition without generic
			outKeys[def.name] = true
		}
	}

	// extract more definitions from given types
	var extractFn func(typ string)
	extractFn = func(typ string) {
		if _, ok := outKeys[typ]; ok {
			return
		}
		at := parseApiType(typ)
		for at.kind == apiArrayKind {
			at = at.array.item
		}
		if at.kind != apiObjectKind {
			return
		}
		// `at` belongs to object
		obj := at.object

		// check object existence and generic parameter
		genDef, ok := allDefMap[obj.typ]
		if !ok {
			panic("Object type `" + at.name + "` not found")
		}
		if len(obj.generics) != len(genDef.generics) {
			panic("Object type `" + at.name + "`'s generic parameter length is not matched")
		}
		if len(obj.generics) == 0 {
			return
		}

		// specific definition need to be added
		specDef := &Definition{
			name:       genDef.name, // TypeName<GenericName, ...>
			desc:       genDef.desc,
			generics:   nil, // empty
			properties: make([]*Property, 0, len(genDef.properties)),
		}
		for _, prop := range genDef.properties {
			specDef.properties = append(specDef.properties, cloneProperty(prop)) // << need to extract type recurrently
		}
		// replace to spec name for new definition
		specNames := make([]string, 0, len(genDef.generics))
		for idx, genName := range genDef.generics {
			specName := obj.generics[idx].name
			specNames = append(specNames, specName)
			for _, prop := range specDef.properties {
				prop.typ = strings.ReplaceAll(prop.typ, genName, specName) // «T» -> XXX, replace directly
			}
		}
		specDef.name += "<" + strings.Join(specNames, ", ") + ">" // TypeName -> TypeName<GenericName, ...>

		// extract recurrently and append to outMap
		for _, prop := range specDef.properties {
			extractFn(prop.typ) // << extract property type recurrently
		}
		out = append(out, specDef)
		outKeys[specDef.name] = true
	}

	// for all types, extract generic parameters to definition list
	for _, typ := range allTypes {
		extractFn(typ)
	}

	// return definition slice
	return out
}
