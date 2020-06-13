# goapidoc

+ A tool written in golang for generating rest api document

### Function

+ [x] Basic swagger information
+ [x] Support definitions
+ [x] Support basic generic type

### Type tips

```text
Struct1<Struct2<Struct3, integer#int64>, string#date-time, Struct4[]>[]

$basicType   := {integer, number, string, boolean, file, object, array}
$basicFormat := {int32, int64, float, double, byte, binary, data, data-time, password}

$type    := $basicType
$type    := $basicType#$basicFormat
$type    := $type[]
$type    := $type<$generic>
$generic := $type
$generic := $generic, $type
```

### Usage

```go
package main

import (
    . "github.com/Aoi-hosizora/goapidoc"
)

func main() {
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
                NewQueryParam("page", INTEGER, false, "current page").WithDefault(1),
                NewQueryParam("total", INTEGER, false, "page size").WithDefault(10),
                NewQueryParam("order", STRING, false, "order string").WithDefault(""),
            ).
            WithResponses(
                NewResponse(200).WithType("_Result<_Page<User>>"),
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
    )
    
    _, _ = GenerateYamlWithSwagger2("./docs/api.yaml")
}
```

### References

+ [OpenAPI Specification 2.0](https://swagger.io/specification/v2/)
