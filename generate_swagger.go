package goapidoc

import (
	"gopkg.in/yaml.v2"
	"strconv"
	"strings"
)

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
	Type        string         `yaml:"type"                  json:"type"`
	Required    []string       `yaml:"required"              json:"required"`
	Description string         `yaml:"description,omitempty" json:"description,omitempty"`
	Properties  *linkedHashMap `yaml:"properties,omitempty"  json:"properties,omitempty"` // map[string]*swagSchema
}

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

type swagItems struct {
	Type    string        `yaml:"type,omitempty"    json:"type,omitempty"`
	Format  string        `yaml:"format,omitempty"  json:"format,omitempty"`
	Default interface{}   `yaml:"default,omitempty" json:"default,omitempty"`
	Enum    []interface{} `yaml:"enum,omitempty"    json:"enum,omitempty"`

	OriginRef string     `yaml:"originRef,omitempty" json:"originRef,omitempty"`
	Ref       string     `yaml:"$ref,omitempty"      json:"$ref,omitempty"`
	Items     *swagItems `yaml:"items,omitempty"     json:"items,omitempty"`
}

func handleSwaggerObject(doc *Document, swagDoc *swagDocument, obj *apiObject) (origin string, ref string) {
	if obj == nil || obj.typ == "" {
		return "", ""
	}

	origin = obj.typ
	if len(obj.generics) != 0 {
		origin += "<"
		for _, g := range obj.generics {
			if g.kind == apiPrimeKind {
				origin += g.prime.typ
			} else if g.kind == apiArrayKind {
				origin += g.name
			} else if g.kind == apiObjectKind {
				newOrigin, _ := handleSwaggerObject(doc, swagDoc, g.object)
				origin += newOrigin
			}
			origin += ","
		}
		origin = origin[:len(origin)-1]
		origin += ">"

		var gdef *Definition
		for _, def := range doc.definitions {
			if def.name == obj.typ && len(def.generics) == len(obj.generics) {
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
						example:    p.example,
						minLength:  p.minLength,
						maxLength:  p.maxLength,
						minimum:    p.minimum,
						maximum:    p.maximum,
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
			for idx, g := range obj.generics {
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
			swagDoc.Definitions[origin] = buildSwaggerDefinition(doc, swagDoc, gdef)
		}
	}
	ref = "#/definitions/" + origin
	return origin, ref
}

func handleSwaggerArray(doc *Document, swagDoc *swagDocument, arr *apiArray) *swagItems {
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
	if arr.item.kind == apiPrimeKind {
		return &swagItems{
			Type:   arr.item.prime.typ,
			Format: arr.item.prime.format,
		}
	} else if arr.item.kind == apiArrayKind {
		return &swagItems{
			Type:  ARRAY,
			Items: handleSwaggerArray(doc, swagDoc, arr.item.array),
		}
	} else if arr.item.kind == apiObjectKind {
		origin, ref := handleSwaggerObject(doc, swagDoc, arr.item.object)
		if origin != "" {
			return &swagItems{
				OriginRef: origin,
				Ref:       ref,
			}
		}
	}
	return nil
}

func buildSwaggerParameterSchema(doc *Document, swagDoc *swagDocument, typ string) (outType string, outFormat string, schema *swagSchema, items *swagItems) {
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
	at := parseApiType(typ)

	if at.kind == apiPrimeKind {
		outType = at.prime.typ
		outFormat = at.prime.format
		return
	}
	if at.kind == apiArrayKind {
		outType = ARRAY
		items = handleSwaggerArray(doc, swagDoc, at.array)
		return
	}
	if at.kind == apiObjectKind {
		origin, ref := handleSwaggerObject(doc, swagDoc, at.object)
		if origin != "" {
			schema = &swagSchema{OriginRef: origin, Ref: ref} // schema
			return
		}
	}

	return
}

func buildSwaggerResponseSchema(doc *Document, swagDoc *swagDocument, typ string, req bool) *swagSchema {
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
	at := parseApiType(typ)

	if at.kind == apiPrimeKind {
		return &swagSchema{Type: at.prime.typ, Format: at.prime.format, Required: req}
	}
	if at.kind == apiArrayKind {
		items := handleSwaggerArray(doc, swagDoc, at.array)
		return &swagSchema{Type: ARRAY, Items: items}
	}
	if at.kind == apiObjectKind {
		origin, ref := handleSwaggerObject(doc, swagDoc, at.object)
		if origin != "" {
			return &swagSchema{OriginRef: origin, Ref: ref}
		}
	}

	return nil
}

func buildSwaggerPropertySchema(doc *Document, swagDoc *swagDocument, typ string) (outType string, outFmt string, origin string, ref string, items *swagItems) {
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
	at := parseApiType(typ)

	if at.kind == apiPrimeKind {
		outType = at.prime.typ
		outFmt = at.prime.format
		return
	}
	if at.kind == apiArrayKind {
		outType = ARRAY
		items = handleSwaggerArray(doc, swagDoc, at.array)
		return
	}
	if at.kind == apiObjectKind {
		origin, ref = handleSwaggerObject(doc, swagDoc, at.object) // ref
		return
	}

	return
}

func buildSwaggerParameters(doc *Document, swagDoc *swagDocument, params []*Param) []*swagParam {
	out := make([]*swagParam, len(params))
	for i, p := range params {
		t, f, schema, items := buildSwaggerParameterSchema(doc, swagDoc, p.typ)
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

func buildSwaggerResponses(doc *Document, swagDoc *swagDocument, responses []*Response) map[string]*swagResponse {
	out := make(map[string]*swagResponse)
	for _, r := range responses {
		headers := make(map[string]*swagHeader)
		for _, h := range r.headers {
			headers[h.name] = &swagHeader{Type: h.typ, Description: h.desc}
		}

		code := strconv.Itoa(r.code)
		out[code] = &swagResponse{
			Description: r.desc,
			Schema:      buildSwaggerResponseSchema(doc, swagDoc, r.typ, r.required),
			Examples:    r.examples,
			Headers:     headers,
		}
	}
	return out
}

func buildSwaggerDefinition(doc *Document, swagDoc *swagDocument, def *Definition) *swagDefinition {
	required := make([]string, 0)
	properties := newLinkedHashMap() // make(map[string]*swagSchema)
	for _, p := range def.properties {
		if p.required {
			required = append(required, p.name)
		}

		t, f, origin, ref, items := buildSwaggerPropertySchema(doc, swagDoc, p.typ)
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

func buildSwaggerPaths(doc *Document, swagDoc *swagDocument) {
	for _, p := range doc.paths {
		_, ok := swagDoc.Paths[p.route]
		if !ok {
			swagDoc.Paths[p.route] = make(map[string]*swagPath)
		}

		method := strings.ToLower(p.method)
		id := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(p.route, "/", "-"), "{", ""), "}", "") + "-" + method
		securities := make([]map[string][]interface{}, 0)
		for _, s := range p.securities {
			securities = append(securities, map[string][]interface{}{s: {}})
		}

		swagDoc.Paths[p.route][method] = &swagPath{
			Summary:     p.summary,
			Description: p.desc,
			OperationId: id,
			Tags:        p.tags,
			Securities:  securities,
			Consumes:    p.consumes,
			Produces:    p.produces,
			Deprecated:  p.deprecated,
			Parameters:  buildSwaggerParameters(doc, swagDoc, p.params),
			Responses:   buildSwaggerResponses(doc, swagDoc, p.responses),
		}
	}
}

func buildSwaggerDefinitions(doc *Document, swagDoc *swagDocument) {
	for _, def := range doc.definitions {
		prehandleGenericName(def)
	}

	for _, def := range doc.definitions {
		if len(def.generics) == 0 {
			swagDoc.Definitions[def.name] = buildSwaggerDefinition(doc, swagDoc, def)
		}
	}
}

func buildSwaggerDocument(d *Document) *swagDocument {
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
		Paths:       map[string]map[string]*swagPath{},
		Definitions: map[string]*swagDefinition{},
	}
	if d.info.license != nil {
		out.Info.License = &swagLicense{Name: d.info.license.name, Url: d.info.license.url}
	}
	if d.info.contact != nil {
		out.Info.Contact = &swagContact{Name: d.info.contact.name, Url: d.info.contact.url, Email: d.info.contact.email}
	}
	if len(d.tags) > 0 {
		tags := make([]*swagTag, 0)
		for _, t := range d.tags {
			tags = append(tags, &swagTag{Name: t.name, Description: t.desc})
		}
		out.Tags = tags
	}
	if len(d.securities) > 0 {
		securities := make(map[string]*swagSecurity)
		for _, s := range d.securities {
			securities[s.title] = &swagSecurity{Type: s.typ, Name: s.name, In: s.in}
		}
		out.Securities = securities
	}

	// definitions
	buildSwaggerDefinitions(d, out)

	// paths
	buildSwaggerPaths(d, out)

	return out
}

func (d *Document) GenerateSwaggerYaml(path string) ([]byte, error) {
	swagDoc := buildSwaggerDocument(d)
	bs, err := yaml.Marshal(swagDoc)
	if err != nil {
		return nil, err
	}

	err = saveFile(path, bs)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (d *Document) GenerateSwaggerJson(path string) ([]byte, error) {
	swagDoc := buildSwaggerDocument(d)
	bs, err := jsonMarshal(swagDoc)
	if err != nil {
		return nil, err
	}

	err = saveFile(path, bs)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func GenerateSwaggerYaml(path string) ([]byte, error) {
	return _document.GenerateSwaggerYaml(path)
}

func GenerateSwaggerJson(path string) ([]byte, error) {
	return _document.GenerateSwaggerJson(path)
}
