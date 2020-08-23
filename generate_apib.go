package goapidoc

import (
	"fmt"
	"strings"
)

func buildApibPath(path *RoutePath) string {
	consume := "application/json"
	if len(path.consumes) >= 1 {
		consume = path.consumes[0]
	}
	produce := "application/json"
	if len(path.produces) >= 1 {
		consume = path.produces[0]
	}

	requestStr := "+ Request (" + consume + ")\n\n    + Body"
	responseStr := "+ Response 200 (" + produce + ")\n\n    + Body"
	return requestStr + "\n\n" + responseStr
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

	for _, group := range groups {
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

			for method, route := range methods {
				method = strings.ToUpper(method)
				methodStr := fmt.Sprintf("### %s [%s]", route.summary, method)
				if route.desc != "" {
					methodStr += "\n\n" + route.desc
				}
				methodStr += "\n\n" + buildApibPath(route)
				methodStrings = append(methodStrings, methodStr)
			}

			pathStr := fmt.Sprintf("## XXX [%s]\n\n%s", path, strings.Join(methodStrings, "\n\n"))
			pathStrings = append(pathStrings, pathStr)
		}

		groupStr := fmt.Sprintf("# Group XXX\n\n%s", strings.Join(pathStrings, "\n\n"))
		groupStrings = append(groupStrings, groupStr)
	}

	cmt := "<!-- GROUPS -->"
	return fmt.Sprintf("%s\n\n%s", cmt, strings.Join(groupStrings, "\n\n"))
}

func buildApibDefinitions(d *Document) string {
	cmt := "<!-- DEFINITIONS -->"
	return cmt
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
