package goapidoc

import (
	"fmt"
	"log"
	"testing"
)

func fail(t *testing.T, msg string) {
	fmt.Println(msg)
	t.Fail()
}

func testPanic(t *testing.T, want bool, fn func()) {
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
		fail(t, fmt.Sprintf("Test case want no panic but panic with `%s`.", msg))
	} else if !didPanic && want {
		fail(t, "Test case want panic but no panic happened.")
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
			})
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
	} {
		t.Run(tc.give, func(t *testing.T) {
			testPanic(t, tc.wantPanic, func() {
				if at := parseApiType(tc.give); !tc.checkFn(at) {
					fail(t, "Test case failed")
				}
			})
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
					if len(prehandled.generics) != len(tc.wantGenerics) || len(prehandled.properties) != len(tc.wantPropTypes) {
						fail(t, "Test case failed")
					}
					for idx := range tc.wantGenerics {
						if prehandled.generics[idx] != tc.wantGenerics[idx] {
							fail(t, "Test case failed")
						}
					}
					for idx := range tc.wantPropTypes {
						if prehandled.properties[idx].typ != tc.wantPropTypes[idx] {
							log.Println(prehandled.properties[idx].typ, tc.wantPropTypes[idx])
							fail(t, "Test case failed")
						}
					}
				}
			})
		})
	}
}

func TestPrehandleDefinitionList(t *testing.T) {
	definitions := []*Definition{
		{name: "User", properties: []*Property{}},
		{name: "Login", properties: []*Property{}},
		{name: "String", properties: []*Property{}},
		{name: "Result", generics: []string{"T"}, properties: []*Property{{name: "code", typ: "number"}, {name: "data", typ: "T"}}},
		{name: "Page", generics: []string{"T"}, properties: []*Property{{name: "code", typ: "number"}, {name: "data", typ: "T[]"}}},
		{name: "Result2", generics: []string{"T", "U"}, properties: []*Property{{name: "a", typ: "T"}, {name: "b", typ: "U[]"}}},
		{name: "Result3", generics: []string{"T", "U", "V"}, properties: []*Property{{name: "a", typ: "T"}, {name: "b", typ: "U[][]"}, {name: "c", typ: "Result<V>"}}},
	}
	newDefinitions := make([]*Definition, 0, len(definitions))
	for _, definition := range definitions {
		newDefinitions = append(newDefinitions, prehandleDefinition(definition))
	}
	newDefs := prehandleDefinitionList(newDefinitions, []string{
		"Result<Page<User>>",
		"Result3<User, Page<Result2<Login, Page<Login>>>, String[]>",
		"Integer",
		"Result2<String, Result2<String, String>>",
	})

	if len(newDefs) != 12 {
		t.Fatal()
	}

	contain := func(definition *Definition) bool {
		ok := false
		for _, newDef := range newDefs {
			if newDef.name != definition.name || len(newDef.properties) != len(definition.properties) {
				continue
			}
			if len(newDef.properties) == 0 {
				ok = true
				break
			}

			ok2 := true
			for idx, newProp := range newDef.properties {
				prop := definition.properties[idx]
				if newProp.name != prop.name || newProp.typ != prop.typ {
					ok2 = false
					break
				}
			}
			if ok2 {
				ok = true
				break
			}
		}
		return ok
	}

	// 0: User | Login | String
	// 1: Page<User> | Page<Login> | Result<String[]> | Result2<String, String>
	// 2: Result<Page<User>> | Result2<Login, Page<Login>> | Result2<String, Result2<String, String>>
	// 3: Page<Result2<Login, Page<Login>>>
	// 4: Result3<User, Page<Result2<Login, Page<Login>>>, String[]>

	for idx, ok := range []bool{
		contain(&Definition{name: "User", properties: []*Property{}}),
		contain(&Definition{name: "Login", properties: []*Property{}}),
		contain(&Definition{name: "String", properties: []*Property{}}),

		contain(&Definition{name: "Page<User>", properties: []*Property{{name: "code", typ: "number"}, {name: "data", typ: "User[]"}}}),
		contain(&Definition{name: "Page<Login>", properties: []*Property{{name: "code", typ: "number"}, {name: "data", typ: "Login[]"}}}),
		contain(&Definition{name: "Result<String[]>", properties: []*Property{{name: "code", typ: "number"}, {name: "data", typ: "String[]"}}}),
		contain(&Definition{name: "Result2<String, String>", properties: []*Property{{name: "a", typ: "String"}, {name: "b", typ: "String[]"}}}),

		contain(&Definition{name: "Result<Page<User>>", properties: []*Property{{name: "code", typ: "number"}, {name: "data", typ: "Page<User>"}}}),
		contain(&Definition{name: "Result2<Login, Page<Login>>", properties: []*Property{{name: "a", typ: "Login"}, {name: "b", typ: "Page<Login>[]"}}}),
		contain(&Definition{name: "Result2<String, Result2<String, String>>", properties: []*Property{{name: "a", typ: "String"}, {name: "b", typ: "Result2<String, String>[]"}}}),

		contain(&Definition{name: "Page<Result2<Login, Page<Login>>>", properties: []*Property{{name: "code", typ: "number"}, {name: "data", typ: "Result2<Login, Page<Login>>[]"}}}),

		contain(&Definition{name: "Result3<User, Page<Result2<Login, Page<Login>>>, String[]>", properties: []*Property{{name: "a", typ: "User"}, {name: "b", typ: "Page<Result2<Login, Page<Login>>>[][]"}, {name: "c", typ: "Result<String[]>"}}}),
	} {
		if !ok {
			t.Fatal(idx)
		}
	}
}
