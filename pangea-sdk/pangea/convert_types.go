package pangea

import (
	b64 "encoding/base64"
	"encoding/json"
	"strings"
	"time"

	pu "github.com/pangeacyber/pangea-go/pangea-sdk/v4/internal/pangeautil"
)

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }

// StringValue is a helper routine that returns the value of a string pointer or a default value if nil
func StringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

// IntValue is a helper routine that returns the value of a int pointer or a default value if nil
func IntValue(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

// BoolValue is a helper routine that returns the value of a bool pointer or a default value if nil
func BoolValue(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

// Stringify returns the string representation of a json object.
func Stringify(obj interface{}) string {
	b := new(strings.Builder)
	if err := json.NewEncoder(b).Encode(obj); err != nil {
		return ""
	}
	return b.String()
}

// Stringify returns the string representation of a json object.
func StringifyIndented(obj interface{}) string {
	b := new(strings.Builder)
	enc := json.NewEncoder(b)
	enc.SetIndent("", "  ")
	if err := enc.Encode(obj); err != nil {
		return ""
	}
	return b.String()
}

// Time is a helper routine that allocates a new time.Time value
// to store v and returns a pointer to it.
func Time(v time.Time) *time.Time {
	return &v
}

// Time is a helper routine that allocates a new PangeaTimestamp value
// to store v and returns a pointer to it.
func PangeaTime(v pu.PangeaTimestamp) *pu.PangeaTimestamp {
	return &v
}

func StrToB64(dec string) string {
	return b64.StdEncoding.EncodeToString([]byte(dec))
}

func B64ToStr(enc string) ([]byte, error) {
	return b64.StdEncoding.DecodeString(enc)
}
