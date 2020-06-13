package goapidoc

import (
	"fmt"
	"log"
	"testing"
)

func TestGenerateYaml(t *testing.T) {
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
				NewResponse(200).WithDescription("success").WithExamples(map[string]string{JSON: "{\n\t\"ping\": \"pong\"\n}"}),
			),
		NewPath(GET, "/api/v1/user", "get users").
			WithTags("user").
			WithConsumes(JSON).
			WithProduces(JSON).
			WithSecurities("jwt").
			WithParams(
				NewParam("page", QUERY, INTEGER, false, "current page").WithDefault(1),
				NewParam("total", QUERY, INTEGER, false, "page size").WithDefault(10),
				NewParam("order", QUERY, STRING, false, "order string").WithDefault(""),
			).
			WithResponses(
				NewResponse(200).WithType("_Result<_Page<User>>"),
				// NewResponse(200).WithSchema(RefSchema("_Result", "data", RefSchema("_Page", "data", "User"))),
			),
		NewPath(GET, "/api/v1/user/{id}", "get a user").
			WithTags("user").
			WithConsumes(JSON).
			WithProduces(JSON).
			WithParams(NewParam("id", PATH, INTEGER, true, "user id")).
			WithResponses(
				NewResponse(200).WithType("_Result<User>"),
				// NewResponse(200).WithSchema(RefSchema("_Result", "data", "User")),
			),
		NewPath(PUT, "/api/v1/user/{id}", "update user (ugly api)").
			WithTags("user").
			WithConsumes(JSON).
			WithProduces(JSON).
			WithSecurities("jwt").
			WithParams(
				NewParam("id", PATH, INTEGER, true, "user id"),
				// NewParam("body", BODY, OBJECT, true, "request body").WithSchema(RefSchema("User")),
				NewParam("body", BODY, "User", true, "request body"),
			).
			WithResponses(
				// NewResponse(200).WithDescription("success").WithSchema(RefSchema("Result")),
				NewResponse(200).WithType("Result").WithDescription("success"),
				NewResponse(404).WithDescription("not found").WithHeaders(NewHeader("Content-Kind", STRING, "demo")),
				// NewResponse(400).WithDescription("bad request").WithSchema(NewSchema(STRING, true)).WithExamples(map[string]string{JSON: "bad request"}),
				NewResponse(400).WithType(STRING).WithDescription("bad request").WithExamples(map[string]string{JSON: "bad request"}),
			),
		NewPath(HEAD, "/api/v1/test", "test path").
			WithParams(
				// NewParam("arr", QUERY, ARRAY, true, "test").WithItems(ArrItems(NewItems(INTEGER).SetFormat(INT64))),
				// NewParam("ref", QUERY, ARRAY, true, "test").WithItems(RefItems("User")),
				// NewParam("enum", QUERY, STRING, true, "test").WithEnum("male", "female"),
				// NewParam("option1", QUERY, ARRAY, true, "test").WithItems(RefItems("Result", "code", NewSchema(STRING, true))),
				// NewParam("option2", QUERY, ARRAY, true, "test").WithItems(RefItems("Result", "code", NewItems(STRING))),
				// NewParam("arr2", BODY, ARRAY, true, "test").WithSchema(ArrSchema(NewItems(INTEGER))),
				NewParam("arr", QUERY, "integer#int64[]", true, "test"),
				NewParam("ref", QUERY, "User[]", true, "test"),
				NewParam("enum", QUERY, STRING, true, "test").WithEnum("male", "female"),
				NewParam("option1", QUERY, "Result<string>[]", true, "test"),
				NewParam("option2", QUERY, "Result<string[]>[]", true, "test"),
				NewParam("test", BODY, "_ResultPage<User>", true, "test"),
				NewParam("arr2", BODY, INTEGER, true, "test"),
			),
	)

	AddDefinitions(
		NewDefinition("Result", "global response").WithProperties(
			NewProperty("code", INTEGER, true, "status code"),
			NewProperty("message", STRING, true, "status message"),
		),
		NewDefinition("User", "user response").WithProperties(
			NewProperty("id", INTEGER, true, "user id"),
			NewProperty("name", STRING, true, "user name"),
			NewProperty("profile", STRING, false, "user profile").WithAllowEmptyValue(true),
			NewProperty("gender", STRING, true, "user gender").WithEnum("male", "female"),
			NewProperty("create_at", "string#date-time", true, "user register time"),
			NewProperty("birthday", "string#date", true, "user birthday"),
			NewProperty("scores", "number[]", true, "user scores"),
		),
		NewDefinition("!Page<User>", "user response").WithProperties(
			NewProperty("page", INTEGER, true, "current page"),
			NewProperty("total", INTEGER, true, "data count"),
			NewProperty("limit", INTEGER, true, "page size"),
			NewProperty("data", "User[]", true, "page data"),
		),
		NewDefinition("!Result<Page<User>>", "user response").WithProperties(
			NewProperty("code", INTEGER, true, "status code"),
			NewProperty("message", STRING, true, "status message"),
			NewProperty("data", "!Page<User>", true, "result data"),
		),

		NewDefinition("_Result", "global response").WithGenerics("T").WithProperties(
			NewProperty("code", INTEGER, true, "status code"),
			NewProperty("message", INTEGER, true, "status message"),
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
	)

	// doc, _ := yaml.Marshal(appendKvs(buildDocument(_document), map[string]interface{}{"swagger": "2.0"}))
	// fmt.Println(string(doc))

	doc, _ := jsonMarshal(appendKvs(buildDocument(_document), map[string]interface{}{"swagger": "2.0"}))
	fmt.Println(string(doc))

	err := GenerateYaml("./docs/api.yaml", map[string]interface{}{"swagger": "2.0"})
	log.Println(err)

	err = GenerateJson("./docs/api.json", map[string]interface{}{"swagger": "2.0"})
	log.Println(err)
}
