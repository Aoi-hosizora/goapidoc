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

// region map-func

func mapSchema(doc *innerDocument, schema *Schema) *innerSchema {
	if schema == nil {
		return nil
	}
	if schema.Ref != "" {
		if len(schema.Options) == 0 {
			return &innerSchema{OriginRef: schema.Ref, Ref: "#/definitions/" + schema.Ref, Required: schema.Required}
		}
		newRef := mapRefOptions(doc, schema.Ref, schema.Options)
		return &innerSchema{OriginRef: newRef, Ref: "#/definitions/" + newRef, Required: schema.Required}
	}
	return &innerSchema{
		Type:            schema.Type,
		Required:        schema.Required,
		Description:     schema.Description,
		Format:          schema.Format,
		AllowEmptyValue: schema.AllowEmptyValue,
		Default:         schema.Default,
		Enum:            schema.Enum,
		Items:           mapItems(doc, schema.Items),
	}
}

func mapItems(doc *innerDocument, items *Items) *innerItems {
	if items == nil {
		return nil
	}
	if items.Ref != "" {
		if len(items.Options) == 0 {
			return &innerItems{OriginRef: items.Ref, Ref: "#/definitions/" + items.Ref}
		}
		newRef := mapRefOptions(doc, items.Ref, items.Options)
		return &innerItems{OriginRef: newRef, Ref: "#/definitions/" + newRef}
	}
	return &innerItems{
		Type:    items.Type,
		Format:  items.Format,
		Default: items.Default,
		Enum:    items.Enum,
		Items:   mapItems(doc, items.Items),
	}
}

func mapParams(doc *innerDocument, params []*Param) []*innerParam {
	out := make([]*innerParam, len(params))
	for i, p := range params {
		out[i] = &innerParam{
			Name:            p.Name,
			In:              p.In,
			Required:        p.Required,
			Description:     p.Description,
			Type:            p.Type,
			Format:          p.Format,
			AllowEmptyValue: p.AllowEmptyValue,
			Default:         p.Default,
			Enum:            p.Enum,
			Schema:          mapSchema(doc, p.Schema),
			Items:           mapItems(doc, p.Items),
		}
	}
	return out
}

func mapResponses(doc *innerDocument, responses []*Response) map[string]*innerResponse {
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
			Schema:      mapSchema(doc, r.Schema),
			Examples:    r.Examples,
			Headers:     headers,
		}
	}
	return out
}

func mapDefinition(doc *innerDocument, definition *Definition) *innerDefinition {
	required := make([]string, 0)
	schemas := make(map[string]*innerSchema)
	for _, p := range definition.Properties {
		if p.Required {
			required = append(required, p.Title)
		}
		schemas[p.Title] = mapSchema(doc, p.Schema)
	}

	return &innerDefinition{
		Type:        "object",
		Description: definition.Description,
		Required:    required,
		Properties:  schemas,
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
		out.Definitions[d.Name] = mapDefinition(out, d)
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
			Parameters:  mapParams(out, p.Params),
			Responses:   mapResponses(out, p.Responses),
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
