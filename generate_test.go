package goapidoc

import (
	"log"
	"testing"
)

func TestParseInnerType(t *testing.T) {
	str := "aType<T1<>, T2<TT1<integer[]>>, T3<TT1, TT2>, T4>[][]"
	res := parseApiType(str)
	obj := res.outArray.typ.outArray.typ
	log.Println(res.name, res.outArray)
	log.Println(res.outArray.typ.name, res.outArray.typ.outArray)
	log.Println(obj.name, obj.outObject.typ, obj.outObject.generic)

	g0 := obj.outObject.generic[0]   // T1<>
	g1 := obj.outObject.generic[1]   // T2<TT1<integer[]>>
	g10 := g1.outObject.generic[0]   // TT1<integer[]>
	g100 := g10.outObject.generic[0] // integer[]
	g1000 := g100.outArray.typ       // integer
	g2 := obj.outObject.generic[2]   // T3<TT1, TT2>
	g20 := g2.outObject.generic[0]   // TT1
	g21 := g2.outObject.generic[1]   // TT2
	g3 := obj.outObject.generic[3]   // T4
	log.Println(g0.name, g0.outObject.typ)
	log.Println(g1.name, g1.outObject.typ, g10.outObject.typ, g100.name, g1000.name)
	log.Println(g2.name, g2.outObject.typ, g20.outObject.typ, g21.outObject.typ)
	log.Println(g3.name, g3.outObject.typ)
}

func TestPreHandleGeneric(t *testing.T) {
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
	preHandleGeneric(def)
	log.Println(def.generics)
	for _, p := range def.properties {
		log.Println(p.typ)
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
		NewRoutePath(PUT, "/api/v1/user/{id}", "update user (ugly api)").
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

		NewDefinition("TestGeneric", "test generic").Generics("T", "U", "W").Properties(
			NewProperty("t", "T", true, "t"),
			NewProperty("t2", "T[]", true, "t2"),
			NewProperty("t3", "T[][]", true, "t3"),
			NewProperty("u", "_Result<U>", true, "u"),
			NewProperty("u2", "_Result<U[]>", true, "u2"),
			NewProperty("u3", "_Result<U[]>[]", true, "u3"),
			NewProperty("w", "W", true, "w"),
		),
	)

	doc, err := GenerateSwaggerYaml("./docs/api.yaml")
	log.Println(err)

	doc, err = GenerateSwaggerJson("./docs/api.json")
	log.Println(string(doc), err)

	doc, err = GenerateApib("./docs/api.apib")
	log.Println(string(doc), err)
}
