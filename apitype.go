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
	if typ[l-2:l] == "<>" {
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
func preHandleGenerics(def *Definition) {
	for _, prop := range def.properties {
		for _, gen := range def.generics { // T -> «T»
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

	for idx := range def.generics {
		def.generics[idx] = "«" + def.generics[idx] + "»"
	}
}
