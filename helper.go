package goapidoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func saveFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0644)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
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

func (l *orderedMap) Keys() []string {
	return l.i
}

func (l *orderedMap) Has(key string) bool {
	_, exist := l.m[key]
	return exist
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
	ov := make([]interface{}, len(l.i))
	for idx, field := range l.i {
		ov[idx] = l.m[field]
	}

	buf := &bytes.Buffer{}
	buf.WriteString("{")
	for idx, field := range l.i {
		b, err := json.Marshal(ov[idx])
		if err != nil {
			return nil, err
		}
		buf.WriteString(fmt.Sprintf("  \"%s\": %s", field, string(b)))
		if idx < len(l.i)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString("}")

	return buf.Bytes(), nil
}

func (l *orderedMap) MarshalYAML() (interface{}, error) {
	return l.m, nil
}
