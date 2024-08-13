package decorators

import (
	"github.com/B9O2/Inspector/core"
	"github.com/gookit/color"
)

var (
	Cyan    = core.NewColorDecorator(color.New(color.FgCyan))
	Red     = core.NewColorDecorator(color.New(color.FgRed))
	Blue    = core.NewColorDecorator(color.New(color.FgBlue))
	Green   = core.NewColorDecorator(color.New(color.FgGreen))
	Yellow  = core.NewColorDecorator(color.New(color.FgYellow))
	Black   = core.NewColorDecorator(color.New(color.FgBlack))
	Magenta = core.NewColorDecorator(color.New(color.FgMagenta))
	Gray    = core.NewColorDecorator(color.New(color.FgGray))
	White   = core.NewColorDecorator(color.New(color.FgWhite))
)
