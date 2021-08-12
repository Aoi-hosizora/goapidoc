package goapidoc

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
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
	ExternalDoc    string

	Securities    []string
	Schemes       []string
	Consumes      []string
	Produces      []string
	AdditionalDoc string

	GroupsString      string
	DefinitionsString string
}

type apibGroup struct {
	Tag           string
	Description   string
	ExternalDoc   string
	AdditionalDoc string
	Routes        []*apibRoute
}

type apibRoute struct {
	Route         string
	RawRoute      string
	Summary       string
	AdditionalDoc string
	Methods       string
}

type apibMethod struct {
	Route         string
	Method        string
	Summary       string
	Description   string
	Security      string
	Deprecated    bool
	Example       string
	Consume       string
	Produce       string
	ExternalDoc   string
	AdditionalDoc string

	Parameters []string
	Headers    []string
	AttrBody   string
	Forms      []string
	Responses  []*apibResponse
}

type apibResponse struct {
	Code          int
	Description   string
	Produce       string
	AdditionalDoc string

	Headers  []string
	AttrBody string
	Example  string
}

type apibDefinition struct {
	Name       string
	Properties []string
}

type apibSchema struct {
	name     string
	typ      string
	required bool
	desc     string

	allowEmpty       bool
	defaul           interface{}
	example          interface{}
	pattern          string
	enum             []interface{}
	minLength        *int
	maxLength        *int
	minItems         *int
	maxItems         *int
	uniqueItems      bool
	collectionFormat string
	minimum          *float64
	maximum          *float64
	exclusiveMin     bool
	exclusiveMax     bool
	multipleOf       float64
}

// =====================================
// type & schema & externalDoc & example
// =====================================

func buildApibType(typ string) (string, *apiType) {
	at := parseApiType(typ)
	switch at.kind {
	case apiPrimeKind:
		t := at.prime.typ
		if t == FILE {
			logWarning("file type is not supported in API Blueprint, this will be replaced to string.")
			t = STRING
		} else if t == INTEGER {
			t = NUMBER
		}
		return t, at
	case apiArrayKind:
		t, _ := buildApibType(at.array.item.name)
		return fmt.Sprintf("array[%s]", t), at
	case apiObjectKind:
		return typ, at
	default:
		return "", nil // unreachable
	}
}

func buildApibSchema(schema *apibSchema, in string) string {
	typ, at := buildApibType(schema.typ)
	req := "required"
	if !schema.required {
		req = "optional"
	}
	switch in {
	case BODY:
		return typ
	case HEADER:
		if schema.example != nil {
			return fmt.Sprintf("%s: %v", schema.name, schema.example)
		}
		if schema.desc == "" {
			return fmt.Sprintf("%s: (%s, %s)", schema.name, typ, req)
		}
		return fmt.Sprintf("%s: (%s, %s) - %s", schema.name, typ, req, schema.desc)
	case PATH, QUERY, FORM:
		// pass
	}

	/*
		+ <parameter name>: `<example value>` (<type> | enum[<type>], required | optional) - <description>
		    <additional description>
		    + Default: `<default value>`
		    + Members
		        + `<enumeration value 1>`
		        + `<enumeration value 2>`
	*/
	out := strings.Builder{}
	if len(schema.enum) != 0 {
		typ = fmt.Sprintf("enum[%s]", typ)
	}
	if schema.example != nil {
		out.WriteString(fmt.Sprintf("+ %s: `%v` (%s, %s)", schema.name, schema.example, typ, req))
	} else {
		out.WriteString(fmt.Sprintf("+ %s (%s, %s)", schema.name, typ, req))
	}
	if schema.desc != "" {
		out.WriteString(fmt.Sprintf(" - %s", schema.desc))
	}
	options := make([]string, 0, 4) // cap defaults to 4
	if at.kind == apiPrimeKind && at.prime.format != "" {
		options = append(options, fmt.Sprintf("format: %s", at.prime.format))
	}
	if schema.allowEmpty {
		options = append(options, "allow empty value")
	}
	if schema.pattern != "" {
		options = append(options, fmt.Sprintf("pattern: /%s/", schema.pattern))
	}
	if schema.maxLength != nil && schema.minLength != nil {
		options = append(options, fmt.Sprintf("%d <= len <= %d", *schema.minLength, *schema.maxLength))
	} else if schema.minLength != nil {
		options = append(options, fmt.Sprintf("len >= %d", *schema.minLength))
	} else if schema.maxLength != nil {
		options = append(options, fmt.Sprintf("len <= %d", *schema.maxLength))
	}
	if schema.maxItems != nil && schema.minItems != nil {
		options = append(options, fmt.Sprintf("%d <= #items <= %d", *schema.minItems, *schema.maxItems))
	} else if schema.minItems != nil {
		options = append(options, fmt.Sprintf("#items >= %d", *schema.minItems))
	} else if schema.maxItems != nil {
		options = append(options, fmt.Sprintf("#items <= %d", *schema.maxItems))
	}
	if schema.uniqueItems {
		options = append(options, "unique items")
	}
	if schema.collectionFormat != "" {
		options = append(options, fmt.Sprintf("collection format: %s", schema.collectionFormat))
	}
	ltSign, gtSign, gtSignR := "<=", ">=", "<="
	if schema.exclusiveMax {
		ltSign = "<"
	}
	if schema.exclusiveMin {
		gtSign, gtSignR = ">", "<"
	}
	if schema.maximum != nil && schema.minimum != nil {
		options = append(options, fmt.Sprintf("%.3f %s val %s %.3f", *schema.minimum, gtSignR, ltSign, *schema.maximum))
	} else if schema.minimum != nil {
		options = append(options, fmt.Sprintf("val %s %.3f", gtSign, *schema.minimum))
	} else if schema.maximum != nil {
		options = append(options, fmt.Sprintf("val %s %.3f", ltSign, *schema.maximum))
	}
	if schema.multipleOf != 0 {
		options = append(options, fmt.Sprintf("multiple of %.3f", schema.multipleOf))
	}
	if len(options) > 0 {
		out.WriteString(fmt.Sprintf("\n    (%s)", strings.Join(options, ", ")))
	}
	if schema.defaul != nil {
		out.WriteString(fmt.Sprintf("\n    + Default: `%v`", schema.defaul))
	}
	if len(schema.enum) != 0 {
		out.WriteString("\n    + Members")
		for _, enum := range schema.enum {
			out.WriteString(fmt.Sprintf("\n        + `%v`", enum))
		}
	}

	return out.String()
}

func buildApiExternalDoc(doc *ExternalDoc) string {
	if doc == nil {
		return ""
	}
	desc := doc.desc
	if desc == "" {
		desc = doc.url
	}
	return fmt.Sprintf("[%s](%s)", desc, doc.url)
}

func buildApibExample(e interface{}, produce string) string {
	if e == nil {
		return ""
	}
	str := ""
	switch e.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, uintptr, string, bool:
		str = fmt.Sprintf("%v", e)
	default:
		if produce != JSON {
			logWarning(fmt.Sprintf("Example in %s mime is not supported in API Blueprint yet, this will be shown in json format.", produce))
		}
		bs, err := jsonMarshal(e)
		if err != nil {
			logWarning(fmt.Sprintf("Example in %T type is not supported in API Blueprint yet, this will be ignored.", e))
		} else {
			str = fastBtos(bs)
		}
	}
	return str
}

// ================================
// operation & groups & definitions
// ================================

var apibOperationTemplate = `
### {{ .Summary }} [{{ .Method }}]

` + "> `{{ .Method }} {{ .Route }}`" + `

{{ if .Description }}{{ .Description }}{{ end }}

{{ if .Deprecated }}**Attention: This api is deprecated!**{{ end }}

{{ if .Security }}Security requirement: {{ .Security }}{{ end }}

{{ if .ExternalDoc }}{{ .ExternalDoc }}{{ end }}

{{ if .AdditionalDoc }}{{ .AdditionalDoc }}{{ end }}

{{ if .Parameters }}
+ Parameters

{{ range .Parameters }}{{ . }}
{{ end }}
{{ end }}

+ Request ({{ .Consume }})

{{ if .AttrBody }}
    + Attributes ({{ .AttrBody }})

{{ range .Forms }}{{ . }}
{{ end }}
{{ end }}

{{ if .Headers }}
    + Headers

{{ range .Headers }}{{ . }}
{{ end }}
{{ end }}

    + Body

{{ .Example }}

{{ range .Responses }}

+ Response {{ .Code }} ({{ .Produce }})

    {{ if .Description }}{{ .Description }}{{ end }}

    {{ if .AdditionalDoc }}{{ .AdditionalDoc }}{{ end }}

{{ if .AttrBody }}
    + Attributes ({{ .AttrBody }})
{{ end }}

{{ if .Headers }}
    + Headers

{{ range .Headers }}{{ . }}
{{ end }}
{{ end }}

    + Body

{{ .Example }}

{{ end }}
`

func buildApibOperation(op *Operation, securities map[string]*Security) ([]byte, error) {
	// prehandle operation fields
	consume := JSON
	if len(op.consumes) >= 1 {
		consume = op.consumes[0]
	}
	produce := JSON
	if len(op.produces) >= 1 {
		produce = op.produces[0]
	}
	params := op.params
	secReqName := ""
	if len(op.securities) > 0 {
		secReqName = op.securities[0] // only support single security
		if s, ok := securities[secReqName]; ok {
			desc := fmt.Sprintf("%s, %s", secReqName, s.typ)
			if s.desc != "" {
				desc = fmt.Sprintf("%s (%s), %s", secReqName, s.desc, s.typ)
			}
			if s.typ == APIKEY {
				params = append(params, &Param{name: s.name, in: s.in, typ: "string", desc: desc})
			} else if s.typ == BASIC {
				params = append(params, &Param{name: "Authorization", in: "header", typ: "string", desc: desc})
			} else if s.typ == OAUTH2 {
				params = append(params, &Param{name: "Authorization", in: "header", typ: "string", desc: desc}) // <<<
			}
		}
	}

	// render operation to apibMethod
	out := &apibMethod{
		Route:         op.route,
		Method:        strings.ToUpper(op.method),
		Summary:       op.summary,
		Description:   op.desc,
		Security:      secReqName,
		Deprecated:    op.deprecated,
		Example:       spaceIndent(3, buildApibExample(op.reqExample, produce)),
		Consume:       consume,
		Produce:       produce,
		ExternalDoc:   buildApiExternalDoc(op.externalDoc),
		AdditionalDoc: op.additionalDoc,
		Parameters:    make([]string, 0, 2),
		Headers:       make([]string, 0, 2),
		AttrBody:      "",
		Forms:         make([]string, 0),
		Responses:     make([]*apibResponse, 0, 1),
	}
	for _, p := range params {
		s := buildApibSchema(&apibSchema{
			name:     p.name,
			typ:      p.typ,
			required: p.required,
			desc:     p.desc,

			allowEmpty:       p.allowEmpty,
			defaul:           p.defaul,
			example:          p.example,
			pattern:          p.pattern,
			enum:             p.enum,
			minLength:        p.minLength,
			maxLength:        p.maxLength,
			minItems:         p.minItems,
			maxItems:         p.maxItems,
			uniqueItems:      p.uniqueItems,
			collectionFormat: p.collectionFormat,
			minimum:          p.minimum,
			maximum:          p.maximum,
			exclusiveMin:     p.exclusiveMin,
			exclusiveMax:     p.exclusiveMax,
			multipleOf:       p.multipleOf,
		}, p.in)
		switch p.in {
		case PATH, QUERY:
			out.Parameters = append(out.Parameters, spaceIndent(1, s))
		case FORM:
			out.Forms = append(out.Forms, spaceIndent(2, s))
			out.AttrBody = "object"
		case HEADER:
			out.Headers = append(out.Headers, spaceIndent(3, s))
		case BODY:
			out.AttrBody = s
		}
	}
	for _, r := range op.responses {
		desc := r.desc
		if desc == "" {
			desc = strconv.Itoa(r.code) + " " + http.StatusText(r.code)
		}
		headers := make([]string, 0, len(r.headers))
		for _, h := range r.headers {
			s := buildApibSchema(&apibSchema{name: h.name, typ: h.typ, desc: h.desc, required: true, example: h.example}, HEADER)
			headers = append(headers, spaceIndent(3, s))
		}
		example := ""
		for _, e := range r.examples {
			if e.mime == produce {
				example = spaceIndent(3, buildApibExample(e.example, produce))
				break
			}
		}
		out.Responses = append(out.Responses, &apibResponse{
			Code:          r.code,
			Description:   desc,
			Produce:       produce,
			AdditionalDoc: r.additionalDoc,
			Headers:       headers,
			AttrBody:      buildApibSchema(&apibSchema{typ: r.typ}, BODY),
			Example:       example,
		})
	}

	return renderTemplate(apibOperationTemplate, out)
}

var apibGroupsTemplate = `
{{ range . }}
# Group {{ .Tag }}

{{ if .Description }}{{ .Description }}{{ end }}

{{ if .ExternalDoc }}{{ .ExternalDoc }}{{ end }}

{{ if .AdditionalDoc }}{{ .AdditionalDoc }}{{ end }}

{{ range .Routes }}
## {{ .Summary }} [{{ .Route }}]

` + "> `{{ .RawRoute }}`" + `

{{ if .AdditionalDoc }}{{ .AdditionalDoc }}{{ end }}

{{ .Methods }}

{{ end }}

{{ end }}`

func buildApibGroups(doc *Document) ([]byte, error) {
	// get tags and securities from document.option
	var allTags []*Tag
	var securities map[string]*Security
	var routesOptions map[string]*RoutesOption
	if opt := doc.option; opt != nil {
		allTags = opt.tags
		securities = make(map[string]*Security, len(opt.securities))
		for _, sec := range opt.securities {
			securities[sec.title] = sec
		}
		routesOptions = make(map[string]*RoutesOption, len(opt.routesOptions))
		for _, ro := range opt.routesOptions {
			routesOptions[ro.route] = ro
		}
	}

	// put all operations to trmoMap splitting by tag, route and method
	// trmo: tag - route - method - Operation
	trmoMap := newOrderedMap(len(allTags)) // map[string]map[string]map[string]*Operation
	for _, tag := range allTags {
		trmoMap.Set(tag.name, newOrderedMap(2)) // cap defaults to 2
	}
	for _, op := range doc.operations {
		// get and prehandle operation's tag, route, method
		tag := "Default"
		if len(op.tags) > 0 {
			tag = op.tags[0]
		}
		queries := make([]string, 0, 3) // cap defaults to 3
		for _, param := range op.params {
			if param.in == "query" {
				queries = append(queries, param.name)
			}
		}
		route := op.route
		if len(queries) > 0 {
			route += fmt.Sprintf("{?%s}", strings.Join(queries, ","))
		}
		method := strings.ToUpper(op.method)

		// rmo: route - method - Operation
		rmoMap, ok := trmoMap.Get(tag) // map[string]map[string]*Operation
		if !ok {
			// new tag, need to append
			rmoMap = newOrderedMap(2)
			allTags = append(allTags, &Tag{name: tag, desc: ""})
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
	for _, tag := range allTags {
		rmoMap := trmoMap.MustGet(tag.name).(*orderedMap)
		outRoutes := make([]*apibRoute, 0, rmoMap.Length())
		for _, route := range rmoMap.Keys() {
			moMap := rmoMap.MustGet(route).(*orderedMap) // map[string]*Operation
			rawRoute := route
			summaries := make([]string, 0, moMap.Length())
			moStrings := make([]string, 0, moMap.Length())
			for _, method := range moMap.Keys() {
				op := moMap.MustGet(method).(*Operation)
				rawRoute = op.route
				summaries = append(summaries, op.summary)
				bs, err := buildApibOperation(op, securities)
				if err != nil {
					return nil, err
				}
				moStrings = append(moStrings, fastBtos(bs))
			}
			summary := strings.Join(summaries, " | ")
			additionalDoc := ""
			if opt, ok := routesOptions[route]; ok {
				if opt.summary != "" {
					summary = opt.summary
				}
				additionalDoc = opt.additionalDoc
			}
			outRoutes = append(outRoutes, &apibRoute{
				Route:         route,
				RawRoute:      rawRoute,
				Summary:       summary,
				AdditionalDoc: additionalDoc,
				Methods:       strings.Join(moStrings, "\n\n"),
			})
		}
		out = append(out, &apibGroup{
			Tag:           tag.name,
			Description:   tag.desc,
			ExternalDoc:   buildApiExternalDoc(tag.externalDoc),
			AdditionalDoc: tag.additionalDoc,
			Routes:        outRoutes,
		})
	}

	return renderTemplate(apibGroupsTemplate, out)
}

var apibDefinitionTemplate = `
# Data Structures

{{ range . }}
## {{ .Name }} (object)

{{ range .Properties }}{{ . }}
{{ end }}

{{ end }}`

func buildApibDefinitions(doc *Document) ([]byte, error) {
	// prehandle definition list
	allSpecTypes := collectAllSpecTypes(doc)
	clonedDefinitions := make([]*Definition, 0, len(doc.definitions))
	for _, definition := range doc.definitions {
		clonedDefinitions = append(clonedDefinitions, prehandleDefinition(definition)) // with generic name checked
	}
	newDefinitionList := prehandleDefinitionList(clonedDefinitions, allSpecTypes)

	// render definitions to apibDefinition slice
	out := make([]*apibDefinition, 0, len(newDefinitionList))
	for _, def := range newDefinitionList {
		props := make([]string, 0, len(def.properties))
		for _, p := range def.properties {
			props = append(props, buildApibSchema(&apibSchema{
				name:     p.name,
				typ:      p.typ,
				required: p.required,
				desc:     p.desc,

				allowEmpty:       p.allowEmpty,
				defaul:           p.defaul,
				example:          p.example,
				pattern:          p.pattern,
				enum:             p.enum,
				minLength:        p.minLength,
				maxLength:        p.maxLength,
				minItems:         p.minItems,
				maxItems:         p.maxItems,
				uniqueItems:      p.uniqueItems,
				collectionFormat: p.collectionFormat,
				minimum:          p.minimum,
				maximum:          p.maximum,
				exclusiveMin:     p.exclusiveMin,
				exclusiveMax:     p.exclusiveMax,
				multipleOf:       p.multipleOf,
			}, PATH))
		}
		out = append(out, &apibDefinition{Name: def.name, Properties: props})
	}

	return renderTemplate(apibDefinitionTemplate, out)
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

{{ if .ExternalDoc }}{{ .ExternalDoc }}{{ end }}

{{ if .Securities }}## Securities defined

{{ range .Securities }}+ {{ . }}
{{ end }}{{ end }}

{{ if .Schemes }}## Supported schemes

{{ range .Schemes }}+ {{ . }}
{{ end }}{{ end }}

{{ if .Consumes }}## Available consumes

{{ range .Consumes }}+ {{ . }}
{{ end }}{{ end }}

{{ if .Produces }}## Available produces

{{ range .Produces }}+ {{ . }}
{{ end }}{{ end }}

{{ if .AdditionalDoc }}{{ .AdditionalDoc }}{{ end }}

{{ .GroupsString }}

{{ .DefinitionsString }}
`

func buildApibDocument(doc *Document) ([]byte, error) {
	// check
	checkDocument(doc)

	// info
	out := &apibDocument{
		Host:        doc.host,
		BasePath:    doc.basePath,
		Title:       doc.info.title,
		Version:     doc.info.version,
		Description: doc.info.desc,
	}
	if doc.info.termsOfService != "" {
		out.TermsOfService = fmt.Sprintf("[Terms of service](%s)", doc.info.termsOfService)
	}
	if license := doc.info.license; license != nil {
		if license.url == "" {
			out.License = fmt.Sprintf("License: %s", license.name)
		} else {
			out.License = fmt.Sprintf("[License: %s](%s)", license.name, license.url)
		}
	}
	if contact := doc.info.contact; contact != nil {
		name := contact.name
		if name == "" {
			name = "the developer"
		}
		if contact.url != "" {
			out.ContactUrl = fmt.Sprintf("[%s - Website](%s)", name, contact.url)
		}
		if contact.email != "" {
			out.ContactEmail = fmt.Sprintf("[Send email to %s](mailto:%s)", name, contact.email)
		}
	}
	if opt := doc.option; opt != nil {
		securities := make([]string, 0, len(opt.securities))
		for _, sec := range opt.securities {
			s := ""
			opts := make([]string, 0, 2)
			if sec.typ == APIKEY {
				s = fmt.Sprintf("%s: apiKey", sec.title)
				opts = append(opts, fmt.Sprintf("+ name: %s", sec.name))
				opts = append(opts, fmt.Sprintf("+ in: %s", sec.in))
			} else if sec.typ == BASIC {
				s = fmt.Sprintf("%s: basic", sec.title)
			} else if sec.typ == OAUTH2 {
				s = fmt.Sprintf("%s: oauth2", sec.title)
				opts = append(opts, fmt.Sprintf("+ flow: %s", sec.flow))
				if sec.authorizationUrl != "" {
					opts = append(opts, fmt.Sprintf("+ authUrl: %s", sec.authorizationUrl))
				}
				if sec.tokenUrl != "" {
					opts = append(opts, fmt.Sprintf("+ tokenUrl: %s`", sec.tokenUrl))
				}
				scopes := make([][2]string, 0, len(sec.scopes))
				for _, c := range sec.scopes {
					scopes = append(scopes, [2]string{c.scope, c.desc})
				}
				sort.Slice(scopes, func(i, j int) bool { return scopes[i][0] < scopes[j][0] })
				for _, kv := range scopes {
					opts = append(opts, fmt.Sprintf("+ scope: %s - %s", kv[0], kv[1]))
				}
			} else {
				continue // unreachable
			}
			if sec.desc != "" {
				s += fmt.Sprintf(" - %s", sec.desc)
			}
			if len(opts) != 0 {
				s += fmt.Sprintf("\n    %s", strings.Join(opts, "\n    "))
			}
			securities = append(securities, s)
		}

		out.ExternalDoc = buildApiExternalDoc(opt.externalDoc)
		out.Securities = securities
		out.Schemes = opt.schemes
		out.Consumes = opt.consumes
		out.Produces = opt.produces
		out.AdditionalDoc = opt.additionalDoc
	}

	// definitions & operations
	s1, err1 := buildApibDefinitions(doc)
	s2, err2 := buildApibGroups(doc)
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	out.DefinitionsString = fastBtos(s1)
	out.GroupsString = fastBtos(s2)

	// execute template and format
	bs, err := renderTemplate(apibDocumentTemplate, out)
	if err != nil {
		return nil, err
	}
	bs = regexp.MustCompile(`[ \t]+\n`).ReplaceAll(bs, []byte("\n"))
	bs = regexp.MustCompile(`\n{3,}`).ReplaceAll(bs, []byte("\n\n"))
	bs = regexp.MustCompile(`\n+$`).ReplaceAll(bs, []byte("\n"))
	return bs, nil
}
