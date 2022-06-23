package pangeautil

import (
	"encoding/json"
	"strings"
)

// Stringify returns the string representation of a json object.
func Stringify(obj interface{}) string {
	b := new(strings.Builder)
	if err := json.NewEncoder(b).Encode(obj); err != nil {
		return ""
	}
	return b.String()
}
