package local

import "context"

type localKey struct {
	name string
}

// Get retrieves a request-scoped value, set by [Set].
func Get(c context.Context, name string) any {
	return c.Value(localKey{name: name})
}

// Local sets a request-scoped value, which can be retrieved using [Get] later on.
func Set(c context.Context, name string, value any) {
	c = context.WithValue(c, localKey{name: name}, value)
}
