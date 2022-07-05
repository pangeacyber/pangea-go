package pangeautil

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// CanonicalizeJSONMarshall is not a true canoni
func CanonicalizeJSONMarshall(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	canonicalizeJSONMarshall(reflect.ValueOf(v), buf)
	return buf.Bytes(), nil
}

func canonicalizeJSONMarshall(v reflect.Value, buf *bytes.Buffer) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			buf.WriteString("null")
			return
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		buf.WriteString("{")
		jsonTags := []string{}
		tagKeyRealtion := make(map[string]string, 0)
		for i := 0; i < v.Type().NumField(); i++ {
			field := v.Type().Field(i)
			if field.Name[0:1] == strings.ToLower(field.Name[0:1]) {
				continue // ignore unexported fields
			}
			if tag := field.Tag.Get("json"); tag != "" {
				name, _, _ := strings.Cut(tag, ",")
				jsonTags = append(jsonTags, name)
				tagKeyRealtion[name] = field.Name
			} else {
				continue // ignore non json tagged fields
			}
		}
		sort.Strings(jsonTags)
		for i, n := range jsonTags {
			val := v.FieldByName(tagKeyRealtion[n])
			fmt.Fprintf(buf, `"%v":`, n)
			canonicalizeJSONMarshall(val, buf)
			if i < len(jsonTags)-1 {
				buf.WriteString(",")
			}
		}
		buf.WriteString("}")
	default:
		fmt.Fprintf(buf, `"%v"`, v.Interface())
	}
}
