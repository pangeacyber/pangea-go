package pangeautil

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"reflect"
	"sort"
	"strings"
	"time"
)

func CanonicalizeStruct(v interface{}) ([]byte, error) {
	var smap map[string]interface{}
	ebytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	// Order keys
	json.Unmarshal(ebytes, &smap)
	mbytes, err := json.Marshal(smap)
	if err != nil {
		return nil, err
	}
	return mbytes, nil
}

type PangeaTimestamp time.Time

const ptLayout_Z = "2006-01-02T15:04:05.000000Z"
const ptLayout_noZ = "2006-01-02T15:04:05.000000"

// UnmarshalJSON Parses the json string in the custom format
func (ct *PangeaTimestamp) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(ptLayout_Z, s)
	if err == nil {
		*ct = PangeaTimestamp(nt)
		return nil
	}
	nt, err = time.Parse(ptLayout_noZ, s)
	if err == nil {
		*ct = PangeaTimestamp(nt)
		return
	}
	return err
}

// MarshalJSON writes a quoted string in the custom format
func (ct PangeaTimestamp) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

// String returns the time in the custom format
func (ct *PangeaTimestamp) String() string {
	t := time.Time(*ct)
	utc := t.UTC()
	return fmt.Sprintf("%q", utc.Format(ptLayout_Z))
}

// CanonicalizeJSONMarshall is not a true canoni
func CanonicalizeJSONMarshall(v interface{}) []byte {
	buf := new(bytes.Buffer)
	canonicalizeJSONMarshall(reflect.ValueOf(v), buf)
	return buf.Bytes()
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
				val := v.FieldByName(field.Name)
				if val.Kind() == reflect.Ptr && val.IsNil() {
					continue
				}
				jsonTags = append(jsonTags, name)
				tagKeyRealtion[name] = field.Name
			} else {
				continue // ignore non json tagged fields
			}
		}
		sort.Strings(jsonTags)
		for i, n := range jsonTags {
			val := v.FieldByName(tagKeyRealtion[n])
			if val.Kind() == reflect.Ptr && val.IsNil() {
				continue
			}
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

// Sleep t time, but also listen to ctx Done() signal
// Return true if timeout was reached correctly, false if it was interrupted by ctx signal
func Sleep(t time.Duration, ctx context.Context) bool {
	select {
	case <-ctx.Done(): //context cancelled
		return false
	case <-time.After(t): //timeout
		return true
	}
}

func GetBoundary(contentType string) (string, error) {
	// Parse the Content-Type header
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", err
	}

	// Check if it's a multipart media type
	if !strings.HasPrefix(mediaType, "multipart/") {
		return "", err
	}

	// Extract the boundary parameter
	boundary, ok := params["boundary"]
	if !ok {
		return "", errors.New("Boundary parameter not found in Content-Type")
	}

	return boundary, nil
}

func GetFilename(contentDisposition string) (string, error) {
	if contentDisposition == "" {
		return "", fmt.Errorf("Content-Disposition header is empty")
	}

	// Split the header into parts
	parts := strings.Split(contentDisposition, ";")

	// Search for the "filename" parameter
	for _, part := range parts {
		if strings.Contains(part, "filename") {
			// Extract the filename
			filenamePart := strings.Split(part, "=")
			if len(filenamePart) == 2 {
				return strings.Trim(filenamePart[1], "\" "), nil
			}
		}
	}

	return "", fmt.Errorf("Filename not found in Content-Disposition header")
}
