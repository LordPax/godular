package core

// inArray checks if a string is present in a slice of strings.
func inArray(s string, arr []string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}
