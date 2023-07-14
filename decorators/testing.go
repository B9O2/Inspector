package decorators

import "github.com/B9O2/Inspector/inspect"

func NewTesting(title string, testFunc func(*inspect.Value) (bool, error)) *inspect.Decorator {
	return inspect.NewDecoration("["+title+"].testing", func(v *inspect.Value) interface{} {
		success, err := testFunc(v)
		detail := ""
		if err != nil {
			detail = err.Error()
		}
		return inspect.TestingReport{
			Success: success,
			Title:   title,
			Detail:  detail,
		}
	})
}
