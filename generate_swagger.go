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

func buildSwaggerParameterSchema(typ string) (outType string, outFormat string, schema *swagSchema, items *swagItems) {
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
		items = buildSwaggerItems(at.array)
		return
	}
	if at.kind == apiObjectKind {
		origin := at.name
		ref := "#/definitions/" + origin
		schema = &swagSchema{OriginRef: origin, Ref: ref} // schema
		return
	}

	return
}

func buildSwaggerResponseSchema(typ string) *swagSchema {
	/*
		"schema": {
		  "type": "string"
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
		return &swagSchema{Type: at.prime.typ, Format: at.prime.format}
	}
	if at.kind == apiArrayKind {
		items := buildSwaggerItems(at.array)
		return &swagSchema{Type: ARRAY, Items: items}
	}
	if at.kind == apiObjectKind {
		origin := at.name
		ref := "#/definitions/" + origin
		return &swagSchema{OriginRef: origin, Ref: ref}
	}

	return nil
}

func buildSwaggerPropertySchema(typ string) (outType string, outFmt string, origin string, ref string, items *swagItems) {
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

func buildSwaggerParameters(params []*Param) []*swagParam {
	out := make([]*swagParam, len(params))
	for i, p := range params {
		typ, format, schema, items := buildSwaggerParameterSchema(p.typ)
		param := &swagParam{
			Name:            p.name,
			In:              p.in,
			Required:        p.required,
			Description:     p.desc,
			Type:            typ,
			Format:          format,
			AllowEmptyValue: p.allowEmpty,
			Default:         p.dft,
			Enum:            p.enums,
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
			if param.Schema != nil {
				origin = param.Schema.OriginRef
				ref = param.Schema.Ref
			}
			param.Schema = &swagSchema{
				Type:      param.Type,
				Format:    param.Format,
				OriginRef: origin,
				Ref:       ref,
				Items:     param.Items,
			}
			param.Type = ""
			param.Format = ""
			param.Items = nil
		}

		out[i] = param
	}
	return out
}

func buildSwaggerResponses(responses []*Response) map[string]*swagResponse {
	out := make(map[string]*swagResponse)
	for _, r := range responses {
		headers := make(map[string]*swagHeader)
		for _, h := range r.headers {
			headers[h.name] = &swagHeader{Type: h.typ, Description: h.desc}
		}

		code := strconv.Itoa(r.code)
		out[code] = &swagResponse{
			Description: r.desc,
			Schema:      buildSwaggerResponseSchema(r.typ),
			Examples:    r.examples,
			Headers:     headers,
		}
	}
	return out
}

func buildSwaggerDefinition(dft *Definition) *swagDefinition {
	required := make([]string, 0)
	properties := newLinkedHashMap() // make(map[string]*swagSchema)
	for _, p := range dft.properties {
		if p.required {
			required = append(required, p.name)
		}

		typ, format, origin, ref, items := buildSwaggerPropertySchema(p.typ)
		properties.Set(p.name, &swagSchema{
			Required:        p.required,
			Description:     p.desc,
			Type:            typ,
			Format:          format,
			AllowEmptyValue: p.allowEmpty,
			Default:         p.dft,
			Example:         p.example,
			Enum:            p.enums,
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
		Description: dft.desc,
		Required:    required,
		Properties:  properties,
	}
}

func buildSwaggerPaths(doc *Document) map[string]map[string]*swagPath {
	out := make(map[string]map[string]*swagPath)
	for _, p := range doc.paths {
		_, ok := out[p.route]
		if !ok {
			out[p.route] = make(map[string]*swagPath)
		}

		method := strings.ToLower(p.method)
		id := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(p.route, "/", "-"), "{", ""), "}", "") + "-" + method
		securities := make([]map[string][]interface{}, 0)
		for _, s := range p.securities {
			securities = append(securities, map[string][]interface{}{s: {}})
		}

		out[p.route][method] = &swagPath{
			Summary:     p.summary,
			Description: p.desc,
			OperationId: id,
			Tags:        p.tags,
			Securities:  securities,
			Consumes:    p.consumes,
			Produces:    p.produces,
			Deprecated:  p.deprecated,
			Parameters:  buildSwaggerParameters(p.params),
			Responses:   buildSwaggerResponses(p.responses),
		}
	}
	return out
}

func buildSwaggerDefinitions(doc *Document) map[string]*swagDefinition {
	propertyTypes := make([]string, 0)
	for _, definition := range doc.definitions {
		prehandleGenericName(definition) // new name

		if len(definition.generics) == 0 && len(definition.properties) > 0 {
			for _, property := range definition.properties {
				propertyTypes = append(propertyTypes, property.typ)
			}
		}
	}
	for _, path := range doc.paths {
		for _, param := range path.params {
			propertyTypes = append(propertyTypes, param.typ)
		}
		for _, response := range path.responses {
			propertyTypes = append(propertyTypes, response.typ)
		}
	}
	definitions := prehandleGenericList(doc.definitions, propertyTypes) // new list

	out := make(map[string]*swagDefinition)
	for _, dft := range definitions {
		if len(dft.generics) == 0 {
			out[dft.name] = buildSwaggerDefinition(dft)
		}
	}
	return out
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
	out.Definitions = buildSwaggerDefinitions(d)

	// paths
	out.Paths = buildSwaggerPaths(d)

	return out
}

// GenerateSwaggerYaml generates swagger yaml script and writes into file.
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

// GenerateSwaggerJson generates swagger json script and writes into file.
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

// GenerateSwaggerYaml generates swagger yaml script and writes into file.
func GenerateSwaggerYaml(path string) ([]byte, error) {
	return _document.GenerateSwaggerYaml(path)
}

// GenerateSwaggerJson generates swagger json script and writes into file.
func GenerateSwaggerJson(path string) ([]byte, error) {
	return _document.GenerateSwaggerJson(path)
}
