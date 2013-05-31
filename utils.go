package postmaster

import (
	"fmt"
	"net/url"
	"strings"
	"reflect"
)


// urlencode joins parameters from map[string]string with ampersand (&), and
// also escapes their values
func urlencode(params map[string]string) string {
	arr := make([]string, len(params))
	for k, v := range params {
		if fmt.Sprintf("%s", v) != "" {
			arr = append(arr, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
		}
	}
	return strings.Join(arr, "&")
}


// MapStruct converts struct to map[string]string, using fields' names as keys
// and fields' values as values.
// It also automagically converts any nested structures.
func MapStruct(s interface{}) map[string]string {
	return mapStructNested(s, "")
}


// mapStructNested does all the dirty job that mapStruct was too lazy to do.
func mapStructNested(s interface{}, baseName string) map[string]string {
	result := make(map[string]string)
	// Is s a pointer? We don't want any of those here
	if reflect.TypeOf(s).Kind() == reflect.Ptr {
		s = reflect.TypeOf(s).Elem()
	}
	fields := reflect.TypeOf(s).NumField()
	for i := 0; i < fields; i++ {
		t := reflect.TypeOf(s).Field(i)
		v := reflect.ValueOf(s).Field(i)
		// Name is important
		var name string
		if json := t.Tag.Get("json"); json != "" {
			name = json
		} else {
			name = strings.ToLower(t.Name)
		}
		if baseName != "" {
			name = fmt.Sprintf("%s[%s]", baseName, name)
		}
		// I wonder whether this is a nested object
		if v.Kind() == reflect.Struct { // Nested, activate recursion!
			m := mapStructNested(v.Interface(), name)
			for mk, mv := range m {
				result[mk] = mv
			}
		} else { // Not nested
			value := fmt.Sprintf("%v", v.Interface())
			if value != "" {
				result[name] = value
			}
		}
	}
	return result
}
