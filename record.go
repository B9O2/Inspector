package Inspect

import (
	"errors"
	"fmt"
	"github.com/B9O2/Inspector/decorators"
	colors "github.com/gookit/color"
	"reflect"
	"strings"
)

type Value struct {
	typeLabel       string
	formatter       func(interface{}) string
	extraDecorators []*decorators.Decorator
	tags            []decorators.Tag
	data            interface{}
}

func (v Value) Label() string {
	return v.typeLabel
}

func (v Value) Data() interface{} {
	return v.data
}

type Record []*Value

func (r Record) String() string {
	return r.ToString(" ")
}

func (r Record) decorate(text string, tags []decorators.Tag) (string, colors.Color, []decorators.TestingReport) {
	var color colors.Color
	var patchTags []decorators.Tag
	var testingResults []decorators.TestingReport

	appendErrorTag := func(label string, errText string) {
		patchTags = append(patchTags, decorators.NewTag(
			label+".error",
			errors.New(errText)))
	}

	for _, tag := range tags {
		suffix := ""
		tagParts := strings.Split(tag.Label(), ".")
		if len(tagParts) > 1 {
			suffix = tagParts[len(tagParts)-1]
		}
		data := tag.Data()
		if data == nil {
			continue
		}
		switch suffix {
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
			if res, ok := data.(decorators.TestingReport); ok {
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

func (r Record) genTestingReport(label string, reports []decorators.TestingReport) string {
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
	var reports []decorators.TestingReport
	var part string
	var parts, valueReports []string
	for _, v := range r {
		if v == nil {
			continue
		}
		part = ""
		color = 0
		if v.formatter != nil {
			part = v.formatter(v.data)
			part, color, reports = r.decorate(part, v.tags)
			if len(reports) > 0 {
				valueReports = append(valueReports, r.genTestingReport(v.Label(), reports))
			}
		} else {
			part = fmt.Sprintf("{%s error: no fomatter. value: %v}", v.typeLabel, v.data)
			color = colors.Red
		}

		if len(part) > 0 {
			var lines []string
			for _, line := range strings.Split(part, "\n") {
				lines = append(lines, color.Text(line))
			}
			parts = append(parts, strings.Join(lines, "\n"))
		}
	}
	//strings.Join(valueReports, "\n")
	return strings.Join(parts, sep)
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
