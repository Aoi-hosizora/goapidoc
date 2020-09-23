# goapidoc

[![Build Status](https://www.travis-ci.org/Aoi-hosizora/goapidoc.svg?branch=master)](https://www.travis-ci.org/Aoi-hosizora/goapidoc)

+ A tool written in golang for generating api document (swagger2 & apib)

### Function

+ [x] Support basic information, route and definition
+ [x] Support basic generic type
+ [x] Support basic Swagger 2
+ [x] Support basic API Blueprint 1A (compatible with swagger2)

### Type tips

```text
Struct1<Struct2<Struct3, integer#int64>, string#date-time, Struct4[]>[]

$basicType   := {integer, number, string, boolean, file, object, array}
$basicFormat := {int32, int64, float, double, byte, binary, date, date-time, password}

$type    := $basicType
$type    := $basicType#$basicFormat
$type    := $type[]
$type    := $type<$generic>
$generic := $type
$generic := $generic, $type
```

### Usage

+ See [generate_test,go](./generate_test.go)

```go
package main

import (
    . "github.com/Aoi-hosizora/goapidoc"
)

func main() {
	SetDocument(
		"localhost:65530", "/",
		NewInfo("goapidoc", "goapidoc test api", "1.0").
			TermsOfService("https://xxx.yyy.zzz.com/").
			License(NewLicense("MIT", "https://xxx.yyy.zzz.com/")).
			Contact(NewContact("xxx", "https://xxx.yyy.zzz.com/", "xxx@yyy.zzz")),
	)
	SetTags(
		NewTag("Authorization", "auth-controller"),
		NewTag("User", "user-controller"),
		NewTag("Test", "test-controller"),
	)
	SetSecurities(
		NewSecurity("Jwt", "header", "Authorization"),
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

	_, _ = GenerateSwaggerYaml("./docs/api.yaml")

	_, _ = GenerateSwaggerJson("./docs/api.json")

	_, _ = GenerateApib("./docs/api.apib")
}
```

### References

+ [OpenAPI Specification 2.0](https://swagger.io/specification/v2/)
+ [API Blueprint Specification](https://apiblueprint.org/documentation/specification.html)
