package goapidoc

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unsafe"
)

func spaceIndent(indent int, s string) string {
	spaces := strings.Repeat("    ", indent)
	return spaces + strings.ReplaceAll(s, "\n", "\n"+spaces)
}

func fastBtos(bs []byte) string {
	// unsafe !!!
	return *(*string)(unsafe.Pointer(&bs))
}

func renderTemplate(t string, object interface{}) []byte {
	tmpl, err := template.New("template").Parse(t)
	if err != nil {
		panic("Template rendering error: " + err.Error())
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, object)
	if err != nil {
		panic("Template rendering error: " + err.Error())
	}
	return buf.Bytes()
}

func yamlMarshal(t interface{}) ([]byte, error) {
	return yaml.Marshal(t)
}

func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func xmlMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := xml.NewEncoder(buffer)
	encoder.Indent("", "  ")
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func bsErrToStrErr(bs []byte, err error) (string, error) {
	return string(bs), err
}

func saveFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0644)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(path, data, 0644)
}

// orderedMap represents an ordered hashmap, is used to replace map[K]V.
type orderedMap struct {
	m map[string]interface{}
	i []string
}

// newOrderedMap creates an empty orderedMap with cap.
func newOrderedMap(cap int) *orderedMap {
	return &orderedMap{
		m: make(map[string]interface{}, cap),
		i: make([]string, 0, cap),
	}
}

func (l *orderedMap) Length() int {
	return len(l.i)
}

func (l *orderedMap) Keys() []string {
	return l.i
}

func (l *orderedMap) Set(key string, value interface{}) {
	_, exist := l.m[key]
	l.m[key] = value
	if !exist {
		l.i = append(l.i, key)
	}
}

func (l *orderedMap) Get(key string) (interface{}, bool) {
	v, ok := l.m[key]
	return v, ok
}

func (l *orderedMap) MustGet(key string) interface{} {
	val, ok := l.Get(key)
	if !ok {
		panic("key " + key + " is not found.")
	}
	return val
}

func (l *orderedMap) MarshalJSON() ([]byte, error) {
	length := len(l.i)
	buf := &bytes.Buffer{}
	buf.WriteString("{")
	for idx, k := range l.i {
		b, err := jsonMarshal(l.m[k])
		if err != nil {
			return nil, err
		}
		buf.WriteString(fmt.Sprintf("  \"%s\": %s", k, string(b)))
		if idx < length-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (l *orderedMap) MarshalYAML() (interface{}, error) {
	ms := yaml.MapSlice{}
	for _, k := range l.i {
		ms = append(ms, yaml.MapItem{Key: k, Value: l.m[k]})
	}
	return ms, nil
}
