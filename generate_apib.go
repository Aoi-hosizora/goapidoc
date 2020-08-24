package goapidoc

import (
	"fmt"
	"strings"
)

func buildApibParam(param *Param) string {
	req := "required"
	if !param.required {
		req = "optional"
	}
	typ := param.typ
	typ = strings.ReplaceAll(strings.ReplaceAll(typ, "<", "«"), ">", "»")
	typ = strings.ReplaceAll(typ, "[]", "||")
	if len(param.enum) != 0 {
		typ = "enum[" + typ + "]"
	}

	paramStr := fmt.Sprintf("+ %s (%s, %s) - %s", param.name, typ, req, param.desc) // center
	if param.example != nil {
		paramStr = fmt.Sprintf("+ %s: `%v` (%s, %s) - %s", param.name, param.example, param.typ, req, param.desc)
	}

	options := make([]string, 0)

	if param.allowEmpty {
		options = append(options, "allow empty")
	}
	if param.maxLength != 0 && param.minLength != 0 {
		options = append(options, fmt.Sprintf("len in \\[%d, %d\\]", param.minLength, param.maxLength))
	} else if param.minLength != 0 {
		options = append(options, fmt.Sprintf("len >= %d", param.minLength))
	} else {
		options = append(options, fmt.Sprintf("len <= %d", param.maxLength))
	}
	if param.maximum != 0 && param.minimum != 0 {
		options = append(options, fmt.Sprintf("val in \\[%d, %d\\]", param.minimum, param.maximum))
	} else if param.minLength != 0 {
		options = append(options, fmt.Sprintf("val >= %d", param.minimum))
	} else {
		options = append(options, fmt.Sprintf("val <= %d", param.maximum))
	}

	if len(options) != 0 {
		paramStr += fmt.Sprintf("\n\n        (%s)", strings.Join(options, ", "))
	}

	if param.def != nil {
		paramStr += fmt.Sprintf("\n\n        Default: `%v`", param.def)
	}

	if len(param.enum) != 0 {
		paramStr += "\n\n        + Members"
		for _, enum := range param.enum {
			paramStr += fmt.Sprintf("\n            + `%v`", enum)
		}
	}

	return paramStr
}

func buildApibPath(path *RoutePath) string {
	// request
	parameterStrings := make([]string, 0)
	attributeStrings := make([]string, 0)
	for _, param := range path.params {
		paramStr := buildApibParam(param)
		if param.in == "path" || param.in == "query" {
			parameterStrings = append(parameterStrings, "    "+paramStr)
		} else { // body || formData
			paramStr = strings.ReplaceAll(paramStr, "\n", "\n    ")
			attributeStrings = append(attributeStrings, "        "+paramStr)
		}
	}
	consume := "application/json"
	if len(path.consumes) >= 1 {
		consume = path.consumes[0]
	}

	requestStr := fmt.Sprintf("+ Request (%s)", consume) // center
	if len(parameterStrings) != 0 {
		parameterStrings := strings.Join(parameterStrings, "\n\n")
		requestStr = fmt.Sprintf("+ Parameter\n\n%s\n\n%s", parameterStrings, requestStr)
	}
	if len(attributeStrings) != 0 {
		attributeStrings := strings.Join(attributeStrings, "\n\n")
		requestStr += fmt.Sprintf("\n\n    + Attributes\n\n%s", attributeStrings)
	}
	requestStr += "\n\n    + Body"

	// response
	responseStrings := make([]string, 0)
	produce := "application/json"
	if len(path.produces) >= 1 {
		consume = path.produces[0]
	}

	for _, response := range path.responses {
		responseStr := fmt.Sprintf("+ Response %d (%s)", response.code, produce) // center
		if response.desc != "" {
			responseStr += fmt.Sprintf("\n\n    %s", response.desc)
		}
		headerStrings := make([]string, 0)
		for _, header := range response.headers {
			headerStr := fmt.Sprintf("            %s: %s (%s)", header.name, header.desc, header.typ)
			headerStrings = append(headerStrings, headerStr)
		}
		if len(headerStrings) != 0 {
			responseStr += fmt.Sprintf("\n\n    + Headers\n\n%s", strings.Join(headerStrings, "\n"))
		}
		responseStr += "\n\n    + Body"
		if ex, ok := response.examples[produce]; ok {
			ex = "\n" + ex
			ex = strings.ReplaceAll(ex, "\n", "\n            ")
			responseStr += "\n" + ex
		}
		responseStrings = append(responseStrings, responseStr)
	}

	return requestStr + "\n\n" + strings.Join(responseStrings, "\n\n")
}

func buildApibGroups(d *Document) string {
	// [#, #]
	groupStrings := make([]string, 0)

	groups := map[string][]*RoutePath{}
	for _, path := range d.paths {
		tags := path.tags
		if len(tags) == 0 {
			tags = []string{"Default"}
		}
		tag := tags[0]
		if _, ok := groups[tag]; !ok {
			groups[tag] = make([]*RoutePath, 0)
		}
		groups[tag] = append(groups[tag], path)
	}
	tags := make(map[string]string)
	for _, tag := range d.tags {
		tags[tag.name] = tag.desc
	}

	for tag, group := range groups {
		// [##, ##]
		pathStrings := make([]string, 0)

		paths := map[string]map[string]*RoutePath{}
		for _, path := range group {
			if _, ok := paths[path.route]; !ok {
				paths[path.route] = map[string]*RoutePath{}
			}
			paths[path.route][path.method] = path
		}

		for path, methods := range paths {
			// [###, ###]
			methodStrings := make([]string, 0)

			summaries := make([]string, 0)
			for method, route := range methods {
				method = strings.ToUpper(method)
				summaries = append(summaries, route.summary)
				methodStr := fmt.Sprintf("### %s [%s]\n\n`%s %s`", route.summary, method, method, route.route)
				if route.desc != "" {
					methodStr += "\n\n" + route.desc
				}
				methodStr += "\n\n" + buildApibPath(route)
				methodStrings = append(methodStrings, methodStr)
			}

			summary := strings.Join(summaries, ", ")
			pathStr := fmt.Sprintf("## %s [%s]\n\n%s", summary, path, strings.Join(methodStrings, "\n\n"))
			pathStrings = append(pathStrings, pathStr)
		}

		groupStr := fmt.Sprintf("# Group %s", tag)
		if tagDesc, ok := tags[tag]; ok {
			groupStr += fmt.Sprintf("\n\n%s", tagDesc)
		}
		groupStr += fmt.Sprintf("\n\n%s", strings.Join(pathStrings, "\n\n"))

		groupStrings = append(groupStrings, groupStr)
	}

	return fmt.Sprintf("<!-- GROUPS -->\n\n%s", strings.Join(groupStrings, "\n\n"))
}

func buildApibDefinitions(d *Document) string {
	return "<!-- DEFINITIONS -->"
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

func (d *Document) GenerateApib(path string) ([]byte, error) {
	doc := buildApibDocument(d)

	err := saveFile(path, doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func GenerateApib(path string) ([]byte, error) {
	return _document.GenerateApib(path)
}
