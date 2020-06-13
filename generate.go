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
	_Additions  map[string]*innerDefinition      `yaml:"-"                             json:"-"`
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
	Enum            []interface{} `yaml:"enum,omitempty"            json:"enum,omitempty"`
	Schema          *innerSchema  `yaml:"schema,omitempty"          json:"schema,omitempty"`
	Items           *innerItems   `yaml:"items,omitempty"           json:"items,omitempty"`
}

// !!!
type innerDefinition struct {
	Type        string                  `json:"type"                  json:"type"`
	Required    []string                `json:"required"              json:"required"`
	Description string                  `json:"description,omitempty" json:"description,omitempty"`
	Properties  map[string]*innerSchema `json:"properties,omitempty"  json:"properties,omitempty"`
}

// !!! (include Schema and Property)
type innerSchema struct {
	Type            string        `yaml:"type,omitempty"            json:"type,omitempty"`
	Required        bool          `yaml:"required,omitempty"        json:"required,omitempty"`
	Description     string        `yaml:"description,omitempty"     json:"description,omitempty"`
	Format          string        `yaml:"format,omitempty"          json:"format,omitempty"`
	AllowEmptyValue bool          `yaml:"allowEmptyValue,omitempty" json:"allowEmptyValue,omitempty"`
	Default         interface{}   `yaml:"default,omitempty"         json:"default,omitempty"`
	Enum            []interface{} `yaml:"enum,omitempty"            json:"enum,omitempty"`

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

func handleInnerObject(obj *innerObject) (origin string, ref string) {
	if obj == nil || obj.Type == "" {
		return "", ""
	}
	/*
		"schema": {
		  "originRef": "User",
		  "$ref": "#/definitions/User"
		}
	*/
	origin = obj.Type
	ref = "#/definitions/" + obj.Type
	return origin, ref
}

func handleInnerArray(arr *innerArray) *innerItems {
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
	if arr.Type.Type == innerBaseType {
		return &innerItems{
			Type:   arr.Type.OutType,
			Format: defaultFormat(arr.Type.OutType),
		}
	} else if arr.Type.Type == innerArrayType {
		return &innerItems{
			Type:  ARRAY,
			Items: handleInnerArray(arr.Type.OutItems),
		}
	} else if arr.Type.Type == innerObjectType {
		origin, ref := handleInnerObject(arr.Type.OutSchema)
		if origin != "" {
			return &innerItems{
				OriginRef: origin,
				Ref:       ref,
			}
		}
	}
	return nil
}

func mapParameterSchema(t string) (string, *innerSchema, *innerItems) {
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
	var items *innerItems
	var schema *innerSchema

	if it.Type == innerBaseType {
		typeStr = it.OutType
	} else if it.Type == innerArrayType {
		typeStr = ARRAY
		items = handleInnerArray(it.OutItems)
	} else if it.Type == innerObjectType {
		typeStr = ""
		origin, ref := handleInnerObject(it.OutSchema)
		if origin != "" {
			schema = &innerSchema{
				OriginRef: origin,
				Ref:       ref,
			}
		}
	}
	return typeStr, schema, items
}

func mapResponseSchema(t string, req bool) *innerSchema {
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
	if it.Type == innerBaseType {
		return &innerSchema{
			Type:     it.OutType,
			Required: req,
		}
	} else if it.Type == innerArrayType {
		items := handleInnerArray(it.OutItems)
		return &innerSchema{
			Type:  ARRAY,
			Items: items,
		}
	} else if it.Type == innerObjectType {
		origin, ref := handleInnerObject(it.OutSchema)
		if origin != "" {
			return &innerSchema{
				OriginRef: origin,
				Ref:       ref,
			}
		}

	}
	return nil
}

func mapPropertySchema(t string) (outType string, origin string, ref string, items *innerItems) {
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
	if it.Type == innerBaseType && it.OutType != "" {
		outType = it.OutType
	} else if it.Type == innerArrayType && it.OutItems != nil {
		outType = ARRAY
		items = handleInnerArray(it.OutItems)
	} else if it.Type == innerObjectType && it.OutSchema != nil {
		outType = ""
		origin, ref = handleInnerObject(it.OutSchema)
	}
	return
}

// endregion

// region map-func

func mapParams(params []*Param) []*innerParam {
	out := make([]*innerParam, len(params))
	for i, p := range params {
		t, schema, items := mapParameterSchema(p.Type)
		out[i] = &innerParam{
			Name:            p.Name,
			In:              p.In,
			Required:        p.Required,
			Description:     p.Description,
			Type:            t,
			Format:          p.Format,
			AllowEmptyValue: p.AllowEmptyValue,
			Default:         p.Default,
			Enum:            p.Enum,
			Schema:          schema,
			Items:           items,
		}
	}
	return out
}

func mapResponses(responses []*Response) map[string]*innerResponse {
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
			Schema:      mapResponseSchema(r.Type, r.Required),
			Examples:    r.Examples,
			Headers:     headers,
		}
	}
	return out
}

func mapDefinition(def *Definition) *innerDefinition {
	required := make([]string, 0)
	properties := make(map[string]*innerSchema)
	for _, p := range def.Properties {
		if p.Required {
			required = append(required, p.Name)
		}
		t, origin, ref, items := mapPropertySchema(p.Type)
		properties[p.Name] = &innerSchema{
			Required:        p.Required,
			Description:     p.Description,
			Type:            t,
			Format:          p.Format,
			AllowEmptyValue: p.AllowEmptyValue,
			Default:         p.Default,
			Enum:            p.Enum,
			OriginRef:       origin,
			Ref:             ref,
			Items:           items,
		}
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
	if len(d.Securities) == 0 {
		d.Securities = nil
	} else {
		for _, s := range d.Securities {
			out.Securities[s.Title] = &innerSecurity{
				Type: s.Type,
				Name: s.Name,
				In:   s.In,
			}
		}
	}

	// models
	for _, d := range d.Definitions {
		out.Definitions[d.Name] = mapDefinition(d)
	}

	// paths
	for _, p := range d.Paths {
		_, ok := out.Paths[p.Route]
		if !ok {
			out.Paths[p.Route] = map[string]*innerPath{}
		}
		p.Method = strings.ToLower(p.Method)
		id := strings.ReplaceAll(p.Route, "/", "-")
		id = strings.ReplaceAll(strings.ReplaceAll(id, "{", "-"), "}", "-")
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
			Parameters:  mapParams(p.Params),
			Responses:   mapResponses(p.Responses),
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

func (d *Document) GenerateYaml(path string, kvs map[string]interface{}) error {
	out := appendKvs(buildDocument(d), kvs)
	doc, err := yaml.Marshal(out)
	if err != nil {
		return err
	}

	err = saveFile(path, doc)
	if err != nil {
		return err
	}
	return nil
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

func (d *Document) GenerateJson(path string, kvs map[string]interface{}) error {
	out := appendKvs(buildDocument(d), kvs)
	doc, err := jsonMarshal(out)
	if err != nil {
		return err
	}

	err = saveFile(path, doc)
	if err != nil {
		return err
	}
	return nil
}

func GenerateYaml(path string, kvs map[string]interface{}) error {
	return _document.GenerateYaml(path, kvs)
}

func GenerateJson(path string, kvs map[string]interface{}) error {
	return _document.GenerateJson(path, kvs)
}
