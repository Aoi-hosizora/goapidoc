package goapidoc

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type LinkedHashMap struct {
	m map[string]interface{}
	i []string
}

func NewLinkedHashMap() *LinkedHashMap {
	return &LinkedHashMap{
		m: make(map[string]interface{}),
		i: make([]string, 0),
	}
}

func (l *LinkedHashMap) Set(key string, value interface{}) {
	_, exist := l.m[key]
	l.m[key] = value
	if !exist {
		l.i = append(l.i, key)
	}
}

func (l *LinkedHashMap) MarshalJSON() ([]byte, error) {
	ov := make([]interface{}, len(l.i))
	for idx, field := range l.i {
		ov[idx] = l.m[field]
	}

	buf := &bytes.Buffer{}
	buf.Write([]byte{'{'})
	for idx, field := range l.i {
		b, err := json.Marshal(ov[idx])
		if err != nil {
			return []byte{}, err
		}
		buf.Write([]byte(fmt.Sprintf("  \"%s\": %s", field, string(b))))
		if idx < len(l.i)-1 {
			buf.Write([]byte(","))
		}
	}
	buf.Write([]byte{'}'})
	return []byte(buf.String()), nil
}
