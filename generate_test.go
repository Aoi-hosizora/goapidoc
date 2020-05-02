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
				NewResponse(200).SetDescription("success").SetExamples(map[string]string{"application/json": `{"ping": "pong"}`}),
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
				NewParam("body", "body", "object", true, "request body").SetSchema("User"),
			).
			SetResponses(
				NewResponse(200).SetSchema("Result"),
				NewResponse(404).SetDescription("not found"),
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
			NewProperty("data", "page data", "array", true).SetSchema("User"),
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

/*

=== RUN   TestGenerateYaml
swagger: "2.0"
host: localhost:10086
basePath: /
info:
  title: test-api
  description: a demo description
  version: "1.0"
  termsOfService: http://xxx.yyy.zzz
  license:
    name: MIT
    url: http://xxx.yyy.zzz
  contact:
    name: author
    url: http://xxx.yyy.zzz
    email: xxx@yyy.zzz
tags:
- name: ping
  description: ping-controller
- name: user
  description: user-controller
securityDefinitions:
  jwt:
    type: apiKey
    name: Authorization
    in: header
paths:
  /api/v1/ping:
    get:
      summary: ping
      operationId: -api-v1-ping-get
      description: ping the server
      tags:
      - ping
      consumes:
      - application/json
      produces:
      - application/json
      responses:
        "200":
          description: success
          examples:
            application/json: '{"ping": "pong"}'
  /api/v1/user:
    get:
      summary: get user
      operationId: -api-v1-user-get
      description: get user from database
      tags:
      - user
      consumes:
      - application/json
      produces:
      - application/json
      security:
      - jwt
      parameters:
      - name: page
        in: query
        type: integer
        required: false
        description: current page
        default: 1
      - name: total
        in: query
        type: integer
        required: false
        description: page size
        default: 10
      - name: order
        in: query
        type: string
        required: false
        description: order string
        default: ""
      responses:
        "200":
          schema:
            $ref: '#/definitions/Result<Page<User>>'
  /api/v1/user/{id}:
    put:
      summary: update user (ugly api)
      operationId: -api-v1-user-id-put
      description: update user to database
      tags:
      - user
      consumes:
      - application/json
      produces:
      - application/json
      security:
      - jwt
      parameters:
      - name: id
        in: path
        type: integer
        required: true
        description: user id
      - name: body
        in: body
        type: object
        required: true
        description: request body
        schema:
          $ref: '#/definitions/User'
      responses:
        "200":
          schema:
            $ref: '#/definitions/Result'
        "404":
          description: not found
definitions:
  Page<User>:
    title: Page<User>
    type: object
    required:
    - page
    - total
    - limit
    - data
    description: user response
    properties:
      data:
        type: array
        description: page data
        items:
          $ref: '#/definitions/User'
      limit:
        type: integer
        description: page size
      page:
        type: integer
        description: current page
      total:
        type: integer
        description: data count
  Result:
    title: Result
    type: object
    required:
    - code
    - message
    description: global response
    properties:
      code:
        type: integer
        description: status code
      message:
        type: string
        description: status message
  Result<Page<User>>:
    title: Result<Page<User>>
    type: object
    required:
    - code
    - message
    - data
    description: user response
    properties:
      code:
        type: integer
        description: status code
      data:
        type: object
        description: result data
        $ref: '#/definitions/Page<User>'
      message:
        type: string
        description: status message
  User:
    title: User
    type: object
    required:
    - id
    - name
    - gender
    - create_at
    - birthday
    description: user response
    properties:
      birthday:
        type: date
        description: user birthday
        format: yyyy-MM-dd
      create_at:
        type: datetime
        description: user register time
        format: yyyy-MM-dd HH:mm:ss
      gender:
        type: string
        description: user gender
        enum:
        - male
        - female
      id:
        type: integer
        description: user id
      name:
        type: string
        description: user name
      profile:
        type: string
        description: user profile

--- PASS: TestGenerateYaml (0.03s)
PASS


 */