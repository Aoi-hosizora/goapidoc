package goapidoc

import (
	"gopkg.in/yaml.v2"
	"reflect"
	"strconv"
	"strings"
)

// region inner-type

type innerDocument struct {
	Host        string                           `yaml:"host"                          json:"host"`
	BasePath    string                           `yaml:"basePath"                      json:"basePath"`
	Info        *innerInfo                       `yaml:"info"                          json:"info"`
	Tags        []*innerTag                      `yaml:"tags,omitempty"                json:"tags,omitempty"`
	Securities  map[string]*innerSecurity        `yaml:"securityDefinitions,omitempty" json:"securityDefinitions,omitempty"`
	Paths       map[string]map[string]*innerPath `yaml:"paths,omitempty"               json:"paths,omitempty"`
	Definitions map[string]*innerDefinition      `yaml:"definitions,omitempty"         json:"definitions,omitempty"`
}

type innerTag struct {
	Name        string `yaml:"name"                  json:"name"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

type innerLicense struct {
	Name string `yaml:"name"          json:"name"`
	Url  string `yaml:"url,omitempty" json:"url,omitempty"`
}

type innerContact struct {
	Name  string `yaml:"name"            json:"name"`
	Url   string `yaml:"url,omitempty"   json:"url,omitempty"`
	Email string `yaml:"email,omitempty" json:"email,omitempty"`
}

type innerInfo struct {
	Title          string        `yaml:"title"                    json:"title"`
	Description    string        `yaml:"description"              json:"description"`
	Version        string        `yaml:"version"                  json:"version"`
	TermsOfService string        `yaml:"termsOfService,omitempty" json:"termsOfService,omitempty"`
	License        *innerLicense `yaml:"license,omitempty"        json:"license,omitempty"`
	Contact        *innerContact `yaml:"contact,omitempty"        json:"contact,omitempty"`
}

type innerSecurity struct {
	Type string `yaml:"type" json:"type"`
	Name string `yaml:"name" json:"name"`
	In   string `yaml:"in"   json:"in"`
}

// !!!
type innerPath struct {
	Summary     string                     `yaml:"summary"               json:"summary"`
	OperationId string                     `yaml:"operationId"           json:"operationId"`
	Description string                     `yaml:"description,omitempty" json:"description,omitempty"`
	Tags        []string                   `yaml:"tags,omitempty"        json:"tags,omitempty"`
	Consumes    []string                   `yaml:"consumes,omitempty"    json:"consumes,omitempty"`
	Produces    []string                   `yaml:"produces,omitempty"    json:"produces,omitempty"`
	Securities  []map[string][]interface{} `yaml:"security,omitempty"    json:"security,omitempty"`
	Deprecated  bool                       `yaml:"deprecated,omitempty"  json:"deprecated,omitempty"`
	Parameters  []*innerParam              `yaml:"parameters,omitempty"  json:"parameters,omitempty"`
	Responses   map[string]*innerResponse  `yaml:"responses,omitempty"   json:"responses,omitempty"`
}

// !!!
type innerResponse struct {
	Description string                  `yaml:"description,omitempty" json:"description,omitempty"`
	Headers     map[string]*innerHeader `yaml:"headers,omitempty"     json:"headers,omitempty"`
	Examples    map[string]string       `yaml:"examples,omitempty"    json:"examples,omitempty"`
	Schema      *innerSchema            `yaml:"schema,omitempty"      json:"schema,omitempty"`
}

type innerHeader struct {
	Type        string      `yaml:"type,omitempty"        json:"type,omitempty"`
	Description string      `yaml:"description,omitempty" json:"description,omitempty"`
	Format      string      `yaml:"format,omitempty"      json:"format,omitempty"`
	Default     interface{} `yaml:"default,omitempty"     json:"default,omitempty"`
}

// !!!
type innerParam struct {
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
	Schema          *innerSchema  `yaml:"schema,omitempty"          json:"schema,omitempty"`
	Items           *innerItems   `yaml:"items,omitempty"           json:"items,omitempty"`
}

// !!!
type innerDefinition struct {
	Type        string         `yaml:"type"                  json:"type"`
	Required    []string       `yaml:"required"              json:"required"`
	Description string         `yaml:"description,omitempty" json:"description,omitempty"`
	Properties  *LinkedHashMap `yaml:"properties,omitempty"  json:"properties,omitempty"` // map[string]*innerSchema
}

// !!! (include Schema and Property)
type innerSchema struct {
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

	OriginRef string      `yaml:"originRef,omitempty" json:"originRef,omitempty"`
	Ref       string      `yaml:"$ref,omitempty"      json:"$ref,omitempty"`
	Items     *innerItems `yaml:"items,omitempty"     json:"items,omitempty"`
}

// !!!
type innerItems struct {
	Type    string        `yaml:"type,omitempty"    json:"type,omitempty"`
	Format  string        `yaml:"format,omitempty"  json:"format,omitempty"`
	Default interface{}   `yaml:"default,omitempty" json:"default,omitempty"`
	Enum    []interface{} `yaml:"enum,omitempty"    json:"enum,omitempty"`

	OriginRef string      `yaml:"originRef,omitempty" json:"originRef,omitempty"`
	Ref       string      `yaml:"$ref,omitempty"      json:"$ref,omitempty"`
	Items     *innerItems `yaml:"items,omitempty"     json:"items,omitempty"`
}

// endregion

// region handle-type

func handleInnerObject(doc *Document, innerDoc *innerDocument, obj *apiObject) (origin string, ref string) {
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
				newOrigin, _ := handleInnerObject(doc, innerDoc, g.outObject)
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
						defaultVal: p.defaultVal,
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
			innerDoc.Definitions[origin] = mapDefinition(doc, innerDoc, gdef)
		}
	}
	ref = "#/definitions/" + origin
	return origin, ref
}

func handleInnerArray(doc *Document, innerDoc *innerDocument, arr *apiArray) *innerItems {
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
		return &innerItems{
			Type:   arr.typ.outPrime.typ,
			Format: arr.typ.outPrime.format,
		}
	} else if arr.typ.kind == apiArrayKind {
		return &innerItems{
			Type:  ARRAY,
			Items: handleInnerArray(doc, innerDoc, arr.typ.outArray),
		}
	} else if arr.typ.kind == apiObjectKind {
		origin, ref := handleInnerObject(doc, innerDoc, arr.typ.outObject)
		if origin != "" {
			return &innerItems{
				OriginRef: origin,
				Ref:       ref,
			}
		}
	}
	return nil
}

func mapParameterSchema(doc *Document, innerDoc *innerDocument, typ string) (string, string, *innerSchema, *innerItems) {
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
	var items *innerItems
	var schema *innerSchema

	if it.kind == apiPrimeKind {
		typeStr = it.outPrime.typ
		formatStr = it.outPrime.format
	} else if it.kind == apiArrayKind {
		typeStr = ARRAY
		items = handleInnerArray(doc, innerDoc, it.outArray)
	} else if it.kind == apiObjectKind {
		origin, ref := handleInnerObject(doc, innerDoc, it.outObject)
		if origin != "" {
			schema = &innerSchema{
				OriginRef: origin,
				Ref:       ref,
			}
		}
	}
	return typeStr, formatStr, schema, items
}

func mapResponseSchema(doc *Document, innerDoc *innerDocument, typ string, req bool) *innerSchema {
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
		return &innerSchema{
			Type:     it.outPrime.typ,
			Format:   it.outPrime.format,
			Required: req,
		}
	} else if it.kind == apiArrayKind {
		items := handleInnerArray(doc, innerDoc, it.outArray)
		return &innerSchema{
			Type:  ARRAY,
			Items: items,
		}
	} else if it.kind == apiObjectKind {
		origin, ref := handleInnerObject(doc, innerDoc, it.outObject)
		if origin != "" {
			return &innerSchema{
				OriginRef: origin,
				Ref:       ref,
			}
		}

	}
	return nil
}

func mapPropertySchema(doc *Document, innerDoc *innerDocument, typ string) (outType string, outFmt string, origin string, ref string, items *innerItems) {
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
		items = handleInnerArray(doc, innerDoc, it.outArray)
	} else if it.kind == apiObjectKind {
		outType = ""
		origin, ref = handleInnerObject(doc, innerDoc, it.outObject)
	}
	return
}

// endregion

// region map-func

func mapParams(doc *Document, innerDoc *innerDocument, params []*Param) []*innerParam {
	out := make([]*innerParam, len(params))
	for i, p := range params {
		t, f, schema, items := mapParameterSchema(doc, innerDoc, p.typ)
		out[i] = &innerParam{
			Name:            p.name,
			In:              p.in,
			Required:        p.required,
			Description:     p.desc,
			Type:            t,
			Format:          f,
			AllowEmptyValue: p.allowEmpty,
			Default:         p.defaultVal,
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
			out[i].Schema = &innerSchema{
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

func mapResponses(doc *Document, innerDoc *innerDocument, responses []*Response) map[string]*innerResponse {
	out := make(map[string]*innerResponse)
	for _, r := range responses {
		headers := map[string]*innerHeader{}
		for _, h := range r.headers {
			headers[h.name] = &innerHeader{
				Type:        h.typ,
				Description: h.desc,
			}
		}

		out[strconv.Itoa(r.code)] = &innerResponse{
			Description: r.desc,
			Schema:      mapResponseSchema(doc, innerDoc, r.typ, r.required),
			Examples:    r.examples,
			Headers:     headers,
		}
	}
	return out
}

func mapDefinition(doc *Document, innerDoc *innerDocument, def *Definition) *innerDefinition {
	required := make([]string, 0)
	properties := NewLinkedHashMap() // make(map[string]*innerSchema)
	for _, p := range def.properties {
		if p.required {
			required = append(required, p.name)
		}
		t, f, origin, ref, items := mapPropertySchema(doc, innerDoc, p.typ)
		properties.Set(p.name, &innerSchema{
			Required:        p.required,
			Description:     p.desc,
			Type:            t,
			Format:          f,
			AllowEmptyValue: p.allowEmpty,
			Default:         p.defaultVal,
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

	return &innerDefinition{
		Type:        "object",
		Description: def.desc,
		Required:    required,
		Properties:  properties,
	}
}

// endregion

func buildDocument(d *Document) *innerDocument {
	out := &innerDocument{
		Host:     d.host,
		BasePath: d.basePath,
		Info: &innerInfo{
			Title:          d.info.title,
			Description:    d.info.desc,
			Version:        d.info.version,
			TermsOfService: d.info.termsOfService,
		},
		Tags:        []*innerTag{},
		Securities:  map[string]*innerSecurity{},
		Paths:       map[string]map[string]*innerPath{},
		Definitions: map[string]*innerDefinition{},
	}
	if d.info.license != nil {
		out.Info.License = &innerLicense{Name: d.info.license.name, Url: d.info.license.url}
	}
	if d.info.contact != nil {
		out.Info.Contact = &innerContact{Name: d.info.contact.name, Url: d.info.contact.url, Email: d.info.contact.email}
	}
	for _, t := range d.tags {
		out.Tags = append(out.Tags, &innerTag{
			Name:        t.name,
			Description: t.desc,
		})
	}
	for _, s := range d.securities {
		out.Securities[s.title] = &innerSecurity{
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
			out.Paths[p.route] = map[string]*innerPath{}
		}
		p.method = strings.ToLower(p.method)
		id := strings.ReplaceAll(p.route, "/", "-")
		id = strings.ReplaceAll(strings.ReplaceAll(id, "{", ""), "}", "")
		id += "-" + p.method
		securities := make([]map[string][]interface{}, len(p.securities))
		for i, s := range p.securities {
			securities[i] = map[string][]interface{}{s: {}}
		}

		out.Paths[p.route][p.method] = &innerPath{
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

func appendKvs(d *innerDocument, kvs map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	if kvs != nil {
		for k, v := range kvs {
			out[k] = v
		}
	}

	innerValue := reflect.ValueOf(d).Elem()
	innerType := innerValue.Type()
	for i := 0; i < innerType.NumField(); i++ {
		field := innerType.Field(i)

		tag := field.Tag.Get("yaml")
		omitempty := false
		if tag == "" {
			tag = field.Name
		} else if strings.Index(tag, ",omitempty") != -1 {
			omitempty = true
		}

		name := strings.TrimSpace(strings.Split(tag, ",")[0])
		if name != "-" && name != "" {
			value := innerValue.Field(i).Interface()
			if !omitempty || (value != nil && value != "") {
				out[name] = value
			}
		}
	}

	return out
}

func (d *Document) GenerateYaml(path string, kvs map[string]interface{}) ([]byte, error) {
	out := appendKvs(buildDocument(d), kvs)
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

func (d *Document) GenerateJson(path string, kvs map[string]interface{}) ([]byte, error) {
	out := appendKvs(buildDocument(d), kvs)
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

func (d *Document) GenerateYamlWithSwagger2(path string) ([]byte, error) {
	return d.GenerateYaml(path, map[string]interface{}{"swagger": "2.0"})
}

func (d *Document) GenerateJsonWithSwagger2(path string) ([]byte, error) {
	return d.GenerateJson(path, map[string]interface{}{"swagger": "2.0"})
}

func GenerateYaml(path string, kvs map[string]interface{}) ([]byte, error) {
	return _document.GenerateYaml(path, kvs)
}

func GenerateJson(path string, kvs map[string]interface{}) ([]byte, error) {
	return _document.GenerateJson(path, kvs)
}

func GenerateYamlWithSwagger2(path string) ([]byte, error) {
	return _document.GenerateJsonWithSwagger2(path)
}

func GenerateJsonWithSwagger2(path string) ([]byte, error) {
	return _document.GenerateJsonWithSwagger2(path)
}
