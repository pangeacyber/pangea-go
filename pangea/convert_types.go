package pangea

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
