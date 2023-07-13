package decorators

import "github.com/B9O2/Inspector/inspect"

func Testing(title string, testFunc func(i interface{}) (bool, error)) *inspect.Decorator {
	return inspect.NewDecoration("["+title+"].testing", func(i interface{}) interface{} {
		success, err := testFunc(i)
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
