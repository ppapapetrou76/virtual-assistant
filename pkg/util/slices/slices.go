package slices

type StringSlice []string

// HasString iterates over the string slice to check if the string is present.
func (ss StringSlice) HasString(s string) bool {
	for _, item := range ss {
		if item == s {
			return true
		}
	}
	return false
}

func (ss StringSlice) IsEmpty() bool {
	return len(ss) == 0
}
