package yamldoc

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

type innerDocument struct {
	Host       string                           `yaml:"host"`
	BasePath   string                           `yaml:"basePath"`
	Info       *innerInfo                       `yaml:"info"`
	Tags       []*innerTag                      `yaml:"tags,omitempty"`
	Securities map[string]*innerSecurity        `yaml:"securityDefinitions,omitempty"`
	Paths      map[string]map[string]*innerPath `yaml:"paths,omitempty"`
	Models     map[string]*innerModel           `yaml:"definitions,omitempty"`
}

type innerTag struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
}

type innerLicense struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url,omitempty"`
}

type innerContact struct {
	Name  string `yaml:"name"`
	Url   string `yaml:"url,omitempty"`
	Email string `yaml:"email,omitempty"`
}

type innerInfo struct {
	Title          string        `yaml:"title"`
	Description    string        `yaml:"description"`
	Version        string        `yaml:"version"`
	TermsOfService string        `yaml:"termsOfService,omitempty"`
	License        *innerLicense `yaml:"license,omitempty"`
	Contact        *innerContact `yaml:"contact,omitempty"`
}

type innerSecurity struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
	In   string `yaml:"in"`
}

// !
type innerPath struct {
	Summary     string                    `yaml:"summary"`
	OperationId string                    `yaml:"operationId"`
	Description string                    `yaml:"description,omitempty"`
	Tags        []string                  `yaml:"tags,omitempty"`
	Consumes    []string                  `yaml:"consumes,omitempty"`
	Produces    []string                  `yaml:"produces,omitempty"`
	Securities  []string                  `yaml:"security,omitempty"`
	Deprecated  bool                      `yaml:"deprecated,omitempty"`
	Parameters  []*innerParam             `yaml:"parameters,omitempty"`
	Responses   map[string]*innerResponse `yaml:"responses,omitempty"`
}

type innerResponse struct {
	Description string                  `yaml:"description,omitempty"`
	Headers     map[string]*innerHeader `yaml:"headers,omitempty"`
	Examples    map[string]string       `yaml:"examples,omitempty"`
	Schema      *innerSchema            `yaml:"schema,omitempty"`
}

type innerHeader struct {
	Type        string      `yaml:"type,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Format      string      `yaml:"format,omitempty"`
	Default     interface{} `yaml:"default,omitempty"`
}

// !
type innerParam struct {
	Name            string        `yaml:"name"`
	In              string        `yaml:"in"`
	Required        bool          `yaml:"required"`
	Type            string        `yaml:"type,omitempty"`
	Description     string        `yaml:"description,omitempty"`
	Format          string        `yaml:"format,omitempty"`
	AllowEmptyValue bool          `yaml:"allowEmptyValue,omitempty"`
	Default         interface{}   `yaml:"default,omitempty"`
	Enum            []interface{} `yaml:"enum,omitempty"`
	Schema          *innerSchema  `yaml:"schema,omitempty"`
	Items           *innerItems   `yaml:"items,omitempty"`
}

// !
type innerModel struct {
	Type        string                  `json:"type"`
	Required    []string                `json:"required"`
	Description string                  `json:"description,omitempty"`
	Properties  map[string]*innerSchema `json:"properties,omitempty"`
}

// !!
type innerSchema struct {
	Title       string `yaml:"-"`
	Type        string `yaml:"type,omitempty"`
	Required    bool   `yaml:"required,omitempty"`
	Description string `yaml:"description,omitempty"`

	Format          string        `yaml:"format,omitempty"`
	AllowEmptyValue bool          `yaml:"allowEmptyValue,omitempty"`
	Default         interface{}   `yaml:"default,omitempty"`
	Enum            []interface{} `yaml:"enum,omitempty"`

	OriginRef string      `yaml:"originRef,omitempty"`
	Ref       string      `yaml:"$ref,omitempty"`
	Items     *innerItems `yaml:"items,omitempty"`
}

// !!
type innerItems struct {
	Type    string        `yaml:"type,omitempty"`
	Format  string        `yaml:"format,omitempty"`
	Default interface{}   `yaml:"default,omitempty"`
	Enum    []interface{} `yaml:"enum,omitempty"`

	OriginRef string      `yaml:"originRef,omitempty"`
	Ref       string      `yaml:"$ref,omitempty"`
	Items     *innerItems `yaml:"items,omitempty"`
}

func mapSchema(schema *Schema) *innerSchema {
	if schema == nil {
		return nil
	}
	if schema.Ref != "" {
		return &innerSchema{
			OriginRef: schema.Ref,
			Ref:       "#/definitions/" + schema.Ref,
			// Title:       schema.Title,
			// Type:        schema.Type,
			// Required:    schema.Required,
			// Description: schema.Description,
		}
	}
	return &innerSchema{
		Title:           schema.Title,
		Type:            schema.Type,
		Required:        schema.Required,
		Description:     schema.Description,
		Format:          schema.Format,
		AllowEmptyValue: schema.AllowEmptyValue,
		Default:         schema.Default,
		Enum:            schema.Enum,
		Items:           mapItems(schema.Items),
	}
}

func mapItems(items *Items) *innerItems {
	if items == nil {
		return nil
	}
	if items.Ref != "" {
		return &innerItems{
			OriginRef: items.Ref,
			Ref:       "#/definitions/" + items.Ref,
		}
	}
	return &innerItems{
		Type:    items.Type,
		Format:  items.Format,
		Default: items.Default,
		Enum:    items.Enum,
		Items:   mapItems(items.Items),
	}
}

func mapParam(params []*Param) []*innerParam {
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
			Schema:          mapSchema(p.Schema),
			Items:           mapItems(p.Items),
		}
	}
	return out
}

func mapResponse(responses []*Response) map[string]*innerResponse {
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
			Schema:      mapSchema(r.Schema),
			Examples:    r.Examples,
			Headers:     headers,
		}
	}
	return out
}

func buildDocument(d *Document) *innerDocument {
	out := &innerDocument{
		Host:     d.Host,
		BasePath: d.BasePath,
		Info: &innerInfo{
			Title:          d.Info.Title,
			Description:    d.Info.Description,
			Version:        d.Info.Version,
			TermsOfService: d.Info.TermsOfService,
			License:        &innerLicense{Name: d.Info.License.Name, Url: d.Info.License.Url},
			Contact:        &innerContact{Name: d.Info.Contact.Name, Url: d.Info.Contact.Url, Email: d.Info.Contact.Email},
		},
		Tags:       []*innerTag{},
		Securities: map[string]*innerSecurity{},
		Paths:      map[string]map[string]*innerPath{},
		Models:     map[string]*innerModel{},
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

		out.Paths[p.Route][p.Method] = &innerPath{
			Summary:     p.Summary,
			Description: p.Description,
			OperationId: id,
			Tags:        p.Tags,
			Consumes:    p.Consumes,
			Produces:    p.Produces,
			Securities:  p.Securities,
			Deprecated:  p.Deprecated,
			Parameters:  mapParam(p.Params),
			Responses:   mapResponse(p.Responses),
		}
	}

	// models
	for _, m := range d.Models {
		required := make([]string, 0)
		schemas := make(map[string]*innerSchema)
		for _, p := range m.Properties {
			name := p.Title
			// p.Title = ""
			if p.Required {
				required = append(required, name)
			}
			schemas[name] = mapSchema(p)
		}

		out.Models[m.Name] = &innerModel{
			Type:        "object",
			Description: m.Description,
			Required:    required,
			Properties:  schemas,
		}
	}

	return out
}

func appendKvs(d *innerDocument, kvs map[string]interface{}) *yaml.MapSlice {
	out := &yaml.MapSlice{}
	if kvs != nil {
		for k, v := range kvs {
			*out = append(*out, yaml.MapItem{Key: k, Value: v})
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
		value := innerValue.Field(i).Interface()

		if name != "-" && name != "" {
			if !omitempty || (value != nil && value != "") {
				*out = append(*out, yaml.MapItem{Key: name, Value: value})
			}
		}
	}

	return out
}

func (d *Document) GenerateYaml(path string, kvs map[string]interface{}) error {
	out := appendKvs(buildDocument(d), kvs)

	doc, err := yaml.Marshal(out)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, doc, 0777)
	if err != nil {
		return err
	}
	return nil
}

func GenerateYaml(path string, kvs map[string]interface{}) error {
	return _document.GenerateYaml(path, kvs)
}
