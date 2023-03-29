package goapidoc

import (
	"errors"
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
				ExternalDoc(NewExternalDoc("Find out more", "http://swagger.io")),
			NewTag("store", "Access to Petstore orders"),
			NewTag("user", "Operations about user").
				ExternalDoc(NewExternalDoc("Find out more about our store", "http://swagger.io")),
		).
		Securities(
			NewOAuth2Security("petstore_auth", IMPLICIT_FLOW).
				AuthorizationUrl("http://petstore.swagger.io/oauth/dialog").
				Scopes(
					NewSecurityScope("write:pets", "modify pets in your account"),
					NewSecurityScope("read:pets", "read your pets"),
				),
			NewApiKeySecurity("api_key", HEADER, "api_key"),
			NewBasicSecurity("b").Desc("A demo basic security definition"),
		).
		ExternalDoc(NewExternalDoc("Find out more about Swagger", "http://swagger.io")),
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
			SetSecurityScopes("petstore_auth", "write:pets", "read:pets"),

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
			SetSecurityScopes("petstore_auth", "write:pets", "read:pets"),

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
			SetSecurityScopes("petstore_auth", "write:pets", "read:pets"),

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
			SetSecurityScopes("petstore_auth", "write:pets", "read:pets").
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
			Securities("api_key", "b"),

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
			SetSecurityScopes("petstore_auth", "write:pets", "read:pets"),

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
			SetSecurityScopes("petstore_auth", "write:pets", "read:pets"),

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
			SetSecurityScopes("petstore_auth", "write:pets", "read:pets"),
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
			).
			Securities("b"),

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
			).
			Securities("b"),

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
			).
			Securities("b"),
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
					NewResponseHeader("X-Rate-Limit", "integer#int32", "calls per hour allowed by the user"),
					NewResponseHeader("X-Expires-After", "string#date-time", "date in UTC when token expires"),
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

	_generate(t, "api1")
}

func TestGenerate2(t *testing.T) {
	// https://raw.githubusercontent.com/apiaryio/api-blueprint/master/examples/Gist%20Fox%20API%20%2B%20Auth.md
	// https://editor.docs.apiary.io/
	CleanupDocument()
	SetDocument("api.gistfox.com", "/",
		NewInfo("Gist Fox API", "Gist Fox API is a **pastes service** similar to [GitHub's Gist](http://gist.github.com).", "1.0.0"),
	)

	SetOption(NewOption().
		Tags(
			NewTag("Gist", "Gist-related resources of *Gist Fox API*."),
			NewTag("Access Authorization and Control", "Access and Control of *Gist Fox API* OAuth token."),
			NewTag("Test more functions", "Operations in this group is only used for testing.").AdditionalDoc("## Example\n\nGET /test is used to test more parameters style in api blueprint."),
		).
		AdditionalDoc("## Authentication\n*Gist Fox API* uses OAuth Authorization. First you create a new (or acquire existing) OAuth token using Basic Authentication. After you have acquired your token you can use it to access other resources within token' scope.\n\n## Media Types\nWhere applicable this API uses the [HAL+JSON](https://github.com/mikekelly/hal_specification/blob/master/hal_specification.md) media-type to represent resources states and affordances.\n\nRequests with a message-body are using plain JSON to set or update resource states.\n\n## Error States\nThe common [HTTP Response Status Codes](https://github.com/for-GET/know-your-http-well/blob/master/status-codes.md) are used.").
		RoutesOptions(
			NewRoutesOption("/").
				Summary("Gist Fox API Root").
				AdditionalDoc("Gist Fox API entry point.\n\nThis resource does not have any attributes. Instead it offers the initial API affordances in the form of the HTTP Link header and\nHAL links."),
			NewRoutesOption("/gists/{id}{?access_token}").
				Summary("Gist").
				AdditionalDoc("A single Gist object. The Gist resource is the central resource in the Gist Fox API. It represents one paste - a single text note.\n\nThe Gist resource has the following attributes:\n\n+ id\n+ created_at\n+ description\n+ content\n\nThe states *id* and *created_at* are assigned by the Gist Fox API at the moment of creation."),
			NewRoutesOption("/gists{?since,access_token}").
				Summary("Gists Collection").
				AdditionalDoc("Collection of all Gists.\n\nThe Gist Collection resource has the following attribute:\n\n+ total\n\nIn addition it **embeds** *Gist Resources* in the Gist Fox API."),
			NewRoutesOption("/gists/{id}/star{?access_token}").
				Summary("Star").
				AdditionalDoc("Star resource represents a Gist starred status.\n\nThe Star resource has the following attribute:\n\n+ starred"),
			NewRoutesOption("/authorization").
				Summary("Authorization").
				AdditionalDoc("Authorization Resource represents an authorization granted to the user. You can **only** access your own authorization, and only through **Basic Authentication**.\n\nThe Authorization Resource has the following attribute:\n\n+ token\n+ scopes\n\nWhere *token* represents an OAuth token and *scopes* is an array of scopes granted for the given authorization. At this moment the only available scope is `gist_write`."),
		),
	)

	AddOperations(
		NewGetOperation("/", "Retrieve the Entry Point").
			Produces("application/hal+json").
			Responses(
				NewResponse(200, "").
					Headers(
						NewResponseHeader("Link", "string", "").Example(`<http:/api.gistfox.com/>;rel="self",<http:/api.gistfox.com/gists>;rel="gists",<http:/api.gistfox.com/authorization>;rel="authorization"`),
					).
					Examples(
						NewResponseExample("application/hal+json", "{\n    \"_links\": {\n        \"self\": { \"href\": \"/\" },\n        \"gists\": { \"href\": \"/gists?{since}\", \"templated\": true },\n        \"authorization\": { \"href\": \"/authorization\"}\n    }\n}"),
					),
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
						NewResponseHeader("Link", "string", "").Example(`<http:/api.gistfox.com/gists/42>;rel="self", <http:/api.gistfox.com/gists/42/star>;rel="star"`),
					).
					Examples(
						NewResponseExample("application/hal+json", "{\n    \"_links\": {\n        \"self\": { \"href\": \"/gists/42\" },\n        \"star\": { \"href\": \"/gists/42/star\" },\n    },\n    \"id\": \"42\",\n    \"created_at\": \"2014-04-14T02:15:15Z\",\n    \"description\": \"Description of Gist\",\n    \"content\": \"String contents\"\n}"),
					).
					AdditionalDoc("HAL+JSON representation of Gist Resource. In addition to representing its state in the JSON form it offers affordances in the form of the HTTP Link header and HAL links."),
			),

		NewPatchOperation("/gists/{id}", "Edit a Gist").
			Tags("Gist").
			Desc("To update a Gist send a JSON with updated value for one or more of the Gist resource attributes. All attributes values (states) from the previous version of this Gist are carried over by default if not included in the hash.").
			Consumes(JSON).
			Params(
				NewPathParam("id", "string", true, "ID of the Gist in the form of a hash."),
				NewQueryParam("access_token", "string", false, "Gist Fox API access token."),
			).
			RequestExample("{\n    \"content\": \"Updated file contents\"\n}").
			Produces("application/hal+json").
			Responses(
				NewResponse(200, "").
					Headers(
						NewResponseHeader("Link", "string", "").Example(`<http:/api.gistfox.com/gists/42>;rel="self", <http:/api.gistfox.com/gists/42/star>;rel="star"`),
					).
					Examples(
						NewResponseExample("application/hal+json", "{\n    \"_links\": {\n        \"self\": { \"href\": \"/gists/42\" },\n        \"star\": { \"href\": \"/gists/42/star\" },\n    },\n    \"id\": \"42\",\n    \"created_at\": \"2014-04-14T02:15:15Z\",\n    \"description\": \"Description of Gist\",\n    \"content\": \"String contents\"\n}"),
					).
					AdditionalDoc("HAL+JSON representation of Gist Resource. In addition to representing its state in the JSON form it offers affordances in the form of the HTTP Link header and HAL links."),
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

	AddOperations(
		NewGetOperation("/gists", "List All Gists").
			Tags("Gist").
			Params(
				NewQueryParam("since", "string", false, "Timestamp in ISO 8601 format: `YYYY-MM-DDTHH:MM:SSZ` Only gists updated at or after this time are returned."),
				NewQueryParam("access_token", "string", false, "Gist Fox API access token."),
			).
			Produces("application/hal+json").
			Responses(
				NewResponse(200, "").
					Headers(
						NewResponseHeader("Link", "string", "").Example(" <http:/api.gistfox.com/gists>;rel=\"self\""),
					).
					Examples(
						NewResponseExample("application/hal+json", "{\n    \"_links\": {\n        \"self\": { \"href\": \"/gists\" }\n    },\n    \"_embedded\": {\n        \"gists\": [\n            {\n                \"_links\" : {\n                    \"self\": { \"href\": \"/gists/42\" }\n                },\n                \"id\": \"42\",\n                \"created_at\": \"2014-04-14T02:15:15Z\",\n                \"description\": \"Description of Gist\"\n            }\n        ]\n    },\n    \"total\": 1\n}"),
					).
					AdditionalDoc("HAL+JSON representation of Gist Collection Resource. The Gist resources in collections are embedded. Note the embedded Gists resource are incomplete representations of the Gist in question. Use the respective Gist link to retrieve its full representation."),
			),

		NewPostOperation("/gists", "Create a Gist").
			Tags("Gist").
			Desc("To create a new Gist simply provide a JSON hash of the *description* and *content* attributes for the new Gist.\n\nThis action requires an `access_token` with `gist_write` scope.").
			Params(
				NewQueryParam("since", "string", false, "Timestamp in ISO 8601 format: `YYYY-MM-DDTHH:MM:SSZ` Only gists updated at or after this time are returned."),
				NewQueryParam("access_token", "string", false, "Gist Fox API access token."),
			).
			Consumes(JSON).
			RequestExample("{\n    \"description\": \"Description of Gist\",\n    \"content\": \"String content\"\n}").
			Produces("application/hal+json").
			Responses(
				NewResponse(201, "").
					Headers(
						NewResponseHeader("Link", "string", "").Example(`<http:/api.gistfox.com/gists/42>;rel="self", <http:/api.gistfox.com/gists/42/star>;rel="star"`),
					).
					Examples(
						NewResponseExample("application/hal+json", "{\n    \"_links\": {\n        \"self\": { \"href\": \"/gists/42\" },\n        \"star\": { \"href\": \"/gists/42/star\" },\n    },\n    \"id\": \"42\",\n    \"created_at\": \"2014-04-14T02:15:15Z\",\n    \"description\": \"Description of Gist\",\n    \"content\": \"String contents\"\n}"),
					).
					AdditionalDoc("HAL+JSON representation of Gist Resource. In addition to representing its state in the JSON form it offers affordances in the form of the HTTP Link header and HAL links."),
			),
	)

	AddOperations(
		NewPutOperation("/gists/{id}/star", "Star a Gist").
			Tags("Gist").
			Desc("This action requires an `access_token` with `gist_write` scope.").
			Params(
				NewPathParam("id", "string", true, "ID of the gist in the form of a hash"),
				NewQueryParam("access_token", "string", false, "Gist Fox API access token"),
			).
			Produces("application/hal+json").
			Responses(
				NewResponse(204, ""),
			),

		NewDeleteOperation("/gists/{id}/star", "Unstar a Gist").
			Tags("Gist").
			Desc("This action requires an `access_token` with `gist_write` scope.").
			Params(
				NewPathParam("id", "string", true, "ID of the gist in the form of a hash"),
				NewQueryParam("access_token", "string", false, "Gist Fox API access token"),
			).
			Produces("application/hal+json").
			Responses(
				NewResponse(204, ""),
			),

		NewGetOperation("/gists/{id}/star", "Check if a Gist is Starred").
			Tags("Gist").
			Params(
				NewPathParam("id", "string", true, "ID of the gist in the form of a hash"),
				NewQueryParam("access_token", "string", false, "Gist Fox API access token"),
			).
			Produces("application/hal+json").
			Responses(
				NewResponse(200, "").
					Headers(
						NewResponseHeader("Link", "string", "").Example("<http:/api.gistfox.com/gists/42/star>;rel=\"self\""),
					).
					Examples(
						NewResponseExample("application/hal+json", "{\n    \"_links\": {\n        \"self\": { \"href\": \"/gists/42/star\" },\n    },\n    \"starred\": true\n}"),
					).
					AdditionalDoc("HAL+JSON representation of Star Resource."),
			),
	)

	AddOperations(
		NewGetOperation("/authorization", "Retrieve Authorization").
			Tags("Access Authorization and Control").
			Params(
				NewHeaderParam("Authorization", "string", true, "").Example("Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ=="),
			).
			Produces("application/hal+json").
			Responses(
				NewResponse(200, "").
					Headers(
						NewResponseHeader("Link", "string", "").Example("<http:/api.gistfox.com/authorizations/1>;rel=\"self\""),
					).
					Examples(
						NewResponseExample("application/hal+json", "{\n    \"_links\": {\n        \"self\": { \"href\": \"/authorizations\" },\n    },\n    \"scopes\": [\n        \"gist_write\"\n    ],\n    \"token\": \"abc123\"\n}"),
					),
			),

		NewPostOperation("/authorization", "Create Authorization").
			Tags("Access Authorization and Control").
			Params(
				NewHeaderParam("Authorization", "string", true, "").Example("Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ=="),
			).
			RequestExample("{\n    \"scopes\": [\n        \"gist_write\"\n    ]\n}").
			Produces("application/hal+json").
			Responses(
				NewResponse(201, "").
					Headers(
						NewResponseHeader("Link", "string", "").Example("<http:/api.gistfox.com/authorizations/1>;rel=\"self\""),
					).
					Examples(
						NewResponseExample("application/hal+json", "{\n    \"_links\": {\n        \"self\": { \"href\": \"/authorizations\" },\n    },\n    \"scopes\": [\n        \"gist_write\"\n    ],\n    \"token\": \"abc123\"\n}"),
					),
			),

		NewDeleteOperation("/authorization", "Remove an Authorization").
			Tags("Access Authorization and Control").
			Params(
				NewHeaderParam("Authorization", "string", true, "").Example("Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ=="),
			).
			Responses(
				NewResponse(204, ""),
			),
	)

	// only for test
	AddOperations(
		NewPostOperation("/test", "Test the most difficult operation").
			Tags("Test more functions").
			ExternalDoc(NewExternalDoc("", "https://apiblueprint.org/documentation/specification.html")).
			RequestExample(struct{ A int }{1}).
			Consumes(MPFD).
			Params(
				NewFormParam("query1", "string#password[]", false, "").AllowEmpty(true).Pattern("^.+$").UniqueItems(true).
					CollectionFormat("csv").ValueRange(0, 10).ItemsRange(0, 5).LengthRange(0, 8),
				NewFormParam("query2", "number", true, "").Minimum(0).MinItems(0).MinLength(0).ExclusiveMin(true).MultipleOf(3.3),
				NewFormParam("query3", "number[][]", true, "").Maximum(10).MaxItems(5).MaxLength(8).ExclusiveMax(true).
					ItemOption(NewItemOption().ValueRange(-5, 5).ExclusiveMax(true).ExclusiveMin(true)),
				NewFormParam("query4", "integer#int64", false, "some desc").ValueRange(0, 10).MaxLength(8).ExclusiveMax(true).ExclusiveMin(true),
			).
			Responses(
				NewResponse(200, ""),
			),
	)

	_generate(t, "api2")
}

func TestGenerate3(t *testing.T) {
	CleanupDocument()
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
			NewHeaderParam("X-Special-Flag", "string", false, "a special flag in header"),
		).
		ParamTemplates(
			NewParamTemplate("page_limit_query").
				Params(
					NewQueryParam("page", "integer#int32", false, "query page").Default(1),
					NewQueryParam("limit", "integer#int32", false, "page size").Default(20),
				),
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
				NewQueryParam("force_refresh", "boolean", false, "force refresh flag for querying users").Default(false),
				NewHeaderParam("X-Special-Flag", "string", true, "a special flag in header, which must be set for querying users"),
			).
			ParamTemplateNames("page_limit_query").
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

	_generate(t, "api3")
}

func _generate(t *testing.T, name string) {
	if _, err := GenerateSwaggerYaml(); err != nil {
		failNow(t, fmt.Sprintf("GenerateSwaggerYaml (%s) error: %v", name, err))
	}
	if _, err := GenerateSwaggerJson(); err != nil {
		failNow(t, fmt.Sprintf("GenerateSwaggerJson (%s) error: %v", name, err))
	}
	EnableWarningLogger()
	if _, err := GenerateApib(); err != nil {
		failNow(t, fmt.Sprintf("GenerateApib (%s) error: %v", name, err))
	}
	DisableWarningLogger()
	if _, err := SaveSwaggerYaml("./docs/" + name + ".yaml"); err != nil {
		failNow(t, fmt.Sprintf("SaveSwaggerYaml (%s) error: %v", name, err))
	}
	if _, err := SaveSwaggerJson("./docs/" + name + ".json"); err != nil {
		failNow(t, fmt.Sprintf("SaveSwaggerJson (%s) error: %v", name, err))
	}
	if _, err := SaveApib("./docs/" + name + ".apib"); err != nil {
		failNow(t, fmt.Sprintf("SaveApib (%s) error: %v", name, err))
	}
}

func TestCheckDocument(t *testing.T) {
	demoInfo := func() *Info { return NewInfo("title", "", "1.0.0") }
	demoOp := func() *Operation { return NewOperation("get", "/", "summary").Responses(NewResponse(200, "")) }
	demoDoc := func() *Document { return NewDocument("host", "/", NewInfo("title", "", "1.0")) }
	demoResp := func() *Response { return NewResponse(200, "") }

	for _, tc := range []struct {
		name    string
		giveDoc *Document
	}{
		{"doc.host", NewDocument("", "", nil).Operations(demoOp())},
		{"doc.basePath", NewDocument("host", "", nil).Operations(demoOp())},
		{"doc.info", NewDocument("host", "/", nil).Operations(demoOp())},
		{"doc.info.title", NewDocument("host", "/", NewInfo("", "", "")).Operations(demoOp())},
		{"doc.info.version", NewDocument("host", "/", NewInfo("title", "", "")).Operations(demoOp())},
		{"success", NewDocument("host", "/", NewInfo("title", "", "v1.0.0")).Operations(demoOp())},
		{"doc.info.license.name", NewDocument("host", "/", demoInfo().
			License(NewLicense("", ""))).Operations(demoOp())},
		{"success", NewDocument("host", "/", demoInfo().
			License(NewLicense("name", ""))).Operations(demoOp())},
		{"doc.info.option.tags.name", NewDocument("host", "/", demoInfo()).
			Option(NewOption().Tags(NewTag("", ""))).Operations(demoOp())},
		{"doc.info.option.tags.externalDoc.url", NewDocument("host", "/", demoInfo()).
			Option(NewOption().Tags(NewTag("tag", "").ExternalDoc(NewExternalDoc("", "")))).Operations(demoOp())},
		{"success", NewDocument("host", "/", demoInfo()).
			Option(NewOption().Tags(NewTag("tag", "").ExternalDoc(NewExternalDoc("", "url")))).Operations(demoOp())},
		{"doc.info.option.securities.title", NewDocument("host", "/", demoInfo()).
			Option(NewOption().Securities(NewSecurity("", ""))).Operations(demoOp())},
		{"doc.info.option.securities.typ", NewDocument("host", "/", demoInfo()).
			Option(NewOption().Securities(NewSecurity("title", ""))).Operations(demoOp())},
		{"doc.info.option.securities.name", NewDocument("host", "/", demoInfo()).
			Option(NewOption().Securities(NewSecurity("title", "apiKey").Name("").InLoc(""))).Operations(demoOp())},
		{"doc.info.option.securities.in", NewDocument("host", "/", demoInfo()).
			Option(NewOption().Securities(NewSecurity("title", "apiKey").Name("name").InLoc(""))).Operations(demoOp())},
		{"success", NewDocument("host", "/", demoInfo()).
			Option(NewOption().Securities(NewSecurity("title", "apiKey").Name("name").InLoc("in"))).Operations(demoOp())},
		{"success", NewDocument("host", "/", demoInfo()).
			Option(NewOption().Securities(NewSecurity("title", "basic"))).Operations(demoOp())},
		{"doc.info.option.securities.flow", NewDocument("host", "/", demoInfo()).
			Option(NewOption().Securities(NewSecurity("title", "oauth2").Flow(""))).Operations(demoOp())},
		{"doc.info.option.securities.authorizationUrl", NewDocument("host", "/", demoInfo()).
			Option(NewOption().Securities(NewSecurity("title", "oauth2").Flow("implicit").AuthorizationUrl(""))).Operations(demoOp())},
		{"doc.info.option.securities.tokenUrl", NewDocument("host", "/", demoInfo()).
			Option(NewOption().Securities(NewSecurity("title", "oauth2").Flow("password").TokenUrl(""))).Operations(demoOp())},
		{"doc.info.option.securities.scopes", NewDocument("host", "/", demoInfo()).
			Option(NewOption().Securities(NewSecurity("title", "oauth2").Flow("password").TokenUrl("token").Scopes())).Operations(demoOp())},
		{"doc.info.option.securities.scopes.name", NewDocument("host", "/", demoInfo()).
			Option(NewOption().Securities(NewSecurity("title", "oauth2").Flow("application").TokenUrl("token").Scopes(NewSecurityScope("", "")))).Operations(demoOp())},
		{"success", NewDocument("host", "/", demoInfo()).
			Option(NewOption().Securities(NewSecurity("title", "oauth2").Flow("accessCode").AuthorizationUrl("auth").TokenUrl("token").Scopes(NewSecurityScope("scope", "")))).Operations(demoOp())},
		{"doc.info.option.externalDoc.url", NewDocument("host", "/", demoInfo()).
			Option(NewOption().ExternalDoc(NewExternalDoc("", ""))).Operations(demoOp())},
		{"success", NewDocument("host", "/", demoInfo()).
			Option(NewOption().ExternalDoc(NewExternalDoc("", "url"))).Operations(demoOp())},
		{"doc.option.routesOptions.route", NewDocument("host", "/", demoInfo()).
			Option(NewOption().RoutesOptions(NewRoutesOption(""))).Operations(demoOp())},
		{"success", NewDocument("host", "/", demoInfo()).
			Option(NewOption().RoutesOptions(NewRoutesOption("/"))).Operations(demoOp())},

		{"doc.operations", demoDoc()},
		{"doc.operations.method", demoDoc().Operations(NewOperation("", "", ""))},
		{"doc.operations.route", demoDoc().Operations(NewOperation("get", "", ""))},
		{"doc.operations.summary", demoDoc().Operations(NewOperation("get", "/", ""))},
		{"doc.operations.response", demoDoc().Operations(NewOperation("get", "/", "summary").Responses())},
		{"doc.operations.response.code", demoDoc().Operations(NewOperation("get", "/", "summary").Responses(NewResponse(0, "")))},
		{"doc.operations.response.headers.name", demoDoc().Operations(NewOperation("get", "/", "summary").
			Responses(NewResponse(200, "").Headers(NewResponseHeader("", "", ""))))},
		{"doc.operations.response.headers.type", demoDoc().Operations(NewOperation("get", "/", "summary").
			Responses(NewResponse(200, "").Headers(NewResponseHeader("name", "", ""))))},
		{"doc.operations.response.example.mime", demoDoc().Operations(NewOperation("get", "/", "summary").
			Responses(NewResponse(200, "").Examples(NewResponseExample("", nil))))},
		{"success", demoDoc().Operations(NewOperation("get", "/", "summary").
			Responses(NewResponse(200, "").Headers(NewResponseHeader("name", "integer", ""))))},
		{".externalDoc.url", demoDoc().Operations(NewOperation("get", "/", "summary").Responses(demoResp()).
			ExternalDoc(NewExternalDoc("", "")))},
		{"success", demoDoc().Operations(NewOperation("get", "/", "summary").Responses(demoResp()).
			ExternalDoc(NewExternalDoc("", "url")))},
		{"doc.operations.params.name", demoDoc().Operations(NewOperation("get", "/", "summary").Responses(demoResp()).
			Params(NewParam("", "", "", true, "")))},
		{"doc.operations.params.in", demoDoc().Operations(NewOperation("get", "/", "summary").Responses(demoResp()).
			Params(NewParam("name", "", "", true, "")))},
		{"doc.operations.params.required", demoDoc().Operations(NewOperation("get", "/", "summary").Responses(demoResp()).
			Params(NewParam("name", "path", "", false, "")))},
		{"doc.operations.params.allowEmpty", demoDoc().Operations(NewOperation("get", "/", "summary").Responses(demoResp()).
			Params(NewParam("name", "path", "", true, "").AllowEmpty(true)))},
		{"doc.operations.params.typ", demoDoc().Operations(NewOperation("get", "/", "summary").Responses(demoResp()).
			Params(NewParam("name", "path", "", true, "")))},
		{"success", demoDoc().Operations(NewOperation("get", "/", "summary").Responses(demoResp()).
			Params(NewParam("name", "path", "string", true, "")))},

		{"doc.definitions.name", demoDoc().Operations(demoOp().Responses(demoResp())).Definitions(NewDefinition("", ""))},
		{"success", demoDoc().Operations(demoOp().Responses(demoResp())).Definitions(NewDefinition("name", ""))},
		{"doc.definitions.properties.name", demoDoc().Operations(demoOp().Responses(demoResp())).Definitions(NewDefinition("name", "").
			Properties(NewProperty("", "", true, "")))},
		{"doc.definitions.properties.typ", demoDoc().Operations(demoOp().Responses(demoResp())).Definitions(NewDefinition("name", "").
			Properties(NewProperty("name", "", true, "")))},
		{"success", demoDoc().Operations(demoOp().Responses(demoResp())).Definitions(NewDefinition("name", "").
			Properties(NewProperty("name", "integer", true, "")))},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testPanic(t, tc.name != "success", func() {
				_, _ = tc.giveDoc.GenerateSwaggerJson()
			}, tc.name)
		})
	}

}

type customType struct{}

func (c customType) MarshalYAML() (interface{}, error) {
	return nil, errors.New("demo")
}

func TestRenderAndSaveError(t *testing.T) {
	d1 := NewDocument("host", "/", NewInfo("title", "", "1.0.0")).
		Operations(NewOperation("get", "/", "summary").
			Responses(NewResponse(200, ""))).Definitions(NewDefinition("name", "").
		Properties(NewProperty("name", "integer", true, "").Example(func() {})))
	d2 := NewDocument("host", "/", NewInfo("title", "", "1.0.0")).
		Operations(NewOperation("get", "/", "summary").
			Responses(NewResponse(200, "").Examples(NewResponseExample("mime", customType{}))))
	d3 := NewDocument("host", "/", NewInfo("title", "", "1.0.0")).
		Operations(NewOperation("get", "/", "summary").Responses(NewResponse(200, "")))

	t.Run("json marshal error", func(t *testing.T) {
		_, err := d1.SaveSwaggerJson("docs/error.json")
		if err == nil {
			failNow(t, "SaveSwaggerJson should return error but no error returned 1")
		}
		_, err = d3.SaveSwaggerJson("docs")
		if err == nil {
			failNow(t, "SaveSwaggerJson should return error but no error returned 2")
		}
	})

	t.Run("yaml marshal error", func(t *testing.T) {
		_, err := d2.SaveSwaggerYaml("docs/error.yaml")
		if err == nil {
			failNow(t, "SaveSwaggerYaml should return error but no error returned 1")
		}
		_, err = d3.SaveSwaggerYaml("docs")
		if err == nil {
			failNow(t, "SaveSwaggerYaml should return error but no error returned 2")
		}
	})

	t.Run("apib template error", func(t *testing.T) {
		_bak := apibDocumentTemplate
		defer func() { apibDocumentTemplate = _bak }()
		apibDocumentTemplate = `{{ . }`
		_, err := d1.SaveApib("docs/error.yaml")
		if err == nil {
			failNow(t, "SaveApib should return error but no error returned 1")
		}
		apibDocumentTemplate = `{{ . }}`
		_, err = d3.SaveApib("docs")
		if err == nil {
			failNow(t, "SaveApib should return error but no error returned 2")
		}
	})
}

func TestGenerateError(t *testing.T) {
	// swagger
	for _, tc := range []struct {
		name string
		give *Document
	}{
		{"file array param type", NewDocument("host", "/", NewInfo("title", "", "1.0.0")).AddOperations(NewOperation("get", "/", "s").
			AddResponses(NewResponse(200, "")).Params(NewParam("param", "formData", "file[]", true, "")))},
		{"file response type", NewDocument("host", "/", NewInfo("title", "", "1.0.0")).AddOperations(NewOperation("get", "/", "s").
			AddResponses(NewResponse(200, "file")))},
		{"object non-body param type", NewDocument("host", "/", NewInfo("title", "", "1.0.0")).AddOperations(NewOperation("get", "/", "s").
			Params(NewParam("query", "query", "Obj", true, "")).AddResponses(NewResponse(200, "file"))).AddDefinitions(NewDefinition("Obj", "").AddProperties(NewProperty("name", "string", true, "")))},
		{"non-primitive header type", NewDocument("host", "/", NewInfo("title", "", "1.0.0")).AddOperations(NewOperation("get", "/", "s").
			AddResponses(NewResponse(200, "").AddHeaders(NewResponseHeader("head", "string[]", ""))))},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testPanic(t, tc.name != "success", func() {
				_, _ = tc.give.GenerateSwaggerJson()
			}, tc.name)
		})
	}

	// apib
	tempBak := ""
	for _, tc := range []struct {
		name   string
		give   *Document
		prevFn func()
		nextFn func()
	}{
		{"success", NewDocument("host", "/", NewInfo("title", "", "1.0.0")).AddOperations(NewOperation("get", "/", "s").
			AddResponses(NewResponse(200, "integer"))), nil, nil},
		{"success", NewDocument("host", "/", NewInfo("title", "", "1.0.0")).AddOperations(NewOperation("get", "/", "s").
			AddResponses(NewResponse(200, "").Examples(NewResponseExample("m", func() {}))).Produces("m")), nil, nil},
		{"apibOperationTemplate error", NewDocument("host", "/", NewInfo("title", "", "1.0.0")).AddOperations(NewOperation("get", "/", "s").
			AddResponses(NewResponse(200, ""))), func() { tempBak = apibOperationTemplate; apibOperationTemplate = "{{ . }" }, func() { apibOperationTemplate = tempBak }},
		{"apibGroupsTemplate error", NewDocument("host", "/", NewInfo("title", "", "1.0.0")).AddOperations(NewOperation("get", "/", "s").
			AddResponses(NewResponse(200, ""))), func() { tempBak = apibGroupsTemplate; apibGroupsTemplate = "{{ . }" }, func() { apibGroupsTemplate = tempBak }},
		{"apibDefinitionTemplate error", NewDocument("host", "/", NewInfo("title", "", "1.0.0")).AddOperations(NewOperation("get", "/", "s").
			AddResponses(NewResponse(200, ""))), func() { tempBak = apibDefinitionTemplate; apibDefinitionTemplate = "{{ . }" }, func() { apibDefinitionTemplate = tempBak }},
		{"apibDocumentTemplate error", NewDocument("host", "/", NewInfo("title", "", "1.0.0")).AddOperations(NewOperation("get", "/", "s").
			AddResponses(NewResponse(200, ""))), func() { tempBak = apibDocumentTemplate; apibDocumentTemplate = "{{ . }" }, func() { apibDocumentTemplate = tempBak }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.prevFn != nil {
				tc.prevFn()
			}
			defer func() {
				if tc.nextFn != nil {
					tc.nextFn()
				}
			}()
			wantErr := tc.name != "success"
			_, err := tc.give.GenerateApib()
			if (err == nil && wantErr) || (err != nil && !wantErr) {
				failNow(t, fmt.Sprintf("GenerateApib should return error but no error returned for %v", tc.name))
			}
		})
	}
}
