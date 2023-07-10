package decorators

var (
	Invisible = NewDecoration("invisible.text", func(i interface{}) interface{} {
		return ""
	})
)
