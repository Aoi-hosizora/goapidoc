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

// Obj<integer#int64, string[]>[]
type apiType struct {
	name string
	kind apiTypeKind

	prime  *apiPrime
	object *apiObject
	array  *apiArray
}

// integer#int64
type apiPrime struct {
	typ    string
	format string
}

// xxx<yyy,...>
type apiObject struct {
	generics []*apiType // yyy,...
	typ      string     // xxx
}

// xxx[][]
type apiArray struct {
	item *apiType // xxx[]
}

func parseApiType(typ string) *apiType {
	typ = strings.TrimSpace(typ)
	l := len(typ)

	// array: X[][]
	if l >= 3 && typ[l-2:] == "[]" {
		item := parseApiType(typ[:l-2]) // X[]
		return &apiType{
			name: typ,
			kind: apiArrayKind,
			array: &apiArray{
				item: item,
			},
		}
	}

	// base: X#Y
	for _, prime := range []string{INTEGER, NUMBER, STRING, BOOLEAN, FILE, ARRAY, OBJECT} {
		primeLen := len(prime)
		if typ == prime || (l > primeLen+1 && strings.HasPrefix(typ, prime) && typ[primeLen] == '#') {
			format := defaultFormat(prime)
			if typ != prime {
				format = typ[primeLen+1:] // Y
			}
			return &apiType{
				name: typ,
				kind: apiPrimeKind,
				prime: &apiPrime{
					typ:    prime,
					format: format,
				},
			}
		}
	}

	// object with generic: X<Y,Z<A,B<C>>>
	start := strings.Index(typ, "<")
	if l > 3 && start >= 0 && typ[l-1] == '>' && typ[l-2:l] != "<>" {
		paramStr := typ[start+1 : l-1]         // Y,Z<A,B<C>>
		params := strings.Split(paramStr, ",") // Y | Z<A | B<C>>
		generics := make([]string, 0)
		for idx, param := range params {
			if strings.Count(param, "<") != strings.Count(param, ">") { // Z<A
				if len(params) > idx+1 {
					params[idx+1] = params[idx] + "," + params[idx+1] // Z<A,B<C>>
				} else {
					panic("Failed to parse type of: `" + paramStr + "`")
				}
			} else {
				generics = append(generics, param)
			}
		}
		out := &apiObject{
			typ:      typ[:start],
			generics: []*apiType{},
		}
		for _, generic := range generics {
			out.generics = append(out.generics, parseApiType(generic))
		}
		return &apiType{
			name:   typ,
			kind:   apiObjectKind,
			object: out,
		}
	}

	// object without generic: X | X<>
	object := &apiObject{
		typ:      typ,
		generics: []*apiType{},
	}
	if l >= 2 && typ[l-2:l] == "<>" {
		object.typ = typ[:l-2]
	}
	return &apiType{
		name:   typ,
		kind:   apiObjectKind,
		object: object,
	}
}

// get default format from type
func defaultFormat(typ string) string {
	if typ == INTEGER {
		return INT32
	} else if typ == NUMBER {
		return DOUBLE
	}
	return ""
}

// parse generic param
func prehandleGenericName(def *Definition) {
	// change generic type in properties
	for _, prop := range def.properties {
		for _, gen := range def.generics { // T -> «T»
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
	for idx := range def.generics {
		gen := def.generics[idx]
		if !strings.HasPrefix(gen, "«") || !strings.HasSuffix(gen, "»") {
			def.generics[idx] = "«" + gen + "»"
		}
	}
}

// generate related generic list
func prehandleGenericList(definitions []*Definition, allTypes []string) []*Definition {
	genericDefs := make(map[string]*Definition)
	normalDefs := make(map[string]*Definition)
	for _, def := range definitions {
		if len(def.generics) == 0 {
			normalDefs[def.name] = def
		} else {
			genericDefs[def.name] = def
		}
	}

	addedDefs := newLinkedHashMap() // make(map[string]*Definition, 0) // new definitions to add

	// preHandle
	var preHandle func(string)
	preHandle = func(typ string) { // string
		at := parseApiType(typ) // *apiType
		for at.kind == apiArrayKind {
			at = at.array.item
		}
		if at.kind != apiObjectKind || len(at.object.generics) == 0 {
			return
		}

		genDef, ok := genericDefs[at.object.typ] // *Definition
		if !ok {
			return
		}

		// new motoDef
		properties := make([]*Property, 0)
		for _, prop := range genDef.properties {
			properties = append(properties, &Property{
				name:       prop.name,
				typ:        prop.typ, // << need parse
				required:   prop.required,
				desc:       prop.desc,
				allowEmpty: prop.allowEmpty,
				def:        prop.def,
				example:    prop.example,
				enum:       prop.enum,
				minLength:  prop.minLength,
				maxLength:  prop.maxLength,
				minimum:    prop.minimum,
				maximum:    prop.maximum,
			})
		}
		addedDef := &Definition{ // *Definition
			name:       genDef.name, // << need parse
			desc:       genDef.desc,
			generics:   []string{},
			properties: properties,
		}

		specNames := make([]string, 0)
		for idx, genName := range genDef.generics {
			if len(at.object.generics) < idx {
				break
			}
			specName := at.object.generics[idx].name
			specNames = append(specNames, specName)
			for _, prop := range addedDef.properties {
				prop.typ = strings.ReplaceAll(prop.typ, genName, specName)
			}
		}
		addedDef.name += "<" + strings.Join(specNames, ", ") + ">"
		addedDefs.Set(addedDef.name, addedDef)
		// log.Println("added", addedDef.name)

		for _, prop := range addedDef.properties {
			if !addedDefs.Has(prop.typ) {
				preHandle(prop.typ) // << preHandle properties types
			}
		}
	}

	for _, typ := range allTypes {
		preHandle(typ)
	}

	out := make([]*Definition, 0)
	for _, def := range normalDefs {
		out = append(out, def)
	}
	for _, key := range addedDefs.Keys() {
		val, _ := addedDefs.Get(key)
		out = append(out, val.(*Definition))
	}

	return out
}
