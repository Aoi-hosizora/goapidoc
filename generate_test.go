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
			SetConsumes("application/json").
			SetProduces("application/json"),
		NewPath("GET", "/api/v1/user/{id}", "get user").
			SetDescription("get user from database").
			SetTags("user").
			SetConsumes("application/json").
			SetProduces("application/json").
			SetSecurities("jwt").
			SetParams(
				NewParam("id", "path", "integer", true, "user id"),
				NewParam("page", "query", "integer", false, "current page").SetDefault(1),
				NewParam("total", "query", "integer", false, "page size").SetDefault(10),
				NewParam("order", "query", "string", false, "order string").SetDefault(""),
			),
	)

	doc, _ := yaml.Marshal(appendKvs(mapToInnerDocument(_document), map[string]interface{}{"swagger": "2.0"}))
	fmt.Println(string(doc))
}
