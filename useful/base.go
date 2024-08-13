package useful

import (
	"fmt"
	"time"

	"github.com/B9O2/Inspector/types"
)

var (
	Text = types.NewType("text", func(a any) string {
		return fmt.Sprint(a)
	})

	Time = types.NewType("time", func(a any) string {
		return a.(time.Time).Format("2006/01/02 15:04:05")
	})
)
