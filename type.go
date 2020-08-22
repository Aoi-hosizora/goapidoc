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

func parseApiType(t string) *apiType {
	t = strings.TrimSpace(t)

	// array
	if len(t) >= 3 && t[len(t)-2:] == "[]" {
		return &apiType{
			name: t, kind: apiArrayKind,
			outArray: &apiArray{typ: parseApiType(t[:len(t)-2])},
		}
	}

	// base
	for _, tp := range []string{INTEGER, NUMBER, STRING, BOOLEAN, FILE, ARRAY, OBJECT} {
		if t == tp || (len(t) > len(tp)+1 && strings.HasPrefix(t, tp) && t[len(tp)] == '#') {
			format := defaultFormat(tp)
			if t != tp {
				format = t[len(tp)+1:]
			}
			out := &apiPrime{typ: tp, format: format}
			return &apiType{name: t, kind: apiPrimeKind, outPrime: out}
		}
	}

	// object with generic
	start := strings.Index(t, "<")
	if len(t) >= 3 && start != -1 && t[len(t)-1] == '>' {
		param := t[start+1 : len(t)-1]
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
		out := &apiObject{typ: t[:start]}
		for _, generic := range generics {
			out.generic = append(out.generic, parseApiType(generic))
		}
		return &apiType{name: t, kind: apiObjectKind, outObject: out}
	}

	// object without generic
	return &apiType{name: t, kind: apiObjectKind, outObject: &apiObject{typ: t}}
}

// get default format from type
func defaultFormat(t string) string {
	if t == INTEGER {
		return INT32
	} else if t == NUMBER {
		return DOUBLE
	}
	return ""
}

// parse generic param before mapDefinition()
func preHandleGeneric(def *Definition) {
	for _, prop := range def.Properties {
		for _, gen := range def.Generics { // T -> «T»
			newGen := "«" + gen + "»"
			re, err := regexp.Compile("(^|[, <])" + gen + "([, >\\[]|$)") // (^|[, <])x([, >]|$)
			if err != nil {
				panic("failed to compile generic parameter: " + err.Error())
			}
			prop.Type = strings.ReplaceAll(prop.Type, " ", "")
			prop.Type = re.ReplaceAllString(prop.Type, "$1"+newGen+"$2")
			prop.Type = strings.ReplaceAll(prop.Type, ",", ", ")
		}
	}
	for idx := range def.Generics {
		def.Generics[idx] = "«" + def.Generics[idx] + "»"
	}
}
