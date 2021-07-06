package goapidoc

import (
	"log"
	"testing"
)

func TestGetSet(t *testing.T) {
	for _, tc := range []struct {
		name   string
		giveFn func() interface{}
		testFn func(interface{}) bool
	}{
		{"default document",
			func() interface{} { return NewDocument("", "", nil) },
			func(i interface{}) bool { return i.(*Document).GetHost() == "" && i.(*Document).GetBasePath() == "" && i.(*Document).GetInfo() == nil }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			itf := tc.giveFn()
			if !tc.testFn(itf) {
				failNow(t, "Not matched object")
			}
		})
	}
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
				NewQueryParam("q1", "string#date-time", true, "q1").Enums(0, 1, 2),
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
				NewProperty("gender", "string", true, "user gender").Enums("secret", "male", "female"),
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
				NewProperty("gender", "string", true, "user gender").Enums("secret", "male", "female"),
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
