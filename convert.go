package gojson

import (
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"reflect"
	"sort"
	"strings"
	"unicode"
)

// Given a JSON string representation of an object and a name structName,
// attemp to generate a struct definition
func generate(input io.Reader, structName, pkgName string) ([]byte, error) {
	var iresult interface{}
	var result map[string]interface{}
	if err := json.NewDecoder(input).Decode(&iresult); err != nil {
		return nil, err
	}

	switch iresult := iresult.(type) {
	case map[string]interface{}:
		result = iresult
	case []map[string]interface{}:
		if len(iresult) > 0 {
			result = iresult[0]
		} else {
			return nil, fmt.Errorf("empty array")
		}
	default:
		return nil, fmt.Errorf("unexpected type: %T", iresult)
	}

	src := fmt.Sprintf("package %s\ntype %s %s}",
		pkgName,
		structName,
		generateTypes(result, 0))
	formatted, err := format.Source([]byte(src))
	if err != nil {
		err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	}
	return formatted, err
}

// Generate go struct entries for a map[string]interface{} structure
func generateTypes(obj map[string]interface{}, depth int) string {
	structure := "struct {"

	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := obj[key]
		valueType := typeForValue(value)

		//If a nested value, recurse
		switch value := value.(type) {
		case []map[string]interface{}:
			valueType = "[]" + generateTypes(value[0], depth+1) + "}"
		case map[string]interface{}:
			valueType = generateTypes(value, depth+1) + "}"
		}

		fieldName := fmtFieldName(key)
		structure += fmt.Sprintf("\n%s %s `json:\"%s\"`",
			fieldName,
			valueType,
			key)
	}
	return structure
}

var uppercaseFixups = map[string]bool{"id": true, "url": true}

// fmtFieldName formats a string as a struct key
//
// Example:
// 	fmtFieldName("foo_id")
// Output: FooID
func fmtFieldName(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	if len(parts) > 0 {
		last := parts[len(parts)-1]
		if uppercaseFixups[strings.ToLower(last)] {
			parts[len(parts)-1] = strings.ToUpper(last)
		}
	}
	assembled := strings.Join(parts, "")
	runes := []rune(assembled)
	for i, c := range runes {
		ok := unicode.IsLetter(c) || unicode.IsDigit(c)
		if i == 0 {
			ok = unicode.IsLetter(c)
		}
		if !ok {
			runes[i] = '_'
		}
	}
	return string(runes)
}

// generate an appropriate struct type entry
func typeForValue(value interface{}) string {
	//Check if this is an array
	if objects, ok := value.([]interface{}); ok {
		types := make(map[reflect.Type]bool, 0)
		for _, o := range objects {
			types[reflect.TypeOf(o)] = true
		}
		if len(types) == 1 {
			return "[]" + typeForValue(objects[0])
		}
		return "[]interface{}"
	} else if object, ok := value.(map[string]interface{}); ok {
		return generateTypes(object, 0) + "}"
	} else if reflect.TypeOf(value) == nil {
		return "interface{}"
	}
	return reflect.TypeOf(value).Name()
}
