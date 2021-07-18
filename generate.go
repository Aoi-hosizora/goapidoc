package goapidoc

import (
	"fmt"
	"strings"
	"sync/atomic"
)

func checkDocument(doc *Document) {
	if doc.host == "" {
		panic("Host is required")
	}
	if !strings.HasPrefix(doc.basePath, "/") {
		panic("BasePath must begin with a slash")
	}
	if doc.info == nil {
		panic("Info is required")
	}
	if doc.info.title == "" {
		panic("Info title is required")
	}
	if doc.info.version == "" {
		panic("Info version is required")
	}
	if doc.info.license != nil {
		if doc.info.license.name == "" {
			panic("License name is required")
		}
	}

	if doc.option != nil {
		for _, t := range doc.option.tags {
			if t.name == "" {
				panic("Tag name is required")
			}
			if t.externalDocs != nil && t.externalDocs.url == "" {
				panic("Empty externalDocs url is not allowed")
			}
		}
		for _, s := range doc.option.securities {
			if s.typ == APIKEY {
				if s.title == "" {
					panic("Security title is required")
				}
				if s.name == "" {
					panic("Security name is required")
				}
				if s.in == "" {
					panic("Security in-location is required")
				}
			} else if s.typ == BASIC {
				// pass
			} else if s.typ == OAUTH2 {
				if s.flow == "" {
					panic("Security flow is required")
				}
				if (s.flow == IMPLICIT_FLOW || s.flow == ACCESSCODE_FLOW) && s.authorizationUrl == "" {
					panic("Security authorizationUrl is required")
				}
				if (s.flow == PASSWORD_FLOW || s.flow == APPLICATION_FLOW || s.flow == ACCESSCODE_FLOW) && s.tokenUrl == "" {
					panic("Security tokenUrl is required")
				}
				if len(s.scopes) == 0 {
					panic("Empty security scopes is not allowed")
				}
			} else {
				panic("Security type `" + s.typ + "` is not supported")
			}
		}
		if doc.option.externalDocs != nil && doc.option.externalDocs.url == "" {
			panic("Empty externalDocs url is not allowed")
		}
	}

	if len(doc.operations) == 0 {
		panic("Empty operations is not allowed")
	}
	for _, op := range doc.operations {
		if op.method == "" || op.route == "" {
			panic("Operation method and route path is required")
		}
		if !strings.HasPrefix(op.route, "/") {
			panic("Operation route path must begin with a slash")
		}
		if op.summary == "" {
			panic("Operation summary is required")
		}

		for _, p := range op.params {
			if p.name == "" {
				panic("Request param name is required")
			}
			if p.in == "" {
				panic("Request param in-location is required")
			}
			if p.in == PATH && !p.required && !p.allowEmpty {
				panic("Path param's must be non-optional and non-empty")
			}
			if p.typ == "" {
				panic("Request param type is required")
			}
		}

		if len(op.responses) == 0 {
			panic("Empty operation response is not allowed")
		}
		for _, r := range op.responses {
			if r.code == 0 {
				panic("Response code is required")
			}
			for _, h := range r.headers {
				if h.name == "" {
					panic("Response header field name is required")
				}
				if h.typ == "" {
					panic("Response header type is required")
				}
			}
		}
		if op.externalDocs != nil && op.externalDocs.url == "" {
			panic("Empty externalDocs url is not allowed")
		}
	}

	for _, def := range doc.definitions {
		if def.name == "" {
			panic("Definition name is required")
		}
		for _, p := range def.properties {
			if p.name == "" {
				panic("Definition property name is required")
			}
			if p.typ == "" {
				panic("Definition property type is required")
			}
		}
	}
}

// GenerateSwaggerYaml generates swagger yaml script and returns byte array.
// TODO BREAK CHANGES
func (d *Document) GenerateSwaggerYaml() ([]byte, error) {
	doc := buildSwagDocument(d)
	return yamlMarshal(doc)
}

// GenerateSwaggerJson generates swagger json script and returns byte array.
// TODO BREAK CHANGES
func (d *Document) GenerateSwaggerJson() ([]byte, error) {
	doc := buildSwagDocument(d)
	return jsonMarshal(doc)
}

// GenerateApib generates apib script and returns byte array.
// TODO BREAK CHANGES
func (d *Document) GenerateApib() ([]byte, error) {
	return buildApibDocument(d)
}

// SaveSwaggerYaml generates swagger yaml script and saves into file.
func (d *Document) SaveSwaggerYaml(path string) ([]byte, error) {
	bs, err := d.GenerateSwaggerYaml()
	if err != nil {
		return nil, err
	}
	err = saveFile(path, bs)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// SaveSwaggerJson generates swagger json script and saves into file.
func (d *Document) SaveSwaggerJson(path string) ([]byte, error) {
	bs, err := d.GenerateSwaggerJson()
	if err != nil {
		return nil, err
	}
	err = saveFile(path, bs)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// SaveApib generates apib script and saves into file.
func (d *Document) SaveApib(path string) ([]byte, error) {
	bs, err := d.GenerateApib()
	if err != nil {
		return nil, err
	}
	err = saveFile(path, bs)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// GenerateSwaggerYaml generates swagger yaml script and returns byte array.
func GenerateSwaggerYaml() ([]byte, error) {
	return _document.GenerateSwaggerYaml()
}

// GenerateSwaggerJson generates swagger json script and returns byte array.
func GenerateSwaggerJson() ([]byte, error) {
	return _document.GenerateSwaggerJson()
}

// GenerateApib generates apib script and returns byte array.
func GenerateApib() ([]byte, error) {
	return _document.GenerateApib()
}

// SaveSwaggerYaml generates swagger yaml script and saves into file.
func SaveSwaggerYaml(path string) ([]byte, error) {
	return _document.SaveSwaggerYaml(path)
}

// SaveSwaggerJson generates swagger json script and saves into file.
func SaveSwaggerJson(path string) ([]byte, error) {
	return _document.SaveSwaggerJson(path)
}

// SaveApib generates apib script and saves into file.
func SaveApib(path string) ([]byte, error) {
	return _document.SaveApib(path)
}

// _warningLogger is a global switcher for logger when warning.
var _warningLogger atomic.Value

func init() {
	_warningLogger.Store(true)
}

// DisableWarningLogger disables the warning logger switcher.
func DisableWarningLogger() {
	_warningLogger.Store(false)
}

// EnableWarningLogger enables the warning logger switcher.
func EnableWarningLogger() {
	_warningLogger.Store(true)
}

// logWarning logs the warning message.
func logWarning(s string) {
	if _warningLogger.Load().(bool) {
		fmt.Printf("Warning: %s\n", s)
	}
}
