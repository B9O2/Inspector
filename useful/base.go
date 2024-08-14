package useful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/B9O2/Inspector/decorators"
	"github.com/B9O2/Inspector/types"
)

var (
	Text = types.NewType("text", func(a any) string {
		return fmt.Sprint(a)
	})

	Time = types.NewType("time", func(a any) string {
		return a.(time.Time).Format("2006/01/02 15:04:05")
	})

	// Json 序列化对象并生成美化后的Json字符串
	Json = types.NewType("json", func(i interface{}) string {
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

	// Error 错误，默认红色
	Error = types.NewType("error", func(i interface{}) string {
		return i.(error).Error()
	}, decorators.Red)
)
