package goapidoc

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

func handleInnerObject(doc *Document, innerDoc *innerDocument, obj *innerObject) (origin string, ref string) {
	if obj == nil || obj.Type == "" {
		return "", ""
	}

	origin = obj.Type
	if len(obj.Generic) != 0 {
		origin += "<"
		for _, g := range obj.Generic {
			if g.Kind == innerPrimeKind {
				origin += g.OutPrime.Type
			} else if g.Kind == innerArrayKind {
				origin += g.Name
			} else if g.Kind == innerObjectKind {
				newOrigin, _ := handleInnerObject(doc, innerDoc, g.OutObject)
				origin += newOrigin
			}
			origin += ","
		}
		origin = origin[:len(origin)-1]
		origin += ">"

		var gdef *Definition
		for _, def := range doc.Definitions {
			if def.Name == obj.Type && len(def.Generics) == len(obj.Generic) {
				props := make([]*Property, len(def.Properties))
				for idx, p := range def.Properties {
					props[idx] = &Property{
						Name:            p.Name,
						Type:            p.Type,
						Required:        p.Required,
						Description:     p.Description,
						AllowEmptyValue: p.AllowEmptyValue,
						Default:         p.Default,
						Enum:            p.Enum,
					}
				}
				gdef = &Definition{
					Name:        def.Name,
					Description: def.Description,
					Generics:    def.Generics,
					Properties:  props,
				}
				break
			}
		}
		if gdef != nil {
			for idx, g := range obj.Generic {
				gActual := g.Name
				gtype := gdef.Generics[idx]
				for _, prop := range gdef.Properties {
					if prop.Type == gtype { // T -> Type
						prop.Type = gActual
					} else if strings.Contains(prop.Type, gtype+"[]") { // T[] -> Type[]
						prop.Type = strings.ReplaceAll(prop.Type, gtype+"[]", gActual+"[]")
					} else if strings.Contains(prop.Type, "<"+gtype+">") { // <T> -> <Type>
						prop.Type = strings.ReplaceAll(prop.Type, "<"+gtype+">", "<"+gActual+">")
					}
				}
			}
			innerDoc.Definitions[origin] = mapDefinition(doc, innerDoc, gdef)
		}
	}
	ref = "#/definitions/" + origin
	return origin, ref
}

func handleInnerArray(doc *Document, innerDoc *innerDocument, arr *innerArray) *innerItems {
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
	if arr.Type.Kind == innerPrimeKind {
		return &innerItems{
			Type:   arr.Type.OutPrime.Type,
			Format: arr.Type.OutPrime.Format,
		}
	} else if arr.Type.Kind == innerArrayKind {
		return &innerItems{
			Type:  ARRAY,
			Items: handleInnerArray(doc, innerDoc, arr.Type.OutArray),
		}
	} else if arr.Type.Kind == innerObjectKind {
		origin, ref := handleInnerObject(doc, innerDoc, arr.Type.OutObject)
		if origin != "" {
			return &innerItems{
				OriginRef: origin,
				Ref:       ref,
			}
		}
	}
	return nil
}

func mapParameterSchema(doc *Document, innerDoc *innerDocument, t string) (string, string, *innerSchema, *innerItems) {
	it := parseInnerType(t)
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

	if it.Kind == innerPrimeKind {
		typeStr = it.OutPrime.Type
		formatStr = it.OutPrime.Format
	} else if it.Kind == innerArrayKind {
		typeStr = ARRAY
		items = handleInnerArray(doc, innerDoc, it.OutArray)
	} else if it.Kind == innerObjectKind {
		origin, ref := handleInnerObject(doc, innerDoc, it.OutObject)
		if origin != "" {
			schema = &innerSchema{
				OriginRef: origin,
				Ref:       ref,
			}
		}
	}
	return typeStr, formatStr, schema, items
}

func mapResponseSchema(doc *Document, innerDoc *innerDocument, t string, req bool) *innerSchema {
	it := parseInnerType(t)
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
	if it.Kind == innerPrimeKind {
		return &innerSchema{
			Type:     it.OutPrime.Type,
			Format:   it.OutPrime.Format,
			Required: req,
		}
	} else if it.Kind == innerArrayKind {
		items := handleInnerArray(doc, innerDoc, it.OutArray)
		return &innerSchema{
			Type:  ARRAY,
			Items: items,
		}
	} else if it.Kind == innerObjectKind {
		origin, ref := handleInnerObject(doc, innerDoc, it.OutObject)
		if origin != "" {
			return &innerSchema{
				OriginRef: origin,
				Ref:       ref,
			}
		}

	}
	return nil
}

func mapPropertySchema(doc *Document, innerDoc *innerDocument, t string) (outType string, outFmt string, origin string, ref string, items *innerItems) {
	it := parseInnerType(t)
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
	if it.Kind == innerPrimeKind {
		outType = it.OutPrime.Type
		outFmt = it.OutPrime.Format
	} else if it.Kind == innerArrayKind {
		outType = ARRAY
		items = handleInnerArray(doc, innerDoc, it.OutArray)
	} else if it.Kind == innerObjectKind {
		outType = ""
		origin, ref = handleInnerObject(doc, innerDoc, it.OutObject)
	}
	return
}

// endregion

// region map-func

func mapParams(doc *Document, innerDoc *innerDocument, params []*Param) []*innerParam {
	out := make([]*innerParam, len(params))
	for i, p := range params {
		t, f, schema, items := mapParameterSchema(doc, innerDoc, p.Type)
		out[i] = &innerParam{
			Name:            p.Name,
			In:              p.In,
			Required:        p.Required,
			Description:     p.Description,
			Type:            t,
			Format:          f,
			AllowEmptyValue: p.AllowEmptyValue,
			Default:         p.Default,
			Enum:            p.Enum,
			Example:         p.Example,
			Maximum:         p.Maximum,
			Minimum:         p.Minimum,
			MaxLength:       p.MaxLength,
			MinLength:       p.MinLength,
			Schema:          schema,
			Items:           items,
		}
		if p.In == BODY { // must put in schema
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
		for _, h := range r.Headers {
			headers[h.Name] = &innerHeader{
				Type:        h.Type,
				Description: h.Description,
			}
		}

		out[strconv.Itoa(r.Code)] = &innerResponse{
			Description: r.Description,
			Schema:      mapResponseSchema(doc, innerDoc, r.Type, r.Required),
			Examples:    r.Examples,
			Headers:     headers,
		}
	}
	return out
}

func mapDefinition(doc *Document, innerDoc *innerDocument, def *Definition) *innerDefinition {
	required := make([]string, 0)
	properties := NewLinkedHashMap() // make(map[string]*innerSchema)
	for _, p := range def.Properties {
		if p.Required {
			required = append(required, p.Name)
		}
		t, f, origin, ref, items := mapPropertySchema(doc, innerDoc, p.Type)
		properties.Set(p.Name, &innerSchema{
			Required:        p.Required,
			Description:     p.Description,
			Type:            t,
			Format:          f,
			AllowEmptyValue: p.AllowEmptyValue,
			Default:         p.Default,
			Example:         p.Example,
			Enum:            p.Enum,
			Maximum:         p.Maximum,
			Minimum:         p.Minimum,
			MaxLength:       p.MaxLength,
			MinLength:       p.MinLength,
			OriginRef:       origin,
			Ref:             ref,
			Items:           items,
		})
	}

	return &innerDefinition{
		Type:        "object",
		Description: def.Description,
		Required:    required,
		Properties:  properties,
	}
}

// endregion

func buildDocument(d *Document) *innerDocument {
	out := &innerDocument{
		Host:     d.Host,
		BasePath: d.BasePath,
		Info: &innerInfo{
			Title:          d.Info.Title,
			Description:    d.Info.Description,
			Version:        d.Info.Version,
			TermsOfService: d.Info.TermsOfService,
		},
		Tags:        []*innerTag{},
		Securities:  map[string]*innerSecurity{},
		Paths:       map[string]map[string]*innerPath{},
		Definitions: map[string]*innerDefinition{},
	}
	if d.Info.License != nil {
		out.Info.License = &innerLicense{Name: d.Info.License.Name, Url: d.Info.License.Url}
	}
	if d.Info.Contact != nil {
		out.Info.Contact = &innerContact{Name: d.Info.Contact.Name, Url: d.Info.Contact.Url, Email: d.Info.Contact.Email}
	}
	for _, t := range d.Tags {
		out.Tags = append(out.Tags, &innerTag{
			Name:        t.Name,
			Description: t.Description,
		})
	}
	for _, s := range d.Securities {
		out.Securities[s.Title] = &innerSecurity{
			Type: s.Type,
			Name: s.Name,
			In:   s.In,
		}
	}
	if len(out.Securities) == 0 {
		out.Securities = nil
	}

	// models
	for _, def := range d.Definitions {
		preHandleGeneric(def)
	}
	for idx := 0; idx < len(d.Definitions); idx++ {
		def := d.Definitions[idx]
		ok := true // contain generic
		for _, prop := range def.Properties {
			for _, g := range def.Generics {
				if prop.Type == g || strings.Contains(prop.Type, g+"[]") || strings.Contains(prop.Type, "<"+g+">") {
					ok = false
					break
				}
			}
			if !ok {
				break
			}
		}
		if ok {
			out.Definitions[def.Name] = mapDefinition(d, out, def)
		}
	}

	// paths
	for _, p := range d.Paths {
		_, ok := out.Paths[p.Route]
		if !ok {
			out.Paths[p.Route] = map[string]*innerPath{}
		}
		p.Method = strings.ToLower(p.Method)
		id := strings.ReplaceAll(p.Route, "/", "-")
		id = strings.ReplaceAll(strings.ReplaceAll(id, "{", ""), "}", "")
		id += "-" + p.Method
		securities := make([]map[string][]interface{}, len(p.Securities))
		for i, s := range p.Securities {
			securities[i] = map[string][]interface{}{s: {}}
		}

		out.Paths[p.Route][p.Method] = &innerPath{
			Summary:     p.Summary,
			Description: p.Description,
			OperationId: id,
			Tags:        p.Tags,
			Consumes:    p.Consumes,
			Produces:    p.Produces,
			Securities:  securities,
			Deprecated:  p.Deprecated,
			Parameters:  mapParams(d, out, p.Params),
			Responses:   mapResponses(d, out, p.Responses),
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

func saveFile(path string, data []byte) error {
	err := ioutil.WriteFile(path, data, 0777)
	if err != nil {
		return err
	}
	return nil
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

// stop json to escape
func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
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

// noinspection GoUnusedExportedFunction
func GenerateYaml(path string, kvs map[string]interface{}) ([]byte, error) {
	return _document.GenerateYaml(path, kvs)
}

// noinspection GoUnusedExportedFunction
func GenerateJson(path string, kvs map[string]interface{}) ([]byte, error) {
	return _document.GenerateJson(path, kvs)
}

// noinspection GoUnusedExportedFunction
func GenerateYamlWithSwagger2(path string) ([]byte, error) {
	return _document.GenerateJsonWithSwagger2(path)
}

// noinspection GoUnusedExportedFunction
func GenerateJsonWithSwagger2(path string) ([]byte, error) {
	return _document.GenerateJsonWithSwagger2(path)
}
