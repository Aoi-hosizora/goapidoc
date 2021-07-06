package goapidoc

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
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

func buildApibSchema(schema *apibSchema) string {
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

func buildApibOperation(securities map[string]*Security, op *Operation) string {
	// meta
	method := strings.ToUpper(op.method)
	metaStr := fmt.Sprintf("### %s [%s]\n\n`%s %s`", op.summary, method, method, op.route)
	if op.desc != "" {
		metaStr += "\n\n" + op.desc
	}
	if op.deprecated {
		metaStr += "\n\nAttention: This api is deprecated."
	}

	// request
	consume := "application/json"
	if len(op.consumes) >= 1 {
		consume = op.consumes[0]
	}
	params := op.params
	bodyParamString := ""
	if len(op.securities) >= 1 {
		securityString := op.securities[0] // only support one authentication
		if security, ok := securities[securityString]; ok {
			desc := securityString + " (" + security.desc + "), " + security.typ
			if security.desc == "" {
				desc = securityString + ", " + security.typ
			}
			if security.typ == "apiKey" {
				params = append(params, &Param{name: security.name, in: security.in, typ: "string", desc: desc})
			} else if security.typ == "basic" {
				params = append(params, &Param{name: "Authorization", in: "header", typ: "string", desc: desc})
			}
		}
	}
	parameterStrings := make([]string, 0)
	attributeStrings := make([]string, 0)
	headerStrings := make([]string, 0)
	for _, param := range params {
		paramStr := buildApibSchema(&apibSchema{
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
		})
		paramStr = strings.ReplaceAll(paramStr, "\n", "\n    ")
		if param.in == "path" || param.in == "query" {
			parameterStrings = append(parameterStrings, "    "+paramStr)
		} else if param.in == "formData" {
			paramStr = strings.ReplaceAll(paramStr, "\n", "\n    ")
			attributeStrings = append(attributeStrings, "        "+paramStr)
		} else if param.in == "body" {
			bodyParamString = param.typ // allow one
		} else if param.in == "header" {
			req := "required"
			if !param.required {
				req = "optional"
			}
			headerStr := fmt.Sprintf("            %s: %s (%s, %s)", param.name, param.desc, param.typ, req)
			headerStrings = append(headerStrings, headerStr)
		}
	}

	// + Request <----------------- center
	requestStr := fmt.Sprintf("+ Request (%s)", consume)
	// + Parameters
	if len(parameterStrings) != 0 {
		parameterStr := strings.Join(parameterStrings, "\n")
		requestStr = fmt.Sprintf("+ Parameters\n\n%s\n\n%s", parameterStr, requestStr)
	}
	// + Attributes (req)
	requestStr += "\n\n    + Attributes"
	if bodyParamString != "" && len(attributeStrings) != 0 {
		attributeStr := strings.Join(attributeStrings, "\n")
		requestStr += fmt.Sprintf(" (%s)\n\n%s", bodyParamString, attributeStr)
	} else if len(attributeStrings) != 0 {
		attributeStr := strings.Join(attributeStrings, "\n")
		requestStr += fmt.Sprintf("\n\n%s", attributeStr)
	} else if bodyParamString != "" {
		requestStr += fmt.Sprintf(" (%s)", bodyParamString)
	}
	// + Headers (req)
	requestStr += "\n\n    + Headers"
	if len(headerStrings) != 0 {
		headerStr := strings.Join(headerStrings, "\n")
		requestStr += fmt.Sprintf("\n\n%s", headerStr)
	}
	// + Body (req)
	requestStr += "\n\n    + Body"

	// response
	produce := "application/json"
	if len(op.produces) >= 1 {
		produce = op.produces[0]
	}
	responseStrings := make([]string, 0)
	for _, response := range op.responses {
		// + Response <----------------- center
		responseStr := fmt.Sprintf("+ Response %d (%s)", response.code, produce)
		if response.desc != "" {
			responseStr += fmt.Sprintf("\n\n    %s", response.desc)
		}
		// + Attributes (resp)
		responseStr += "\n\n    + Attributes"
		if response.typ != "" {
			typ, _ := buildApibType(response.typ)
			responseStr += fmt.Sprintf(" (%s)", typ)
		}
		// + Headers (resp)
		responseStr += "\n\n    + Headers"
		headerStrings := make([]string, 0)
		for _, header := range response.headers {
			headerStr := fmt.Sprintf("            %s: %s (%s)", header.name, header.desc, header.typ)
			headerStrings = append(headerStrings, headerStr)
		}
		if len(headerStrings) != 0 {
			headerStr := strings.Join(headerStrings, "\n")
			responseStr += fmt.Sprintf("\n\n%s", headerStr)
		}
		// + Body (resp)
		responseStr += "\n\n    + Body"
		if ex, ok := response.examples[produce]; ok {
			ex = "\n" + ex
			ex = strings.ReplaceAll(ex, "\n", "\n            ")
			responseStr += "\n" + ex
		}
		responseStrings = append(responseStrings, responseStr)
	}

	responseStr := strings.Join(responseStrings, "\n\n")
	return fmt.Sprintf("%s\n\n%s\n\n%s", metaStr, requestStr, responseStr)
}

func buildApibGroups(doc *Document) string {
	var tags map[string]string
	var securities map[string]*Security
	if doc.option != nil {
		tags = make(map[string]string, len(doc.option.tags))
		for _, tag := range doc.option.tags {
			tags[tag.name] = tag.desc
		}
		securities = make(map[string]*Security, len(doc.option.securities))
		for _, security := range doc.option.securities {
			securities[security.title] = security
		}
	}

	// tag - Operation (route & method)
	groups := newOrderedMap(len(tags)) // map[string][]*Operation{}
	for name := range tags {
		groups.Set(name, make([]*Operation, 0))
	}
	for _, op := range doc.operations {
		tag := "Default"
		if len(op.tags) > 0 {
			tag = op.tags[0]
		}
		ops, ok := groups.Get(tag)
		if !ok {
			ops = make([]*Operation, 0)
		}
		ops = append(ops.([]*Operation), op)
		groups.Set(tag, ops)
	}

	// [#, #]
	groupStrings := make([]string, 0, groups.Length())

	for _, tag := range groups.Keys() {
		group := groups.MustGet(tag).([]*Operation)
		// route - method - Operation
		operations := newOrderedMap(len(group)) // map[string]map[string]*Operation{}
		for _, op := range group {
			route := op.route
			query := make([]string, 0)
			for _, param := range op.params {
				if param.in == "query" {
					query = append(query, param.name)
				}
			}
			if len(query) != 0 {
				route += fmt.Sprintf("{?%s}", strings.Join(query, ","))
			}

			methods, ok := operations.Get(route)
			if !ok {
				methods = newOrderedMap(0) // map[string]*Operation
			}
			methods.(*orderedMap).Set(op.method, op)
			operations.Set(route, methods)
		}

		// [##, ##]
		opStrings := make([]string, 0, operations.Length())
		for _, opKey := range operations.Keys() {
			methods := operations.MustGet(opKey).(*orderedMap)

			// [###, ###]
			methodStrings := make([]string, 0, methods.Length())
			summaries := make([]string, 0, methods.Length())
			for _, methodKey := range methods.Keys() {
				operation := methods.MustGet(methodKey).(*Operation)
				summaries = append(summaries, operation.summary)
				methodStr := buildApibOperation(securities, operation)
				methodStrings = append(methodStrings, methodStr)
			}

			summary := strings.Join(summaries, ", ")
			opStr := fmt.Sprintf("## %s [%s]\n\n%s", summary, opKey, strings.Join(methodStrings, "\n\n"))
			opStrings = append(opStrings, opStr)
		}

		groupStr := fmt.Sprintf("# Group %s", tag)
		if tagDesc, ok := tags[tag]; ok {
			groupStr += fmt.Sprintf("\n\n%s", tagDesc)
		}
		if len(opStrings) > 0 {
			groupStr += fmt.Sprintf("\n\n%s", strings.Join(opStrings, "\n\n"))
		}

		groupStrings = append(groupStrings, groupStr)
	}

	return strings.Join(groupStrings, "\n\n")
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
			})
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

var apibTemplate = `FORMAT: 1A
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
	if doc.info == nil {
		panic("Nil document info")
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
	out.GroupsString = buildApibGroups(doc)

	// execute template
	t := template.Must(template.New("apibTemplate").Parse(apibTemplate))
	buf := &bytes.Buffer{}
	err := t.Execute(buf, out)
	if err != nil {
		panic("Internal error: " + err.Error())
	}

	// remove extra newlines
	re := regexp.MustCompile(`\n{3,}`)
	return re.ReplaceAll(buf.Bytes(), []byte("\n\n"))
}
