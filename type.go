package goapidoc

import (
	"regexp"
	"strings"
)

type innerTypeKind int

const (
	innerPrimeKind innerTypeKind = iota
	innerObjectKind
	innerArrayKind
)

type innerType struct {
	Name string
	Kind innerTypeKind

	OutPrime  *innerPrime  // prime
	OutObject *innerObject // object
	OutArray  *innerArray  // array
}

// Obj<integer#int64, string[]>[]
func parseInnerType(t string) *innerType {
	t = strings.TrimSpace(t)

	// array
	if len(t) >= 3 && t[len(t)-2:] == "[]" {
		return &innerType{
			Name: t, Kind: innerArrayKind,
			OutArray: &innerArray{Type: parseInnerType(t[:len(t)-2])},
		}
	}

	// base
	for _, tp := range []string{INTEGER, NUMBER, STRING, BOOLEAN, FILE, ARRAY, OBJECT} {
		if t == tp || (len(t) > len(tp)+1 && strings.HasPrefix(t, tp) && t[len(tp)] == '#') {
			format := defaultFormat(tp)
			if t != tp {
				format = t[len(tp)+1:]
			}
			out := &innerPrime{Type: tp, Format: format}
			return &innerType{Name: t, Kind: innerPrimeKind, OutPrime: out}
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
		out := &innerObject{Type: t[:start]}
		for _, generic := range generics {
			out.Generic = append(out.Generic, parseInnerType(generic))
		}
		return &innerType{Name: t, Kind: innerObjectKind, OutObject: out}
	}

	// object without generic
	return &innerType{Name: t, Kind: innerObjectKind, OutObject: &innerObject{Type: t}}
}

// integer#int64
type innerPrime struct {
	Type   string
	Format string
}

// xxx<yyy,...>
type innerObject struct {
	Generic []*innerType // yyy,...
	Type    string       // xxx
}

// xxx[][]
type innerArray struct {
	Type *innerType // xxx[]
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
				continue
			}
			prop.Type = re.ReplaceAllString(prop.Type, "$1"+newGen+"$2")
		}
	}
	for idx := range def.Generics {
		def.Generics[idx] = "«" + def.Generics[idx] + "»"
	}
}
