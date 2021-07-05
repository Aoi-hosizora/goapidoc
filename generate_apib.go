package goapidoc

import (
	"fmt"
	"strings"
)

func buildApibType(typ string) (string, *apiType) {
	at := parseApiType(typ)
	// typ = strings.ReplaceAll(strings.ReplaceAll(typ, "<", "«"), ">", "»")
	if at.kind == apiPrimeKind {
		t := at.prime.typ
		if t == "integer" {
			t = "number"
		}
		return t, at
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

func buildApibParameter(param *Param) string {
	req := "required"
	if !param.required {
		req = "optional"
	}
	typ, at := buildApibType(param.typ)
	if len(param.enums) != 0 {
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
		options = append(options, fmt.Sprintf("val in \\[%.3f, %.3f\\]", param.minimum, param.maximum))
	} else if param.minimum != 0 {
		options = append(options, fmt.Sprintf("val >= %.3f", param.minimum))
	} else if param.maximum != 0 {
		options = append(options, fmt.Sprintf("val <= %.3f", param.maximum))
	}

	if len(options) != 0 {
		paramStr += fmt.Sprintf("\n    (%s)", strings.Join(options, ", "))
	}

	if param.defaul != nil {
		paramStr += fmt.Sprintf("\n    + Default: `%v`", param.defaul)
	}

	if len(param.enums) != 0 {
		paramStr += "\n    + Members"
		for _, enum := range param.enums {
			paramStr += fmt.Sprintf("\n        + `%v`", enum)
		}
	}

	return paramStr
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
		securityString := path.securities[0] // only support one authentication
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
		paramStr := buildApibParameter(param)
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

	// tag - RoutePath (route&method)
	groups := newOrderedMap(len(tags)) // map[string][]*RoutePath{}
	for name := range tags {
		groups.Set(name, make([]*RoutePath, 0))
	}
	for _, path := range doc.paths {
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

	// [#, #]
	groupStrings := make([]string, 0, groups.Length())

	for _, tag := range groups.Keys() {
		group := groups.MustGet(tag).([]*RoutePath)
		// route - method - RoutePath
		paths := newOrderedMap(len(group)) // map[string]map[string]*RoutePath{}
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
				methods = newOrderedMap(0) // map[string]*RoutePath
			}
			methods.(*orderedMap).Set(path.method, path)
			paths.Set(route, methods)
		}

		// [##, ##]
		pathStrings := make([]string, 0, paths.Length())
		for _, pathKey := range paths.Keys() {
			methods := paths.MustGet(pathKey).(*orderedMap)

			// [###, ###]
			methodStrings := make([]string, 0, methods.Length())
			summaries := make([]string, 0, methods.Length())
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

func buildApibDefinitions(doc *Document) string {
	// check and collect type names
	allSpecTypes := collectAllSpecTypes(doc)

	// prehandle cloned definition list
	clonedDefinitions := make([]*Definition, 0, len(doc.definitions))
	for _, definition := range doc.definitions {
		clonedDefinitions = append(clonedDefinitions, prehandleDefinition(definition)) // with generic name checked
	}
	newDefinitions := prehandleDefinitionList(clonedDefinitions, allSpecTypes)

	// render result string
	definitionStrings := make([]string, 0)
	for _, definition := range newDefinitions {
		if len(definition.generics) > 0 || len(definition.properties) == 0 {
			continue
		}

		propertyStrings := make([]string, 0)
		for _, property := range definition.properties {
			propertyStr := buildApibParameter(cloneParamFromProperty(property))
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

func buildApibDocument(doc *Document) []byte {
	if doc.info == nil {
		panic("Nil document info")
	}

	// header
	template := fmt.Sprintf(apibTemplate, doc.host, doc.basePath, doc.info.title, doc.info.version, doc.info.desc, "%s", "%s", "%s")
	infoArray := make([]string, 0, 4)
	if doc.info.termsOfService != "" {
		infoArray = append(infoArray, fmt.Sprintf("[Terms of service](%s)", doc.info.termsOfService))
	}
	if license := doc.info.license; license != nil {
		infoArray = append(infoArray, fmt.Sprintf("[License: %s](%s)", license.name, license.url))
	}
	if contact := doc.info.contact; contact != nil {
		infoArray = append(infoArray, fmt.Sprintf("[%s - Website](%s)", contact.name, contact.url))
		if contact.email != "" {
			infoArray = append(infoArray, fmt.Sprintf("[Send email to %s](mailto:%s)", contact.name, contact.email))
		}
	}
	infoString := strings.Join(infoArray, "\n\n")
	template = fmt.Sprintf(template, infoString, "%s", "%s")

	// definition
	definitionsString := buildApibDefinitions(doc)
	template = fmt.Sprintf(template, "%s", definitionsString)

	// route path
	groupsString := buildApibGroups(doc)
	template = fmt.Sprintf(template, groupsString)

	return []byte(template)
}
