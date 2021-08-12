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
	Schemes     []string                             `yaml:"schemes,omitempty"             json:"schemes,omitempty"`
	Consumes    []string                             `yaml:"consumes,omitempty"            json:"consumes,omitempty"`
	Produces    []string                             `yaml:"produces,omitempty"            json:"produces,omitempty"`
	Tags        []*swagTag                           `yaml:"tags,omitempty"                json:"tags,omitempty"`
	Securities  map[string]*swagSecurity             `yaml:"securityDefinitions,omitempty" json:"securityDefinitions,omitempty"`
	ExternalDoc *swagExternalDoc                     `yaml:"externalDocs,omitempty"        json:"externalDocs,omitempty"`
	Operations  map[string]map[string]*swagOperation `yaml:"paths,omitempty"               json:"paths,omitempty"`
	Definitions map[string]*swagDefinition           `yaml:"definitions,omitempty"         json:"definitions,omitempty"`
}

type swagInfo struct {
	Title          string       `yaml:"title"                    json:"title"`
	Version        string       `yaml:"version"                  json:"version"`
	Description    string       `yaml:"description,omitempty"    json:"description,omitempty"`
	TermsOfService string       `yaml:"termsOfService,omitempty" json:"termsOfService,omitempty"`
	License        *swagLicense `yaml:"license,omitempty"        json:"license,omitempty"`
	Contact        *swagContact `yaml:"contact,omitempty"        json:"contact,omitempty"`
}

type swagLicense struct {
	Name string `yaml:"name"          json:"name"`
	Url  string `yaml:"url,omitempty" json:"url,omitempty"`
}

type swagContact struct {
	Name  string `yaml:"name,omitempty"  json:"name,omitempty"`
	Url   string `yaml:"url,omitempty"   json:"url,omitempty"`
	Email string `yaml:"email,omitempty" json:"email,omitempty"`
}

type swagTag struct {
	Name        string           `yaml:"name"                   json:"name"`
	Description string           `yaml:"description,omitempty"  json:"description,omitempty"`
	ExternalDoc *swagExternalDoc `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
}

type swagSecurity struct {
	Type             string            `yaml:"type"                       json:"type"`
	Description      string            `yaml:"description,omitempty"      json:"description,omitempty"`
	Name             string            `yaml:"name,omitempty"             json:"name,omitempty"`
	In               string            `yaml:"in,omitempty"               json:"in,omitempty"`
	Flow             string            `yaml:"flow,omitempty"             json:"flow,omitempty"`
	AuthorizationUrl string            `yaml:"authorizationUrl,omitempty" json:"authorizationUrl,omitempty"`
	TokenUrl         string            `yaml:"tokenUrl,omitempty"         json:"tokenUrl,omitempty"`
	Scopes           map[string]string `yaml:"scopes,omitempty"           json:"scopes,omitempty"`
}

type swagExternalDoc struct {
	Url         string `yaml:"url"                   json:"url"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

type swagOperation struct {
	Summary     string                   `yaml:"summary"                json:"summary"`
	OperationId string                   `yaml:"operationId"            json:"operationId"`
	Description string                   `yaml:"description,omitempty"  json:"description,omitempty"`
	Schemes     []string                 `yaml:"schemes,omitempty"      json:"schemes,omitempty"`
	Consumes    []string                 `yaml:"consumes,omitempty"     json:"consumes,omitempty"`
	Produces    []string                 `yaml:"produces,omitempty"     json:"produces,omitempty"`
	Tags        []string                 `yaml:"tags,omitempty"         json:"tags,omitempty"`
	Securities  []map[string][]string    `yaml:"security,omitempty"     json:"security,omitempty"`
	Deprecated  bool                     `yaml:"deprecated,omitempty"   json:"deprecated,omitempty"`
	ExternalDoc *swagExternalDoc         `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
	Parameters  []*swagParam             `yaml:"parameters,omitempty"   json:"parameters,omitempty"`
	Responses   map[string]*swagResponse `yaml:"responses,omitempty"    json:"responses,omitempty"`
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
	MaxLength        *int          `yaml:"maxLength,omitempty"        json:"maxLength,omitempty"`
	MinLength        *int          `yaml:"minLength,omitempty"        json:"minLength,omitempty"`
	MaxItems         *int          `yaml:"maxItems,omitempty"         json:"maxItems,omitempty"`
	MinItems         *int          `yaml:"minItems,omitempty"         json:"minItems,omitempty"`
	UniqueItems      bool          `yaml:"uniqueItems,omitempty"      json:"uniqueItems,omitempty"`
	CollectionFormat string        `yaml:"collectionFormat,omitempty" json:"collectionFormat,omitempty"`
	Maximum          *float64      `yaml:"maximum,omitempty"          json:"maximum,omitempty"`
	Minimum          *float64      `yaml:"minimum,omitempty"          json:"minimum,omitempty"`
	ExclusiveMin     bool          `yaml:"exclusiveMinimum,omitempty" json:"exclusiveMinimum,omitempty"`
	ExclusiveMax     bool          `yaml:"exclusiveMaximum,omitempty" json:"exclusiveMaximum,omitempty"`
	MultipleOf       float64       `yaml:"multipleOf,omitempty"       json:"multipleOf,omitempty"`
	XMLRepr          *swagXMLRepr  `yaml:"xml,omitempty"              json:"xml,omitempty"`

	Items  *swagItems  `yaml:"items,omitempty"  json:"items,omitempty"`
	Schema *swagSchema `yaml:"schema,omitempty" json:"schema,omitempty"`
}

type swagResponse struct {
	Description string                         `yaml:"description"        json:"description"`
	Headers     map[string]*swagResponseHeader `yaml:"headers,omitempty"  json:"headers,omitempty"`
	Examples    map[string]interface{}         `yaml:"examples,omitempty" json:"examples,omitempty"`
	Schema      *swagResponseSchema            `yaml:"schema,omitempty"   json:"schema,omitempty"`
}

type swagResponseHeader struct {
	Type        string      `yaml:"type"                  json:"type"`
	Format      string      `yaml:"format,omitempty"      json:"format,omitempty"`
	Description string      `yaml:"description,omitempty" json:"description,omitempty"`
	Example     interface{} `yaml:"example,omitempty"     json:"example,omitempty"` // ?
}

type swagResponseSchema struct {
	Type      string     `yaml:"type,omitempty"   json:"type,omitempty"`
	Format    string     `yaml:"format,omitempty" json:"format,omitempty"`
	Items     *swagItems `yaml:"items,omitempty"  json:"items,omitempty"`
	OriginRef string     `yaml:"-"                json:"-"`
	Ref       string     `yaml:"$ref,omitempty"   json:"$ref,omitempty"`
}

type swagDefinition struct {
	Type        string       `yaml:"type"                  json:"type"`
	Required    []string     `yaml:"required,omitempty"    json:"required,omitempty"`
	Description string       `yaml:"description,omitempty" json:"description,omitempty"`
	XMLRepr     *swagXMLRepr `yaml:"xml,omitempty"         json:"xml,omitempty"`
	Properties  *orderedMap  `yaml:"properties,omitempty"  json:"properties,omitempty"` // map[string]*swagSchema
}

type swagXMLRepr struct {
	Name      string `yaml:"name,omitempty"      json:"name,omitempty"`
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Prefix    string `yaml:"prefix,omitempty"    json:"prefix,omitempty"`
	Attribute bool   `yaml:"attribute,omitempty" json:"attribute,omitempty"`
	Wrapped   bool   `yaml:"wrapped,omitempty"   json:"wrapped,omitempty"`
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
	MaxLength        *int          `yaml:"maxLength,omitempty"        json:"maxLength,omitempty"`
	MinLength        *int          `yaml:"minLength,omitempty"        json:"minLength,omitempty"`
	MaxItems         *int          `yaml:"maxItems,omitempty"         json:"maxItems,omitempty"`
	MinItems         *int          `yaml:"minItems,omitempty"         json:"minItems,omitempty"`
	UniqueItems      bool          `yaml:"uniqueItems,omitempty"      json:"uniqueItems,omitempty"`
	CollectionFormat string        `yaml:"collectionFormat,omitempty" json:"collectionFormat,omitempty"` // ?
	Maximum          *float64      `yaml:"maximum,omitempty"          json:"maximum,omitempty"`
	Minimum          *float64      `yaml:"minimum,omitempty"          json:"minimum,omitempty"`
	ExclusiveMin     bool          `yaml:"exclusiveMinimum,omitempty" json:"exclusiveMinimum,omitempty"`
	ExclusiveMax     bool          `yaml:"exclusiveMaximum,omitempty" json:"exclusiveMaximum,omitempty"`
	MultipleOf       float64       `yaml:"multipleOf,omitempty"       json:"multipleOf,omitempty"`
	XMLRepr          *swagXMLRepr  `yaml:"xml,omitempty"              json:"xml,omitempty"`

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
	MaxLength        *int          `yaml:"maxLength,omitempty"        json:"maxLength,omitempty"`
	MinLength        *int          `yaml:"minLength,omitempty"        json:"minLength,omitempty"`
	MaxItems         *int          `yaml:"maxItems,omitempty"         json:"maxItems,omitempty"`
	MinItems         *int          `yaml:"minItems,omitempty"         json:"minItems,omitempty"`
	UniqueItems      bool          `yaml:"uniqueItems,omitempty"      json:"uniqueItems,omitempty"`
	CollectionFormat string        `yaml:"collectionFormat,omitempty" json:"collectionFormat,omitempty"`
	Maximum          *float64      `yaml:"maximum,omitempty"          json:"maximum,omitempty"`
	Minimum          *float64      `yaml:"minimum,omitempty"          json:"minimum,omitempty"`
	ExclusiveMin     bool          `yaml:"exclusiveMinimum,omitempty" json:"exclusiveMinimum,omitempty"`
	ExclusiveMax     bool          `yaml:"exclusiveMaximum,omitempty" json:"exclusiveMaximum,omitempty"`
	MultipleOf       float64       `yaml:"multipleOf,omitempty"       json:"multipleOf,omitempty"`
	XMLRepr          *swagXMLRepr  `yaml:"xml,omitempty"              json:"xml,omitempty"`

	Items     *swagItems `yaml:"items,omitempty" json:"items,omitempty"`
	OriginRef string     `yaml:"-"               json:"-"`
	Ref       string     `yaml:"$ref,omitempty"  json:"$ref,omitempty"`
}

// ======================================
// items & schema & externalDoc & xmlRepr
// ======================================

func buildSwagItems(arr *apiArray, opt *ItemOption) *swagItems {
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
	var items *swagItems
	if opt == nil {
		items = &swagItems{}
	} else {
		items = &swagItems{
			AllowEmpty:       opt.allowEmpty, // ?
			Default:          opt.defaul,
			Example:          opt.example, // ?
			Pattern:          opt.pattern,
			Enum:             opt.enum,
			MaxLength:        opt.maxLength,
			MinLength:        opt.minLength,
			MaxItems:         opt.maxItems,
			MinItems:         opt.minItems,
			UniqueItems:      opt.uniqueItems,
			CollectionFormat: opt.collectionFormat,
			Maximum:          opt.maximum,
			Minimum:          opt.minimum,
			ExclusiveMin:     opt.exclusiveMin,
			ExclusiveMax:     opt.exclusiveMax,
			MultipleOf:       opt.multipleOf,
			XMLRepr:          buildSwagXMLRepr(opt.xmlRepr),
		}
	}

	switch arr.item.kind {
	case apiPrimeKind:
		prime := arr.item.prime
		if prime.typ == FILE {
			panic("Invalid file type used in non-request parameter")
		}
		items.Type = prime.typ
		items.Format = prime.format
		return items
	case apiArrayKind:
		items.Type = ARRAY
		var o *ItemOption
		if opt != nil {
			o = opt.itemOption
		}
		items.Items = buildSwagItems(arr.item.array, o)
		return items
	case apiObjectKind:
		origin := arr.item.name
		ref := "#/definitions/" + origin
		return &swagItems{OriginRef: origin, Ref: ref}
	default:
		return nil // unreachable
	}
}

func buildSwagSchema(typ string, option *ItemOption, allowFile bool) (outType, outFmt, origin, ref string, items *swagItems) {
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
		if at.prime.typ == FILE && !allowFile {
			panic("Invalid file type used in non-request parameter")
		}
		outType = at.prime.typ
		outFmt = at.prime.format
		return
	case apiArrayKind:
		outType = ARRAY
		items = buildSwagItems(at.array, option)
		return
	case apiObjectKind:
		origin = at.name
		ref = "#/definitions/" + origin // ref
		return
	default:
		return // unreachable
	}
}

func buildSwagExternalDoc(doc *ExternalDoc) *swagExternalDoc {
	if doc == nil {
		return nil
	}
	return &swagExternalDoc{Url: doc.url, Description: doc.desc}
}

func buildSwagXMLRepr(xml *XMLRepr) *swagXMLRepr {
	if xml == nil {
		return nil
	}
	return &swagXMLRepr{
		Name:      xml.name,
		Namespace: xml.namespace,
		Prefix:    xml.prefix,
		Attribute: xml.attribute,
		Wrapped:   xml.wrapped,
	}
}

// ===============================
// params & responses & definition
// ===============================

func buildSwagParams(params []*Param) []*swagParam {
	out := make([]*swagParam, 0, len(params))
	for _, p := range params {
		var param *swagParam
		typ, format, origin, ref, items := buildSwagSchema(p.typ, p.itemOption, true)
		if p.in != BODY {
			// cannot use schema
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
				XMLRepr:          buildSwagXMLRepr(p.xmlRepr),
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
					XMLRepr:          buildSwagXMLRepr(p.xmlRepr),
					Items:            items,
				}
			}
		}

		out = append(out, param)
	}
	return out
}

func buildSwagResponses(responses []*Response) map[string]*swagResponse {
	out := make(map[string]*swagResponse, len(responses))
	for _, r := range responses {
		desc := r.desc
		if desc == "" {
			desc = strconv.Itoa(r.code) + " " + http.StatusText(r.code)
		}
		headers := make(map[string]*swagResponseHeader, len(r.headers))
		for _, h := range r.headers {
			typ, format, _, ref, items := buildSwagSchema(h.typ, nil, false)
			if ref != "" || items != nil {
				panic("Invalid type `" + h.typ + "` used in response header") // only allow primitive
			}
			headers[h.name] = &swagResponseHeader{
				Type:        typ,
				Format:      format,
				Description: h.desc,
				Example:     h.example, // ?
				// ignore other fields, see https://swagger.io/specification/v2/#headerObject
			}
		}
		examples := make(map[string]interface{}, len(r.examples))
		for _, e := range r.examples {
			examples[e.mime] = e.example
		}

		resp := &swagResponse{
			Description: desc,
			Headers:     headers,
			Examples:    examples,
		}
		if r.typ != "" {
			typ, format, origin, ref, items := buildSwagSchema(r.typ, nil, false)
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

func buildSwagDefinition(definition *Definition) *swagDefinition {
	required := make([]string, 0, len(definition.properties)/2)
	properties := newOrderedMap(len(definition.properties)) // map[string]*swagSchema
	for _, p := range definition.properties {
		if p.required {
			required = append(required, p.name)
		}

		var schema *swagSchema
		typ, format, origin, ref, items := buildSwagSchema(p.typ, p.itemOption, false)
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
				XMLRepr:          buildSwagXMLRepr(p.xmlRepr),
				Items:            items,
			}
		}
		properties.Set(p.name, schema)
	}

	return &swagDefinition{
		Type:        OBJECT, // fixed schema type to object
		Required:    required,
		Description: definition.desc,
		XMLRepr:     buildSwagXMLRepr(definition.xmlRepr),
		Properties:  properties,
	}
}

// ========================
// operations & definitions
// ========================

func buildSwagOperations(doc *Document) map[string]map[string]*swagOperation {
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
			secReq := map[string][]string{s: {}}
			if scopes, ok := op.secsScopes[s]; ok {
				secReq[s] = scopes
			}
			securities = append(securities, secReq)
		}

		_, ok := out[op.route]
		if !ok {
			out[op.route] = make(map[string]*swagOperation, 4) // cap defaults to 4
		}
		out[op.route][method] = &swagOperation{
			Summary:     op.summary,
			OperationId: operationId,
			Description: op.desc,
			Schemes:     op.schemes,
			Consumes:    op.consumes,
			Produces:    op.produces,
			Tags:        op.tags,
			Securities:  securities,
			Deprecated:  op.deprecated,
			ExternalDoc: buildSwagExternalDoc(op.externalDoc),
			Parameters:  buildSwagParams(op.params),
			Responses:   buildSwagResponses(op.responses),
		}
	}
	return out
}

func buildSwagDefinitions(doc *Document) map[string]*swagDefinition {
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
		out[definition.name] = buildSwagDefinition(definition)
	}
	return out
}

// ========
// document
// ========

func buildSwagDocument(doc *Document) *swagDocument {
	// check
	checkDocument(doc)

	// info
	out := &swagDocument{
		Swagger:  "2.0",
		Host:     doc.host,
		BasePath: doc.basePath,
		Info: &swagInfo{
			Title:          doc.info.title,
			Version:        doc.info.version,
			Description:    doc.info.desc,
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
	if opt := doc.option; opt != nil {
		tags := make([]*swagTag, 0, len(opt.tags))
		for _, t := range opt.tags {
			tags = append(tags, &swagTag{Name: t.name, Description: t.desc, ExternalDoc: buildSwagExternalDoc(t.externalDoc)})
		}
		securities := make(map[string]*swagSecurity, len(opt.securities))
		for _, s := range opt.securities {
			if s.typ == APIKEY {
				securities[s.title] = &swagSecurity{Type: APIKEY, Description: s.desc, Name: s.name, In: s.in}
			} else if s.typ == BASIC {
				securities[s.title] = &swagSecurity{Type: BASIC, Description: s.desc}
			} else if s.typ == OAUTH2 {
				scopes := make(map[string]string, len(s.scopes))
				for _, c := range s.scopes {
					scopes[c.scope] = c.desc
				}
				securities[s.title] = &swagSecurity{Type: OAUTH2, Description: s.desc, Flow: s.flow, AuthorizationUrl: s.authorizationUrl, TokenUrl: s.tokenUrl, Scopes: scopes}
			}
		}

		out.Schemes = opt.schemes
		out.Consumes = opt.consumes
		out.Produces = opt.produces
		out.Tags = tags
		out.Securities = securities
		out.ExternalDoc = buildSwagExternalDoc(opt.externalDoc)
	}

	// definitions & operations
	out.Definitions = buildSwagDefinitions(doc)
	out.Operations = buildSwagOperations(doc)

	return out
}
