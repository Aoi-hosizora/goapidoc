package goapidoc

import (
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
	Paths       map[string]map[string]*swagOperation `yaml:"paths,omitempty"               json:"paths,omitempty"`
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
	Type             string        `yaml:"type,omitempty"             json:"type,omitempty"`
	Description      string        `yaml:"description,omitempty"      json:"description,omitempty"`
	Format           string        `yaml:"format,omitempty"           json:"format,omitempty"`
	AllowEmpty       bool          `yaml:"allowEmptyValue,omitempty"  json:"allowEmptyValue,omitempty"`
	Default          interface{}   `yaml:"default,omitempty"          json:"default,omitempty"`
	Example          interface{}   `yaml:"example,omitempty"          json:"example,omitempty"`
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
	ExclusiveMin     bool          `yaml:"exclusiveMin,omitempty"     json:"exclusiveMin,omitempty"`
	ExclusiveMax     bool          `yaml:"exclusiveMax,omitempty"     json:"exclusiveMax,omitempty"`
	MultipleOf       float64       `yaml:"multipleOf,omitempty"       json:"multipleOf,omitempty"`

	Schema *swagSchema `yaml:"schema,omitempty" json:"schema,omitempty"`
	Items  *swagItems  `yaml:"items,omitempty"  json:"items,omitempty"`
}

type swagResponse struct {
	Description string                 `yaml:"description,omitempty" json:"description,omitempty"`
	Headers     map[string]*swagHeader `yaml:"headers,omitempty"     json:"headers,omitempty"`
	Examples    map[string]string      `yaml:"examples,omitempty"    json:"examples,omitempty"`
	Schema      *swagSchema            `yaml:"schema,omitempty"      json:"schema,omitempty"`
}

type swagHeader struct {
	Type        string      `yaml:"type,omitempty"        json:"type,omitempty"`
	Description string      `yaml:"description,omitempty" json:"description,omitempty"`
	Format      string      `yaml:"format,omitempty"      json:"format,omitempty"`
	Default     interface{} `yaml:"default,omitempty"     json:"default,omitempty"`
}

type swagDefinition struct {
	Type        string      `yaml:"type"                  json:"type"`
	Required    []string    `yaml:"required"              json:"required"`
	Description string      `yaml:"description,omitempty" json:"description,omitempty"`
	Properties  *orderedMap `yaml:"properties,omitempty"  json:"properties,omitempty"` // map[string]*swagSchema
}

type swagSchema struct {
	Type        string        `yaml:"type,omitempty"            json:"type,omitempty"`
	Required    bool          `yaml:"required,omitempty"        json:"required,omitempty"`
	Description string        `yaml:"description,omitempty"     json:"description,omitempty"`
	Format      string        `yaml:"format,omitempty"          json:"format,omitempty"`
	AllowEmpty  bool          `yaml:"allowEmptyValue,omitempty" json:"allowEmptyValue,omitempty"`
	Default     interface{}   `yaml:"default,omitempty"         json:"default,omitempty"`
	Example     interface{}   `yaml:"example,omitempty"         json:"example,omitempty"`
	Enum        []interface{} `yaml:"enum,omitempty"            json:"enum,omitempty"`
	Maximum     float64       `yaml:"maximum,omitempty"         json:"maximum,omitempty"`
	Minimum     float64       `yaml:"minimum,omitempty"         json:"minimum,omitempty"`
	MaxLength   int           `yaml:"maxLength,omitempty"       json:"maxLength,omitempty"`
	MinLength   int           `yaml:"minLength,omitempty"       json:"minLength,omitempty"`

	OriginRef string     `yaml:"-"               json:"-"`
	Ref       string     `yaml:"$ref,omitempty"  json:"$ref,omitempty"`
	Items     *swagItems `yaml:"items,omitempty" json:"items,omitempty"`
}

type swagItems struct {
	Type    string        `yaml:"type,omitempty"    json:"type,omitempty"`
	Format  string        `yaml:"format,omitempty"  json:"format,omitempty"`
	Default interface{}   `yaml:"default,omitempty" json:"default,omitempty"`
	Enum    []interface{} `yaml:"enum,omitempty"    json:"enum,omitempty"`

	OriginRef string     `yaml:"-"               json:"-"`
	Ref       string     `yaml:"$ref,omitempty"  json:"$ref,omitempty"`
	Items     *swagItems `yaml:"items,omitempty" json:"items,omitempty"`
}

// ==============
// items & schema
// ==============

func buildSwaggerItems(arr *apiArray) *swagItems {
	/*
		"items": {
		  "type": "integer",
		  "format": "int64"
		}
		"items": {
		  "type": "array",
		  "items": {}
		}
		"items": {
		  "originRef": "User",
		  "$ref": "#/definitions/User"
		}
	*/
	if arr == nil {
		return nil
	}
	if arr.item.kind == apiPrimeKind {
		prime := arr.item.prime
		return &swagItems{Type: prime.typ, Format: prime.format}
	}
	if arr.item.kind == apiArrayKind {
		return &swagItems{Type: ARRAY, Items: buildSwaggerItems(arr.item.array)}
	}
	if arr.item.kind == apiObjectKind {
		origin := arr.item.name
		ref := "#/definitions/" + origin
		return &swagItems{OriginRef: origin, Ref: ref}
	}

	return nil
}

func buildSwaggerSchema(typ string) (outType, outFmt, origin, ref string, items *swagItems) {
	/*
		{
		  "type": "string",
		  "format": "password"
		}
		{
		  "type": "array",
		  "items": {}
		},
		{
		  "schema": {
		    "originRef": "User",
		    "$ref": "#/definitions/User"
		  }
		}
	*/
	at := parseApiType(typ)

	if at.kind == apiPrimeKind {
		outType = at.prime.typ
		outFmt = at.prime.format
		return
	}
	if at.kind == apiArrayKind {
		outType = ARRAY
		items = buildSwaggerItems(at.array)
		return
	}
	if at.kind == apiObjectKind {
		origin = at.name
		ref = "#/definitions/" + origin // ref
		return
	}

	return
}

// ===============================
// params & responses & definition
// ===============================

func buildSwaggerParams(params []*Param) []*swagParam {
	out := make([]*swagParam, 0, len(params))
	for _, p := range params {
		if p.name == "" {
			panic("Param name is required in swagger 2.0")
		}
		if p.in == "" {
			panic("Param in-location is required in swagger 2.0")
		}
		if p.in == PATH && !p.required {
			panic("Path param's required option must be true in swagger 2.0")
		}

		var param *swagParam
		typ, format, origin, ref, items := buildSwaggerSchema(p.typ)
		if p.in != BODY {
			if ref != "" {
				panic("Invalid type `" + p.typ + "` used in non-body parameter")
			}
			param = &swagParam{
				Name:             p.name,
				In:               p.in,
				Required:         p.required,
				Description:      p.desc,
				Type:             typ,
				Format:           format,
				Items:            items,
				AllowEmpty:       p.allowEmpty,
				Default:          p.defaul,
				Example:          p.example,
				Pattern:          p.pattern,
				Enum:             p.enums,
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
					Type:       typ,
					Format:     format,
					Items:      items,
					AllowEmpty: p.allowEmpty,
					Default:    p.defaul,
					Example:    p.example,
					// Pattern:          p.pattern,
					Enum:      p.enums,
					MaxLength: p.maxLength,
					MinLength: p.minLength,
					// MaxItems:         p.maxItems,
					// MinItems:         p.minItems,
					// UniqueItems:      p.uniqueItems,
					// CollectionFormat: p.collectionFormat,
					Maximum: p.maximum,
					Minimum: p.minimum,
					// ExclusiveMin:     p.exclusiveMin,
					// ExclusiveMax:     p.exclusiveMax,
					// MultipleOf:       p.multipleOf,
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
		if r.code == 0 {
			panic("Response code is required in swagger 2.0")
		}
		headers := make(map[string]*swagHeader, len(r.headers))
		for _, h := range r.headers {
			if h.name == "" {
				panic("Response header field name is required in swagger 2.0")
			}
			if h.typ == "" {
				panic("Response header field type is required in swagger 2.0")
			}
			headers[h.name] = &swagHeader{
				Type:        h.typ, // use type directly temporarily
				Description: h.desc,
				// ignore other fields, see https://swagger.io/specification/v2/#headerObject
			}
		}

		typ, format, origin, ref, items := buildSwaggerSchema(r.typ)
		out[strconv.Itoa(r.code)] = &swagResponse{
			Description: r.desc,
			Schema: &swagSchema{
				Type:      typ,
				Format:    format,
				OriginRef: origin,
				Ref:       ref,
				Items:     items,
				// ignore other fields, see https://swagger.io/specification/v2/#schemaObject
			}, // TODO schema
			Examples: r.examples,
			Headers:  headers,
		}
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

		typ, format, origin, ref, items := buildSwaggerSchema(p.typ)
		schema := &swagSchema{
			Required:    p.required,
			Description: p.desc,
			Type:        typ,
			Format:      format,
			AllowEmpty:  p.allowEmpty,
			Default:     p.defaul,
			Example:     p.example,
			Enum:        p.enums,
			Maximum:     p.maximum,
			Minimum:     p.minimum,
			MaxLength:   p.maxLength,
			MinLength:   p.minLength,
			OriginRef:   origin,
			Ref:         ref,
			Items:       items,
		} // TODO schema
		properties.Set(p.name, schema)
	}

	return &swagDefinition{
		Type:        OBJECT,
		Description: definition.desc,
		Required:    required,
		Properties:  properties,
	}
}

// ===================
// paths & definitions
// ===================

func buildSwaggerPaths(doc *Document) map[string]map[string]*swagOperation {
	// path - method - operation
	out := make(map[string]map[string]*swagOperation)
	for _, p := range doc.paths {
		if p.method == "" || p.route == "" {
			panic("Route path is required in swagger 2.0")
		}
		if len(p.responses) == 0 {
			panic("Empty route path response is not allowed in swagger 2.0")
		}
		if !strings.HasPrefix(p.route, "/") {
			panic("Route path must begin with a slash in swagger 2.0")
		}

		method := strings.ToLower(p.method)
		operationId := p.operationId
		if operationId == "" {
			operationId = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
				p.route, "/", "-"), "{", ":"), "}", "") + "-" + method
		}
		securities := make([]map[string][]string, 0, len(p.securities))
		for _, s := range p.securities {
			securities = append(securities, map[string][]string{s: {}}) // support apiKey and basic
		}

		_, ok := out[p.route]
		if !ok {
			out[p.route] = make(map[string]*swagOperation, 2) // cap defaults to 2
		}
		out[p.route][method] = &swagOperation{
			Summary:     p.summary,
			Description: p.desc,
			OperationId: operationId,
			Schemas:     p.schemas,
			Consumes:    p.consumes,
			Produces:    p.produces,
			Tags:        p.tags,
			Securities:  securities,
			Deprecated:  p.deprecated,
			Parameters:  buildSwaggerParams(p.params),
			Responses:   buildSwaggerResponses(p.responses),
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

	// get result map
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
	if doc.host == "" {
		panic("Host is required in swagger 2.0")
	}
	if doc.info == nil {
		panic("Info is required in swagger 2.0")
	}
	if doc.info.title == "" {
		panic("Info.title is required in swagger 2.0")
	}
	if doc.info.version == "" {
		panic("Info.version is required in swagger 2.0")
	}
	if len(doc.paths) == 0 {
		panic("Empty route paths is not allowed in swagger 2.0")
	}

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
			if t.name == "" {
				panic("Tag name is required in swagger 2.0")
			}
			tags = append(tags, &swagTag{Name: t.name, Description: t.desc})
		}
		securities := make(map[string]*swagSecurity, len(doc.option.securities))
		for _, s := range doc.option.securities {
			if s.typ == "apiKey" {
				if s.title == "" {
					panic("Security title is required in swagger 2.0")
				}
				if s.name == "" {
					panic("Security name is required in swagger 2.0")
				}
				if s.in == "" {
					panic("Security in-location is required in swagger 2.0")
				}
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

	// definitions & paths
	out.Definitions = buildSwaggerDefinitions(doc)
	out.Paths = buildSwaggerPaths(doc)

	return out
}
