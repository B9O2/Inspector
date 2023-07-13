package inspect

import (
	"errors"
	"fmt"
	colors "github.com/gookit/color"
	"reflect"
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

type Record []*Value

func (r Record) String() string {
	return r.ToString(" ")
}

func (r Record) CalCondition(conditions ...Condition) bool {
	if len(conditions) <= 0 {
		return true
	}
	for _, v := range r {
		for _, condition := range conditions {
			if condition(v) {
				return true
			}
		}
	}
	return false
}

func (r Record) decorate(text string, tags []Tag) (string, colors.Color, []TestingReport) {
	var color colors.Color
	var patchTags []Tag
	var testingResults []TestingReport

	appendErrorTag := func(label string, errText string) {
		patchTags = append(patchTags, NewTag(
			label,
			"error",
			errors.New(errText)))
	}

	for _, tag := range tags {
		data := tag.Data()
		if data == nil {
			continue
		}
		switch tag.tagType {
		case "color":
			switch data.(type) {
			case colors.Color:
				color = data.(colors.Color)
			case colors.RGBColor:
				color = data.(colors.RGBColor).Color()
			default:
				appendErrorTag(tag.Label(), "unknown color type '"+reflect.TypeOf(data).String()+"'")
			}
		case "testing":
			if res, ok := data.(TestingReport); ok {
				testingResults = append(testingResults, res)
			} else {
				appendErrorTag(tag.Label(), "unknown testing type '"+reflect.TypeOf(data).String()+"'")
			}
		case "text":
			if res, ok := data.(string); ok {
				text = res
			} else {
				appendErrorTag(tag.Label(), "unknown text type '"+reflect.TypeOf(data).String()+"'")
			}
		case "prefix":
			if res, ok := data.(string); ok {
				text = res + text
			} else {
				appendErrorTag(tag.Label(), "unknown prefix type '"+reflect.TypeOf(data).String()+"'")
			}
		case "suffix":
			if res, ok := data.(string); ok {
				text += res
			} else {
				appendErrorTag(tag.Label(), "unknown suffix type '"+reflect.TypeOf(data).String()+"'")
			}
		case "error":
			color = colors.Red
			if err, ok := data.(error); ok {
				text += fmt.Sprintf("{@%s error: %s}", tag.Label(), err)
			} else {
				text += fmt.Sprintf("{@%s error: %s}", tag.Label(), "unknown error type '"+reflect.TypeOf(data).String()+"'")
			}
		default:
			continue
		}
	}
	if len(patchTags) > 0 {
		return r.decorate(text, patchTags)
	} else {
		return text, color, testingResults
	}
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

func (r Record) ToString(sep string) string {
	var color colors.Color
	var reports []TestingReport
	var valueStr string
	var parts, valueReports []string
	for _, v := range r {
		if v == nil {
			continue
		}
		valueStr = ""
		color = 0
		if v.formatter != nil {
			valueStr = v.formatter(v.data)
			valueStr, color, reports = r.decorate(valueStr, v.tags)
			if len(reports) > 0 {
				valueReports = append(valueReports, r.genTestingReport(v.Label(), reports))
			}
		} else {
			valueStr = fmt.Sprintf("{%s error: no fomatter. value: %v}", v.typeLabel, v.data)
			color = colors.Red
		}

		if len(valueStr) > 0 || v.Label() == "_end" {
			var lines []string
			for _, line := range strings.Split(valueStr, "\n") {
				lines = append(lines, color.Text(line))
			}
			parts = append(parts, strings.Join(lines, "\n"))
		}
	}
	//strings.Join(valueReports, "\n")
	length := len(parts)
	if length > 0 {
		return strings.Join(parts[:length-1], sep) + parts[length-1]
	} else {
		return ""
	}
}

func (r Record) StringWithVType(sep string) string {
	var parts []string
	part := ""
	for _, v := range r {
		if v == nil {
			continue
		}
		if v.formatter != nil {
			part = fmt.Sprintf("%v", v.data)
		} else {
			part = v.formatter(v.data)
		}
		part += "<" + v.typeLabel + ">"
		parts = append(parts, part)
	}
	return strings.Join(parts, sep)
}
