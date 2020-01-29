package ptr

// FromInt returns a pointer to the argument
func FromInt(i int) *int { return &i }

// FromString returns a pointer to the argument
func FromString(s string) *string { return &s }

// Int returns the dereferenced value of the pointer, if the pointer is nil, 'or' is returned
func IntOr(p *int, or int) int {
	if p == nil {
		return or
	}
	return *p
}

// Int returns the dereferenced value of the pointer, if the pointer is nil, 0 is returned
func Int(p *int) int { return IntOr(p, 0) }

// String returns the dereferenced value of the pointer, if the pointer is nil, "" is returned
func String(p *string) string { return StringOr(p, "") }

// String returns the dereferenced value of the pointer, if the pointer is nil, 'or' is returned
func StringOr(p *string, or string) string {
	if p == nil {
		return or
	}
	return *p
}
