package goapidoc

import (
	"fmt"
	"testing"
)

func failNow(t *testing.T, msg string) {
	fmt.Println(msg)
	t.FailNow()
}

func testPanic(t *testing.T, want bool, fn func(), message string) {
	didPanic := false
	var msg interface{}
	func() {
		defer func() {
			if msg = recover(); msg != nil {
				didPanic = true
			}
		}()
		fn()
	}()

	if didPanic && !want {
		failNow(t, fmt.Sprintf("Test case for %s want no panic but panic with `%s`", message, msg))
	} else if !didPanic && want {
		failNow(t, fmt.Sprintf("Test case for %s want panic but no panic happened", message))
	}
}

func testMatchElements(t *testing.T, s1, s2 []string, s1Msg, s2Msg string) {
	if len(s1) != len(s2) {
		failNow(t, fmt.Sprintf("Two slice's lengths (%s and %s) is not same", s1Msg, s2Msg))
	}
	for _, i1 := range s1 {
		contained := false
		for _, i2 := range s2 {
			if i1 == i2 {
				contained = true
				break
			}
		}
		if !contained {
			failNow(t, fmt.Sprintf("There are some items in %s that are not found in %s", s1Msg, s2Msg))
		}
	}
	for _, i2 := range s2 {
		contained := false
		for _, i1 := range s1 {
			if i1 == i2 {
				contained = true
				break
			}
		}
		if !contained {
			failNow(t, fmt.Sprintf("There are some items in %s that are not found in %s", s2Msg, s1Msg))
		}
	}
}

func TestCheckTypeName(t *testing.T) {
	for _, tc := range []struct {
		give      string
		wantPanic bool
	}{
		{"", true},
		{"$", true},
		{"#", true},
		{"integer", false},
		{"integer#", false},
		{"integer#int64", false},
		{"integer[]", false},
		{"integer#[]", false},
		{"integer#int64[]", false},
		{"integer#[][][]", false},
		{"integer[]#", true},
		{"integer##", true},
		{"integer#a#b", true},
		{"integer[][", true},
		{"integer[0]", true},
		{"int$$eger#int64", true},
		{"int%eger[]", true},
		{"int#eger[]", false},
		{"int#eger$int64[]", true},
		{"Object", false},
		{"Object#xxx", false},
		{"Object<>", true},
		{"Object<T>", false},
		{"Object<T", true},
		{"Object<T<", true},
		{"Object<T<<X>>", true},
		{"Object<T1,T2>", false},
		{"Object<T1, T2>", false},
		{"Object<T1, T2<T3>>", false},
		{"Object<T>#fmt", true},
		{"Object<T#fmt>", false},
		{"Object<T1<T2<T3<T4>>>>", false},
		{"Object<T1<T2<T3<T4[]>[]>[]>[]>[]", false},
		{"Object<T1, T2<T3>>[]", false},
		{"Object<T1, T2<T3, T4>>", false},
		{"Object<T1, T2<T3, T4<>", true},
		{"Object<T1, T2<TT1<integer#int64[]>>, T3<TT2, TT3<TT4, TT5<number#>>[]>[], string#date-time>[][]", false},
	} {
		t.Run(tc.give, func(t *testing.T) {
			testPanic(t, tc.wantPanic, func() {
				checkTypeName(tc.give)
			}, "checkTypeName")
		})
	}
}

func TestParseApiType(t *testing.T) {
	for _, tc := range []struct {
		give      string
		wantPanic bool
		checkFn   func(*apiType) bool
	}{
		{"integer", false, func(at *apiType) bool {
			return at.prime.typ == "integer" && at.prime.format == "int32"
		}},
		{"number", false, func(at *apiType) bool {
			return at.prime.typ == "number" && at.prime.format == "double"
		}},
		{"integer#", false, func(at *apiType) bool {
			return at.prime.typ == "integer" && at.prime.format == ""
		}},
		{"string#date-time", false, func(at *apiType) bool {
			return at.prime.typ == "string" && at.prime.format == "date-time"
		}},
		{"boolean[]", false, func(at *apiType) bool {
			return at.array.item.prime.typ == "boolean"
		}},
		{"boolean[][][]", false, func(at *apiType) bool {
			return at.array.item.name == "boolean[][]" && at.array.item.array.item.array.item.prime.typ == "boolean"
		}},
		{"string#date[]", false, func(at *apiType) bool {
			return at.array.item.prime.typ == "string" && at.array.item.prime.format == "date"
		}},
		{"Object", false, func(at *apiType) bool {
			return at.object.typ == "Object" && len(at.object.generics) == 0
		}},
		{"Object[][]", false, func(at *apiType) bool {
			return at.array.item.array.item.object.typ == "Object" && len(at.array.item.array.item.object.generics) == 0
		}},
		{"Object<integer>", false, func(at *apiType) bool {
			return at.object.typ == "Object" && at.object.generics[0].prime.typ == "integer" && at.object.generics[0].prime.format == "int32"
		}},
		{"Object<integer#int64>", false, func(at *apiType) bool {
			return at.object.typ == "Object" && at.object.generics[0].prime.typ == "integer" && at.object.generics[0].prime.format == "int64"
		}},
		{"Object<integer[]>", false, func(at *apiType) bool {
			return at.object.typ == "Object" && at.object.generics[0].array.item.prime.typ == "integer" && at.object.generics[0].array.item.prime.format == "int32"
		}},
		{"Object<T>", false, func(at *apiType) bool {
			return at.object.typ == "Object" && at.object.generics[0].object.typ == "T" && len(at.object.generics[0].object.generics) == 0
		}},
		{"Object<boolean, T>", false, func(at *apiType) bool {
			return at.object.typ == "Object" && at.object.generics[0].prime.typ == "boolean" &&
				at.object.generics[1].object.typ == "T" && len(at.object.generics[1].object.generics) == 0
		}},
		{"Object<boolean, T<U1, U2<U3[]>>>", false, func(at *apiType) bool {
			return at.object.typ == "Object" && at.object.generics[0].prime.typ == "boolean" && at.object.generics[1].object.typ == "T" &&
				at.object.generics[1].object.generics[0].object.typ == "U1" && len(at.object.generics[1].object.generics[0].object.generics) == 0 &&
				at.object.generics[1].object.generics[1].object.typ == "U2" && at.object.generics[1].object.generics[1].object.generics[0].array.item.object.typ == "U3"
		}},
		{"Object<T1, T2<TT1<integer#int64[]>>, T3<TT2, TT3<TT4, TT5<number#>>[]>[], string#date-time>[][]", false, func(at *apiType) bool {
			at = at.array.item.array.item
			return at.object.typ == "Object" &&
				at.object.generics[0].object.typ == "T1" && len(at.object.generics[0].object.generics) == 0 &&
				at.object.generics[1].object.typ == "T2" && at.object.generics[1].object.generics[0].object.typ == "TT1" &&
				at.object.generics[1].object.generics[0].object.generics[0].array.item.prime.typ == "integer" &&
				at.object.generics[1].object.generics[0].object.generics[0].array.item.prime.format == "int64" &&
				at.object.generics[2].array.item.object.typ == "T3" &&
				at.object.generics[2].array.item.object.generics[0].object.typ == "TT2" &&
				at.object.generics[2].array.item.object.generics[1].array.item.object.typ == "TT3" &&
				at.object.generics[2].array.item.object.generics[1].array.item.object.generics[0].object.typ == "TT4" &&
				at.object.generics[2].array.item.object.generics[1].array.item.object.generics[1].object.typ == "TT5" &&
				at.object.generics[2].array.item.object.generics[1].array.item.object.generics[1].object.generics[0].prime.typ == "number" &&
				at.object.generics[2].array.item.object.generics[1].array.item.object.generics[1].object.generics[0].prime.format == "" &&
				at.object.generics[3].prime.typ == "string" && at.object.generics[3].prime.format == "date-time"
		}},

		{"integer<Object>", true, nil},
		{"Object#xxx", true, nil},
		{"Object<xxx,xxx<xxx>", true, nil},
		{"array", true, nil},
		{"array[]", true, nil},
		{"object#", true, nil},
		{"Object<object>", true, nil},
	} {
		t.Run(tc.give, func(t *testing.T) {
			testPanic(t, tc.wantPanic, func() {
				if at := parseApiType(tc.give); !tc.checkFn(at) {
					failNow(t, "Parse and get a wrong ApiType")
				}
			}, "parseApiType")
		})
	}
}

func TestPrehandleDefinition(t *testing.T) {
	for _, tc := range []struct {
		name          string
		giveGenerics  []string
		givePropTypes []string
		wantPanic     bool
		wantGenerics  []string
		wantPropTypes []string
	}{
		{"empty", []string{}, []string{}, false, []string{}, []string{}},

		{"check1", []string{"T-T"}, []string{}, true, []string{}, []string{}},
		{"check2", []string{"T#"}, []string{}, true, []string{}, []string{}},
		{"check3", []string{"T<T>"}, []string{}, true, []string{}, []string{}},
		{"check4", []string{"T_T"}, []string{}, false, []string{"«T_T»"}, []string{}},
		{"check5", []string{"T123"}, []string{}, false, []string{"«T123»"}, []string{}},

		{"parse1", []string{"T"}, []string{"U"},
			false, []string{"«T»"}, []string{"U"}},
		{"parse2", []string{"T"}, []string{"T", "T[][]", "Obj<T>", "T<Obj>"},
			false, []string{"«T»"}, []string{"«T»", "«T»[][]", "Obj<«T»>", "«T»<Obj>"}},
		{"parse3", []string{"T", "T"}, []string{"Obj<U, T, T[][], T[]>[]", "inTeger#T"},
			false, []string{"«T»"}, []string{"Obj<U, «T», «T»[][], «T»[]>[]", "inTeger#T"}},
		{"parse4", []string{"T", "U"}, []string{"Obj<T, U[], TT>[]"},
			false, []string{"«T»", "«U»"}, []string{"Obj<«T», «U»[], TT>[]"}},
		{"parse5", []string{"T", "U", "V"}, []string{"inT[]", "ObjT<inT[], TV[], U<V>>"},
			false, []string{"«T»", "«U»", "«V»"}, []string{"inT[]", "ObjT<inT[], TV[], «U»<«V»>>"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testPanic(t, tc.wantPanic, func() {
				def := &Definition{generics: tc.giveGenerics, properties: make([]*Property, 0, len(tc.givePropTypes))}
				for _, typ := range tc.givePropTypes {
					def.properties = append(def.properties, &Property{typ: typ})
				}
				prehandled := prehandleDefinition(def)

				if !tc.wantPanic {
					testMatchElements(t, prehandled.generics, tc.wantGenerics, "prehandledGenerics", "wantGenerics")
					prehandledProperties := make([]string, 0, len(prehandled.properties))
					for _, prop := range prehandled.properties {
						prehandledProperties = append(prehandledProperties, prop.typ)
					}
					testMatchElements(t, prehandledProperties, tc.wantPropTypes, "prehandledPropTypes", "wantPropTypes")
				}
			}, "prehandleDefinition")
		})
	}
}

func TestPrehandleDefinitionList(t *testing.T) {
	definitions := []*Definition{
		{name: "Result", generics: []string{"T"}, properties: []*Property{{name: "code", typ: "integer"}, {name: "data", typ: "T"}}},
		{name: "Result2", generics: []string{"T", "U"}, properties: []*Property{{name: "code", typ: "integer"}, {name: "data", typ: "T"}, {name: "error", typ: "U"}}},
		{name: "Page", generics: []string{"T"}, properties: []*Property{{name: "total", typ: "integer"}, {name: "data", typ: "T[]"}}},
		{name: "Page2", generics: []string{"T"}, properties: []*Property{{name: "next_max_id", typ: "integer"}, {name: "total", typ: "integer"}, {name: "data", typ: "T[]"}}},

		{name: "UserDto", properties: []*Property{{name: "uid", typ: "integer"}, {name: "name", typ: "string"}}},
		{name: "ErrorDto", properties: []*Property{{name: "type", typ: "string"}, {name: "detail", typ: "string"}}},
	}
	prehandledDefinitions := make([]*Definition, 0, len(definitions))
	for _, definition := range definitions {
		prehandledDefinitions = append(prehandledDefinitions, prehandleDefinition(definition))
	}

	t.Run("dup definition", func(t *testing.T) {
		testPanic(t, true, func() {
			prehandleDefinitionList([]*Definition{{name: "UserDto"}, {name: "UserDto"}}, []string{})
		}, "prehandleDefinitionList")
		testPanic(t, true, func() {
			prehandleDefinitionList([]*Definition{{name: "Result", generics: []string{"T"}}, {name: "Result"}}, []string{})
		}, "prehandleDefinitionList")
	})
	for _, tc := range []struct {
		name            string
		giveTypes       []string
		wantPanic       bool
		wantObjectNames []string
		wantPropNames   [][]string
		wantPropTypes   [][]string
	}{
		{"empty", []string{}, false, []string{"UserDto", "ErrorDto"},
			[][]string{{"uid", "name"}, {"type", "detail"}},
			[][]string{{"integer", "string"}, {"string", "string"}}},
		{"UserDto, ErrorDto", []string{"UserDto", "ErrorDto"}, false, []string{"UserDto", "ErrorDto"},
			[][]string{{"uid", "name"}, {"type", "detail"}},
			[][]string{{"integer", "string"}, {"string", "string"}}},
		{"integer", []string{"integer"}, false, []string{"UserDto", "ErrorDto"},
			[][]string{{"uid", "name"}, {"type", "detail"}},
			[][]string{{"integer", "string"}, {"string", "string"}}},
		{"Result<integer>", []string{"Result<integer>"}, false, []string{"UserDto", "ErrorDto", "Result<integer>"},
			[][]string{{"uid", "name"}, {"type", "detail"}, {"code", "data"}},
			[][]string{{"integer", "string"}, {"string", "string"}, {"integer", "integer"}}},
		{"Result<UserDto>", []string{"Result<UserDto>"}, false, []string{"UserDto", "ErrorDto", "Result<UserDto>"},
			[][]string{{"uid", "name"}, {"type", "detail"}, {"code", "data"}},
			[][]string{{"integer", "string"}, {"string", "string"}, {"integer", "UserDto"}}},
		{"Result2<UserDto, ErrorDto>", []string{"Result2<UserDto, ErrorDto>"}, false, []string{"UserDto", "ErrorDto", "Result2<UserDto, ErrorDto>"},
			[][]string{{"uid", "name"}, {"type", "detail"}, {"code", "data", "error"}},
			[][]string{{"integer", "string"}, {"string", "string"}, {"integer", "UserDto", "ErrorDto"}}},
		{"Result<Page<UserDto>>", []string{"Result<Page<UserDto>>"}, false, []string{"UserDto", "ErrorDto", "Page<UserDto>", "Result<Page<UserDto>>"},
			[][]string{{"uid", "name"}, {"type", "detail"}, {"total", "data"}, {"code", "data"}},
			[][]string{{"integer", "string"}, {"string", "string"}, {"integer", "UserDto[]"}, {"integer", "Page<UserDto>"}}},
		{"UserDto[][] | Page<UserDto[]>", []string{"UserDto[][]", "Result<UserDto[]>"}, false, []string{"UserDto", "ErrorDto", "Result<UserDto[]>"},
			[][]string{{"uid", "name"}, {"type", "detail"}, {"code", "data"}},
			[][]string{{"integer", "string"}, {"string", "string"}, {"integer", "UserDto[]"}}},
		{"Page<UserDto> | Page<UserDto> | Result<Page<UserDto>>", []string{"Page<UserDto>", "Page<UserDto>", "Result<Page<UserDto>>"}, false,
			[]string{"UserDto", "ErrorDto", "Page<UserDto>", "Result<Page<UserDto>>"},
			[][]string{{"uid", "name"}, {"type", "detail"}, {"total", "data"}, {"code", "data"}},
			[][]string{{"integer", "string"}, {"string", "string"}, {"integer", "UserDto[]"}, {"integer", "Page<UserDto>"}}},
		{"Result<UserDto> | Result<Result<UserDto>> | Result<Result<Result<UserDto>>>", []string{"Result<UserDto>", "Result<Result<UserDto>>", "Result<Result<Result<UserDto>>>"}, false,
			[]string{"UserDto", "ErrorDto", "Result<UserDto>", "Result<Result<UserDto>>", "Result<Result<Result<UserDto>>>"},
			[][]string{{"uid", "name"}, {"type", "detail"}, {"code", "data"}, {"code", "data"}, {"code", "data"}},
			[][]string{{"integer", "string"}, {"string", "string"}, {"integer", "UserDto"}, {"integer", "Result<UserDto>"}, {"integer", "Result<Result<UserDto>>"}}},
		{"Page<UserDto> | Result<Result<UserDto>>[] | Result<Result<Result<UserDto>>>",
			[]string{"Page<UserDto>", "Result<UserDto>", "Result<Result<UserDto>>[]", "Result<Result<Result<UserDto>>>"}, false,
			[]string{"UserDto", "ErrorDto", "Page<UserDto>", "Result<UserDto>", "Result<Result<UserDto>>", "Result<Result<Result<UserDto>>>"},
			[][]string{{"uid", "name"}, {"type", "detail"}, {"total", "data"}, {"code", "data"}, {"code", "data"}, {"code", "data"}},
			[][]string{{"integer", "string"}, {"string", "string"}, {"integer", "UserDto[]"}, {"integer", "UserDto"}, {"integer", "Result<UserDto>"}, {"integer", "Result<Result<UserDto>>"}}},
		{"Result2<Result<UserDto>, ErrorDto> | Result<Result<UserDto[]>> | Page<Result<Result<UserDto>>>",
			[]string{"Result2<Result<UserDto>, ErrorDto>", "Result<Result<UserDto[]>>", "Page<Result<Result<UserDto>>>"}, false,
			[]string{"UserDto", "ErrorDto", "Result<UserDto>", "Result2<Result<UserDto>, ErrorDto>", "Result<UserDto[]>", "Result<Result<UserDto[]>>", "Result<Result<UserDto>>", "Page<Result<Result<UserDto>>>"},
			[][]string{{"uid", "name"}, {"type", "detail"}, {"code", "data"}, {"code", "data", "error"}, {"code", "data"}, {"code", "data"}, {"code", "data"}, {"total", "data"}},
			[][]string{{"integer", "string"}, {"string", "string"}, {"integer", "UserDto"}, {"integer", "Result<UserDto>", "ErrorDto"}, {"integer", "UserDto[]"}, {"integer", "Result<UserDto[]>"}, {"integer", "Result<UserDto>"}, {"integer", "Result<Result<UserDto>>[]"}}},

		{"blank type", []string{""}, true, nil,nil,nil},
		{"not found", []string{"xxx"}, true, nil,nil,nil},
		{"not found sub type", []string{"Result<xxx>"}, true, nil,nil,nil},
		{"integer<integer>", []string{"integer<integer>"}, true, nil,nil,nil},
		{"UserDto<integer>", []string{"UserDto<integer>"}, true, nil,nil,nil},
		{"Result<UserDto, ErrorDto>", []string{"Result<UserDto, ErrorDto>"}, true, nil,nil,nil},
		{"Result2<UserDto>", []string{"Result2<UserDto>"}, true, nil,nil,nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testPanic(t, tc.wantPanic, func() {
				newDefinitions := prehandleDefinitionList(prehandledDefinitions, tc.giveTypes)
				newObjectNames := make([]string, 0, len(newDefinitions))
				for _, def := range newDefinitions {
					newObjectNames = append(newObjectNames, def.name)
				}
				testMatchElements(t, newObjectNames, tc.wantObjectNames, "newObjectNames", "wantObjectNames")
				for idx, def := range newDefinitions {
					newPropNames := make([]string, 0, len(def.properties))
					newPropTypes := make([]string, 0, len(def.properties))
					for _, prop := range def.properties {
						newPropNames = append(newPropNames, prop.name)
						newPropTypes = append(newPropTypes, prop.typ)
					}
					testMatchElements(t, newPropNames, tc.wantPropNames[idx], "newPropNames", "wantPropNames")
					testMatchElements(t, newPropTypes, tc.wantPropTypes[idx], "newPropTypes", "wantPropTypes")
				}
			}, "prehandleDefinitionList")
		})
	}
}
