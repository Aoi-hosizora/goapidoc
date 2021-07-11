package goapidoc

import (
	"fmt"
	"testing"
)

func TestGenerate1(t *testing.T) {
	// https://editor.swagger.io/

	t.Cleanup(func() { CleanupDocument() })
	SetDocument("petstore.swagger.io", "/v2",
		NewInfo("Swagger Petstore", "This is a sample server Petstore server.", "1.0.0").
			TermsOfService("http://swagger.io/terms/").
			Contact(NewContact("", "", "apiteam@swagger.io")).
			License(NewLicense("Apache 2.0", "http://www.apache.org/licenses/LICENSE-2.0.html")),
	)
	SetOption(NewOption().
		Schemes("https", "http").
		Tags(
			NewTag("pet", "Everything about your Pets"),
			NewTag("store", "Access to Petstore orders"),
			NewTag("user", "Operations about user"),
		).
		Securities(
			NewApiKeySecurity("api_key", HEADER, "api_key"),
			NewBasicSecurity("b"),
		),
	)
	AddOperations(
		NewPostOperation("/pet", "Add a new pet to the store").
			Tags("pet").
			OperationId("addPet").
			Consumes(JSON, XML).
			Produces(XML, JSON).
			Params(NewBodyParam("body", "Pet[]", true, "Pet object that needs to be added to the store")).
			Responses(NewResponse(405, "").Desc("Invalid input")).
			Securities("b"),
		NewPutOperation("/pet", "Update an existing pet").
			Tags("pet").
			OperationId("updatePet").
			Consumes(JSON, XML).
			Produces(XML, JSON).
			Params(NewBodyParam("body", "Pet[]", true, "Pet object that needs to be added to the store")).
			Responses(
				NewResponse(400, "string").Desc("Invalid ID supplied"),
				NewResponse(404, "string").Desc("Pet not found"),
				NewResponse(405, "string").Desc("Validation exception"),
			), // oauth2
		NewGetOperation("/pet/findByStatus", "Finds Pets by status").
			Tags("pet").
			Desc("Multiple status values can be provided with comma separated strings").
			OperationId("findPetsByStatus").
			Produces(JSON).
			Params(
				NewQueryParam("status", "string[]", true, "Status values that need to be considered for filter").CollectionFormat(MULTI).
					ItemOption(NewItemOption().Enum("available", "pending", "sold").Default("available")),
			).
			Responses(
				NewResponse(200, "Pet[]").Desc("successful operation").AddExample(JSON, "a"),
				NewResponse(400, "").Desc("Invalid status value"),
			), // oauth2
		NewGetOperation("/pet/findByTags", "Finds Pets by tags").
			Tags("pet").
			Desc("Multiple tags can be provided with comma separated strings.").
			OperationId("findPetsByTags").
			Produces(XML, JSON).
			Params(NewQueryParam("tags", "string[]", true, "Tags to filter by").CollectionFormat(MULTI)).
			Responses(
				NewResponse(200, "Pet[]").Desc("successful operation").AddExample(JSON, "a"),
				NewResponse(400, "").Desc("Invalid status value"),
			).
			Deprecated(true), // oauth2
	)
	AddDefinitions(
		NewDefinition("Pet", "").
			Properties(
				NewProperty("id", "integer#int64", false, ""),
				NewProperty("category", "Category", false, ""),
				NewProperty("name", "string", true, "").Example("doggie"),
				NewProperty("photoUrls", "string[]", true, "").ItemOption(NewItemOption().Pattern("^[123]*$")),
				NewProperty("tags", "Tag[]", false, ""),
				NewProperty("status", "string", false, "pet status in the store").Enum("available", "pending", "sold"),
			),
		NewDefinition("Category", "").
			Properties(
				NewProperty("id", "integer#int64", false, ""),
				NewProperty("name", "string", false, ""),
			),
		NewDefinition("Tag", "").
			Properties(
				NewProperty("id", "integer#int64", false, ""),
				NewProperty("name", "string", false, ""),
			),
	)

	if _, err := GenerateSwaggerYaml(); err != nil {
		failNow(t, fmt.Sprintf("GenerateSwaggerYaml error: %v", err))
	}
	if _, err := GenerateSwaggerJson(); err != nil {
		failNow(t, fmt.Sprintf("GenerateSwaggerJson error: %v", err))
	}
	if _, err := GenerateApib(); err != nil {
		failNow(t, fmt.Sprintf("GenerateApib error: %v", err))
	}
	if _, err := SaveSwaggerYaml("./docs/api.yaml"); err != nil {
		failNow(t, fmt.Sprintf("SaveSwaggerYaml error: %v", err))
	}
	if _, err := SaveSwaggerJson("./docs/api.json"); err != nil {
		failNow(t, fmt.Sprintf("SaveSwaggerJson error: %v", err))
	}
	if _, err := SaveApib("./docs/api.apib"); err != nil {
		failNow(t, fmt.Sprintf("SaveApib error: %v", err))
	}

	// if _document.definitions[1].GetGenerics()[0] != "T" {
	// 	t.Fatal(`GetDefinitions()[1].GetGenerics()[0] != "T"`)
	// }
	// if _document.definitions[1].GetProperties()[2].GetType() != "T" {
	// 	t.Fatal(`GetDefinitions()[1].GetProperties()[2].GetType() != "T"`)
	// }
}

func TestGenerate2(t *testing.T) {
	// https://github.com/apiaryio/api-blueprint/blob/master/examples/Gist%20Fox%20API%20%2B%20Auth.md
	// https://editor.docs.apiary.io/
}
