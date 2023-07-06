package Inspect

import (
	"fmt"
	"github.com/B9O2/Inspector/decorators"
	colors "github.com/gookit/color"
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

func (r Record) ToString(sep string) string {
	var color colors.Color
	var part string
	var parts []string
	for _, v := range r {
		if v == nil {
			continue
		}
		part = ""
		color = 0
		if v.formatter != nil {
			part = v.formatter(v.data)
			for _, tag := range v.tags {
				switch strings.SplitN(tag.Label(), ".", 2)[1] {
				case "color":
					if c, ok := tag.Data().(colors.Color); ok {
						color = c
					}
				default:
					continue
				}
			}
		} else {
			part = fmt.Sprintf("{%s error: no fomatter. value: %v}", v.typeLabel, v.data)
			color = colors.Red
		}
		parts = append(parts, color.Text(part))
	}
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
