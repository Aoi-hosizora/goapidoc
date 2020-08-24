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

func TestPreHandleGenerics(t *testing.T) {
	def := &Definition{
		generics: []string{"T", "U", "V"},
		properties: []*Property{
			{typ: "inT[]"},
			{typ: "O<inT[], T[], T, inT<int>>"},
			{typ: "T"},
			{typ: "tT<T<tT>[][], T>[]"},
			{typ: "TtT<T,tT[],T[][]>[]"},
		},
	}
	preHandleDefinitionForGeneric(def)

	if def.generics[0] != "«T»" {
		t.Fatal("def.generics[0]")
	}
	if def.generics[1] != "«U»" {
		t.Fatal("def.generics[1]")
	}
	if def.generics[2] != "«V»" {
		t.Fatal("def.generics[2]")
	}

	p0 := def.properties[0]
	p1 := def.properties[1]
	p2 := def.properties[2]
	p3 := def.properties[3]
	p4 := def.properties[4]
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

func TestGenerate(t *testing.T) {
	SetDocument(
		"localhost:10086", "/",
		NewInfo("test-api", "a demo description", "1.0").
			TermsOfService("http://xxx.yyy.zzz").
			License(NewLicense("MIT", "http://xxx.yyy.zzz")).
			Contact(NewContact("author", "http://xxx.yyy.zzz", "xxx@yyy.zzz")),
	)
	SetTags(
		NewTag("ping", "ping-controller"),
		NewTag("user", "user-controller"),
	)
	SetSecurities(
		NewSecurity("jwt", HEADER, "Authorization"),
	)

	AddPaths(
		NewRoutePath(GET, "/api/v1/ping", "ping").
			Desc("ping the server").
			Tags("ping").
			Consumes(JSON).
			Produces(JSON).
			Responses(
				NewResponse(200).Desc("success").Examples(map[string]string{JSON: "{\n    \"ping\": \"pong\"\n}"}),
			),
		NewRoutePath(GET, "/api/v1/user", "get users").
			Tags("user").
			Consumes(JSON).
			Produces(JSON).
			Securities("jwt").
			Params(
				NewQueryParam("page", INTEGER, false, "current page").Default(1).Minimum(1).Maximum(50),
				NewQueryParam("total", INTEGER, false, "page size").Default(10).Example(20),
				NewQueryParam("order", STRING, false, "order string").Default("").MinLength(1).MaxLength(50),
			).
			Responses(
				NewResponse(200).Type("_Result<_Page<User>>"),
			),
		NewRoutePath(GET, "/api/v1/user/{id}", "get a user").
			Tags("user").
			Consumes(JSON).
			Produces(JSON).
			Params(NewPathParam("id", INTEGER, true, "user id")).
			Responses(
				NewResponse(200).Type("_Result<User>"),
			),
		NewRoutePath(PUT, "/api/v1/user/{id}", "ugly update user").
			Deprecated(true).
			Tags("user").
			Consumes(JSON).
			Produces(JSON).
			Securities("jwt").
			Params(
				NewPathParam("id", INTEGER, true, "user id"),
				NewBodyParam("body", "User", true, "request body"),
			).
			Responses(
				NewResponse(200).Type("Result").Desc("success"),
				NewResponse(404).Desc("not found").Headers(NewHeader("Content-Kind", STRING, "demo")),
				NewResponse(400).Type(STRING).Desc("bad request").Examples(map[string]string{JSON: "bad request"}),
			),
		NewRoutePath(HEAD, "/api/v1/test", "test path").
			Params(
				NewQueryParam("arr", "integer#int64[]", true, "test"),
				NewQueryParam("ref", "User[]", true, "test"),
				NewQueryParam("enum", STRING, true, "test").Enum("male", "female"),
				NewQueryParam("option1", "_Result<string>[]", true, "test"),
				NewQueryParam("option2", "_Result<string[]>[]", true, "test"),
				NewBodyParam("test", "_ResultPage<User>", true, "test"),
				NewBodyParam("arr2", INTEGER, true, "test"),
			).
			Responses(NewResponse(200).Type("TestGeneric<integer, User, string>")),
	)

	AddDefinitions(
		NewDefinition("Result", "global response").Properties(
			NewProperty("code", INTEGER, true, "status code"),
			NewProperty("message", STRING, true, "status message"),
		),
		NewDefinition("User", "user response").Properties(
			NewProperty("id", INTEGER, true, "user id").Minimum(1).Maximum(65535),
			NewProperty("name", STRING, true, "user name"),
			NewProperty("profile", STRING, false, "user profile").AllowEmpty(true).MinLength(1).MaxLength(255),
			NewProperty("gender", STRING, true, "user gender").Enum("male", "female").Example("female"),
			NewProperty("create_at", "string#date-time", true, "user register time"),
			NewProperty("birthday", "string#date", true, "user birthday"),
			NewProperty("scores", "number[]", true, "user scores"),
		),
		NewDefinition("_Result", "global response").Generics("T").Properties(
			NewProperty("code", INTEGER, true, "status code"),
			NewProperty("message", STRING, true, "status message"),
			NewProperty("data", "T", true, "response data"),
		),
		NewDefinition("_Page", "global page response").Generics("T").Properties(
			NewProperty("page", INTEGER, true, "current page"),
			NewProperty("total", INTEGER, true, "data count"),
			NewProperty("limit", INTEGER, true, "page size"),
			NewProperty("data", "T[]", true, "page data"),
		),
		NewDefinition("_ResultPage", "global response").Generics("T").Properties(
			NewProperty("code", INTEGER, true, "status code"),
			NewProperty("message", INTEGER, true, "status message"),
			NewProperty("data", "_Page<T>", true, "response data"),
		),

		NewDefinition("TestGeneric", "test generics").Generics("T", "U", "W").Properties(
			NewProperty("t", "T", true, "t"),
			NewProperty("t2", "T[]", true, "t2"),
			NewProperty("t3", "T[][]", true, "t3"),
			NewProperty("u", "_Result<U>", true, "u"),
			NewProperty("u2", "_Result<U[]>", true, "u2"),
			NewProperty("u3", "_Result<U[]>[]", true, "u3"),
			NewProperty("w", "W", true, "w"),
		),
	)

	_, err := GenerateSwaggerYaml("./docs/api.yaml")
	log.Println(err)

	_, err = GenerateSwaggerJson("./docs/api.json")
	log.Println(err)

	_, err = GenerateApib("./docs/api.apib")
	log.Println(err)
}
