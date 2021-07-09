package goapidoc

import (
	"testing"
)

func TestSpaceIndent(t *testing.T) {
	for _, tc := range []struct {
		name       string
		giveIndent int
		giveString string
		want       string
	}{
		{"empty", 0, "", ""},
		{"test_0", 0, "test", "test"},
		{"test_1", 1, "test", "    test"},
		{"test_4", 4, "test", "                test"},
		{"test|test_1", 1, "test\ntest", "    test\n    test"},
		{"test|----x_1", 1, "test\n    x", "    test\n        x"},
		{"test|x|_2", 2, "test\nx\n", "        test\n        x\n        "},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if spaceIndent(tc.giveIndent, tc.giveString) != tc.want {
				failNow(t, "spaceIndent get an unexpected result")
			}
		})
	}
}

func TestFastBtos(t *testing.T) {
	for _, tc := range []struct {
		name string
		give []byte
		want string
	}{
		{"empty", []byte{}, ""},
		{"a", []byte{'a'}, "a"},
		{"hello", []byte{'h', 'e', 'l', 'l', 'o'}, "hello"},
		{"a b c", []byte{'a', ' ', 'b', ' ', 'c'}, "a b c"},
		{"测试", []byte("测试"), "测试"},
		{"テス", []byte("テス"), "テス"},
		{"тест", []byte("тест"), "тест"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if fastBtos(tc.give) != tc.want {
				failNow(t, "fastBtos get an unexpected result")
			}
		})
	}
}

func TestRenderTemplate(t *testing.T) {
	for _, tc := range []struct {
		name         string
		giveTemplate string
		giveObject   interface{}
		wantPanic    bool
		want         string
	}{
		{"wrong template", "{{ /* xxx */ }}", nil, true, ""},
		{"wrong object 1", "{{ .Test }}", "test", true, ""},
		{"wrong object 2", "{{ .Test }}", struct{ X string }{"test"}, true, ""},
		{"nil object", "{{ .Test }}", nil, false, "<no value>"},
		{"normal value", "{{ . }}", "hello world", false, "hello world"},
		{"normal struct", "Test: {{ .Test }}", struct{ Test string }{"hello world"}, false, "Test: hello world"},
		{"normal array", "{{ range . }}{{ .Test }}{{ end }}", []struct{ Test string }{{"a"}, {"b"}, {"c"}}, false, "abc"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testPanic(t, tc.wantPanic, func() {
				if string(renderTemplate(tc.giveTemplate, tc.giveObject)) != tc.want {
					failNow(t, "renderTemplate get an unexpected result")
				}
			}, "testPanic")
		})
	}
}

// func TestYamlMarshal(t *testing.T) {
//
// }
//
// func TestJsonMarshal(t *testing.T) {
//
// }
