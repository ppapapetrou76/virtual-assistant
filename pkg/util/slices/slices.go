package slices

// StringSlice wraps a string slice to provide methods on top of it
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

// IsEmpty returns true if a string slice is empty, else false
func (ss StringSlice) IsEmpty() bool {
	return len(ss) == 0
}

// ContainsAny returns true if a string slice contains any of the given elements
func (ss StringSlice) ContainsAny(elements ...string) bool {
	for _, e := range elements {
		if ss.HasString(e) {
			return true
		}
	}
	return false
}

// Add adds the given element to the slice and returns a new string slice
func (ss StringSlice) Add(e string) StringSlice {
	return append(ss, e)
}
