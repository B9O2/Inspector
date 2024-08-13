package decorators

import (
	"github.com/B9O2/Inspector/types"
	"github.com/gookit/color"
)

var (
	Cyan    = types.NewColorDecorator(color.New(color.FgCyan))
	Red     = types.NewColorDecorator(color.New(color.FgRed))
	Blue    = types.NewColorDecorator(color.New(color.FgBlue))
	Green   = types.NewColorDecorator(color.New(color.FgGreen))
	Yellow  = types.NewColorDecorator(color.New(color.FgYellow))
	Black   = types.NewColorDecorator(color.New(color.FgBlack))
	Magenta = types.NewColorDecorator(color.New(color.FgMagenta))
	Gray    = types.NewColorDecorator(color.New(color.FgGray))
	White   = types.NewColorDecorator(color.New(color.FgWhite))
)
