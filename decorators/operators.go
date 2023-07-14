package decorators

import (
	"github.com/B9O2/Inspector/inspect"
	colors "github.com/gookit/color"
)

func And(c1, c2 *inspect.Decorator) *inspect.Decorator {
	title := c1.Label() + " AND " + c2.Label()
	return inspect.NewDecoration(title, func(i *inspect.Value) interface{} {
		tag1 := c1.Decorate(i)
		tag2 := c2.Decorate(i)
		switch tag1.(type) {
		case inspect.TestingTag:
			switch tag2.(type) {
			case inspect.TestingTag:
				return tag1.(inspect.TestingTag).IsSuccess() && tag2.(inspect.TestingTag).IsSuccess()
			}
		case inspect.ColorTag:
			switch tag2.(type) {
			case inspect.ColorTag:
				return colors.New(
					tag1.(inspect.ColorTag).Color().ToFg(),
					tag2.(inspect.ColorTag).Color().ToBg(),
				)
			}
		}
		return false
	})
}

func Or(c1, c2 *inspect.Decorator) *inspect.Decorator {
	title := c1.Label() + " OR " + c2.Label()
	return inspect.NewDecoration(title, func(i *inspect.Value) interface{} {
		tag1 := c1.Decorate(i)
		tag2 := c2.Decorate(i)
		switch tag1.(type) {
		case inspect.TestingTag:
			switch tag2.(type) {
			case inspect.TestingTag:
				return tag1.(inspect.TestingTag).IsSuccess() || tag2.(inspect.TestingTag).IsSuccess()
			}
		}
		return false
	})
}

func Not(c *inspect.Decorator) *inspect.Decorator {
	title := "NOT " + c.Label()
	return inspect.NewDecoration(title, func(i *inspect.Value) interface{} {
		tag := c.Decorate(i)
		switch tag.(type) {
		case inspect.TestingTag:
			switch tag.(type) {
			case inspect.TestingTag:
				return !tag.(inspect.TestingTag).IsSuccess()
			}
			//todo 颜色运算
		}
		return false
	})
}
