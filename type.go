package goapidoc

import (
	"strings"
)

type innerTypeEnum int

const (
	innerBaseType innerTypeEnum = iota
	innerObjectType
	innerArrayType
)

type innerType struct {
	Name string
	Type innerTypeEnum

	OutType   string       // prime
	OutSchema *innerObject // object
	OutItems  *innerArray  // array
}

func parseInnerType(t string) *innerType {
	t = strings.TrimSpace(t)

	// base
	if t == INTEGER || t == NUMBER || t == STRING || t == BOOLEAN || t == FILE {
		return &innerType{Name: t, Type: innerBaseType, OutType: t}
	}

	// array
	if len(t) >= 3 && t[len(t)-2:] == "[]" {
		return &innerType{
			Name: t, Type: innerArrayType,
			OutItems: &innerArray{Type: parseInnerType(t[:len(t)-2])},
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
		return &innerType{Name: t, Type: innerObjectType, OutSchema: out}
	}

	// object without generic
	return &innerType{Name: t, Type: innerObjectType, OutSchema: &innerObject{Type: t}}
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
