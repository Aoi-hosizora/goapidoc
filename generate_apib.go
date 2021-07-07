package goapidoc

import (
	"fmt"
	"regexp"
	"strings"
)

type apibDocument struct {
	Host     string
	BasePath string

	Title          string
	Description    string
	Version        string
	TermsOfService string
	License        string
	ContactUrl     string
	ContactEmail   string

	GroupsString      string
	DefinitionsString string
}

type apibGroup struct {
	Tag         string
	Description string
	Routes      []*apibRoute
}

type apibRoute struct {
	Route   string
	Summary string
	Methods string
}

type apibMethod struct {
	Route       string
	Method      string
	Summary     string
	Description string
	Deprecated  bool
	Consume     string
	Produce     string

	Parameters []string
	Headers    []string
	Body       string
	Forms      []string

	Responses []*apibResponse
}

type apibResponse struct {
	Code        int
	Description string
	Produce     string

	Headers []string
	Body    string
	Example string
}

type apibSchema struct {
	name       string
	typ        string
	required   bool
	desc       string
	allowEmpty bool
	defaul     interface{}
	example    interface{}
	enum       []interface{}
	minLength  int
	maxLength  int
	minimum    float64
	maximum    float64
}

// =============
// type & schema
// =============

func buildApibType(typ string) (string, *apiType) {
	at := parseApiType(typ)
	// typ = strings.ReplaceAll(strings.ReplaceAll(typ, "<", "«"), ">", "»")
	switch at.kind {
	case apiPrimeKind:
		t := at.prime.typ
		if t == "integer" {
			t = "number"
		}
		return t, at
	case apiArrayKind:
		item, _ := buildApibType(at.array.item.name)
		if typ == "integer" {
			item = "number"
		}
		return fmt.Sprintf("array[%s]", item), at
	case apiObjectKind:
		return typ, at
	default:
		return "", nil
	}
}

func buildApibSchema(schema *apibSchema, in string) string {
	/*
		+ <parameter name>: `<example value>` (<type> | enum[<type>], required | optional) - <description>
		    <additional description>
		    + Default: `<default value>`
		    + Members
		        + `<enumeration value 1>`
		        + `<enumeration value 2>`
	*/
	switch in {
	case HEADER:
		req := "required"
		if !schema.required {
			req = "optional"
		}
		return fmt.Sprintf("%s: %s (%s, %s)", schema.name, schema.desc, schema.typ, req)
	case BODY:
		return schema.typ // <<<
	case PATH, QUERY, FORM:
		// pass
	}
	req := "required"
	if !schema.required {
		req = "optional"
	}
	typ, at := buildApibType(schema.typ)
	if len(schema.enum) != 0 {
		typ = "enum[" + typ + "]"
	}

	paramStr := fmt.Sprintf("+ %s (%s, %s) - %s", schema.name, typ, req, schema.desc) // center
	if schema.example != nil {
		paramStr = fmt.Sprintf("+ %s: `%v` (%s, %s) - %s", schema.name, schema.example, typ, req, schema.desc)
	}

	options := make([]string, 0)

	if at.kind == apiPrimeKind && at.prime.format != "" {
		options = append(options, "format: "+at.prime.format)
	}
	if schema.allowEmpty {
		options = append(options, "allow empty")
	}
	if schema.maxLength != 0 && schema.minLength != 0 {
		options = append(options, fmt.Sprintf("len in \\[%d, %d\\]", schema.minLength, schema.maxLength))
	} else if schema.minLength != 0 {
		options = append(options, fmt.Sprintf("len >= %d", schema.minLength))
	} else if schema.maxLength != 0 {
		options = append(options, fmt.Sprintf("len <= %d", schema.maxLength))
	}
	if schema.maximum != 0 && schema.minimum != 0 {
		options = append(options, fmt.Sprintf("val in \\[%.3f, %.3f\\]", schema.minimum, schema.maximum))
	} else if schema.minimum != 0 {
		options = append(options, fmt.Sprintf("val >= %.3f", schema.minimum))
	} else if schema.maximum != 0 {
		options = append(options, fmt.Sprintf("val <= %.3f", schema.maximum))
	}

	if len(options) != 0 {
		paramStr += fmt.Sprintf("\n    (%s)", strings.Join(options, ", "))
	}

	if schema.defaul != nil {
		paramStr += fmt.Sprintf("\n    + Default: `%v`", schema.defaul)
	}

	if len(schema.enum) != 0 {
		paramStr += "\n    + Members"
		for _, enum := range schema.enum {
			paramStr += fmt.Sprintf("\n        + `%v`", enum)
		}
	}

	return paramStr
}

// =================================
// operations & groups & definitions
// =================================

var apibOperationTemplate = `
### {{ .Summary }} [{{ .Method }}]

` + "`{{ .Method }} {{ .Route }}`" + `

{{ if .Description }}> {{ .Description }}{{ end }}

{{ if .Deprecated }}Attention: This api is deprecated.{{ end }}

{{ if .Parameters }}
+ Parameters

{{ range .Parameters }}{{.}}
{{ end }}
{{ end }}

+ Request ({{ .Consume }})

    + Attributes {{ if .Body }}({{ .Body }}){{ end }}

{{ range .Forms }}{{.}}
{{ end }}

    + Headers

{{ range .Headers }}{{.}}
{{ end }}

    + Body

{{ range .Responses }}

+ Response {{ .Code }} ({{ .Produce }})

    {{ if .Description }}> {{ .Description }}{{ end }}

    + Attributes {{ if .Body }}({{ .Body }}){{ end }}

    + Headers

{{ range .Headers }}{{.}}
{{ end }}

    + Body

{{/* {{ .Example }} */}}

{{ end }}
`

func buildApibOperation(op *Operation, securities map[string]*Security) []byte {
	// prehandle operation fields
	consume := "application/json"
	if len(op.consumes) >= 1 {
		consume = op.consumes[0]
	}
	produce := "application/json"
	if len(op.produces) >= 1 {
		produce = op.produces[0]
	}
	params := op.params
	if len(op.securities) > 0 {
		name := op.securities[0]
		if s, ok := securities[name]; ok {
			desc := fmt.Sprintf("%s, %s", name, s.typ)
			if s.desc != "" {
				desc = fmt.Sprintf("%s (%s), %s", name, s.desc, s.typ)
			}
			if s.typ == "apiKey" {
				params = append(params, &Param{name: s.name, in: s.in, typ: "string", desc: desc})
			} else if s.typ == "basic" {
				params = append(params, &Param{name: "Authorization", in: "header", typ: "string", desc: desc})
			}
		}
	}

	// render operation to apibMethod
	out := &apibMethod{
		Summary:     op.summary,
		Method:      op.method,
		Route:       op.route,
		Description: op.desc,
		Deprecated:  op.deprecated,
		Consume:     consume,
		Produce:     produce,
	}
	for _, param := range params {
		s := buildApibSchema(&apibSchema{
			name:       param.name,
			typ:        param.typ,
			required:   param.required,
			desc:       param.desc,
			allowEmpty: param.allowEmpty,
			defaul:     param.defaul,
			example:    param.example,
			enum:       param.enum,
			minLength:  param.minLength,
			maxLength:  param.maxLength,
			minimum:    param.minimum,
			maximum:    param.maximum,
		}, param.in)
		switch param.in {
		case PATH, QUERY:
			out.Parameters = append(out.Parameters, spaceIndent(1, s))
		case FORM:
			out.Forms = append(out.Headers, spaceIndent(2, s))
		case HEADER:
			out.Headers = append(out.Headers, spaceIndent(3, s))
		case BODY:
			out.Body = s
		}
	}
	for _, resp := range op.responses {
		headers := make([]string, 0, len(resp.headers))
		for _, h := range resp.headers {
			s := buildApibSchema(&apibSchema{name: h.name, typ: h.typ, desc: h.desc, required: true}, HEADER)
			headers = append(headers, spaceIndent(3, s))
		}
		example := ""
		if e, ok := resp.examples[produce]; ok {
			example = e // <<<
		}
		out.Responses = append(out.Responses, &apibResponse{
			Code:        resp.code,
			Description: resp.desc,
			Produce:     produce,
			Headers:     headers,
			Body:        buildApibSchema(&apibSchema{typ: resp.typ}, BODY),
			Example:     example,
		})
	}

	bs, err := renderTemplate(apibOperationTemplate, out)
	if err != nil {
		panic("Internal error: " + err.Error())
	}
	return bs
}

var apibGroupsTemplate = `
{{ range . }}
# Group {{ .Tag }}

{{ if .Description }}> {{ .Description }}{{ end }}

{{ range .Routes }}
## {{ .Summary }} [{{ .Route }}]

{{ .Methods }}

{{ end }}

{{ end }}`

func buildApibGroups(doc *Document) []byte {
	// get tags and securities from document.option
	var tags []*Tag
	var securities map[string]*Security
	if doc.option != nil {
		tags = doc.option.tags
		securities = make(map[string]*Security, len(doc.option.securities))
		for _, sec := range doc.option.securities {
			securities[sec.title] = sec
		}
	}

	// put all operations to trmoMap splitting by tag, route and method
	// trmo: tag - route - method - Operation
	trmoMap := newOrderedMap(len(tags)) // map[string]map[string]map[string]*Operation
	for _, tag := range tags {
		trmoMap.Set(tag.name, newOrderedMap(2)) // cap defaults to 2
	}
	for _, op := range doc.operations {
		// get and prehandle operation's tag, route, method
		tag := "Default"
		if len(op.tags) > 0 {
			tag = op.tags[0]
		}
		queries := make([]string, 0, 3)
		for _, param := range op.params {
			if param.in == "query" {
				queries = append(queries, param.name)
			}
		}
		route := op.route
		if len(queries) > 0 {
			route += fmt.Sprintf("{?%s}", strings.Join(queries, ","))
		}
		method := strings.ToLower(op.method)

		// rmo: route - method - Operation
		rmoMap, ok := trmoMap.Get(tag) // map[string]map[string]*Operation
		if !ok {
			// new tag, need to append
			rmoMap = newOrderedMap(2)
			tags = append(tags, &Tag{name: tag})
		}
		// mo: method - Operation
		moMap, ok := rmoMap.(*orderedMap).Get(route) // map[string]*Operation
		if !ok {
			moMap = newOrderedMap(4) // cap defaults to 4
		}
		moMap.(*orderedMap).Set(method, op)
		rmoMap.(*orderedMap).Set(route, moMap)
		trmoMap.Set(tag, rmoMap)
	}

	// render operations from trmoMap to apibGroup slice
	out := make([]*apibGroup, 0, trmoMap.Length())
	for _, tag := range tags {
		rmoMap := trmoMap.MustGet(tag.name).(*orderedMap)
		outRoutes := make([]*apibRoute, 0, rmoMap.Length())
		for _, route := range rmoMap.Keys() {
			moMap := rmoMap.MustGet(route).(*orderedMap) // map[string]*Operation
			summaries := make([]string, 0, moMap.Length())
			moStrings := make([]string, 0, moMap.Length())
			for _, method := range moMap.Keys() {
				op := moMap.MustGet(method).(*Operation)
				summaries = append(summaries, op.summary)
				moStrings = append(moStrings, string(buildApibOperation(op, securities)))
			}
			outRoutes = append(outRoutes, &apibRoute{
				Route:   route,
				Summary: strings.Join(summaries, ", "),
				Methods: strings.Join(moStrings, "\n\n"),
			})
		}
		out = append(out, &apibGroup{
			Tag:         tag.name,
			Description: tag.desc,
			Routes:      outRoutes,
		})
	}

	bs, err := renderTemplate(apibGroupsTemplate, out)
	if err != nil {
		panic("Internal error: " + err.Error())
	}
	return bs
}

func buildApibDefinitions(doc *Document) string {
	// prehandle definition list
	allSpecTypes := collectAllSpecTypes(doc)
	clonedDefinitions := make([]*Definition, 0, len(doc.definitions))
	for _, definition := range doc.definitions {
		clonedDefinitions = append(clonedDefinitions, prehandleDefinition(definition)) // with generic name checked
	}
	newDefinitionList := prehandleDefinitionList(clonedDefinitions, allSpecTypes)

	// render result string
	definitionStrings := make([]string, 0)
	for _, definition := range newDefinitionList {
		if len(definition.generics) > 0 || len(definition.properties) == 0 {
			continue
		}

		propertyStrings := make([]string, 0)
		for _, property := range definition.properties {
			propertyStr := buildApibSchema(&apibSchema{
				name:       property.name,
				typ:        property.typ,
				required:   property.required,
				desc:       property.desc,
				allowEmpty: property.allowEmpty,
				defaul:     property.defaul,
				example:    property.example,
				enum:       property.enum,
				minLength:  property.minLength,
				maxLength:  property.maxLength,
				minimum:    property.minimum,
				maximum:    property.maximum,
			}, PATH)
			propertyStrings = append(propertyStrings, propertyStr)
		}

		definitionStr := fmt.Sprintf("## %s (object)", definition.name)
		definitionStr += fmt.Sprintf("\n\n%s", strings.Join(propertyStrings, "\n"))
		definitionStrings = append(definitionStrings, definitionStr)
	}

	return strings.Join(definitionStrings, "\n\n")
}

// ========
// document
// ========

var apibDocumentTemplate = `FORMAT: 1A
HOST: {{ .Host }}{{ .BasePath }}

# {{ .Title }} ({{ .Version }})

{{ if .Description }}{{ .Description }}{{ end }}

{{ if .TermsOfService }}{{ .TermsOfService }}{{ end }}

{{ if .License }}{{ .License }}{{ end }}

{{ if .ContactUrl }}{{ .ContactUrl }}{{ end }}

{{ if .ContactEmail }}{{ .ContactEmail }}{{ end }}

{{ .GroupsString }}

# Data Structures

{{ .DefinitionsString }}
`

func buildApibDocument(doc *Document) []byte {
	// check
	if doc.host == "" {
		panic("Host is required in api blueprint 1A")
	}
	if doc.info == nil {
		panic("Info is required in api blueprint 1A")
	}
	if doc.info.title == "" {
		panic("Info.title is required in api blueprint 1A")
	}
	if doc.info.version == "" {
		panic("Info.version is required in api blueprint 1A")
	}
	if len(doc.operations) == 0 {
		panic("Empty operations is not allowed in api blueprint 1A")
	}

	// header
	out := &apibDocument{
		Host:        doc.host,
		BasePath:    doc.basePath,
		Title:       doc.info.title,
		Version:     doc.info.version,
		Description: doc.info.desc,
	}

	// info
	if doc.info.termsOfService != "" {
		out.TermsOfService = fmt.Sprintf("[Terms of service](%s)", doc.info.termsOfService)
	}
	if license := doc.info.license; license != nil {
		out.License = fmt.Sprintf("[License: %s](%s)", license.name, license.url)
	}
	if contact := doc.info.contact; contact != nil {
		out.ContactUrl = fmt.Sprintf("[%s - Website](%s)", contact.name, contact.url)
		if contact.email != "" {
			out.ContactEmail = fmt.Sprintf("[Send email to %s](mailto:%s)", contact.name, contact.email)
		}
	}

	// definition & operation
	out.DefinitionsString = buildApibDefinitions(doc)
	out.GroupsString = string(buildApibGroups(doc))

	// execute template
	bs, err := renderTemplate(apibDocumentTemplate, out)
	if err != nil {
		panic("Internal error: " + err.Error())
	}

	// remove extra newlines
	bs = regexp.MustCompile(`[ \t]+\n`).ReplaceAll(bs, []byte("\n"))
	bs = regexp.MustCompile(`\n{3,}`).ReplaceAll(bs, []byte("\n\n"))
	return bs
}
