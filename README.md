# goapidoc

[![Build Status](https://travis-ci.com/Aoi-hosizora/goapidoc.svg?branch=master)](https://travis-ci.com/Aoi-hosizora/goapidoc)
[![codecov](https://codecov.io/gh/Aoi-hosizora/goapidoc/branch/master/graph/badge.svg)](https://codecov.io/gh/Aoi-hosizora/goapidoc)
[![Go Report Card](https://goreportcard.com/badge/github.com/Aoi-hosizora/goapidoc)](https://goreportcard.com/report/github.com/Aoi-hosizora/goapidoc)
[![License](http://img.shields.io/badge/license-mit-blue.svg)](./LICENSE)
[![Release](https://img.shields.io/github/v/release/Aoi-hosizora/goapidoc)](https://github.com/Aoi-hosizora/goapidoc/releases)

+ A golang library for generating api document, including swagger2 and apib.

### Function

+ [x] Support api, routes and definitions information
+ [x] Support generic definition type
+ [x] Support most of the functions for swagger 2
+ [x] Support basic functions for API Blueprint 1A

### Usage

+ Visit [generate_test.go](./generate_test.go) for detail examples, and [api1.json](./docs/api1.json) / [api2.apib](./docs/api2.apib) for generated documents.

```go
package main

import (
	. "github.com/Aoi-hosizora/goapidoc"
)

func main() {
	SetDocument("localhost:60001", "/",
		NewInfo("Demo api", "This is a demo api only for testing goapidoc.", "1.0.0").
			License(NewLicense("MIT", "")).
			Contact(NewContact("", "https://github.com/Aoi-hosizora", "")),
	)

	SetOption(NewOption().
		Schemes("http").
		Tags(
			NewTag("Authorization", "auth-controller"),
			NewTag("User", "user-controller"),
		).
		Securities(
			NewApiKeySecurity("jwt", HEADER, "Authorization"),
		).
		GlobalParams(
			NewQueryParam("force_refresh", "boolean", false, "force refresh flag").Default(false),
			NewHeaderParam("X-Special-Flag", "string", false, "a special flag in header").Example("token-xxx"),
		),
	)

	AddOperations(
		NewPostOperation("/auth/register", "Sign up").
			Tags("Authorization").
			Params(
				NewBodyParam("param", "RegisterParam", true, "register param"),
			).
			Responses(
				NewResponse(200, "Result"),
			),

		NewPostOperation("/auth/login", "Sign in").
			Tags("Authorization").
			Params(
				NewBodyParam("param", "LoginParam", true, "login param"),
			).
			Responses(
				NewResponse(200, "_Result<LoginDto>"),
			),

		NewGetOperation("/auth/me", "Get the authorized user").
			Tags("Authorization").
			Securities("jwt").
			Responses(
				NewResponse(200, "_Result<UserDto>"),
			),

		NewDeleteOperation("/auth/logout", "Sign out").
			Tags("Authorization").
			Securities("jwt").
			Responses(
				NewResponse(200, "Result"),
			),
	)

	AddOperations(
		NewGetOperation("/user", "Query all users").
			Tags("User").
			Securities("jwt").
			Params(
				NewQueryParam("page", "integer#int32", false, "query page").Default(1),
				NewQueryParam("limit", "integer#int32", false, "page size").Default(20),
				NewQueryParam("force_refresh", "boolean", false, "force refresh flag for querying users").Default(false),
				NewHeaderParam("X-Special-Flag", "string", true, "a special flag in header, which must be set for querying users"),
			).
			Responses(
				NewResponse(200, "_Result<_Page<UserDto>>"),
			),

		NewGetOperation("/user/{id}", "Query the specific user").
			Tags("User").
			Securities("jwt").
			Params(
				NewPathParam("id", "integer#int64", true, "user id"),
			).
			Responses(
				NewResponse(200, "_Result<UserDto>"),
			),

		NewPutOperation("/user", "Update the authorized user").
			Tags("User").
			Securities("jwt").
			Params(
				NewBodyParam("param", "UpdateUserParam", true, "update user param"),
			).
			Responses(
				NewResponse(200, "Result"),
			),

		NewDeleteOperation("/user", "Delete the authorized user").
			Tags("User").
			Securities("jwt").
			Responses(
				NewResponse(200, "Result"),
			),
	)

	AddDefinitions(
		NewDefinition("Result", "Global response").
			Properties(
				NewProperty("code", "integer#int32", true, "status code"),
				NewProperty("message", "string", true, "status message"),
			),

		NewDefinition("_Result", "Global generic response").
			Generics("T").
			Properties(
				NewProperty("code", "integer#int32", true, "status code"),
				NewProperty("message", "string", true, "status message"),
				NewProperty("data", "T", true, "response data"),
			),

		NewDefinition("_Page", "Global generic page response").
			Generics("T").
			Properties(
				NewProperty("page", "integer#int32", true, "current page"),
				NewProperty("limit", "integer#int32", true, "page size"),
				NewProperty("total", "integer#int32", true, "total count"),
				NewProperty("data", "T[]", true, "response data"),
			),

		NewDefinition("LoginParam", "Login parameter").
			Properties(
				NewProperty("username", "string", true, "username"),
				NewProperty("password", "string", true, "password"),
			),

		NewDefinition("RegisterParam", "Register parameter").
			Properties(
				NewProperty("username", "string", true, "username"),
				NewProperty("password", "string", true, "password"),
			),

		NewDefinition("UpdateUserParam", "Update user parameter").
			Properties(
				NewProperty("username", "string", true, "username"),
				NewProperty("bio", "string", true, "user bio"),
				NewProperty("gender", "string", true, "user gender").Enum("Secret", "Male", "Female"),
				NewProperty("birthday", "string#date", true, "user birthday"),
			),

		NewDefinition("LoginDto", "Login response").
			Properties(
				NewProperty("user", "UserDto", true, "authorized user"),
				NewProperty("token", "string", true, "access token"),
			),

		NewDefinition("UserDto", "User response").
			Properties(
				NewProperty("id", "integer#int64", true, "user id"),
				NewProperty("username", "string", true, "username"),
				NewProperty("bio", "string", true, "user bio"),
				NewProperty("gender", "string", true, "user gender").Enum("Secret", "Male", "Female"),
				NewProperty("birthday", "string#date", true, "user birthday"),
			),
	)

	_, _ = SaveSwaggerYaml("./docs/api3.yaml")
	_, _ = SaveSwaggerJson("./docs/api3.json")
	_, _ = SaveApib("./docs/api3.apib")
}
```

### References

+ [OpenAPI Specification 2.0](https://swagger.io/specification/v2/)
+ [API Blueprint Specification](https://apiblueprint.org/documentation/specification.html)
