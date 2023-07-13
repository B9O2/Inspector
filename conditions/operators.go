package conditions

import "github.com/B9O2/Inspector/inspect"

func And(c1, c2 inspect.Condition) inspect.Condition {
	return func(v *inspect.Value) bool {
		return c1(v) && c2(v)
	}
}

func Or(c1, c2 inspect.Condition) inspect.Condition {
	return func(v *inspect.Value) bool {
		return c1(v) || c2(v)
	}
}

func Not(c inspect.Condition) inspect.Condition {
	return func(v *inspect.Value) bool {
		return !c(v)
	}
}
