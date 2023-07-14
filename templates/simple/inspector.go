// Package simple
// 这是一个简单的检查器样板，您可以直接导入使用或基于此自定义符合您项目的检查器。
// /*
package simple

import (
	"bytes"
	"encoding/json"
	"github.com/B9O2/Inspector/decorators"
	"github.com/B9O2/Inspector/inspect"
	colors "github.com/gookit/color"
)

var Insp = inspect.NewInspector("simple", 10000)
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
	}, inspect.NewDecoration("level.color", func(i *inspect.Value) interface{} {
		switch i.Data().(int) {
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

	// Json 序列化对象并生成美化后的Json字符串
	Json, _ = Insp.NewType("json", func(i interface{}) string {
		var out bytes.Buffer
		if obj, ok := i.(string); ok {
			err := json.Indent(&out, []byte(obj), "", "  ")
			if err != nil {
				return err.Error()
			}
			return "\n" + out.String() + "\n"
		} else {
			marshal, err := json.MarshalIndent(i, "", "  ")
			if err != nil {
				return err.Error()
			}
			return "\n" + string(marshal) + "\n"
		}

	}, decorators.Cyan)
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
