package value_types

import (
	"fmt"
	"time"

	"github.com/B9O2/Inspector/core"
)

var (
	Text = core.NewType("text", func(a any) string {
		return fmt.Sprint(a)
	})

	Time = core.NewType("time", func(a any) string {
		return a.(time.Time).Format("2006/01/02 15:04:05")
	})
)
