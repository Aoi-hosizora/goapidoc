package yamldoc

import (
	"fmt"
	"gopkg.in/yaml.v2"
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
	SetSecurities(NewSecurity("jwt", "header", "Authorization"))

	AddPaths(
		NewPath("GET", "/api/v1/ping", "ping").
			SetDescription("ping the server").
			SetTags("ping").
			SetConsumes(JSON).
			SetProduces(JSON).
			SetResponses(
				NewResponse(200).SetDescription("success").SetExamples(map[string]string{JSON: `{"ping": "pong"}`}),
			),
		NewPath("GET", "/api/v1/user", "get user").
			SetDescription("get user from database").
			SetTags("user").
			SetConsumes(JSON).
			SetProduces(JSON).
			SetSecurities("jwt").
			SetParams(
				NewParam("page", "query", "integer", false, "current page").SetDefault(1),
				NewParam("total", "query", "integer", false, "page size").SetDefault(10),
				NewParam("order", "query", "string", false, "order string").SetDefault(""),
			).
			SetResponses(
				NewResponse(200).SetSchema("Result<Page<User>>"),
			),
		NewPath("PUT", "/api/v1/user/{id}", "update user (ugly api)").
			SetDescription("update user to database").
			SetTags("user").
			SetConsumes(JSON).
			SetProduces(JSON).
			SetSecurities("jwt").
			SetParams(
				NewParam("id", "path", "integer", true, "user id"),
				NewParam("body", "body", "", true, "request body").SetSchema("User"),
			).
			SetResponses(
				NewResponse(200).SetDescription("success").SetSchema("Result").
					SetHeaders(NewHeader("Content-Type", "Content-Type", "string").SetDefault(JSON)),
				NewResponse(404).SetDescription("not found"),
			),
		NewPath("HEADER", "/api/v1/test", "test path").
			SetParams(
				NewParam("arr", "query", "array", true, "test").SetItems(NewItems("integer")),
				NewParam("arr-ref", "query", "array", true, "test").SetItems(NewItems("").SetSchema("User")),
				NewParam("enum", "query", "string", true, "test").SetEnum("male", "female"),
			),
	)

	AddModels(
		NewModel("Result", "global response").SetProperties(
			NewProperty("code", "status code", "integer", true),
			NewProperty("message", "status message", "string", true),
		),
		NewModel("User", "user response").SetProperties(
			NewProperty("id", "user id", "integer", true),
			NewProperty("name", "user name", "string", true),
			NewProperty("profile", "user profile", "string", false).SetAllowEmptyValue(true),
			NewProperty("gender", "user gender", "string", true).SetEnum("male", "female"),
			NewProperty("create_at", "user register time", "datetime", true).SetFormat("yyyy-MM-dd HH:mm:ss"),
			NewProperty("birthday", "user birthday", "date", true).SetFormat("yyyy-MM-dd"),
		),
		NewModel("Page<User>", "user response").SetProperties(
			NewProperty("page", "current page", "integer", true),
			NewProperty("total", "data count", "integer", true),
			NewProperty("limit", "page size", "integer", true),
			NewProperty("data", "page data", "array", true).SetItems(NewItems("").SetSchema("User")),
		),
		NewModel("Result<Page<User>>", "user response").SetProperties(
			NewProperty("code", "status code", "integer", true),
			NewProperty("message", "status message", "string", true),
			NewProperty("data", "result data", "object", true).SetSchema("Page<User>"),
		),
	)

	doc, _ := yaml.Marshal(appendKvs(mapToInnerDocument(_document), map[string]interface{}{"swagger": "2.0"}))
	fmt.Println(string(doc))
}
