package goapidoc

import (
	"errors"
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

func TestMarshal(t *testing.T) {
	type DemoStruct struct {
		Code    int
		Message string
		Data    float64
	}

	for _, tc := range []struct {
		name     string
		give     interface{}
		wantJson string
		wantYaml string
		wantErr  bool
	}{
		{"nil", nil, "null\n", "null\n", false},
		{"int", 1025, "1025\n", "1025\n", false},
		{"float", 1.56, "1.56\n", "1.56\n", false},
		{"string", "demo string", "\"demo string\"\n", "demo string\n", false},
		{"string2", "1 < 2", "\"1 < 2\"\n", "1 < 2\n", false},
		{"array_of_int", []int{1, 2, 3},
			"[\n  1,\n  2,\n  3\n]\n", "- 1\n- 2\n- 3\n", false},
		{"array_of_string", []string{"1", "2", "3"},
			"[\n  \"1\",\n  \"2\",\n  \"3\"\n]\n", "- \"1\"\n- \"2\"\n- \"3\"\n", false},
		{"map_of_string_int", map[string]int{"item1": 1, "item2": 2, "item3": 3},
			"{\n  \"item1\": 1,\n  \"item2\": 2,\n  \"item3\": 3\n}\n", "item1: 1\nitem2: 2\nitem3: 3\n", false},
		{"map_of_int_string", map[int]string{1: "item1", 2: "item2", 3: "item3"},
			"{\n  \"1\": \"item1\",\n  \"2\": \"item2\",\n  \"3\": \"item3\"\n}\n", "1: item1\n2: item2\n3: item3\n", false},
		{"empty_struct", struct{}{}, "{}\n", "{}\n", false},
		{"struct", DemoStruct{200, "Success", 3.14159},
			"{\n  \"Code\": 200,\n  \"Message\": \"Success\",\n  \"Data\": 3.14159\n}\n", "code: 200\nmessage: Success\ndata: 3.14159\n", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			j, err := jsonMarshal(tc.give)
			if string(j) != tc.wantJson {
				failNow(t, "jsonMarshal get an unexpected result for value")
			}
			if (err == nil && tc.wantErr) || (err != nil && !tc.wantErr) {
				failNow(t, "jsonMarshal get an unexpected result for error")
			}
			y, err := yamlMarshal(tc.give)
			if string(y) != tc.wantYaml {
				failNow(t, "yamlMarshal get an unexpected result for value")
			}
			if (err == nil && tc.wantErr) || (err != nil && !tc.wantErr) {
				failNow(t, "yamlMarshal get an unexpected result for error")
			}
		})
	}
}

func TestBsErrToStrErr(t *testing.T) {
	for _, tc := range []struct {
		name    string
		giveBs  []byte
		giveErr error
		wantStr string
		wantErr error
	}{
		{"empty string, nil error", nil, nil, "", nil},
		{"string, nil error", []byte("test test"), nil, "test test", nil},
		{"empty string, error", nil, errors.New("error error"), "", errors.New("error error")},
		{"string, error", []byte("test test"), errors.New("error error"), "test test", errors.New("error error")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			str, err := bsErrToStrErr(tc.giveBs, tc.giveErr)
			if str != tc.wantStr {
				failNow(t, "bsErrToStrErr get an unexpected result")
			}
			if (err == nil && tc.wantErr != nil) || (err != nil && tc.wantErr == nil) || (err != nil && err.Error() != tc.wantErr.Error()) {
				failNow(t, "bsErrToStrErr get an unexpected result")
			}
		})
	}
}
