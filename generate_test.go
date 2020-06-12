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
				NewResponse(200).WithSchema(RefSchema("_Result", "data", RefSchema("_Page", "data", "User"))),
			),
		NewPath(GET, "/api/v1/user/{id}", "get a user").
			WithTags("user").
			WithConsumes(JSON).
			WithProduces(JSON).
			WithParams(NewParam("id", PATH, INTEGER, true, "user id")).
			WithResponses(
				NewResponse(200).WithSchema(RefSchema("_Result", "data", "User")),
			),
		NewPath(PUT, "/api/v1/user/{id}", "update user (ugly api)").
			WithTags("user").
			WithConsumes(JSON).
			WithProduces(JSON).
			WithSecurities("jwt").
			WithParams(
				NewParam("id", PATH, INTEGER, true, "user id"),
				NewParam("body", BODY, OBJECT, true, "request body").WithSchema(RefSchema("User")),
			).
			WithResponses(
				NewResponse(200).WithDescription("success").WithSchema(RefSchema("Result")),
				NewResponse(404).WithDescription("not found").WithHeaders(NewHeader("Content-Type", STRING, "demo")),
				NewResponse(400).WithDescription("bad request").WithSchema(NewSchema(STRING, true)).WithExamples(map[string]string{JSON: "bad request"}),
			),
		NewPath(HEAD, "/api/v1/test", "test path").
			WithParams(
				NewParam("arr", QUERY, ARRAY, true, "test").WithItems(ArrItems(NewItems(INTEGER).SetFormat(INT64))),
				NewParam("ref", QUERY, ARRAY, true, "test").WithItems(RefItems("User")),
				NewParam("enum", QUERY, STRING, true, "test").WithEnum("male", "female"),
				NewParam("option1", QUERY, ARRAY, true, "test").WithItems(RefItems("Result", "code", NewSchema(STRING, true))),
				NewParam("option2", QUERY, ARRAY, true, "test").WithItems(RefItems("Result", "code", NewItems(STRING))),
				NewParam("arr2", BODY, ARRAY, true, "test").WithSchema(ArrSchema(NewItems(INTEGER))),
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
			NewProperty("create_at", STRING, true, "user register time").WithFormat(DATETIME),
			NewProperty("birthday", STRING, true, "user birthday").WithFormat(DATE),
			NewProperty("scores", ARRAY, true, "user scores").WithItems(NewItems(NUMBER)),
		),
		NewDefinition("!Page<User>", "user response").WithProperties(
			NewProperty("page", INTEGER, true, "current page"),
			NewProperty("total", INTEGER, true, "data count"),
			NewProperty("limit", INTEGER, true, "page size"),
			NewArrayProperty("data", RefItems("User"), true),
		),
		NewDefinition("!Result<Page<User>>", "user response").WithProperties(
			NewProperty("code", INTEGER, true, "status code"),
			NewProperty("message", STRING, true, "status message"),
			NewObjectProperty("data", "Page<User>", true),
		),

		NewDefinition("_Result", "global response").WithProperties(
			NewProperty("code", INTEGER, true, "status code"),
			NewProperty("message", INTEGER, true, "status message"),
			NewProperty("data", OBJECT, true, "response data"),
		),
		NewDefinition("_Page", "global page response").WithProperties(
			NewProperty("page", INTEGER, true, "current page"),
			NewProperty("total", INTEGER, true, "data count"),
			NewProperty("limit", INTEGER, true, "page size"),
			NewProperty("data", ARRAY, true, "page data"),
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
