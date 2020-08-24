# goapidoc

+ A tool written in golang for generating rest api document (swagger & apib)

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
            TermsOfService("http://xxx.yyy.zzz").
            License(NewLicense("MIT", "http://xxx.yyy.zzz")).
            Contact(NewContact("author", "http://xxx.yyy.zzz", "xxx@yyy.zzz")),
    )
    SetTags(
        NewTag("ping", "ping-controller"),
        NewTag("user", "user-controller"),
    )
    SetSecurities(
        NewSecurity("jwt", HEADER, "Authorization"),
    )
    
    AddPaths(
        NewRoutePath(GET, "/api/v1/ping", "ping").
            Desc("ping the server").
            Tags("ping").
            Consumes(JSON).
            Produces(JSON).
            Responses(
                NewResponse(200).Desc("success").Examples(map[string]string{JSON: "{\n    \"ping\": \"pong\"\n}"}),
            ),
        NewRoutePath(GET, "/api/v1/user", "get users").
            Tags("user").
            Consumes(JSON).
            Produces(JSON).
            Securities("jwt").
            Params(
                NewQueryParam("page", INTEGER, false, "current page").Default(1),
                NewQueryParam("total", INTEGER, false, "page size").Default(10),
                NewQueryParam("order", STRING, false, "order string").Default(""),
            ).
            Responses(
                NewResponse(200).Type("_Result<_Page<User>>"),
            ),
        NewRoutePath(PUT, "/api/v1/user/{id}", "update user (ugly api)").
            Tags("user").
            Consumes(JSON).
            Produces(JSON).
            Securities("jwt").
            Params(
                NewPathParam("id", INTEGER, true, "user id"),
                NewBodyParam("body", "User", true, "request body"),
            ).
            Responses(
                NewResponse(200).Type("Result").Desc("success"),
                NewResponse(404).Desc("not found").Headers(NewHeader("Content-Kind", STRING, "demo")),
                NewResponse(400).Type(STRING).Desc("bad request").Examples(map[string]string{JSON: "bad request"}),
            ),
    )
    
    AddDefinitions(
        NewDefinition("Result", "global response").Properties(
            NewProperty("code", INTEGER, true, "status code"),
            NewProperty("message", STRING, true, "status message"),
        ),
        NewDefinition("User", "user response").Properties(
            NewProperty("id", INTEGER, true, "user id"),
            NewProperty("name", STRING, true, "user name"),
            NewProperty("profile", STRING, false, "user profile").AllowEmpty(true),
            NewProperty("gender", STRING, true, "user gender").Enum("male", "female"),
            NewProperty("create_at", "string#date-time", true, "user register time"),
            NewProperty("birthday", "string#date", true, "user birthday"),
            NewProperty("scores", "number[]", true, "user scores"),
        ),
        NewDefinition("_Result", "global response").Generics("T").Properties(
            NewProperty("code", INTEGER, true, "status code"),
            NewProperty("message", STRING, true, "status message"),
            NewProperty("data", "T", true, "response data"),
        ),
        NewDefinition("_Page", "global page response").Generics("T").Properties(
            NewProperty("page", INTEGER, true, "current page"),
            NewProperty("total", INTEGER, true, "data count"),
            NewProperty("limit", INTEGER, true, "page size"),
            NewProperty("data", "T[]", true, "page data"),
        ),
    )
    
    _, _ = GenerateSwaggerYaml("./docs/api.yaml")
}
```

### References

+ [OpenAPI Specification 2.0](https://swagger.io/specification/v2/)
+ [API Blueprint Specification](https://apiblueprint.org/documentation/specification.html)
