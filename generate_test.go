package goapidoc

import (
	"log"
	"testing"
)

func TestParseApiType(t *testing.T) {
	str := "aType<T1<>, T2<TT1<integer#int32[]>>, T3<TT1, TT2<TT3, TT4<TT5>>>, T4#FF>[][]"
	res := parseApiType(str)

	if res.name != str {
		t.Fatal("res.name")
	}
	if res.array.item.name != "aType<T1<>, T2<TT1<integer#int32[]>>, T3<TT1, TT2<TT3, TT4<TT5>>>, T4#FF>[]" {
		t.Fatal("res.array.item.name")
	}
	if res.array.item.array.item.name != "aType<T1<>, T2<TT1<integer#int32[]>>, T3<TT1, TT2<TT3, TT4<TT5>>>, T4#FF>" {
		t.Fatal("res.array.item.array.item.name")
	}

	obj := res.array.item.array.item
	if len(obj.object.generics) != 4 {
		t.Fatal("len(obj.object.generics)")
	}

	gen := obj.object.generics
	if gen[0].name != "T1<>" {
		t.Fatal("gen[0].name")
	}
	if gen[1].name != "T2<TT1<integer#int32[]>>" {
		t.Fatal("gen[1].name")
	}
	if gen[2].name != "T3<TT1, TT2<TT3, TT4<TT5>>>" {
		t.Fatal("gen[2].name")
	}
	if gen[3].name != "T4#FF" {
		t.Fatal()
	}

	g0 := gen[0] // T1<>
	if g0.object.typ != "T1" {
		t.Fatal("g0.object.typ")
	}
	if len(g0.object.generics) != 0 {
		t.Fatal("len(g0.object.generics)")
	}

	g1 := gen[1] // T2<TT1<integer[]>>
	if g1.object.typ != "T2" {
		t.Fatal("g1.object.typ")
	}
	g10 := g1.object.generics[0]
	if g10.object.typ != "TT1" {
		t.Fatal("g10.object.typ")
	}
	g100 := g10.object.generics[0]
	if g100.name != "integer#int32[]" {
		t.Fatal("g100.name")
	}
	if g100.array.item.prime.typ != "integer" {
		t.Fatal("g100.array.item.prime.typ")
	}
	if g100.array.item.prime.format != "int32" {
		t.Fatal("g100.array.item.prime.format")
	}

	g2 := gen[2] // T3<TT1, TT2<TT3, TT4<TT5>>>
	g20 := g2.object.generics[0]
	g21 := g2.object.generics[1]
	if g20.object.typ != "TT1" {
		t.Fatal("g20.object.typ")
	}
	if g21.name != "TT2<TT3, TT4<TT5>>" {
		t.Fatal("g21.name")
	}
	if g21.object.typ != "TT2" {
		t.Fatal("g21.object.typ")
	}
	g210 := g21.object.generics[0]
	g211 := g21.object.generics[1] // TT4<TT5>
	if g210.object.typ != "TT3" {
		t.Fatal("g210.object.typ")
	}
	if g211.name != "TT4<TT5>" {
		t.Fatal("g211.name")
	}
	if g211.object.typ != "TT4" {
		t.Fatal("g211.object.typ")
	}
	if g211.object.generics[0].object.typ != "TT5" {
		t.Fatal("g211.object.generics[0].object.typ")
	}

	g3 := gen[3] // T4#FF
	if g3.object.typ != "T4#FF" {
		t.Fatal("g3.object.typ")
	}
	if len(g3.object.generics) != 0 {
		t.Fatal("len(g3.object.generics)")
	}
}

func TestPrehandleGenericName(t *testing.T) {
	definition := &Definition{
		generics: []string{"T", "U", "V"},
		properties: []*Property{
			{typ: "inT[]"},
			{typ: "O<inT[], T[], T, inT<int>>"},
			{typ: "T"},
			{typ: "tT<T<tT>[][], T>[]"},
			{typ: "TtT<T,tT[],T[][]>[]"},
		},
	}
	prehandleGenericName(definition)

	if definition.generics[0] != "«T»" {
		t.Fatal("definition.generics[0]")
	}
	if definition.generics[1] != "«U»" {
		t.Fatal("definition.generics[1]")
	}
	if definition.generics[2] != "«V»" {
		t.Fatal("definition.generics[2]")
	}

	p0 := definition.properties[0]
	p1 := definition.properties[1]
	p2 := definition.properties[2]
	p3 := definition.properties[3]
	p4 := definition.properties[4]
	if p0.typ != "inT[]" {
		t.Fatal("p0.typ")
	}
	if p1.typ != "O<inT[], «T»[], «T», inT<int>>" {
		t.Fatal("p1.typ")
	}
	if p2.typ != "«T»" {
		t.Fatal("p2.typ")
	}
	if p3.typ != "tT<«T»<tT>[][], «T»>[]" {
		t.Fatal("p3.typ")
	}
	if p4.typ != "TtT<«T», tT[], «T»[][]>[]" {
		t.Fatal("p4.typ")
	}
}

func TestPrehandleGenericList(t *testing.T) {
	definitions := []*Definition{
		{name: "User", properties: []*Property{}},
		{name: "Login", properties: []*Property{}},
		{name: "String", properties: []*Property{}},
		{name: "Result", generics: []string{"T"}, properties: []*Property{{name: "code", typ: "number"}, {name: "data", typ: "T"}}},
		{name: "Page", generics: []string{"T"}, properties: []*Property{{name: "code", typ: "number"}, {name: "data", typ: "T[]"}}},
		{name: "Result2", generics: []string{"T", "U"}, properties: []*Property{{name: "a", typ: "T"}, {name: "b", typ: "U[]"}}},
		{name: "Result3", generics: []string{"T", "U", "V"}, properties: []*Property{{name: "a", typ: "T"}, {name: "b", typ: "U[][]"}, {name: "c", typ: "Result<V>"}}},
	}
	for _, definition := range definitions {
		prehandleGenericName(definition)
	}
	newDefs := prehandleGenericList(definitions, []string{
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

func TestGenerate(t *testing.T) {
	SetDocument(
		"localhost:65530", "/",
		NewInfo("goapidoc", "goapidoc test api", "1.0").
			TermsOfService("https://github.com/Aoi-hosizora").
			License(NewLicense("MIT", "https://github.com/Aoi-hosizora")).
			Contact(NewContact("Aoi-hosizora", "https://github.com/Aoi-hosizora", "aoihosizora@hotmail.com")),
	)
	SetTags(
		NewTag("Authorization", "auth-controller"),
		NewTag("User", "user-controller"),
		NewTag("Test", "test-controller"),
	)
	SetSecurities(
		NewSecurity("Jwt", "header", "Authorization"),
	)

	AddRoutePaths(
		NewRoutePath("POST", "/auth/register", "Register").
			Desc("Register.").
			Tags("Authorization").
			Params(NewBodyParam("param", "RegisterParam", true, "register param")).
			Responses(NewResponse(200, "Result")),

		NewRoutePath("POST", "/auth/login", "Login").
			Desc("Login.").
			Tags("Authorization").
			Params(NewBodyParam("param", "LoginParam", true, "login param")).
			Responses(
				NewResponse(200, "_Result<LoginDto>"),
				NewResponse(400, "Result").Examples(map[string]string{"application/json": "{\n  \"code\": 400, \n  \"message\": \"Unauthorized\"\n}"}),
			),

		NewRoutePath("DELETE", "/auth/logout", "Logout").
			Tags("Authorization").
			Securities("Jwt").
			Responses(NewResponse(200, "Result")),

		NewRoutePath("GET", "/user", "Get users").
			Tags("User").
			Securities("Jwt").
			Params(
				NewQueryParam("page", "integer#int32", false, "current page").Default(1).Example(1).Minimum(1),
				NewQueryParam("limit", "integer#int32", false, "page size").Default(20).Example(20).Minimum(15),
			).
			Responses(NewResponse(200, "_Result<_Page<UserDto>>")),

		NewRoutePath("GET", "/user/{username}", "Get a user").
			Tags("User").
			Securities("Jwt").
			Params(NewPathParam("username", "string", true, "username")).
			Responses(NewResponse(200, "_Result<UserDto>")),

		NewRoutePath("PUT", "/user/deprecated", "Update user").
			Tags("User").
			Securities("Jwt").
			Deprecated(true).
			Params(NewBodyParam("param", "UpdateUserParam", true, "update user param")).
			Responses(NewResponse(200, "Result")),

		NewRoutePath("PUT", "/user", "Update user").
			Tags("User").
			Securities("Jwt").
			Params(NewBodyParam("param", "UpdateUserParam", true, "update user param")).
			Responses(NewResponse(200, "Result")),

		NewRoutePath("DELETE", "/user", "Delete user").
			Tags("User").
			Securities("Jwt").
			Responses(NewResponse(200, "Result")),

		NewRoutePath("HEAD", "/test/a", "Test a").
			Tags("Test").
			Securities("WrongSecurity").
			Params(
				NewQueryParam("q1", "string#date-time", true, "q1").Enum(0, 1, 2),
				NewQueryParam("q2", "number", false, "q2").Minimum(-5),
				NewQueryParam("q3", "string#password", true, "q3").AllowEmpty(true).Example("example").Default("default"),
				NewFormParam("f1", "file", true, "f1"),
				NewFormParam("f2", "string", true, "f2").AllowEmpty(true),
				NewHeaderParam("Authorization", "header", false, "authorization"),
			).
			Responses(
				NewResponse(200, "Result").
					Desc("200 Success").
					Headers(
						NewHeader("Content-Type", "string", "content type"),
						NewHeader("X-My-Token", "string", "my token"),
						NewHeader("X-My-Object", "UserDto", "my object"),
					),
				NewResponse(409, "string").Desc("409 Conflict"),
			),
	)

	AddDefinitions(
		NewDefinition("Result", "global response").
			Properties(
				NewProperty("code", "integer#int32", true, "status code").Example("200"),
				NewProperty("message", "string", true, "status message").Example("success"),
			),

		NewDefinition("_Result", "global response").
			Generics("T").
			Properties(
				NewProperty("code", "integer#int32", true, "status code"),
				NewProperty("message", "string", true, "status message"),
				NewProperty("data", "T", true, "response data"),
			),

		NewDefinition("_Page", "global page response").
			Generics("T").
			Properties(
				NewProperty("page", "integer#int32", true, "current page"),
				NewProperty("limit", "integer#int32", true, "page size"),
				NewProperty("total", "integer#int32", true, "total count"),
				NewProperty("data", "T[]", true, "response data"),
			),

		NewDefinition("UserDto", "user response").
			Properties(
				NewProperty("uid", "integer#int64", true, "user id"),
				NewProperty("username", "string", true, "username"),
				NewProperty("nickname", "string", true, "nickname"),
				NewProperty("profile", "string", true, "user profile").AllowEmpty(true),
				NewProperty("gender", "string", true, "user gender").Enum("secret", "male", "female"),
			),

		NewDefinition("LoginDto", "login response").
			Properties(
				NewProperty("user", "UserDto", true, "authorized user"),
				NewProperty("token", "string", true, "access token"),
			),

		NewDefinition("RegisterParam", "register param").
			Properties(
				NewProperty("username", "string", true, "username").MinLength(5).MaxLength(30),
				NewProperty("password", "string", true, "password").MinLength(5).MaxLength(30),
			),

		NewDefinition("LoginParam", "login param").
			Properties(
				NewProperty("parameter", "string", true, "login parameter"),
				NewProperty("password", "string", true, "password"),
			),

		NewDefinition("UpdateUserParam", "update user param").
			Properties(
				NewProperty("username", "string", true, "username"),
				NewProperty("nickname", "string", true, "nickname"),
				NewProperty("profile", "string", true, "user profile").AllowEmpty(true),
				NewProperty("gender", "string", true, "user gender").Enum("secret", "male", "female"),
			),
	)

	_, err := GenerateSwaggerYaml()
	if err != nil {
		log.Println(err)
		t.Fatal("yaml")
	}

	_, err = GenerateSwaggerJson()
	if err != nil {
		log.Println(err)
		t.Fatal("json")
	}

	_, err = GenerateApib()
	if err != nil {
		log.Println(err)
		t.Fatal("apib")
	}

	_, err = SaveSwaggerYaml("./docs/api.yaml")
	if err != nil {
		log.Println(err)
		t.Fatal("yaml")
	}

	_, err = SaveSwaggerJson("./docs/api.json")
	if err != nil {
		log.Println(err)
		t.Fatal("json")
	}

	_, err = SaveApib("./docs/api.apib")
	if err != nil {
		log.Println(err)
		t.Fatal("apib")
	}

	if GetDefinitions()[1].GetGenerics()[0] != "T" {
		t.Fatal(`GetDefinitions()[1].GetGenerics()[0] != "T"`)
	}
	if GetDefinitions()[1].GetProperties()[2].GetType() != "T" {
		t.Fatal(`GetDefinitions()[1].GetProperties()[2].GetType() != "T"`)
	}
}
