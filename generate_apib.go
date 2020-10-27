package goapidoc

import (
	"fmt"
	"strings"
)

func buildApibType(typ string) (string, *apiType) {
	at := parseApiType(typ)
	// typ = strings.ReplaceAll(strings.ReplaceAll(typ, "<", "«"), ">", "»")
	if at.kind == apiPrimeKind {
		typ := at.prime.typ
		if typ == "integer" {
			typ = "number"
		}
		return typ, at
	} else if at.kind == apiObjectKind {
		return typ, at
	} else {
		item, _ := buildApibType(at.array.item.name)
		if typ == "integer" {
			item = "number"
		}
		return fmt.Sprintf("array[%s]", item), at
	}
}

func buildApibParam(param *Param) string {
	req := "required"
	if !param.required {
		req = "optional"
	}
	typ, at := buildApibType(param.typ)
	if len(param.enum) != 0 {
		typ = "enum[" + typ + "]"
	}

	paramStr := fmt.Sprintf("+ %s (%s, %s) - %s", param.name, typ, req, param.desc) // center
	if param.example != nil {
		paramStr = fmt.Sprintf("+ %s: `%v` (%s, %s) - %s", param.name, param.example, typ, req, param.desc)
	}

	options := make([]string, 0)

	if at.kind == apiPrimeKind && at.prime.format != "" {
		options = append(options, "format: "+at.prime.format)
	}
	if param.allowEmpty {
		options = append(options, "allow empty")
	}
	if param.maxLength != 0 && param.minLength != 0 {
		options = append(options, fmt.Sprintf("len in \\[%d, %d\\]", param.minLength, param.maxLength))
	} else if param.minLength != 0 {
		options = append(options, fmt.Sprintf("len >= %d", param.minLength))
	} else if param.maxLength != 0 {
		options = append(options, fmt.Sprintf("len <= %d", param.maxLength))
	}
	if param.maximum != 0 && param.minimum != 0 {
		options = append(options, fmt.Sprintf("val in \\[%d, %d\\]", param.minimum, param.maximum))
	} else if param.minimum != 0 {
		options = append(options, fmt.Sprintf("val >= %d", param.minimum))
	} else if param.maximum != 0 {
		options = append(options, fmt.Sprintf("val <= %d", param.maximum))
	}

	if len(options) != 0 {
		paramStr += fmt.Sprintf("\n    (%s)", strings.Join(options, ", "))
	}

	if param.def != nil {
		paramStr += fmt.Sprintf("\n    + Default: `%v`", param.def)
	}

	if len(param.enum) != 0 {
		paramStr += "\n    + Members"
		for _, enum := range param.enum {
			paramStr += fmt.Sprintf("\n        + `%v`", enum)
		}
	}

	return paramStr
}

func buildApibProperty(prop *Property) string {
	return buildApibParam(&Param{
		name:       prop.name,
		typ:        prop.typ,
		required:   prop.required,
		desc:       prop.desc,
		allowEmpty: prop.allowEmpty,
		def:        prop.def,
		example:    prop.example,
		enum:       prop.enum,
		minLength:  prop.minimum,
		maxLength:  prop.maxLength,
		minimum:    prop.minimum,
		maximum:    prop.maxLength,
	})
}

func buildApibPath(securities map[string]*Security, path *RoutePath) string {
	// meta
	method := strings.ToUpper(path.method)
	metaStr := fmt.Sprintf("### %s [%s]\n\n`%s %s`", path.summary, method, method, path.route)
	if path.desc != "" {
		metaStr += "\n\n" + path.desc
	}
	if path.deprecated {
		metaStr += "\n\nAttention: This api is deprecated."
	}

	// request
	consume := "application/json"
	if len(path.consumes) >= 1 {
		consume = path.consumes[0]
	}
	params := path.params
	bodyParamString := ""
	if len(path.securities) >= 1 {
		securityString := path.securities[0]
		if security, ok := securities[securityString]; ok {
			params = append(params, &Param{name: security.name, in: security.in, typ: "string", desc: securityString + " " + security.typ})
		}
	}
	parameterStrings := make([]string, 0)
	attributeStrings := make([]string, 0)
	headerStrings := make([]string, 0)
	for _, param := range params {
		paramStr := buildApibParam(param)
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
	if len(path.produces) >= 1 {
		produce = path.produces[0]
	}
	responseStrings := make([]string, 0)
	for _, response := range path.responses {
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

func buildApibGroups(d *Document) string {
	// [#, #]
	groupStrings := make([]string, 0)

	tags := make(map[string]string)
	for _, tag := range d.tags {
		tags[tag.name] = tag.desc
	}
	securities := make(map[string]*Security)
	for _, security := range d.securities {
		securities[security.title] = security
	}

	groups := newLinkedHashMap() // map[string][]*RoutePath{}
	for _, tag := range d.tags {
		groups.Set(tag.name, make([]*RoutePath, 0))
	}
	for _, path := range d.paths {
		tag := "Default"
		if len(path.tags) > 0 {
			tag = path.tags[0]
		}
		paths, ok := groups.Get(tag)
		if !ok {
			paths = make([]*RoutePath, 0)
		}
		paths = append(paths.([]*RoutePath), path)
		groups.Set(tag, paths)
	}

	for _, tag := range groups.Keys() {
		// [##, ##]
		pathStrings := make([]string, 0)
		group := groups.MustGet(tag).([]*RoutePath)

		paths := newLinkedHashMap() // map[string]map[string]*RoutePath{}
		for _, path := range group {
			route := path.route
			query := make([]string, 0)
			for _, param := range path.params {
				if param.in == "query" {
					query = append(query, param.name)
				}
			}
			if len(query) != 0 {
				route += fmt.Sprintf("{?%s}", strings.Join(query, ","))
			}

			methods, ok := paths.Get(route)
			if !ok {
				methods = newLinkedHashMap() // make(map[string]*RoutePath)
			}
			methods.(*linkedHashMap).Set(path.method, path)
			paths.Set(route, methods)
		}

		for _, pathKey := range paths.Keys() {
			// [###, ###]
			methodStrings := make([]string, 0)
			methods := paths.MustGet(pathKey).(*linkedHashMap)

			summaries := make([]string, 0)
			for _, methodKey := range methods.Keys() {
				routePath := methods.MustGet(methodKey).(*RoutePath)
				summaries = append(summaries, routePath.summary)
				methodStr := buildApibPath(securities, routePath)
				methodStrings = append(methodStrings, methodStr)
			}

			summary := strings.Join(summaries, ", ")
			pathStr := fmt.Sprintf("## %s [%s]\n\n%s", summary, pathKey, strings.Join(methodStrings, "\n\n"))
			pathStrings = append(pathStrings, pathStr)
		}

		groupStr := fmt.Sprintf("# Group %s", tag)
		if tagDesc, ok := tags[tag]; ok {
			groupStr += fmt.Sprintf("\n\n%s", tagDesc)
		}
		if len(pathStrings) > 0 {
			groupStr += fmt.Sprintf("\n\n%s", strings.Join(pathStrings, "\n\n"))
		}

		groupStrings = append(groupStrings, groupStr)
	}

	return strings.Join(groupStrings, "\n\n")
}

func buildApibDefinitions(d *Document) string {
	propertyTypes := make([]string, 0)
	for _, definition := range d.definitions {
		prehandleGenericName(definition) // new name
		if len(definition.generics) == 0 && len(definition.properties) > 0 {
			for _, property := range definition.properties {
				propertyTypes = append(propertyTypes, property.typ)
			}
		}
	}
	for _, path := range d.paths {
		for _, param := range path.params {
			propertyTypes = append(propertyTypes, param.typ)
		}
		for _, response := range path.responses {
			propertyTypes = append(propertyTypes, response.typ)
		}
	}

	newDefinitions := prehandleGenericList(d.definitions, propertyTypes) // new list
	definitionStrings := make([]string, 0)
	for _, definition := range newDefinitions {
		if len(definition.generics) > 0 || len(definition.properties) == 0 {
			continue
		}

		propertyStrings := make([]string, 0)
		for _, property := range definition.properties {
			propertyStr := buildApibProperty(property)
			propertyStrings = append(propertyStrings, propertyStr)
		}

		definitionStr := fmt.Sprintf("## %s (object)", definition.name)
		definitionStr += fmt.Sprintf("\n\n%s", strings.Join(propertyStrings, "\n"))
		definitionStrings = append(definitionStrings, definitionStr)
	}

	return "# Data Structures\n\n" + strings.Join(definitionStrings, "\n\n")
}

var apibTemplate = `FORMAT: 1A
HOST: %s%s

# %s (%s)

%s

%s

%s

%s
`

func buildApibDocument(d *Document) []byte {
	// header
	template := fmt.Sprintf(apibTemplate, d.host, d.basePath, d.info.title, d.info.version, d.info.desc, "%s", "%s", "%s")
	infoArray := make([]string, 0)
	if d.info.termsOfService != "" {
		infoArray = append(infoArray, fmt.Sprintf("[Terms of service](%s)", d.info.termsOfService))
	}
	if license := d.info.license; license != nil {
		infoArray = append(infoArray, fmt.Sprintf("[License: %s](%s)", license.name, license.url))
	}
	if contact := d.info.contact; contact != nil {
		infoArray = append(infoArray, fmt.Sprintf("[%s - Website](%s)", contact.name, contact.url))
		if contact.email != "" {
			infoArray = append(infoArray, fmt.Sprintf("[Send email to %s](mailto:%s)", contact.name, contact.email))
		}
	}
	infoString := strings.Join(infoArray, "\n\n")
	template = fmt.Sprintf(template, infoString, "%s", "%s")

	// route path
	routePathString := buildApibGroups(d)
	template = fmt.Sprintf(template, routePathString, "%s")

	// definition
	definitionString := buildApibDefinitions(d)
	template = fmt.Sprintf(template, definitionString)

	return []byte(template)
}

// GenerateApib generates apib script and writes into file.
func (d *Document) GenerateApib(path string) ([]byte, error) {
	doc := buildApibDocument(d)

	err := saveFile(path, doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// GenerateApib generates apib script and writes into file.
func GenerateApib(path string) ([]byte, error) {
	return _document.GenerateApib(path)
}
