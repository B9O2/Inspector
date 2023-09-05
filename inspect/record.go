package inspect

import (
	"fmt"
	colors "github.com/gookit/color"
	"strings"
)

type Value struct {
	typeLabel       string
	formatter       func(interface{}) string
	extraDecorators []*Decorator
	tags            []Tag
	data            interface{}
}

func (v Value) Label() string {
	return v.typeLabel
}

func (v Value) Tags() []Tag {
	return v.tags
}

func (v Value) Data() interface{} {
	return v.data
}

type ValueString struct {
	text           string
	color          colors.Style
	testingReports []TestingReport
}

func (vs ValueString) String() string {
	return vs.ToString(false)
}

func (vs ValueString) ToString(disableColor bool) string {
	if disableColor {
		return vs.text
	}
	var lines []string
	for _, line := range strings.Split(vs.text, "\n") {
		lines = append(lines, vs.color.Sprint(line))
	}
	return strings.Join(lines, "\n")
}

type Record []*Value

func (r Record) String() string {
	return r.ToString(" ", false, false)
}

func (r Record) CalCondition(conditions ...*Decorator) bool {
	if len(conditions) <= 0 {
		return true
	}
	for _, v := range r {
		for _, condition := range conditions {
			if c, ok := condition.Decorate(v).(TestingTag); ok {
				if c.IsSuccess() {
					return true
				}
			}
		}
	}
	return false
}

func (r Record) genTestingReport(label string, reports []TestingReport) string {
	finalReportFmt := "\\___[%s] Testing Report___\n    %s\n_____________________(%d/%d)\n\n"
	totalSuccess := 0
	var parts []string
	for _, report := range reports {
		if report.Success {
			totalSuccess += 1
		}
		parts = append(parts, report.String())
	}
	return fmt.Sprintf(
		finalReportFmt,
		label,
		strings.Join(parts, "\n    "),
		totalSuccess,
		len(reports),
	)
}

func (r Record) ToString(sep string, showLabel, disableColor bool) string {
	var parts []string
	var endl string
	for _, v := range r {
		if v == nil || v.formatter == nil {
			continue
		}
		vs := ValueString{
			text: v.formatter(v.data),
		}

		for _, tag := range v.tags {
			if err := tag.Modify(&vs); err != nil {
				//todo 处理错误
			}
		}
		part := vs.ToString(disableColor)
		if showLabel {
			part += v.Label()
		}
		if len(part) > 0 && v.Label() != "_end" {
			parts = append(parts, part)
		} else {
			endl = part
		}
	}
	return strings.Join(parts, sep) + endl
}
