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
		Generics: []string{"T", "U", "V"},
		Properties: []*Property{
			{Type: "inT[]"},
			{Type: "O<inT[], T[], T, inT<int>>"},
			{Type: "T"},
			{Type: "tT<T<tT>[][], T>[]"},
			{Type: "TtT<T,tT[],T[][]>[]"},
		},
	}
	preHandleGeneric(def)
	log.Println(def.Generics)
	for _, p := range def.Properties {
		log.Println(p.Type)
	}
}

func TestGenerate(t *testing.T) {
	SetDocument(
		"localhost:10086", "/",
		NewInfo("test-api", "a demo description", "1.0").
			WithTermsOfService("http://xxx.yyy.zzz").
			WithLicense(NewLicense("MIT", "http://xxx.yyy.zzz")).
			WithContact(NewContact("author", "http://xxx.yyy.zzz", "xxx@yyy.zzz")),
	)
	SetTags(
		NewTag("ping", "ping-controller"),
		NewTag("user", "user-controller"),
	)
	SetSecurities(
		NewSecurity("jwt", HEADER, "Authorization"),
	)

	AddPaths(
		NewPath(GET, "/api/v1/ping", "ping").
			WithDescription("ping the server").
			WithTags("ping").
			WithConsumes(JSON).
			WithProduces(JSON).
			WithResponses(
				NewResponse(200).WithDescription("success").WithExamples(map[string]string{JSON: "{\n    \"ping\": \"pong\"\n}"}),
			),
		NewPath(GET, "/api/v1/user", "get users").
			WithTags("user").
			WithConsumes(JSON).
			WithProduces(JSON).
			WithSecurities("jwt").
			WithParams(
				NewQueryParam("page", INTEGER, false, "current page").WithDefault(1).WithMinimum(1).WithMaximum(50),
				NewQueryParam("total", INTEGER, false, "page size").WithDefault(10).WithExample(20),
				NewQueryParam("order", STRING, false, "order string").WithDefault("").WithMinLength(1).WithMaxLength(50),
			).
			WithResponses(
				NewResponse(200).WithType("_Result<_Page<User>>"),
			),
		NewPath(GET, "/api/v1/user/{id}", "get a user").
			WithTags("user").
			WithConsumes(JSON).
			WithProduces(JSON).
			WithParams(NewPathParam("id", INTEGER, true, "user id")).
			WithResponses(
				NewResponse(200).WithType("_Result<User>"),
			),
		NewPath(PUT, "/api/v1/user/{id}", "update user (ugly api)").
			WithTags("user").
			WithConsumes(JSON).
			WithProduces(JSON).
			WithSecurities("jwt").
			WithParams(
				NewPathParam("id", INTEGER, true, "user id"),
				NewBodyParam("body", "User", true, "request body"),
			).
			WithResponses(
				NewResponse(200).WithType("Result").WithDescription("success"),
				NewResponse(404).WithDescription("not found").WithHeaders(NewHeader("Content-Kind", STRING, "demo")),
				NewResponse(400).WithType(STRING).WithDescription("bad request").WithExamples(map[string]string{JSON: "bad request"}),
			),
		NewPath(HEAD, "/api/v1/test", "test path").
			WithParams(
				NewQueryParam("arr", "integer#int64[]", true, "test"),
				NewQueryParam("ref", "User[]", true, "test"),
				NewQueryParam("enum", STRING, true, "test").WithEnum("male", "female"),
				NewQueryParam("option1", "_Result<string>[]", true, "test"),
				NewQueryParam("option2", "_Result<string[]>[]", true, "test"),
				NewBodyParam("test", "_ResultPage<User>", true, "test"),
				NewBodyParam("arr2", INTEGER, true, "test"),
			).
			WithResponses(NewResponse(200).WithType("TestGeneric<integer, User, string>")),
	)

	AddDefinitions(
		NewDefinition("Result", "global response").WithProperties(
			NewProperty("code", INTEGER, true, "status code"),
			NewProperty("message", STRING, true, "status message"),
		),
		NewDefinition("User", "user response").WithProperties(
			NewProperty("id", INTEGER, true, "user id").WithMinimum(1).WithMaximum(65535),
			NewProperty("name", STRING, true, "user name"),
			NewProperty("profile", STRING, false, "user profile").WithAllowEmptyValue(true).WithMinLength(1).WithMaxLength(255),
			NewProperty("gender", STRING, true, "user gender").WithEnum("male", "female").WithExample("female"),
			NewProperty("create_at", "string#date-time", true, "user register time"),
			NewProperty("birthday", "string#date", true, "user birthday"),
			NewProperty("scores", "number[]", true, "user scores"),
		),
		NewDefinition("_Result", "global response").WithGenerics("T").WithProperties(
			NewProperty("code", INTEGER, true, "status code"),
			NewProperty("message", STRING, true, "status message"),
			NewProperty("data", "T", true, "response data"),
		),
		NewDefinition("_Page", "global page response").WithGenerics("T").WithProperties(
			NewProperty("page", INTEGER, true, "current page"),
			NewProperty("total", INTEGER, true, "data count"),
			NewProperty("limit", INTEGER, true, "page size"),
			NewProperty("data", "T[]", true, "page data"),
		),
		NewDefinition("_ResultPage", "global response").WithGenerics("T").WithProperties(
			NewProperty("code", INTEGER, true, "status code"),
			NewProperty("message", INTEGER, true, "status message"),
			NewProperty("data", "_Page<T>", true, "response data"),
		),

		NewDefinition("TestGeneric", "test generic").WithGenerics("T", "U", "W").WithProperties(
			NewProperty("t", "T", true, "t"),
			NewProperty("t2", "T[]", true, "t2"),
			NewProperty("t3", "T[][]", true, "t3"),
			NewProperty("u", "_Result<U>", true, "u"),
			NewProperty("u2", "_Result<U[]>", true, "u2"),
			NewProperty("u3", "_Result<U[]>[]", true, "u3"),
			NewProperty("w", "W", true, "w"),
		),
	)

	doc, err := GenerateYamlWithSwagger2("./docs/api.yaml")
	log.Println(err)

	doc, err = GenerateJsonWithSwagger2("./docs/api.json")
	log.Println(string(doc), err)
}
