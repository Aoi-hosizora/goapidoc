package goapidoc

import (
	"gopkg.in/yaml.v2"
	"strconv"
	"strings"
)

// region swag-type

type swagDocument struct {
	Swagger     string                          `yaml:"swagger"                       json:"swagger"`
	Host        string                          `yaml:"host"                          json:"host"`
	BasePath    string                          `yaml:"basePath"                      json:"basePath"`
	Info        *swagInfo                       `yaml:"info"                          json:"info"`
	Tags        []*swagTag                      `yaml:"tags,omitempty"                json:"tags,omitempty"`
	Securities  map[string]*swagSecurity        `yaml:"securityDefinitions,omitempty" json:"securityDefinitions,omitempty"`
	Paths       map[string]map[string]*swagPath `yaml:"paths,omitempty"               json:"paths,omitempty"`
	Definitions map[string]*swagDefinition      `yaml:"definitions,omitempty"         json:"definitions,omitempty"`
}

type swagTag struct {
	Name        string `yaml:"name"                  json:"name"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
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

type swagInfo struct {
	Title          string       `yaml:"title"                    json:"title"`
	Description    string       `yaml:"description"              json:"description"`
	Version        string       `yaml:"version"                  json:"version"`
	TermsOfService string       `yaml:"termsOfService,omitempty" json:"termsOfService,omitempty"`
	License        *swagLicense `yaml:"license,omitempty"        json:"license,omitempty"`
	Contact        *swagContact `yaml:"contact,omitempty"        json:"contact,omitempty"`
}

type swagSecurity struct {
	Type string `yaml:"type" json:"type"`
	Name string `yaml:"name" json:"name"`
	In   string `yaml:"in"   json:"in"`
}

// !!!
type swagPath struct {
	Summary     string                     `yaml:"summary"               json:"summary"`
	OperationId string                     `yaml:"operationId"           json:"operationId"`
	Description string                     `yaml:"description,omitempty" json:"description,omitempty"`
	Tags        []string                   `yaml:"tags,omitempty"        json:"tags,omitempty"`
	Consumes    []string                   `yaml:"consumes,omitempty"    json:"consumes,omitempty"`
	Produces    []string                   `yaml:"produces,omitempty"    json:"produces,omitempty"`
	Securities  []map[string][]interface{} `yaml:"security,omitempty"    json:"security,omitempty"`
	Deprecated  bool                       `yaml:"deprecated,omitempty"  json:"deprecated,omitempty"`
	Parameters  []*swagParam               `yaml:"parameters,omitempty"  json:"parameters,omitempty"`
	Responses   map[string]*swagResponse   `yaml:"responses,omitempty"   json:"responses,omitempty"`
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

// !!!
type swagParam struct {
	Name            string        `yaml:"name"                      json:"name"`
	In              string        `yaml:"in"                        json:"in"`
	Required        bool          `yaml:"required"                  json:"required"`
	Type            string        `yaml:"type,omitempty"            json:"type,omitempty"`
	Description     string        `yaml:"description,omitempty"     json:"description,omitempty"`
	Format          string        `yaml:"format,omitempty"          json:"format,omitempty"`
	AllowEmptyValue bool          `yaml:"allowEmptyValue,omitempty" json:"allowEmptyValue,omitempty"`
	Default         interface{}   `yaml:"default,omitempty"         json:"default,omitempty"`
	Example         interface{}   `yaml:"example,omitempty"         json:"example,omitempty"`
	Enum            []interface{} `yaml:"enum,omitempty"            json:"enum,omitempty"`
	Maximum         int           `yaml:"maximum,omitempty"         json:"maximum,omitempty"`
	Minimum         int           `yaml:"minimum,omitempty"         json:"minimum,omitempty"`
	MaxLength       int           `yaml:"maxLength,omitempty"       json:"maxLength,omitempty"`
	MinLength       int           `yaml:"minLength,omitempty"       json:"minLength,omitempty"`
	Schema          *swagSchema   `yaml:"schema,omitempty"          json:"schema,omitempty"`
	Items           *swagItems    `yaml:"items,omitempty"           json:"items,omitempty"`
}

// !!!
type swagDefinition struct {
	Type        string         `yaml:"type"                  json:"type"`
	Required    []string       `yaml:"required"              json:"required"`
	Description string         `yaml:"description,omitempty" json:"description,omitempty"`
	Properties  *linkedHashMap `yaml:"properties,omitempty"  json:"properties,omitempty"` // map[string]*swagSchema
}

// !!!!!!!!! (include schema and property)
type swagSchema struct {
	Type            string        `yaml:"type,omitempty"            json:"type,omitempty"`
	Required        bool          `yaml:"required,omitempty"        json:"required,omitempty"`
	Description     string        `yaml:"description,omitempty"     json:"description,omitempty"`
	Format          string        `yaml:"format,omitempty"          json:"format,omitempty"`
	AllowEmptyValue bool          `yaml:"allowEmptyValue,omitempty" json:"allowEmptyValue,omitempty"`
	Default         interface{}   `yaml:"default,omitempty"         json:"default,omitempty"`
	Example         interface{}   `yaml:"example,omitempty"         json:"example,omitempty"`
	Enum            []interface{} `yaml:"enum,omitempty"            json:"enum,omitempty"`
	Maximum         int           `yaml:"maximum,omitempty"         json:"maximum,omitempty"`
	Minimum         int           `yaml:"minimum,omitempty"         json:"minimum,omitempty"`
	MaxLength       int           `yaml:"maxLength,omitempty"       json:"maxLength,omitempty"`
	MinLength       int           `yaml:"minLength,omitempty"       json:"minLength,omitempty"`

	OriginRef string     `yaml:"originRef,omitempty" json:"originRef,omitempty"`
	Ref       string     `yaml:"$ref,omitempty"      json:"$ref,omitempty"`
	Items     *swagItems `yaml:"items,omitempty"     json:"items,omitempty"`
}

// !!!!!!!!!
type swagItems struct {
	Type    string        `yaml:"type,omitempty"    json:"type,omitempty"`
	Format  string        `yaml:"format,omitempty"  json:"format,omitempty"`
	Default interface{}   `yaml:"default,omitempty" json:"default,omitempty"`
	Enum    []interface{} `yaml:"enum,omitempty"    json:"enum,omitempty"`

	OriginRef string     `yaml:"originRef,omitempty" json:"originRef,omitempty"`
	Ref       string     `yaml:"$ref,omitempty"      json:"$ref,omitempty"`
	Items     *swagItems `yaml:"items,omitempty"     json:"items,omitempty"`
}

// endregion

// region handle-type

func handleSwagObject(doc *Document, swagDoc *swagDocument, obj *apiObject) (origin string, ref string) {
	if obj == nil || obj.typ == "" {
		return "", ""
	}

	origin = obj.typ
	if len(obj.generic) != 0 {
		origin += "<"
		for _, g := range obj.generic {
			if g.kind == apiPrimeKind {
				origin += g.outPrime.typ
			} else if g.kind == apiArrayKind {
				origin += g.name
			} else if g.kind == apiObjectKind {
				newOrigin, _ := handleSwagObject(doc, swagDoc, g.outObject)
				origin += newOrigin
			}
			origin += ","
		}
		origin = origin[:len(origin)-1]
		origin += ">"

		var gdef *Definition
		for _, def := range doc.definitions {
			if def.name == obj.typ && len(def.generics) == len(obj.generic) {
				props := make([]*Property, len(def.properties))
				for idx, p := range def.properties {
					props[idx] = &Property{
						name:       p.name,
						typ:        p.typ,
						required:   p.required,
						desc:       p.desc,
						allowEmpty: p.allowEmpty,
						def:        p.def,
						enum:       p.enum,
					}
				}
				gdef = &Definition{
					name:       def.name,
					desc:       def.desc,
					generics:   def.generics,
					properties: props,
				}
				break
			}
		}
		if gdef != nil {
			for idx, g := range obj.generic {
				gActual := g.name
				gtype := gdef.generics[idx]
				for _, prop := range gdef.properties {
					if prop.typ == gtype { // T -> Type
						prop.typ = gActual
					} else if strings.Contains(prop.typ, gtype+"[]") { // T[] -> Type[]
						prop.typ = strings.ReplaceAll(prop.typ, gtype+"[]", gActual+"[]")
					} else if strings.Contains(prop.typ, "<"+gtype+">") { // <T> -> <Type>
						prop.typ = strings.ReplaceAll(prop.typ, "<"+gtype+">", "<"+gActual+">")
					}
				}
			}
			swagDoc.Definitions[origin] = mapDefinition(doc, swagDoc, gdef)
		}
	}
	ref = "#/definitions/" + origin
	return origin, ref
}

func handleSwagArray(doc *Document, swagDoc *swagDocument, arr *apiArray) *swagItems {
	if arr == nil {
		return nil
	}
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
	if arr.typ.kind == apiPrimeKind {
		return &swagItems{
			Type:   arr.typ.outPrime.typ,
			Format: arr.typ.outPrime.format,
		}
	} else if arr.typ.kind == apiArrayKind {
		return &swagItems{
			Type:  ARRAY,
			Items: handleSwagArray(doc, swagDoc, arr.typ.outArray),
		}
	} else if arr.typ.kind == apiObjectKind {
		origin, ref := handleSwagObject(doc, swagDoc, arr.typ.outObject)
		if origin != "" {
			return &swagItems{
				OriginRef: origin,
				Ref:       ref,
			}
		}
	}
	return nil
}

func mapParameterSchema(doc *Document, swagDoc *swagDocument, typ string) (string, string, *swagSchema, *swagItems) {
	it := parseApiType(typ)
	/*
		{
		  "type": "string"
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
	typeStr := ""
	formatStr := ""
	var items *swagItems
	var schema *swagSchema

	if it.kind == apiPrimeKind {
		typeStr = it.outPrime.typ
		formatStr = it.outPrime.format
	} else if it.kind == apiArrayKind {
		typeStr = ARRAY
		items = handleSwagArray(doc, swagDoc, it.outArray)
	} else if it.kind == apiObjectKind {
		origin, ref := handleSwagObject(doc, swagDoc, it.outObject)
		if origin != "" {
			schema = &swagSchema{
				OriginRef: origin,
				Ref:       ref,
			}
		}
	}
	return typeStr, formatStr, schema, items
}

func mapResponseSchema(doc *Document, swagDoc *swagDocument, typ string, req bool) *swagSchema {
	it := parseApiType(typ)
	/*
		"schema": {
		  "type": "string",
		  "required": true
		}
		"schema": {
		  "type": "array",
		  "items": {}
		}
		"schema": {
		  "originRef": "Result",
		  "$ref": "#/definitions/Result"
		}
	*/
	if it.kind == apiPrimeKind {
		return &swagSchema{
			Type:     it.outPrime.typ,
			Format:   it.outPrime.format,
			Required: req,
		}
	} else if it.kind == apiArrayKind {
		items := handleSwagArray(doc, swagDoc, it.outArray)
		return &swagSchema{
			Type:  ARRAY,
			Items: items,
		}
	} else if it.kind == apiObjectKind {
		origin, ref := handleSwagObject(doc, swagDoc, it.outObject)
		if origin != "" {
			return &swagSchema{
				OriginRef: origin,
				Ref:       ref,
			}
		}

	}
	return nil
}

func mapPropertySchema(doc *Document, swagDoc *swagDocument, typ string) (outType string, outFmt string, origin string, ref string, items *swagItems) {
	it := parseApiType(typ)
	/*
		{
		  "type": "integer"
		}
		{
		  "type": "array",
		  "items": {}
		}
		{
		  "originRef": "Page<User>",
		  "$ref": "#/definitions/Page<User>"
		}
	*/
	if it.kind == apiPrimeKind {
		outType = it.outPrime.typ
		outFmt = it.outPrime.format
	} else if it.kind == apiArrayKind {
		outType = ARRAY
		items = handleSwagArray(doc, swagDoc, it.outArray)
	} else if it.kind == apiObjectKind {
		outType = ""
		origin, ref = handleSwagObject(doc, swagDoc, it.outObject)
	}
	return
}

// endregion

// region map-func

func mapParams(doc *Document, swagDoc *swagDocument, params []*Param) []*swagParam {
	out := make([]*swagParam, len(params))
	for i, p := range params {
		t, f, schema, items := mapParameterSchema(doc, swagDoc, p.typ)
		out[i] = &swagParam{
			Name:            p.name,
			In:              p.in,
			Required:        p.required,
			Description:     p.desc,
			Type:            t,
			Format:          f,
			AllowEmptyValue: p.allowEmpty,
			Default:         p.def,
			Enum:            p.enum,
			Example:         p.example,
			Maximum:         p.maximum,
			Minimum:         p.minimum,
			MaxLength:       p.maxLength,
			MinLength:       p.minLength,
			Schema:          schema,
			Items:           items,
		}
		if p.in == BODY { // must put in schema
			origin := ""
			ref := ""
			if out[i].Schema != nil {
				origin = out[i].Schema.OriginRef
				ref = out[i].Schema.Ref
			}
			out[i].Schema = &swagSchema{
				Type:      out[i].Type,
				Format:    out[i].Format,
				OriginRef: origin,
				Ref:       ref,
				Items:     out[i].Items,
			}
			out[i].Type = ""
			out[i].Format = ""
			out[i].Items = nil
		}
	}
	return out
}

func mapResponses(doc *Document, swagDoc *swagDocument, responses []*Response) map[string]*swagResponse {
	out := make(map[string]*swagResponse)
	for _, r := range responses {
		headers := map[string]*swagHeader{}
		for _, h := range r.headers {
			headers[h.name] = &swagHeader{
				Type:        h.typ,
				Description: h.desc,
			}
		}

		out[strconv.Itoa(r.code)] = &swagResponse{
			Description: r.desc,
			Schema:      mapResponseSchema(doc, swagDoc, r.typ, r.required),
			Examples:    r.examples,
			Headers:     headers,
		}
	}
	return out
}

func mapDefinition(doc *Document, swagDoc *swagDocument, def *Definition) *swagDefinition {
	required := make([]string, 0)
	properties := newLinkedHashMap() // make(map[string]*swagSchema)
	for _, p := range def.properties {
		if p.required {
			required = append(required, p.name)
		}
		t, f, origin, ref, items := mapPropertySchema(doc, swagDoc, p.typ)
		properties.Set(p.name, &swagSchema{
			Required:        p.required,
			Description:     p.desc,
			Type:            t,
			Format:          f,
			AllowEmptyValue: p.allowEmpty,
			Default:         p.def,
			Example:         p.example,
			Enum:            p.enum,
			Maximum:         p.maximum,
			Minimum:         p.minimum,
			MaxLength:       p.maxLength,
			MinLength:       p.minLength,
			OriginRef:       origin,
			Ref:             ref,
			Items:           items,
		})
	}

	return &swagDefinition{
		Type:        "object",
		Description: def.desc,
		Required:    required,
		Properties:  properties,
	}
}

// endregion

func buildDocument(d *Document) *swagDocument {
	out := &swagDocument{
		Swagger:  "2.0",
		Host:     d.host,
		BasePath: d.basePath,
		Info: &swagInfo{
			Title:          d.info.title,
			Description:    d.info.desc,
			Version:        d.info.version,
			TermsOfService: d.info.termsOfService,
		},
		Tags:        []*swagTag{},
		Securities:  map[string]*swagSecurity{},
		Paths:       map[string]map[string]*swagPath{},
		Definitions: map[string]*swagDefinition{},
	}
	if d.info.license != nil {
		out.Info.License = &swagLicense{Name: d.info.license.name, Url: d.info.license.url}
	}
	if d.info.contact != nil {
		out.Info.Contact = &swagContact{Name: d.info.contact.name, Url: d.info.contact.url, Email: d.info.contact.email}
	}
	for _, t := range d.tags {
		out.Tags = append(out.Tags, &swagTag{
			Name:        t.name,
			Description: t.desc,
		})
	}
	for _, s := range d.securities {
		out.Securities[s.title] = &swagSecurity{
			Type: s.typ,
			Name: s.name,
			In:   s.in,
		}
	}
	if len(out.Securities) == 0 {
		out.Securities = nil
	}

	// models
	for _, def := range d.definitions {
		preHandleGeneric(def)
	}
	for idx := 0; idx < len(d.definitions); idx++ {
		def := d.definitions[idx]
		ok := true // contain generic
		for _, prop := range def.properties {
			for _, g := range def.generics {
				if prop.typ == g || strings.Contains(prop.typ, g+"[]") || strings.Contains(prop.typ, "<"+g+">") {
					ok = false
					break
				}
			}
			if !ok {
				break
			}
		}
		if ok {
			out.Definitions[def.name] = mapDefinition(d, out, def)
		}
	}

	// paths
	for _, p := range d.paths {
		_, ok := out.Paths[p.route]
		if !ok {
			out.Paths[p.route] = map[string]*swagPath{}
		}
		p.method = strings.ToLower(p.method)
		id := strings.ReplaceAll(p.route, "/", "-")
		id = strings.ReplaceAll(strings.ReplaceAll(id, "{", ""), "}", "")
		id += "-" + p.method
		securities := make([]map[string][]interface{}, len(p.securities))
		for i, s := range p.securities {
			securities[i] = map[string][]interface{}{s: {}}
		}

		out.Paths[p.route][p.method] = &swagPath{
			Summary:     p.summary,
			Description: p.desc,
			OperationId: id,
			Tags:        p.tags,
			Consumes:    p.consumes,
			Produces:    p.produces,
			Securities:  securities,
			Deprecated:  p.deprecated,
			Parameters:  mapParams(d, out, p.params),
			Responses:   mapResponses(d, out, p.responses),
		}
	}

	return out
}

func (d *Document) GenerateSwaggerYaml(path string) ([]byte, error) {
	out := buildDocument(d)
	doc, err := yaml.Marshal(out)
	if err != nil {
		return nil, err
	}

	err = saveFile(path, doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (d *Document) GenerateSwaggerJson(path string) ([]byte, error) {
	out := buildDocument(d)
	doc, err := jsonMarshal(out)
	if err != nil {
		return nil, err
	}

	err = saveFile(path, doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func GenerateSwaggerYaml(path string) ([]byte, error) {
	return _document.GenerateSwaggerYaml(path)
}

func GenerateSwaggerJson(path string) ([]byte, error) {
	return _document.GenerateSwaggerJson(path)
}
