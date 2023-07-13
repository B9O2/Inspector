package conditions

import (
	"github.com/B9O2/Inspector/inspect"
)

func newColorCondition(colorName string) inspect.Condition {
	return func(v *inspect.Value) bool {
		color := ""
		for _, tag := range v.Tags() {
			if tag.Type() == "color" {
				color = tag.Name()
			}
		}
		return color == colorName
	}
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
