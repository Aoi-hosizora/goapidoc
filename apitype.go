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

	outPrime  *apiPrime  // prime
	outObject *apiObject // object
	outArray  *apiArray  // array
}

// integer#int64
type apiPrime struct {
	typ    string
	format string
}

// xxx<yyy,...>
type apiObject struct {
	generic []*apiType // yyy,...
	typ     string     // xxx
}

// xxx[][]
type apiArray struct {
	typ *apiType // xxx[]
}

func parseApiType(typ string) *apiType {
	typ = strings.TrimSpace(typ)

	// array
	if len(typ) >= 3 && typ[len(typ)-2:] == "[]" {
		return &apiType{
			name: typ, kind: apiArrayKind,
			outArray: &apiArray{typ: parseApiType(typ[:len(typ)-2])},
		}
	}

	// base
	for _, tp := range []string{INTEGER, NUMBER, STRING, BOOLEAN, FILE, ARRAY, OBJECT} {
		if typ == tp || (len(typ) > len(tp)+1 && strings.HasPrefix(typ, tp) && typ[len(tp)] == '#') {
			format := defaultFormat(tp)
			if typ != tp {
				format = typ[len(tp)+1:]
			}
			out := &apiPrime{typ: tp, format: format}
			return &apiType{name: typ, kind: apiPrimeKind, outPrime: out}
		}
	}

	// object with generic
	start := strings.Index(typ, "<")
	if len(typ) >= 3 && start != -1 && typ[len(typ)-1] == '>' {
		param := typ[start+1 : len(typ)-1]
		temp := strings.Split(param, ",")
		generics := make([]string, 0)
		for idx := 0; idx < len(temp); idx++ {
			if strings.Contains(temp[idx], "<") && !strings.Contains(temp[idx], ">") {
				if idx+1 < len(temp) {
					temp[idx+1] = temp[idx] + "," + temp[idx+1]
				} else {
					panic("Failed to parse type of " + param)
				}
			} else {
				generics = append(generics, temp[idx])
			}
		}
		out := &apiObject{typ: typ[:start]}
		for _, generic := range generics {
			out.generic = append(out.generic, parseApiType(generic))
		}
		return &apiType{name: typ, kind: apiObjectKind, outObject: out}
	}

	// object without generic
	return &apiType{name: typ, kind: apiObjectKind, outObject: &apiObject{typ: typ}}
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
func preHandleGeneric(def *Definition) {
	for _, prop := range def.properties {
		for _, gen := range def.generics { // T -> «T»
			newGen := "«" + gen + "»"
			re, err := regexp.Compile("(^|[, <])" + gen + "([, >\\[]|$)") // (^|[, <])x([, >]|$)
			if err != nil {
				panic("failed to compile generic parameter: " + err.Error())
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
