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
			Desc("Multiple status values can be provided with comma separated strings.").
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
			Desc("Returns a single pet.").
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

		NewPostOperation("/pet/{petId}/uploadImage", "Uploads an image").
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

	AddOperations(
		NewPostOperation("/store/order", "Place an order for a pet").
			Tags("store").
			OperationId("placeOrder").
			Produces(XML, JSON).
			Params(
				NewBodyParam("body", "Order", true, "order placed for purchasing the pet"),
			).
			Responses(
				NewResponse(200, "Order").Desc("successful operation"),
				NewResponse(400, "").Desc("Invalid Order"),
			),

		NewGetOperation("/store/order/{orderId}", "Find purchase order by ID").
			Tags("store").
			Desc("For valid response try integer IDs with value >= 1 and <= 10.").
			OperationId("getOrderById").
			Produces(XML, JSON).
			Params(
				NewPathParam("orderId", "integer#int64", true, "ID of pet that needs to be fetched").ValueRange(1.0, 10.0),
			).
			Responses(
				NewResponse(200, "Order").Desc("successful operation"),
				NewResponse(400, "").Desc("Invalid ID supplied"),
				NewResponse(404, "").Desc("Order not found"),
			),

		NewDeleteOperation("/store/order/{orderId}", "Delete purchase order by ID").
			Tags("store").
			Desc("For valid response try integer IDs with positive integer value.").
			OperationId("deleteOrder").
			Produces(XML, JSON).
			Params(
				NewPathParam("orderId", "integer#int64", true, "ID of the order that needs to be deleted").Minimum(1.0),
			).
			Responses(
				NewResponse(400, "").Desc("Invalid ID supplied"),
				NewResponse(404, "").Desc("Order not found"),
			),
	)

	AddOperations(
		NewPostOperation("/user", "Create user").
			Tags("user").
			Desc("This can only be done by the logged in user.").
			OperationId("createUser").
			Produces(XML, JSON).
			Params(
				NewBodyParam("body", "User", true, "Created user object"),
			).
			Responses(
				NewResponse(200, "").Desc("successful operation"),
			),

		NewPostOperation("/user/createWithArray", "Creates list of users with given input array").
			Tags("user").
			OperationId("createUsersWithArrayInput").
			Produces(XML, JSON).
			Params(
				NewBodyParam("body", "User[]", true, "List of user object"),
			).
			Responses(
				NewResponse(200, "").Desc("successful operation"),
			),

		NewGetOperation("/user/login", "Logs user into the system").
			Tags("user").
			OperationId("loginUser").
			Produces(XML, JSON).
			Params(
				NewQueryParam("username", "string", true, "The user name for login"),
				NewQueryParam("password", "string", true, "The password for login in clear text"),
			).
			Responses(
				NewResponse(200, "string").Desc("successful operation").Headers(
					NewHeader("X-Rate-Limit", "integer#int32", "calls per hour allowed by the user"),
					NewHeader("X-Expires-After", "string#date-time", "date in UTC when token expires"),
				),
				NewResponse(400, "").Desc("Invalid username/password supplied"),
			),

		NewGetOperation("/user/logout", "Logs out current logged in user session").
			Tags("user").
			OperationId("logoutUser").
			Produces(XML, JSON).
			Responses(
				NewResponse(200, "").Desc("successful operation"),
			),

		NewGetOperation("/user/{username}", "Get user by user name").
			Tags("user").
			OperationId("getUserByName").
			Produces(XML, JSON).
			Params(
				NewPathParam("username", "string", true, "The name that needs to be fetched. Use user1 for testing."),
			).
			Responses(
				NewResponse(200, "User").Desc("successful operation"),
				NewResponse(400, "").Desc("Invalid username supplied"),
				NewResponse(404, "").Desc("User not found"),
			),

		NewPutOperation("/user/{username}", "Update user").
			Tags("user").
			Desc("This can only be done by the logged in user.").
			OperationId("updateUser").
			Produces(XML, JSON).
			Params(
				NewPathParam("username", "string", true, "name that need to be updated"),
				NewBodyParam("body", "User", true, "Updated user object"),
			).
			Responses(
				NewResponse(400, "").Desc("Invalid user supplied"),
				NewResponse(404, "").Desc("User not found"),
			),

		NewDeleteOperation("/user/{username}", "Delete user").
			Tags("user").
			Desc("This can only be done by the logged in user.").
			OperationId("deleteUser").
			Produces(XML, JSON).
			Params(
				NewPathParam("username", "string", true, "The name that needs to be deleted"),
			).
			Responses(
				NewResponse(400, "").Desc("Invalid username supplied"),
				NewResponse(404, "").Desc("User not found"),
			),
	)

	AddDefinitions(
		NewDefinition("Order", "").
			XMLRepr(NewXMLRepr("Order")).
			Properties(
				NewProperty("id", "integer#int64", false, ""),
				NewProperty("petId", "integer#int64", false, ""),
				NewProperty("quantity", "integer#int32", false, ""),
				NewProperty("shipDate", "string#date-time", false, ""),
				NewProperty("status", "string", false, "Order Status").Enum("placed", "approved", "delivered"),
				NewProperty("complete", "boolean", false, "").Default(false),
			),

		NewDefinition("Category", "").
			XMLRepr(NewXMLRepr("Category")).
			Properties(
				NewProperty("id", "integer#int64", false, ""),
				NewProperty("name", "string", false, ""),
			),

		NewDefinition("User", "").
			XMLRepr(NewXMLRepr("User")).
			Properties(
				NewProperty("id", "integer#int64", false, ""),
				NewProperty("username", "string", false, ""),
				NewProperty("firstName", "string", false, ""),
				NewProperty("lastName", "string", false, ""),
				NewProperty("email", "string", false, ""),
				NewProperty("password", "string", false, ""),
				NewProperty("phone", "string", false, ""),
				NewProperty("userStatus", "integer#int32", false, "User Status"),
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
	if _, err := SaveSwaggerYaml("./docs/api1.yaml"); err != nil {
		failNow(t, fmt.Sprintf("SaveSwaggerYaml error: %v", err))
	}
	if _, err := SaveSwaggerJson("./docs/api1.json"); err != nil {
		failNow(t, fmt.Sprintf("SaveSwaggerJson error: %v", err))
	}
	if _, err := SaveApib("./docs/api1.apib"); err != nil {
		failNow(t, fmt.Sprintf("SaveApib error: %v", err))
	}
}

func TestGenerate2(t *testing.T) {
	// https://github.com/apiaryio/api-blueprint/blob/master/examples/Gist%20Fox%20API%20%2B%20Auth.md
	// https://editor.docs.apiary.io/
	CleanupDocument()
	SetDocument("api.fake.com", "/",
		NewInfo("Gist Fox API", "Gist Fox API is a **pastes service** similar to [GitHub's Gist](http://gist.github.com).", "1.0.0"),
	)

	SetOption(NewOption().
		Tags(
			NewTag("Gist", "Gist-related resources of *Gist Fox API*."),
			NewTag("Access Authorization and Control", "Access and Control of *Gist Fox API* OAuth token."),
		).
		AdditionalDoc(`## Authentication
*Gist Fox API* uses OAuth Authorization. First you create a new (or acquire existing) OAuth token using Basic Authentication. After you have acquired your token you can use it to access other resources within token' scope.

## Media Types
Where applicable this API uses the [HAL+JSON](https://github.com/mikekelly/hal_specification/blob/master/hal_specification.md) media-type to represent resources states and affordances.

Requests with a message-body are using plain JSON to set or update resource states.

## Error States
The common [HTTP Response Status Codes](https://github.com/for-GET/know-your-http-well/blob/master/status-codes.md) are used.`).
		AddRoutesAlias("/", "Gist Fox API Root").
		AddRoutesAlias("/gists/{id}{?access_token}", "Gist").
		AddRoutesAlias("/gists{?access_token,since}", "Gists Collection").
		AddRoutesAlias("/gists/{id}/star{?access_token}", "Star").
		AddRoutesAlias("/authorization", "Authorization").
		AddRoutesAdditionalDoc("/", "This resource does not have any attributes. Instead it offers the initial API affordances in the form of the HTTP Link header and HAL links.").
		AddRoutesAdditionalDoc("/gists/{id}{?access_token}", `A single Gist object. The Gist resource is the central resource in the Gist Fox API. It represents one paste - a single text note.

The Gist resource has the following attributes:

+ id
+ created_at
+ description
+ content

The states *id* and *created_at* are assigned by the Gist Fox API at the moment of creation.`),
	)

	AddOperations(
		NewGetOperation("/", "Retrieve the Entry Point").
			Produces("application/hal+json").
			Responses(
				NewResponse(200, "").
					Headers(
						NewHeader("Link", "string", "").Example(`<http:/api.gistfox.com/>;rel="self",<http:/api.gistfox.com/gists>;rel="gists",<http:/api.gistfox.com/authorization>;rel="authorization"`),
					).
					AddExample("application/hal+json", nil),
			),
	)

	AddOperations(
		NewGetOperation("/gists/{id}", "Retrieve a Single Gist").
			Tags("Gist").
			Params(
				NewPathParam("id", "string", true, "ID of the Gist in the form of a hash."),
				NewQueryParam("access_token", "string", false, "Gist Fox API access token."),
			).
			Produces("application/hal+json").
			Responses(
				NewResponse(200, "").
					Headers(
						NewHeader("Link", "string", "").Example(`<http:/api.gistfox.com/gists/42>;rel="self", <http:/api.gistfox.com/gists/42/star>;rel="star"`),
					).
					AddExample("application/hal+json", nil),
			),

		NewPatchOperation("/gists/{id}", "Edit a Gist").
			Tags("Gist").
			Desc("To update a Gist send a JSON with updated value for one or more of the Gist resource attributes. All attributes values (states) from the previous version of this Gist are carried over by default if not included in the hash.").
			Consumes(JSON).
			Params(
				NewPathParam("id", "string", true, "ID of the Gist in the form of a hash."),
				NewQueryParam("access_token", "string", false, "Gist Fox API access token."),
				NewBodyParam("content", "string", true, "").Example("Updated file contents"), // <<<
			).
			Produces("application/hal+json").
			Responses(
				NewResponse(200, "").
					Headers(
						NewHeader("Link", "string", "").Example(`<http:/api.gistfox.com/gists/42>;rel="self", <http:/api.gistfox.com/gists/42/star>;rel="star"`),
					).
					AddExample("application/hal+json", nil),
			),

		NewDeleteOperation("/gists/{id}", "Delete a Gist").
			Tags("Gist").
			Params(
				NewPathParam("id", "string", true, "ID of the Gist in the form of a hash."),
				NewQueryParam("access_token", "string", false, "Gist Fox API access token."),
			).
			Responses(
				NewResponse(204, ""),
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
	if _, err := SaveSwaggerYaml("./docs/api2.yaml"); err != nil {
		failNow(t, fmt.Sprintf("SaveSwaggerYaml error: %v", err))
	}
	if _, err := SaveSwaggerJson("./docs/api2.json"); err != nil {
		failNow(t, fmt.Sprintf("SaveSwaggerJson error: %v", err))
	}
	if _, err := SaveApib("./docs/api2.apib"); err != nil {
		failNow(t, fmt.Sprintf("SaveApib error: %v", err))
	}
}

func TestGenerate3(t *testing.T) {

}
