package decorators

import "github.com/B9O2/Inspector/inspect"

var (
	Invisible = inspect.NewDecoration("invisible.text", func(*inspect.Value) interface{} {
		return ""
	})
)
