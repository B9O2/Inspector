package core

import (
	"github.com/B9O2/Inspector/utils"
)

type ValueType func(v any, decos ...Decorator) *Value

// NewType 创建类型方法可以声明一种新的类型描述函数(VType)，类型描述函数用于在输出、记录时标记值的类型。 您总是应该在有新类型的值时调用此方法。
func NewType(label string, formatter func(any) string, extraDecos ...Decorator) ValueType {
	if formatter == nil {
		return func(v any, decos ...Decorator) *Value {
			return &Value{
				Label: label,
				Text:  "{#" + label + " init error: formatter is nil}",
				Color: utils.PanicStyle,
			}
		}
	}

	f := func(v any, decos ...Decorator) (value *Value) {
		defer func() {
			if r := recover(); r != nil {
				value.Text = "{#" + label + " panic: " + r.(error).Error() + "}"
				value.Color = utils.PanicStyle
			}
		}()

		value = &Value{
			Label: label,
			Color: utils.DefaultStyle,
		}

		value.Text = formatter(v)

		for _, d := range append(extraDecos, decos...) {
			v, err := d.Decorate(*value)
			if err != nil {
				v.Text = "{" + label + " decorator error: " + err.Error() + "}"
				v.Color = utils.PanicStyle
			}
			value = &v
		}

		return value
	}
	return f
}
