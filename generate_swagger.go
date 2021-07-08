package goapidoc

import (
	"net/http"
	"strconv"
	"strings"
)

type swagDocument struct {
	Swagger     string                               `yaml:"swagger"                       json:"swagger"`
	Host        string                               `yaml:"host"                          json:"host"`
	BasePath    string                               `yaml:"basePath"                      json:"basePath"`
	Info        *swagInfo                            `yaml:"info"                          json:"info"`
	Schemas     []string                             `yaml:"schemas,omitempty"             json:"schemas,omitempty"`
	Consumes    []string                             `yaml:"consumes,omitempty"            json:"consumes,omitempty"`
	Produces    []string                             `yaml:"produces,omitempty"            json:"produces,omitempty"`
	Tags        []*swagTag                           `yaml:"tags,omitempty"                json:"tags,omitempty"`
	Securities  map[string]*swagSecurity             `yaml:"securityDefinitions,omitempty" json:"securityDefinitions,omitempty"`
	Operations  map[string]map[string]*swagOperation `yaml:"paths,omitempty"               json:"paths,omitempty"`
	Definitions map[string]*swagDefinition           `yaml:"definitions,omitempty"         json:"definitions,omitempty"`
}

type swagInfo struct {
	Title          string       `yaml:"title"                    json:"title"`
	Description    string       `yaml:"description"              json:"description"`
	Version        string       `yaml:"version"                  json:"version"`
	TermsOfService string       `yaml:"termsOfService,omitempty" json:"termsOfService,omitempty"`
	License        *swagLicense `yaml:"license,omitempty"        json:"license,omitempty"`
	Contact        *swagContact `yaml:"contact,omitempty"        json:"contact,omitempty"`
}

type swagLicense struct {
	Name string `yaml:"name"          json:"name"`
	Url  string `yaml:"url,omitempty" json:"url,omitempty"`
}

type swagContact struct {
	Name  string `yaml:"name"            json:"name"`
	Url   string `yaml:"url,omitempty"   json:"url,omitempty"`
	Email string `yaml:"email,omitempty" json:"email,omitempty"`
}

type swagTag struct {
	Name        string `yaml:"name"                  json:"name"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

type swagSecurity struct {
	Type        string `yaml:"type"                  json:"type"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Name        string `yaml:"name"                  json:"name"`
	In          string `yaml:"in"                    json:"in"`
}

type swagOperation struct {
	Summary     string                   `yaml:"summary"               json:"summary"`
	Description string                   `yaml:"description,omitempty" json:"description,omitempty"`
	OperationId string                   `yaml:"operationId"           json:"operationId"`
	Schemas     []string                 `yaml:"schemas,omitempty"     json:"schemas,omitempty"`
	Consumes    []string                 `yaml:"consumes,omitempty"    json:"consumes,omitempty"`
	Produces    []string                 `yaml:"produces,omitempty"    json:"produces,omitempty"`
	Tags        []string                 `yaml:"tags,omitempty"        json:"tags,omitempty"`
	Securities  []map[string][]string    `yaml:"security,omitempty"    json:"security,omitempty"`
	Deprecated  bool                     `yaml:"deprecated,omitempty"  json:"deprecated,omitempty"`
	Parameters  []*swagParam             `yaml:"parameters,omitempty"  json:"parameters,omitempty"`
	Responses   map[string]*swagResponse `yaml:"responses,omitempty"   json:"responses,omitempty"`
}

type swagParam struct {
	Name             string        `yaml:"name"                       json:"name"`
	In               string        `yaml:"in"                         json:"in"`
	Required         bool          `yaml:"required"                   json:"required"`
	Description      string        `yaml:"description,omitempty"      json:"description,omitempty"`
	Type             string        `yaml:"type,omitempty"             json:"type,omitempty"`
	Format           string        `yaml:"format,omitempty"           json:"format,omitempty"`
	AllowEmpty       bool          `yaml:"allowEmptyValue,omitempty"  json:"allowEmptyValue,omitempty"`
	Default          interface{}   `yaml:"default,omitempty"          json:"default,omitempty"`
	Example          interface{}   `yaml:"example,omitempty"          json:"example,omitempty"` // ?
	Pattern          string        `yaml:"pattern,omitempty"          json:"pattern,omitempty"`
	Enum             []interface{} `yaml:"enum,omitempty"             json:"enum,omitempty"`
	MaxLength        int           `yaml:"maxLength,omitempty"        json:"maxLength,omitempty"`
	MinLength        int           `yaml:"minLength,omitempty"        json:"minLength,omitempty"`
	MaxItems         int           `yaml:"maxItems,omitempty"         json:"maxItems,omitempty"`
	MinItems         int           `yaml:"minItems,omitempty"         json:"minItems,omitempty"`
	UniqueItems      bool          `yaml:"uniqueItems,omitempty"      json:"uniqueItems,omitempty"`
	CollectionFormat string        `yaml:"collectionFormat,omitempty" json:"collectionFormat,omitempty"`
	Maximum          float64       `yaml:"maximum,omitempty"          json:"maximum,omitempty"`
	Minimum          float64       `yaml:"minimum,omitempty"          json:"minimum,omitempty"`
	ExclusiveMin     bool          `yaml:"exclusiveMinimum,omitempty" json:"exclusiveMinimum,omitempty"`
	ExclusiveMax     bool          `yaml:"exclusiveMaximum,omitempty" json:"exclusiveMaximum,omitempty"`
	MultipleOf       float64       `yaml:"multipleOf,omitempty"       json:"multipleOf,omitempty"`

	Items  *swagItems  `yaml:"items,omitempty"  json:"items,omitempty"`
	Schema *swagSchema `yaml:"schema,omitempty" json:"schema,omitempty"`
}

type swagResponse struct {
	Description string                 `yaml:"description,omitempty" json:"description,omitempty"`
	Headers     map[string]*swagHeader `yaml:"headers,omitempty"     json:"headers,omitempty"`
	Examples    map[string]string      `yaml:"examples,omitempty"    json:"examples,omitempty"`
	Schema      *swagResponseSchema    `yaml:"schema,omitempty"      json:"schema,omitempty"`
}

type swagHeader struct {
	Type        string `yaml:"type,omitempty"        json:"type,omitempty"`
	Format      string `yaml:"format,omitempty"      json:"format,omitempty"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

type swagResponseSchema struct {
	Type      string     `yaml:"type,omitempty"   json:"type,omitempty"`
	Format    string     `yaml:"format,omitempty" json:"format,omitempty"`
	Items     *swagItems `yaml:"items,omitempty"  json:"items,omitempty"`
	OriginRef string     `yaml:"-"                json:"-"`
	Ref       string     `yaml:"$ref,omitempty"   json:"$ref,omitempty"`
}

type swagDefinition struct {
	Type        string      `yaml:"type"                  json:"type"`
	Required    []string    `yaml:"required"              json:"required"`
	Description string      `yaml:"description,omitempty" json:"description,omitempty"`
	Properties  *orderedMap `yaml:"properties,omitempty"  json:"properties,omitempty"` // map[string]*swagSchema
}

type swagSchema struct {
	Type             string        `yaml:"type,omitempty"             json:"type,omitempty"`
	Format           string        `yaml:"format,omitempty"           json:"format,omitempty"`
	Required         bool          `yaml:"required,omitempty"         json:"required,omitempty"` // <<<
	Description      string        `yaml:"description,omitempty"      json:"description,omitempty"`
	AllowEmpty       bool          `yaml:"allowEmptyValue,omitempty"  json:"allowEmptyValue,omitempty"` // ?
	Default          interface{}   `yaml:"default,omitempty"          json:"default,omitempty"`
	Example          interface{}   `yaml:"example,omitempty"          json:"example,omitempty"`
	Pattern          string        `yaml:"pattern,omitempty"          json:"pattern,omitempty"`
	Enum             []interface{} `yaml:"enum,omitempty"             json:"enum,omitempty"`
	MaxLength        int           `yaml:"maxLength,omitempty"        json:"maxLength,omitempty"`
	MinLength        int           `yaml:"minLength,omitempty"        json:"minLength,omitempty"`
	MaxItems         int           `yaml:"maxItems,omitempty"         json:"maxItems,omitempty"`
	MinItems         int           `yaml:"minItems,omitempty"         json:"minItems,omitempty"`
	UniqueItems      bool          `yaml:"uniqueItems,omitempty"      json:"uniqueItems,omitempty"`
	CollectionFormat string        `yaml:"collectionFormat,omitempty" json:"collectionFormat,omitempty"` // ?
	Maximum          float64       `yaml:"maximum,omitempty"          json:"maximum,omitempty"`
	Minimum          float64       `yaml:"minimum,omitempty"          json:"minimum,omitempty"`
	ExclusiveMin     bool          `yaml:"exclusiveMinimum,omitempty" json:"exclusiveMinimum,omitempty"`
	ExclusiveMax     bool          `yaml:"exclusiveMaximum,omitempty" json:"exclusiveMaximum,omitempty"`
	MultipleOf       float64       `yaml:"multipleOf,omitempty"       json:"multipleOf,omitempty"`

	Items     *swagItems `yaml:"items,omitempty" json:"items,omitempty"`
	OriginRef string     `yaml:"-"               json:"-"`
	Ref       string     `yaml:"$ref,omitempty"  json:"$ref,omitempty"`
}

type swagItems struct {
	Type             string        `yaml:"type,omitempty"             json:"type,omitempty"`
	Format           string        `yaml:"format,omitempty"           json:"format,omitempty"`
	AllowEmpty       bool          `yaml:"allowEmptyValue,omitempty"  json:"allowEmptyValue,omitempty"` // ?
	Default          interface{}   `yaml:"default,omitempty"          json:"default,omitempty"`
	Example          interface{}   `yaml:"example,omitempty"          json:"example,omitempty"` // ?
	Pattern          string        `yaml:"pattern,omitempty"          json:"pattern,omitempty"`
	Enum             []interface{} `yaml:"enum,omitempty"             json:"enum,omitempty"`
	MaxLength        int           `yaml:"maxLength,omitempty"        json:"maxLength,omitempty"`
	MinLength        int           `yaml:"minLength,omitempty"        json:"minLength,omitempty"`
	MaxItems         int           `yaml:"maxItems,omitempty"         json:"maxItems,omitempty"`
	MinItems         int           `yaml:"minItems,omitempty"         json:"minItems,omitempty"`
	UniqueItems      bool          `yaml:"uniqueItems,omitempty"      json:"uniqueItems,omitempty"`
	CollectionFormat string        `yaml:"collectionFormat,omitempty" json:"collectionFormat,omitempty"`
	Maximum          float64       `yaml:"maximum,omitempty"          json:"maximum,omitempty"`
	Minimum          float64       `yaml:"minimum,omitempty"          json:"minimum,omitempty"`
	ExclusiveMin     bool          `yaml:"exclusiveMinimum,omitempty" json:"exclusiveMinimum,omitempty"`
	ExclusiveMax     bool          `yaml:"exclusiveMaximum,omitempty" json:"exclusiveMaximum,omitempty"`
	MultipleOf       float64       `yaml:"multipleOf,omitempty"       json:"multipleOf,omitempty"`

	Items     *swagItems `yaml:"items,omitempty" json:"items,omitempty"`
	OriginRef string     `yaml:"-"               json:"-"`
	Ref       string     `yaml:"$ref,omitempty"  json:"$ref,omitempty"`
}

// ==============
// items & schema
// ==============

func buildSwaggerItems(arr *apiArray, option *ItemOption) *swagItems {
	/*
		"items": {
		  "type": "integer",
		  "format": "int64",
		  // ...
		}
		"items": {
		  "type": "array",
		  "items": {},
		  // ...
		}
		"items": {
		  "$ref": "#/definitions/User"
		}
	*/
	if arr == nil {
		return nil
	}

	var items *swagItems
	if option == nil {
		items = &swagItems{}
	} else {
		items = &swagItems{
			AllowEmpty:       option.allowEmpty, // ?
			Default:          option.defaul,
			Example:          option.example, // ?
			Pattern:          option.pattern,
			Enum:             option.enum,
			MaxLength:        option.maxLength,
			MinLength:        option.minLength,
			MaxItems:         option.maxItems,
			MinItems:         option.minItems,
			UniqueItems:      option.uniqueItems,
			CollectionFormat: option.collectionFormat,
			Maximum:          option.maximum,
			Minimum:          option.minimum,
			ExclusiveMin:     option.exclusiveMin,
			ExclusiveMax:     option.exclusiveMax,
			MultipleOf:       option.multipleOf,
		}
	}

	switch arr.item.kind {
	case apiPrimeKind:
		prime := arr.item.prime
		items.Type = prime.typ
		items.Format = prime.format
		return items
	case apiArrayKind:
		items.Type = ARRAY
		var o *ItemOption
		if option != nil {
			o = option.itemOption
		}
		items.Items = buildSwaggerItems(arr.item.array, o)
		return items
	case apiObjectKind:
		origin := arr.item.name
		ref := "#/definitions/" + origin
		return &swagItems{OriginRef: origin, Ref: ref}
	default:
		return nil
	}
}

func buildSwaggerSchema(typ string, option *ItemOption) (outType, outFmt, origin, ref string, items *swagItems) {
	/*
		{
		  "type": "string",
		  "format": "password",
		  // ...
		}
		{
		  "type": "array",
		  "items": {},
		  // ...
		},
		{
		  "$ref": "#/definitions/User"
		}
	*/
	at := parseApiType(typ)

	switch at.kind {
	case apiPrimeKind:
		outType = at.prime.typ
		outFmt = at.prime.format
		return
	case apiArrayKind:
		outType = ARRAY
		items = buildSwaggerItems(at.array, option)
		return
	case apiObjectKind:
		origin = at.name
		ref = "#/definitions/" + origin // ref
		return
	default:
		return
	}
}

// ===============================
// params & responses & definition
// ===============================

func buildSwaggerParams(params []*Param) []*swagParam {
	out := make([]*swagParam, 0, len(params))
	for _, p := range params {
		var param *swagParam
		typ, format, origin, ref, items := buildSwaggerSchema(p.typ, p.itemOption)
		if p.in != BODY {
			if ref != "" {
				panic("Invalid type `" + p.typ + "` used in non-body parameter") // only allowed primitive and array
			}
			param = &swagParam{
				Name:             p.name,
				In:               p.in,
				Required:         p.required,
				Description:      p.desc,
				Type:             typ,
				Format:           format,
				AllowEmpty:       p.allowEmpty,
				Default:          p.defaul,
				Example:          p.example, // ?
				Pattern:          p.pattern,
				Enum:             p.enum,
				MaxLength:        p.maxLength,
				MinLength:        p.minLength,
				MaxItems:         p.maxItems,
				MinItems:         p.minItems,
				UniqueItems:      p.uniqueItems,
				CollectionFormat: p.collectionFormat,
				Maximum:          p.maximum,
				Minimum:          p.minimum,
				ExclusiveMin:     p.exclusiveMin,
				ExclusiveMax:     p.exclusiveMax,
				MultipleOf:       p.multipleOf,
				Items:            items,
			}
		} else {
			// must put in schema
			param = &swagParam{
				Name:        p.name,
				In:          p.in,
				Required:    p.required,
				Description: p.desc,
			}
			if ref != "" {
				param.Schema = &swagSchema{OriginRef: origin, Ref: ref}
			} else {
				param.Schema = &swagSchema{
					Type:             typ,
					Format:           format,
					Required:         p.required,
					Description:      p.desc,
					AllowEmpty:       p.allowEmpty, // ?
					Default:          p.defaul,
					Example:          p.example,
					Pattern:          p.pattern,
					Enum:             p.enum,
					MaxLength:        p.maxLength,
					MinLength:        p.minLength,
					MaxItems:         p.maxItems,
					MinItems:         p.minItems,
					UniqueItems:      p.uniqueItems,
					CollectionFormat: p.collectionFormat, // ?
					Maximum:          p.maximum,
					Minimum:          p.minimum,
					ExclusiveMin:     p.exclusiveMin,
					ExclusiveMax:     p.exclusiveMax,
					MultipleOf:       p.multipleOf,
					Items:            items,
				} // TODO schema
			}
		}

		out = append(out, param)
	}
	return out
}

func buildSwaggerResponses(responses []*Response) map[string]*swagResponse {
	out := make(map[string]*swagResponse, len(responses))
	for _, r := range responses {
		desc := r.desc
		if desc == "" {
			desc = strconv.Itoa(r.code) + " " + http.StatusText(r.code)
		}
		headers := make(map[string]*swagHeader, len(r.headers))
		for _, h := range r.headers {
			typ, format, _, ref, items := buildSwaggerSchema(h.typ, nil)
			if ref != "" || items != nil {
				panic("Invalid type `" + h.typ + "` used in response header") // only allow primitive
			}
			headers[h.name] = &swagHeader{
				Type:        typ,
				Format:      format,
				Description: h.desc,
				// ignore other fields, see https://swagger.io/specification/v2/#headerObject
			}
		}

		resp := &swagResponse{
			Description: desc,
			Examples:    r.examples,
			Headers:     headers,
		}
		if r.typ != "" {
			typ, format, origin, ref, items := buildSwaggerSchema(r.typ, nil)
			resp.Schema = &swagResponseSchema{
				Type:      typ,
				Format:    format,
				Items:     items,
				OriginRef: origin,
				Ref:       ref,
				// ignore other fields, see https://swagger.io/specification/v2/#schemaObject
			}
		}
		out[strconv.Itoa(r.code)] = resp
	}
	return out
}

func buildSwaggerDefinition(definition *Definition) *swagDefinition {
	required := make([]string, 0, len(definition.properties))
	properties := newOrderedMap(len(definition.properties)) // map[string]*swagSchema
	for _, p := range definition.properties {
		if p.required {
			required = append(required, p.name)
		}

		var schema *swagSchema
		typ, format, origin, ref, items := buildSwaggerSchema(p.typ, p.itemOption)
		if ref != "" {
			schema = &swagSchema{OriginRef: origin, Ref: ref}
		} else {
			schema = &swagSchema{
				// Required: p.required,
				Type:             typ,
				Format:           format,
				Description:      p.desc,
				AllowEmpty:       p.allowEmpty, // ?
				Default:          p.defaul,
				Example:          p.example,
				Pattern:          p.pattern,
				Enum:             p.enum,
				MaxLength:        p.maxLength,
				MinLength:        p.minLength,
				MaxItems:         p.maxItems,
				MinItems:         p.minItems,
				UniqueItems:      p.uniqueItems,
				CollectionFormat: p.collectionFormat, // ?
				Maximum:          p.maximum,
				Minimum:          p.minimum,
				ExclusiveMin:     p.exclusiveMin,
				ExclusiveMax:     p.exclusiveMax,
				MultipleOf:       p.multipleOf,
				Items:            items,
			} // TODO schema
		}
		properties.Set(p.name, schema)
	}

	return &swagDefinition{
		Type:        OBJECT, // fixed schema type to object
		Description: definition.desc,
		Required:    required,
		Properties:  properties,
	}
}

// ========================
// operations & definitions
// ========================

func buildSwaggerOperations(doc *Document) map[string]map[string]*swagOperation {
	// route - method - operation
	out := make(map[string]map[string]*swagOperation, 2) // cap defaults to 2
	for _, op := range doc.operations {
		method := strings.ToLower(op.method)
		operationId := op.operationId
		if operationId == "" {
			operationId = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
				op.route, "/", "-"), "{", ":"), "}", "") + "-" + method
		}
		securities := make([]map[string][]string, 0, len(op.securities))
		for _, s := range op.securities {
			securities = append(securities, map[string][]string{s: {}}) // support apiKey and basic
		}

		_, ok := out[op.route]
		if !ok {
			out[op.route] = make(map[string]*swagOperation, 4) // cap defaults to 4
		}
		out[op.route][method] = &swagOperation{
			Summary:     op.summary,
			Description: op.desc,
			OperationId: operationId,
			Schemas:     op.schemas,
			Consumes:    op.consumes,
			Produces:    op.produces,
			Tags:        op.tags,
			Securities:  securities,
			Deprecated:  op.deprecated,
			Parameters:  buildSwaggerParams(op.params),
			Responses:   buildSwaggerResponses(op.responses),
		}
	}
	return out
}

func buildSwaggerDefinitions(doc *Document) map[string]*swagDefinition {
	// prehandle definition list
	allSpecTypes := collectAllSpecTypes(doc)
	clonedDefinitions := make([]*Definition, 0, len(doc.definitions))
	for _, definition := range doc.definitions {
		clonedDefinitions = append(clonedDefinitions, prehandleDefinition(definition)) // with generic name checked
	}
	newDefinitionList := prehandleDefinitionList(clonedDefinitions, allSpecTypes)

	// return result map
	out := make(map[string]*swagDefinition, len(newDefinitionList))
	for _, definition := range newDefinitionList {
		out[definition.name] = buildSwaggerDefinition(definition)
	}
	return out
}

// ========
// document
// ========

func buildSwaggerDocument(doc *Document) *swagDocument {
	// check
	checkDocument(doc)

	// info
	out := &swagDocument{
		Swagger:  "2.0",
		Host:     doc.host,
		BasePath: doc.basePath,
		Info: &swagInfo{
			Title:          doc.info.title,
			Description:    doc.info.desc,
			Version:        doc.info.version,
			TermsOfService: doc.info.termsOfService,
		},
	}
	if doc.info.license != nil {
		out.Info.License = &swagLicense{Name: doc.info.license.name, Url: doc.info.license.url}
	}
	if doc.info.contact != nil {
		out.Info.Contact = &swagContact{Name: doc.info.contact.name, Url: doc.info.contact.url, Email: doc.info.contact.email}
	}

	// option
	if doc.option != nil {
		tags := make([]*swagTag, 0, len(doc.option.tags))
		for _, t := range doc.option.tags {
			tags = append(tags, &swagTag{Name: t.name, Description: t.desc})
		}
		securities := make(map[string]*swagSecurity, len(doc.option.securities))
		for _, s := range doc.option.securities {
			if s.typ == "apiKey" {
				securities[s.title] = &swagSecurity{Type: "apiKey", Description: s.desc, Name: s.name, In: s.in}
			} else if s.typ == "basic" {
				securities[s.title] = &swagSecurity{Type: "basic", Description: s.desc}
			}
		}

		out.Schemas = doc.option.schemas
		out.Consumes = doc.option.consumes
		out.Produces = doc.option.produces
		out.Tags = tags
		out.Securities = securities
	}

	// definitions & operations
	out.Definitions = buildSwaggerDefinitions(doc)
	out.Operations = buildSwaggerOperations(doc)

	return out
}
