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

type linkedHashMap struct {
	m map[string]interface{}
	i []string
}

func newLinkedHashMap() *linkedHashMap {
	return &linkedHashMap{
		m: make(map[string]interface{}),
		i: make([]string, 0),
	}
}

func (l *linkedHashMap) Keys() []string {
	return l.i
}

func (l *linkedHashMap) Has(key string) bool {
	_, exist := l.m[key]
	return exist
}

func (l *linkedHashMap) Set(key string, value interface{}) {
	_, exist := l.m[key]
	l.m[key] = value
	if !exist {
		l.i = append(l.i, key)
	}
}

func (l *linkedHashMap) Get(key string) (interface{}, bool) {
	v, ok := l.m[key]
	return v, ok
}

func (l *linkedHashMap) MustGet(key string) interface{} {
	val, ok := l.Get(key)
	if !ok {
		panic("key " + key + " is not found.")
	}
	return val
}

func (l *linkedHashMap) MarshalJSON() ([]byte, error) {
	ov := make([]interface{}, len(l.i))
	for idx, field := range l.i {
		ov[idx] = l.m[field]
	}

	buf := &bytes.Buffer{}
	buf.WriteString("{")
	for idx, field := range l.i {
		b, err := json.Marshal(ov[idx])
		if err != nil {
			return []byte{}, err
		}
		buf.WriteString(fmt.Sprintf("  \"%s\": %s", field, string(b)))
		if idx < len(l.i)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString("}")

	return buf.Bytes(), nil
}

func (l *linkedHashMap) MarshalYAML() (interface{}, error) {
	return l.m, nil
}
