package yamldoc

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"testing"
)

func TestGenerateYaml(t *testing.T) {
	SetDocument(
		"localhost:10086", "/",
		NewInfo("test-api", "a demo description", "1.0").
			SetTermsOfService("http://xxx.yyy.zzz").
			SetLicense(NewLicense("MIT", "http://xxx.yyy.zzz")).
			SetContact(NewContact("author", "http://xxx.yyy.zzz", "xxx@yyy.zzz")),
	)
	SetTags(NewTag("ping", "ping-controller"), NewTag("user", "user-controller"))
	SetSecurities(NewSecurity("jwt", HEADER, "Authorization"))

	AddPaths(
		NewPath(GET, "/api/v1/ping", "ping").
			SetDescription("ping the server").
			SetTags("ping").
			SetConsumes(JSON).
			SetProduces(JSON).
			SetResponses(
				NewResponse(200).SetDescription("success").SetExamples(map[string]string{JSON: "{\n\t\"ping\": \"pong\"\n}"}),
			),
		NewPath(GET, "/api/v1/user", "get user").
			SetDescription("get user from database").
			SetTags("user").
			SetConsumes(JSON).
			SetProduces(JSON).
			SetSecurities("jwt").
			SetParams(
				NewParam("page", QUERY, INTEGER, false, "current page").SetDefault(1),
				NewParam("total", QUERY, INTEGER, false, "page size").SetDefault(10),
				NewParam("order", QUERY, STRING, false, "order string").SetDefault(""),
			).
			SetResponses(
				NewResponse(200).SetSchema(NewRefSchema("Result<Page<User>>")),
			),
		NewPath(PUT, "/api/v1/user/{id}", "update user (ugly api)").
			SetDescription("update user to database").
			SetTags("user").
			SetConsumes(JSON).
			SetProduces(JSON).
			SetSecurities("jwt").
			SetParams(
				NewParam("id", PATH, "integer", true, "user id"),
				NewParam("body", BODY, "", true, "request body").SetSchema(NewRefSchema("User")),
			).
			SetResponses(
				NewResponse(200).SetDescription("success").SetSchema(NewRefSchema("Result")).
					SetHeaders(NewHeader("Content-Type", STRING, "demo content type").SetDefault(JSON)),
				NewResponse(404).SetDescription("not found"),
			),
		NewPath(HEAD, "/api/v1/test", "test path").
			SetParams(
				NewParam("arr", QUERY, ARRAY, true, "test").SetItems(NewItems(INTEGER).SetFormat(INT64).SetItems(NewItems(INTEGER))),
				NewParam("ref", QUERY, ARRAY, true, "test").SetItems(NewRefItems("User")),
				NewParam("ref2", QUERY, ARRAY, true, "test").SetItems(NewRefItems("User")),
				NewParam("enum", QUERY, STRING, true, "test").SetEnum("male", "female"),
			),
	)

	AddModels(
		NewModel("Result", "global response").SetProperties(
			NewProperty("code", "integer", true, "status code"),
			NewProperty("message", "string", true, "status message"),
		),
		NewModel("User", "user response").SetProperties(
			NewProperty("id", "integer", true, "user id"),
			NewProperty("name", "string", true, "user name"),
			NewProperty("profile", "string", false, "user profile").SetAllowEmptyValue(true),
			NewProperty("gender", "string", true, "user gender").SetEnum("male", "female"),
			NewProperty("create_at", "datetime", true, "user register time").SetFormat("yyyy-MM-dd HH:mm:ss"),
			NewProperty("birthday", "date", true, "user birthday").SetFormat("yyyy-MM-dd"),
		),
		NewModel("Page<User>", "user response").SetProperties(
			NewProperty("page", "integer", true, "current page"),
			NewProperty("total", "integer", true, "data count"),
			NewProperty("limit", "integer", true, "page size"),
			NewProperty("data", "array", true, "page data").SetItems(NewRefItems("User")),
		),
		NewModel("Result<Page<User>>", "user response").SetProperties(
			NewProperty("code", "integer", true, "status code"),
			NewProperty("message", "string", true, "status message"),
			NewRefProperty("data", "Page<User>", true, "result data"),
		),
	)

	doc, _ := yaml.Marshal(appendKvs(buildDocument(_document), map[string]interface{}{"swagger": "2.0"}))
	fmt.Println(string(doc))

	err := GenerateYaml("./docs/api.yaml", map[string]interface{}{"swagger": "2.0"})
	log.Println(err)
}
