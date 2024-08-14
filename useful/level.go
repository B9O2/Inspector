package useful

import (
	"github.com/B9O2/Inspector/decorators"
	"github.com/B9O2/Inspector/types"
)

var (
	Level = types.NewType("level", func(a any) string {
		name := ""
		switch a.(int) {
		case 0:
			name = "ERROR"
		case 1:
			name = "WARN"
		case 2:
			name = "INFO"
		case 3:
			name = "DEBUG"
		default:
			name = "UNKNOW"
		}
		return "[" + name + "]"
	})
)

var (
	ERROR = Level(0, decorators.Red)
	WARN  = Level(1, decorators.Yellow)
	INFO  = Level(2)
	DEBUG = Level(3, decorators.Cyan)
)

func LevelFilter(limit *types.Value) types.Middleware {
	return types.NewBaseMiddleware(func(r types.Record) types.Record {
		show := true
		if l, ok := limit.Raw.(int); ok {
			for _, v := range r {
				if v.Label == "level" {
					if vl, ok := v.Raw.(int); ok && vl > l {
						show = false
						break
					}
				}
			}
		}
		if show {
			return r
		} else {
			return nil
		}
	})
}
