package inspect

import (
	"errors"
	"fmt"
	colors "github.com/gookit/color"
	"reflect"
	"strings"
)

// Decorator 装饰器可以对值（Value）进行装饰
type Decorator struct {
	label string
	//todo 装饰器参数
	decoration func(v *Value) interface{}
}

func (d Decorator) Label() string {
	return d.label
}
func (d Decorator) parseLabel() (string, string) {
	suffix := ""
	parts := strings.Split(d.label, ".")
	if len(parts) > 1 {
		suffix = parts[len(parts)-1]
	}
	return parts[0], suffix
}
func (d Decorator) Decorate(v *Value) (ret Tag) {
	name, tagType := d.parseLabel()
	defer func() {
		if r := recover(); r != nil {
			ret = ErrorTag{
				name: name,
				err:  r.(error),
			}
		}
	}()
	data := d.decoration(v)
	switch data.(type) {
	case colors.Style:
		ret = ColorStyleTag{
			name:  name,
			color: data.(colors.Style),
		}
	case colors.Color:
		ret = ColorTag{
			name:  name,
			color: data.(colors.Color),
		}
	case string:
		ret = StringTag{
			name: name,
			str:  data.(string),
			mode: tagType,
		}
	case bool:
		ret = TestingTag{
			name: name,
			report: TestingReport{
				Success: data.(bool),
				Title:   "Testing <" + name + ">",
				Detail:  "No detail",
			},
		}
	case TestingReport:
		ret = TestingTag{
			name:   name,
			report: data.(TestingReport),
		}
	default:
		ret = ErrorTag{
			name: name,
			err:  errors.New("unknown decorator type '" + reflect.TypeOf(data).String() + "'"),
		}
	}
	return
}

// NewDecoration 初始化一个新的装饰器。标签需要使用后缀标记其类型，例如："MyRed.color"中包含".color"后缀，因此decoration应当返回colors.Color类型。
func NewDecoration(label string, decoration func(v *Value) interface{}) *Decorator {
	return &Decorator{
		label:      label,
		decoration: decoration,
	}
}

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
