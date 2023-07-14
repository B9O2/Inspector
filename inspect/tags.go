package inspect

import (
	"errors"
	colors "github.com/gookit/color"
)

type Tag interface {
	Name() string
	Modify(v *ValueString) error
}

type ColorTag struct {
	name  string
	color colors.Color
}

func (c ColorTag) Name() string {
	return c.name
}

func (c ColorTag) Color() colors.Color {
	return c.color
}

func (c ColorTag) Modify(v *ValueString) error {
	v.color = colors.New(c.color.ToFg())
	return nil
}

type ColorStyleTag struct {
	name  string
	color colors.Style
}

func (c ColorStyleTag) Name() string {
	return c.name
}

func (c ColorStyleTag) ColorStyle() colors.Style {
	return c.color
}

func (c ColorStyleTag) Modify(v *ValueString) error {
	v.color = c.color
	return nil
}

type ErrorTag struct {
	name string
	err  error
}

func (c ErrorTag) Name() string {
	return c.name
}

func (c ErrorTag) Modify(v *ValueString) error {
	v.text += "{@" + c.name + " error:" + c.err.Error() + "}"
	v.color = colors.New(colors.Red)
	return nil
}

type StringTag struct {
	name string
	str  string
	mode string
}

func (c StringTag) Name() string {
	return c.name
}

func (c StringTag) Modify(v *ValueString) error {
	switch c.mode {
	case "text": //text
		v.text = c.str
	case "prefix": //prefix
		v.text = c.str + v.text
	case "suffix": //suffix
		v.text += c.str
	default:
		return errors.New("unknown string tag mode '" + c.mode + "'")
	}
	return nil
}

type TestingTag struct {
	name   string
	report TestingReport
}

func (c TestingTag) Name() string {
	return c.name
}

func (c TestingTag) IsSuccess() bool {
	return c.report.Success
}

func (c TestingTag) Modify(vs *ValueString) error {
	vs.testingReports = append(vs.testingReports, c.report)
	return nil
}
