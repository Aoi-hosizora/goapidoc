package goapidoc

import (
	"fmt"
	"testing"
)

func TestGenerate1(t *testing.T) {
	// https://editor.swagger.io/
	CleanupDocument()
	SetDocument("petstore.swagger.io", "/v2",
		NewInfo("Swagger Petstore", "This is a sample server Petstore server.", "1.0.0").
			TermsOfService("http://swagger.io/terms/").
			Contact(NewContact("", "", "apiteam@swagger.io")).
			License(NewLicense("Apache 2.0", "http://www.apache.org/licenses/LICENSE-2.0.html")),
	)
	SetOption(NewOption().
		Schemes("https", "http").
		Tags(
			NewTag("pet", "Everything about your Pets").
				ExternalDocs(NewExternalDocs("Find out more", "http://swagger.io")),
			NewTag("store", "Access to Petstore orders"),
			NewTag("user", "Operations about user").
				ExternalDocs(NewExternalDocs("Find out more about our store", "http://swagger.io")),
		).
		Securities(
			NewOAuth2Security("petstore_auth", IMPLICIT_FLOW).
				AuthorizationUrl("http://petstore.swagger.io/oauth/dialog").
				AddScope("write:pets", "modify pets in your account").
				AddScope("read:pets", "read your pets"),
			NewApiKeySecurity("api_key", HEADER, "api_key"),
		).
		ExternalDocs(NewExternalDocs("Find out more about Swagger", "http://swagger.io")),
	)
	AddOperations(
		NewPostOperation("/pet", "Add a new pet to the store").
			Tags("pet").
			OperationId("addPet").
			Consumes(JSON, XML).
			Produces(XML, JSON).
			Params(
				NewBodyParam("body", "Pet", true, "Pet object that needs to be added to the store"),
			).
			Responses(
				NewResponse(405, "").Desc("Invalid input"),
			).
			Securities("petstore_auth").
			AddSecurityScopes("petstore_auth", "write:pets", "read:pets"),

		NewPutOperation("/pet", "Update an existing pet").
			Tags("pet").
			OperationId("updatePet").
			Consumes(JSON, XML).
			Produces(XML, JSON).
			Params(
				NewBodyParam("body", "Pet", true, "Pet object that needs to be added to the store"),
			).
			Responses(
				NewResponse(400, "").Desc("Invalid ID supplied"),
				NewResponse(404, "").Desc("Pet not found"),
				NewResponse(405, "").Desc("Validation exception"),
			).
			Securities("petstore_auth").
			AddSecurityScopes("petstore_auth", "write:pets", "read:pets"),

		NewGetOperation("/pet/findByStatus", "Finds Pets by status").
			Tags("pet").
			Desc("Multiple status values can be provided with comma separated strings").
			OperationId("findPetsByStatus").
			Produces(XML, JSON).
			Params(
				NewQueryParam("status", "string[]", true, "Status values that need to be considered for filter").CollectionFormat(MULTI).
					ItemOption(NewItemOption().Enum("available", "pending", "sold").Default("available")),
			).
			Responses(
				NewResponse(200, "Pet[]").Desc("successful operation"),
				NewResponse(400, "").Desc("Invalid status value"),
			).
			Securities("petstore_auth").
			AddSecurityScopes("petstore_auth", "write:pets", "read:pets"),

		NewGetOperation("/pet/findByTags", "Finds Pets by tags").
			Tags("pet").
			Desc("Multiple tags can be provided with comma separated strings.").
			OperationId("findPetsByTags").
			Produces(XML, JSON).
			Params(
				NewQueryParam("tags", "string[]", true, "Tags to filter by").CollectionFormat(MULTI),
			).
			Responses(
				NewResponse(200, "Pet[]").Desc("successful operation"),
				NewResponse(400, "").Desc("Invalid tag value"),
			).
			Securities("petstore_auth").
			AddSecurityScopes("petstore_auth", "write:pets", "read:pets").
			Deprecated(true),

		NewGetOperation("/pet/{petId}", "Find pet by ID").
			Tags("pet").
			Desc("Returns a single pet").
			OperationId("getPetById").
			Produces(XML, JSON).
			Params(
				NewPathParam("petId", "integer#int64", true, "ID of pet to return"),
			).
			Responses(
				NewResponse(200, "Pet").Desc("successful operation"),
				NewResponse(400, "").Desc("Invalid ID supplied"),
				NewResponse(404, "").Desc("Pet not found"),
			).
			Securities("api_key"),

		NewPostOperation("/pet/{petId}", "Updates a pet in the store with form data").
			Tags("pet").
			OperationId("updatePetWithForm").
			Consumes(URL).
			Produces(XML, JSON).
			Params(
				NewPathParam("petId", "integer#int64", true, "ID of pet that needs to be updated"),
				NewFormParam("name", "string", false, "Updated name of the pet"),
				NewFormParam("status", "string", false, "Updated status of the pet"),
			).
			Responses(
				NewResponse(405, "").Desc("Invalid input"),
			).
			Securities("petstore_auth").
			AddSecurityScopes("petstore_auth", "write:pets", "read:pets"),

		NewDeleteOperation("/pet/{petId}", "Deletes a pet").
			Tags("pet").
			OperationId("deletePet").
			Produces(XML, JSON).
			Params(
				NewHeaderParam("api_key", "string", false, ""),
				NewPathParam("petId", "integer#int64", true, "Pet id to delete"),
			).
			Responses(
				NewResponse(400, "").Desc("Invalid ID supplied"),
				NewResponse(404, "").Desc("Pet not found"),
			).
			Securities("petstore_auth").
			AddSecurityScopes("petstore_auth", "write:pets", "read:pets"),

		NewPostOperation("/pet/{petId}/uploadImage", "uploads an image").
			Tags("pet").
			OperationId("uploadFile").
			Consumes(MPFD).
			Produces(JSON).
			Params(
				NewPathParam("petId", "integer#int64", true, "ID of pet to update"),
				NewFormParam("additionalMetadata", "string", false, "Additional data to pass to server"),
				NewFormParam("file", "file", false, "file to upload"),
			).
			Responses(
				NewResponse(200, "ApiResponse").Desc("successful operation"),
			).
			Securities("petstore_auth").
			AddSecurityScopes("petstore_auth", "write:pets", "read:pets"),
	)
	AddDefinitions(
		NewDefinition("Category", "").
			XMLRepr(NewXMLRepr("Category")).
			Properties(
				NewProperty("id", "integer#int64", false, ""),
				NewProperty("name", "string", false, ""),
			),

		NewDefinition("Tag", "").
			XMLRepr(NewXMLRepr("Tag")).
			Properties(
				NewProperty("id", "integer#int64", false, ""),
				NewProperty("name", "string", false, ""),
			),

		NewDefinition("Pet", "").
			XMLRepr(NewXMLRepr("Pet")).
			Properties(
				NewProperty("id", "integer#int64", false, ""),
				NewProperty("category", "Category", false, ""),
				NewProperty("name", "string", true, "").Example("doggie"),
				NewProperty("photoUrls", "string[]", true, "").ItemOption(NewItemOption().Pattern("^[123]*$")),
				NewProperty("tags", "Tag[]", false, ""),
				NewProperty("status", "string", false, "pet status in the store").Enum("available", "pending", "sold"),
			),

		NewDefinition("ApiResponse", "").
			Properties(
				NewProperty("code", "integer#int32", false, ""),
				NewProperty("type", "string", false, ""),
				NewProperty("message", "string", false, ""),
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
