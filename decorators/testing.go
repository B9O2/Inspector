package decorators

import "fmt"

type TestingReport struct {
	Success bool
	Title   string
	Detail  string
}

func (tr TestingReport) String() string {
	checkSign := "○"
	tip := ""
	if tr.Success {
		checkSign = "◉"
		tip = "SUCCESS"
	} else {
		checkSign = "○"
		tip = "FAILED"
	}
	return fmt.Sprintf("%s [%s] %s %s",
		checkSign,
		tr.Title,
		tr.Detail,
		tip,
	)
}

func Testing(title string, testFunc func(i interface{}) (bool, error)) *Decorator {
	return NewDecoration("["+title+"].testing", func(i interface{}) interface{} {
		success, err := testFunc(i)
		detail := ""
		if err != nil {
			detail = err.Error()
		}
		return TestingReport{
			Success: success,
			Title:   title,
			Detail:  detail,
		}
	})
}
