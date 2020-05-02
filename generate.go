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

type innerPath struct {
	Summary     string                    `yaml:"summary"`
	Description string                    `yaml:"description,omitempty"`
	OperationId string                    `yaml:"operationId"`
	Tags        []string                  `yaml:"tags,omitempty"`
	Consumes    []string                  `yaml:"consumes,omitempty"`
	Produces    []string                  `yaml:"produces,omitempty"`
	Securities  []string                  `yaml:"security,omitempty"`
	Parameters  []*innerParam             `yaml:"parameters,omitempty"`
	Responses   map[string]*innerResponse `yaml:"responses,omitempty"`
}

type innerModel struct {
	Title       string                    `json:"title"`
	Description string                    `json:"description,omitempty"`
	Type        string                    `json:"type"`
	Required    []string                  `json:"required"`
	Properties  map[string]*innerProperty `json:"properties,omitempty"`
}

type innerParam struct {
	Name        string        `yaml:"name"`
	In          string        `yaml:"in"`
	Type        string        `yaml:"type"`
	Required    bool          `yaml:"required"`
	Description string        `yaml:"description,omitempty"`
	Format      string        `yaml:"format,omitempty"`
	Schema      *innerSchema  `yaml:"schema,omitempty"`
	Default     interface{}   `yaml:"default,omitempty"`
	Enum        []interface{} `yaml:"enum,omitempty"`
}

type innerResponse struct {
	Description string            `yaml:"description,omitempty"`
	Schema      *innerSchema      `yaml:"schema,omitempty"`
	Examples    map[string]string `yaml:"examples,omitempty"`
}

type innerProperty struct {
	Description string       `yaml:"description,omitempty"`
	Type        string       `yaml:"type"`
	Format      string       `yaml:"format,omitempty"`
	Enum        []string     `yaml:"enum,omitempty"`
	Ref         string       `yaml:"$ref,omitempty"`
	Items       *innerSchema `yaml:"items,omitempty"`
}

type innerSchema struct {
	Ref string `yaml:"$ref"`
}

func getRef(ref string) string {
	if ref == "" {
		return ""
	}
	return "#/definitions/" + ref
}

func getSchema(ref string) *innerSchema {
	if ref == "" {
		return nil
	}
	return &innerSchema{Ref: "#/definitions/" + ref}
}

func mapToInnerParam(params []*Param) []*innerParam {
	out := make([]*innerParam, len(params))
	for i, p := range params {
		out[i] = &innerParam{
			Name:        p.Name,
			In:          p.In,
			Type:        p.Type,
			Required:    p.Required,
			Description: p.Description,
			Format:      p.Format,
			Schema:      getSchema(p.Schema),
			Default:     p.Default,
			Enum:        p.Enum,
		}
	}
	return out
}

func mapToInnerResponse(responses []*Response) map[string]*innerResponse {
	out := make(map[string]*innerResponse)
	for _, r := range responses {
		out[strconv.Itoa(r.Code)] = &innerResponse{
			Description: r.Description,
			Schema:      getSchema(r.Schema),
			Examples:    r.Examples,
		}
	}
	return out
}

func mapToInnerProperty(properties []*Property) map[string]*innerProperty {
	out := make(map[string]*innerProperty)
	for _, p := range properties {
		out[p.Title] = &innerProperty{
			Description: p.Description,
			Type:        p.Type,
			Format:      p.Format,
			Enum:        p.Enum,
		}
		if p.Type == "object" {
			out[p.Title].Ref = getRef(p.Schema)
		} else if p.Type == "array" {
			out[p.Title].Items = getSchema(p.Schema)
		}
	}
	return out
}

func mapToInnerDocument(d *Document) *innerDocument {
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
			Parameters:  mapToInnerParam(p.Params),
			Responses:   mapToInnerResponse(p.Responses),
		}
	}

	// models
	for _, m := range d.Models {
		required := make([]string, 0)
		for _, p := range m.Properties {
			if p.Required {
				required = append(required, p.Title)
			}
		}
		out.Models[m.Title] = &innerModel{
			Title:       m.Title,
			Description: m.Description,
			Type:        m.Type,
			Required:    required,
			Properties:  mapToInnerProperty(m.Properties),
		}
	}

	return out
}

func appendKvs(d *innerDocument, kvs map[string]interface{}) *yaml.MapSlice {
	out := &yaml.MapSlice{}
	for k, v := range kvs {
		*out = append(*out, yaml.MapItem{Key: k, Value: v})
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
	out := appendKvs(mapToInnerDocument(d), kvs)

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
