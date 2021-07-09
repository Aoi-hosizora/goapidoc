package goapidoc

import (
	"log"
	"testing"
)

func TestGetSet(t *testing.T) {
	SetHost("localhost:12334")
	SetBasePath("v1")
	SetInfo(NewInfo("", "", "").
		Title("Test api").
		Desc("This is a test api").
		Version("v0.0.1").
		TermsOfService("ToS").
		License(NewLicense("", "").
			Name("MIT").
			Url("https://opensource.org/licenses/MIT")).
		Contact(NewContact("", "", "").
			Name("Aoi-hosizora").
			Url("https://github.com/Aoi-hosizora").
			Email("a970335605@hotmail.com")))
	SetOption(NewOption().
		Schemas("http").
		AddSchemas("https", "ws", "wss").
		Consumes("application/json").
		AddConsumes("multipart/form-data", "application/protobuf").
		Produces("application/json").
		AddProduces("application/xml", "application/protobuf").
		Tags(NewTag("", "").
			Name("Authorization").
			Desc("auth-controller")).
		AddTags(NewTag("User", "user-controller"),
			NewTag("Resource", "resource-controller")).
		Securities(NewSecurity("", "").
			Title("jwt").
			Type(APIKEY).
			Desc("A apiKey security called jwt").
			InLoc(HEADER).
			Name("Authorization")).
		AddSecurities(NewBasicSecurity(BASIC),
			NewApiKeySecurity("another_jwt", HEADER, "X-JWT")),
	)
	AddOperations(NewOperation("", "", ""))
	SetOperations()
	AddDefinitions(NewDefinition("", ""))
	SetDefinitions()

	if GetHost() != "localhost:12334" {
		failNow(t, "SetHost has a wrong behavior")
	}
	if GetBasePath() != "v1" {
		failNow(t, "SetBasePath has a wrong behavior")
	}
	if GetInfo().GetTitle() != "Test api" {
		failNow(t, "Info.Title has a wrong behavior")
	}
	if GetInfo().GetDesc() != "This is a test api" {
		failNow(t, "Info.Desc has a wrong behavior")
	}
	if GetInfo().GetVersion() != "v0.0.1" {
		failNow(t, "Info.Version has a wrong behavior")
	}
	if GetInfo().GetTermsOfService() != "ToS" {
		failNow(t, "Info.TermsOfService has a wrong behavior")
	}
	if GetInfo().GetLicense().GetName() != "MIT" {
		failNow(t, "Info.License.Name has a wrong behavior")
	}
	if GetInfo().GetLicense().GetUrl() != "https://opensource.org/licenses/MIT" {
		failNow(t, "Info.License.Url has a wrong behavior")
	}
	if GetInfo().GetContact().GetName() != "Aoi-hosizora" {
		failNow(t, "Info.Contact.Name has a wrong behavior")
	}
	if GetInfo().GetContact().GetUrl() != "https://github.com/Aoi-hosizora" {
		failNow(t, "Info.Contact.Url has a wrong behavior")
	}
	if GetInfo().GetContact().GetEmail() != "a970335605@hotmail.com" {
		failNow(t, "Info.Contact.Email has a wrong behavior")
	}
	if s := GetOption().GetSchemas(); s[0] != "http" || s[1] != "https" || s[2] != "ws" || s[3] != "wss" {
		failNow(t, "Option.Schemas or Option.AddSchemas has a wrong behavior")
	}
	if c := GetOption().GetConsumes(); c[0] != "application/json" || c[1] != "multipart/form-data" || c[2] != "application/protobuf" {
		failNow(t, "Option.Consumes or Option.AddConsumes has a wrong behavior")
	}
	if p := GetOption().GetProduces(); p[0] != "application/json" || p[1] != "application/xml" || p[2] != "application/protobuf" {
		failNow(t, "Option.Produces or Option.AddProduces has a wrong behavior")
	}
	if a := GetOption().GetTags(); a[0].GetName() != "Authorization" || a[0].GetDesc() != "auth-controller" ||
		a[1].GetName() != "User" || a[1].GetDesc() != "user-controller" || a[2].GetName() != "Resource" || a[2].GetDesc() != "resource-controller" {
		failNow(t, "Option.Tags or Option.AddTags or Tags.XXX has a wrong behavior")
	}
	if s := GetOption().GetSecurities(); s[0].GetTitle() != "jwt" || s[0].GetType() != "apiKey" || s[0].GetDesc() != "A apiKey security called jwt" || s[0].GetInLoc() != "header" || s[0].GetName() != "Authorization" ||
		s[1].GetTitle() != "basic" || s[1].GetType() != "basic" || s[2].GetTitle() != "another_jwt" || s[2].GetType() != "apiKey" || s[2].GetInLoc() != "header" || s[2].GetName() != "X-JWT" {
		failNow(t, "Option.Securities or Option.AddSecurities or Security.XXX has a wrong behavior")
	}
	if len(GetOperations()) != 0 {
		failNow(t, "Option.AddOperations or Option.SetOperations has a wrong behavior")
	}
	if len(GetDefinitions()) != 0 {
		failNow(t, "Option.AddDefinitions or Option.SetDefinitions has a wrong behavior")
	}

	SetDocument("", "", nil)
}

func TestGenerate(t *testing.T) {
	// https://editor.swagger.io/

	SetDocument(
		"localhost:65530", "/",
		NewInfo("goapidoc", "goapidoc test api", "1.0").
			TermsOfService("https://github.com/Aoi-hosizora").
			License(NewLicense("MIT", "https://github.com/Aoi-hosizora")).
			Contact(NewContact("Aoi-hosizora", "https://github.com/Aoi-hosizora", "a970335605@hotmail.com")),
	)
	SetOption(
		NewOption().Tags(
			NewTag("Authorization", "auth-controller"),
			NewTag("User", "user-controller"),
			NewTag("Test", "test-controller"),
		).Securities(
			NewApiKeySecurity("Jwt", "header", "Authorization"),
		),
	)

	AddOperations(
		NewOperation("POST", "/auth/register", "Register").
			Desc("Register.").
			Tags("Authorization").
			Params(NewBodyParam("param", "RegisterParam", true, "register param")).
			Responses(NewResponse(200, "Result")),

		NewOperation("POST", "/auth/login", "Login").
			Desc("Login.").
			Tags("Authorization").
			Params(NewBodyParam("param", "LoginParam", true, "login param")).
			Responses(
				NewResponse(200, "_Result<LoginDto>"),
				NewResponse(400, "Result").Examples(map[string]string{"application/json": "{\n  \"code\": 400, \n  \"message\": \"Unauthorized\"\n}"}),
			),

		NewOperation("DELETE", "/auth/logout", "Logout").
			Tags("Authorization").
			Securities("Jwt").
			Responses(NewResponse(200, "Result")),

		NewOperation("GET", "/user", "Get users").
			Tags("User").
			Securities("Jwt").
			Params(
				NewQueryParam("page", "integer#int32", false, "current page").Default(1).Example(1).Minimum(1),
				NewQueryParam("limit", "integer#int32", false, "page size").Default(20).Example(20).Minimum(15),
			).
			Responses(NewResponse(200, "_Result<_Page<UserDto>>")),

		NewOperation("GET", "/user/{username}", "Get a user").
			Tags("User").
			Securities("Jwt").
			Params(NewPathParam("username", "string", true, "username")).
			Responses(NewResponse(200, "_Result<UserDto>")),

		NewOperation("PUT", "/user/deprecated", "Update user").
			Tags("User").
			Securities("Jwt").
			Deprecated(true).
			Params(NewBodyParam("param", "UpdateUserParam", true, "update user param")).
			Responses(NewResponse(200, "Result")),

		NewOperation("PUT", "/user", "Update user").
			Tags("User").
			Securities("Jwt").
			Params(NewBodyParam("param", "UpdateUserParam", true, "update user param")).
			Responses(NewResponse(200, "Result")),

		NewOperation("DELETE", "/user", "Delete user").
			Tags("User").
			Securities("Jwt").
			Responses(NewResponse(200, "Result")),

		NewOperation("HEAD", "/test/a", "Test a").
			Tags("Test").
			Securities("WrongSecurity").
			Params(
				NewQueryParam("q1", "string#date-time", true, "q1").Enum(0, 1, 2),
				NewQueryParam("q2", "number", false, "q2").Minimum(-5),
				NewQueryParam("q3", "string#password", true, "q3").AllowEmpty(true).Example("example").Default("default"),
				NewFormParam("f1", "file", true, "f1"),
				NewFormParam("f2", "string", true, "f2").AllowEmpty(true),
				NewHeaderParam("Authorization", "string", false, "authorization"),
			).
			Responses(
				NewResponse(200, "Result").
					Desc("200 Success").
					Headers(
						NewHeader("Content-Type", "string", "content type"),
						NewHeader("X-My-Token", "string", "my token"),
						NewHeader("X-My-Number", "number", "my number"),
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

	if _document.definitions[1].GetGenerics()[0] != "T" {
		t.Fatal(`GetDefinitions()[1].GetGenerics()[0] != "T"`)
	}
	if _document.definitions[1].GetProperties()[2].GetType() != "T" {
		t.Fatal(`GetDefinitions()[1].GetProperties()[2].GetType() != "T"`)
	}
}
