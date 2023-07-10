// Package simple
// 这是一个简单的检查器样板间，您可以直接导入使用或基于此自定义符合您项目的检查器。
// 同
// /*
package simple

import (
	Inspect "github.com/B9O2/Inspector"
	"github.com/B9O2/Inspector/decorators"
	colors "github.com/gookit/color"
)

var Insp = Inspect.NewInspector("simple", 10000)
var (
	// Text 普通本文
	Text, _ = Insp.NewType("text", func(i interface{}) string {
		return i.(string)
	})

	// Path 路径，默认蓝色
	Path, _ = Insp.NewType("path", func(i interface{}) string {
		return "'" + i.(string) + "'"
	}, decorators.Blue)

	// Error 错误，默认红色
	Error, _ = Insp.NewType("error", func(i interface{}) string {
		return i.(error).Error()
	}, decorators.Red)

	// Level 输出级别，根据放入的数字不同显示不同文本。
	Level, _ = Insp.NewType("level", func(i interface{}) string {
		switch i.(int) {
		case 0:
			return "[INFO]"
		case 1:
			return "[WARNING]"
		case 2:
			return "[ERROR]"
		default:
			return "[UNKNOWN]"
		}
		//根据不同数字装饰不同颜色
	}, decorators.NewDecoration("level.color", func(i interface{}) interface{} {
		switch i.(int) {
		case 0:
			return colors.White
		case 1:
			return colors.Yellow
		case 2:
			return colors.Red
		default:
			return colors.Magenta
		}
	}))
)

var (
	//提前声明好三种输出级别
	LEVEL_INFO    = Level(0)
	LEVEL_WARNING = Level(1)
	LEVEL_ERROR   = Level(2)
)

func init() {
	Insp.SetOrders(Level, "_time", Path, Text)
}
