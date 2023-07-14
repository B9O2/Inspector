package decorators

import (
	"github.com/B9O2/Inspector/inspect"
	colors "github.com/gookit/color"
)

var (
	Cyan = inspect.NewDecoration("cyan.color", func(i *inspect.Value) interface{} {
		return colors.Cyan
	})
	Red = inspect.NewDecoration("red.color", func(*inspect.Value) interface{} {
		return colors.Red
	})
	Blue = inspect.NewDecoration("blue.color", func(*inspect.Value) interface{} {
		return colors.Blue
	})
	Green = inspect.NewDecoration("green.color", func(*inspect.Value) interface{} {
		return colors.Green
	})
	Yellow = inspect.NewDecoration("yellow.color", func(*inspect.Value) interface{} {
		return colors.Yellow
	})
	Black = inspect.NewDecoration("black.color", func(*inspect.Value) interface{} {
		return colors.Black
	})
	Magenta = inspect.NewDecoration("magenta.color", func(*inspect.Value) interface{} {
		return colors.Magenta
	})
	Gray = inspect.NewDecoration("gray.color", func(*inspect.Value) interface{} {
		return colors.Gray
	})
	White = inspect.NewDecoration("white.color", func(*inspect.Value) interface{} {
		return colors.White
	})
)

func newColorCondition(colorName string) *inspect.Decorator {
	return NewTesting(colorName+".testing", func(v *inspect.Value) (bool, error) {
		color := ""
		for _, tag := range v.Tags() {
			if t, ok := tag.(inspect.ColorTag); ok {
				color = t.Name()
			}
		}
		return color == colorName, nil
	})
}

var (
	IsRed     = newColorCondition("red")
	IsYellow  = newColorCondition("yellow")
	IsBlue    = newColorCondition("blue")
	IsCyan    = newColorCondition("cyan")
	IsBlack   = newColorCondition("black")
	IsGreen   = newColorCondition("green")
	IsGray    = newColorCondition("gray")
	IsWhite   = newColorCondition("white")
	IsMagenta = newColorCondition("magenta")
)
