package decorators

import (
	colors "github.com/gookit/color"
)

var (
	Cyan = NewDecoration("cyan.color", func(i interface{}) interface{} {
		return colors.Cyan
	})
	Red = NewDecoration("red.color", func(i interface{}) interface{} {
		return colors.Red
	})
	Blue = NewDecoration("blue.color", func(i interface{}) interface{} {
		return colors.Blue
	})
	Green = NewDecoration("green.color", func(i interface{}) interface{} {
		return colors.Green
	})
	Yellow = NewDecoration("yellow.color", func(i interface{}) interface{} {
		return colors.Yellow
	})
	Black = NewDecoration("black.color", func(i interface{}) interface{} {
		return colors.Black
	})
	Magenta = NewDecoration("magenta.color", func(i interface{}) interface{} {
		return colors.Magenta
	})
	Gray = NewDecoration("gray.color", func(i interface{}) interface{} {
		return colors.Gray
	})
	White = NewDecoration("white.color", func(i interface{}) interface{} {
		return colors.White
	})
)
