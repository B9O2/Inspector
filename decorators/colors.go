package decorators

import (
	"github.com/B9O2/Inspector/inspect"
	colors "github.com/gookit/color"
)

var (
	Cyan = inspect.NewDecoration("cyan.color", func(i interface{}) interface{} {
		return colors.Cyan
	})
	Red = inspect.NewDecoration("red.color", func(i interface{}) interface{} {
		return colors.Red
	})
	Blue = inspect.NewDecoration("blue.color", func(i interface{}) interface{} {
		return colors.Blue
	})
	Green = inspect.NewDecoration("green.color", func(i interface{}) interface{} {
		return colors.Green
	})
	Yellow = inspect.NewDecoration("yellow.color", func(i interface{}) interface{} {
		return colors.Yellow
	})
	Black = inspect.NewDecoration("black.color", func(i interface{}) interface{} {
		return colors.Black
	})
	Magenta = inspect.NewDecoration("magenta.color", func(i interface{}) interface{} {
		return colors.Magenta
	})
	Gray = inspect.NewDecoration("gray.color", func(i interface{}) interface{} {
		return colors.Gray
	})
	White = inspect.NewDecoration("white.color", func(i interface{}) interface{} {
		return colors.White
	})
)
