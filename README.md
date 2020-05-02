# go-yaml-doc

+ A tool written in golang for generating yaml-format api document
+ If you want to adapt to swagger 2, you need to add `map[string]interface{}{"swagger": "2.0"}`

### Usage

```go
import . "github.com/Aoi-hosizora/go-yaml-doc"

SetDocument(
    "localhost:10086", "/",
    NewInfo("test-api", "a demo description", "1.0").
        SetTermsOfService("http://xxx.yyy.zzz").
        SetLicense(NewLicense("MIT", "http://xxx.yyy.zzz")).
        SetContact(NewContact("author", "http://xxx.yyy.zzz", "xxx@yyy.zzz")),
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
        SetDescription("ping the server").
        SetTags("ping").
        SetConsumes(JSON).
        SetProduces(JSON).
        SetResponses( // response with example value
            NewResponse(200).SetDescription("success").SetExamples(map[string]string{JSON: `{"ping":"pong"}`}),
        ),
    NewPath(GET, "/api/v1/user", "get user").
        SetDescription("get user from database").
        SetTags("user").
        SetConsumes(JSON).
        SetProduces(JSON).
        SetSecurities("jwt"). // a security setting
        SetParams(
            NewParam("page", QUERY, INTEGER, false, "current page").SetDefault(1), // parameter with default value
            NewParam("total", QUERY, INTEGER, false, "page size").SetDefault(10),
            NewParam("order", QUERY, STRING, false, "order string").SetDefault(""),
        ).
        SetResponses(
            NewResponse(200).SetSchema(NewSchemaRef("Result<Page<User>>")), // object type response
        ),
    NewPath(PUT, "/api/v1/user/{id}", "update user (ugly api)").
        SetDescription("update user to database").
        SetTags("user").
        SetConsumes(JSON).
        SetProduces(JSON).
        SetSecurities("jwt").
        SetParams(
            NewParam("id", PATH, INTEGER, true, "user id"), // normal parameter
            NewParam("body", BODY, OBJECT, true, "body").SetSchema(NewSchemaRef("User")), // object parameter
        ).
        SetResponses(
            NewResponse(200).SetDescription("success").SetSchema(NewSchemaRef("Result")).
                SetHeaders(NewHeader("Content-Type", STRING, "demo")), // response with header
            NewResponse(404).SetDescription("not found"),
            NewResponse(400).SetDescription("bad request").SetSchema(NewSchema(STRING, true)).
                SetExamples(map[string]string{JSON: "bad request"}), // string type response
        ),
    NewPath(HEAD, "/api/v1/test", "test path").
        SetParams(
            NewParam("arr", QUERY, ARRAY, true, "test").
                SetItems(NewItems(INTEGER).SetFormat(INT64).SetItems(NewItems(INTEGER))), // array-array parameter
            NewParam("ref", QUERY, ARRAY, true, "test").SetItems(NewItemsRef("User")), // array-object parameter
            NewParam("enum", QUERY, STRING, true, "test").SetEnum("male", "female"), // enum type parameter
        ),
)

AddDefinitions(
    NewDefinition("Result", "global response").SetProperties( // a normal definition
        NewProperty("code", INTEGER, true, "status code"), // a normal property
        NewProperty("message", STRING, true, "status message"),
    ),
    NewDefinition("User", "user response").SetProperties(
        NewProperty("id", INTEGER, true, "user id"),
        NewProperty("name", STRING, true, "user name"),
        NewProperty("profile", STRING, false, "user profile").SetAllowEmptyValue(true),
        NewProperty("gender", STRING, true, "user gender").SetEnum("male", "female"), // enum type property
        NewProperty("create_at", STRING, true, "user register time").SetFormat(DATETIME), // datetime type property
        NewProperty("birthday", STRING, true, "user birthday").SetFormat(DATE), // date type property
        NewProperty("scores", ARRAY, true, "user scores").SetItems(NewItems(NUMBER)), // array property
    ),
    NewDefinition("Page<User>", "user response").SetProperties(
        NewProperty("page", INTEGER, true, "current page"),
        NewProperty("total", INTEGER, true, "data count"),
        NewProperty("limit", INTEGER, true, "page size"),
        NewArrayProperty("data", NewItemsRef("User"), true, "page data"), // array-object property
    ),
    NewDefinition("Result<Page<User>>", "user response").SetProperties(
        NewProperty("code", INTEGER, true, "status code"),
        NewProperty("message", STRING, true, "status message"),
        NewObjectProperty("data", "Page<User>", true, "result data"), // object type property
    ),
)

err := GenerateYaml("api.yaml", map[string]interface{}{"swagger": "2.0"}) // add some fields to output
```

### References

+ [OpenAPI Specification 2.0](https://swagger.io/specification/v2/)
