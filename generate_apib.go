package goapidoc

import (
	"fmt"
	"strings"
)

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
	routePathString := "<!-- ROUTE PATH -->"
	template = fmt.Sprintf(template, routePathString, "%s")

	// definition
	definitionString := "<!-- DEFINITION -->"
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
