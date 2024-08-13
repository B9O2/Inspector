package inspect

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/B9O2/Inspector/types"
	"github.com/B9O2/Inspector/useful"
	"github.com/B9O2/NStruct/ScrollArray"
)

type LogConifg struct {
	Level int
}

// Inspector 检查器是一切输出、记录、变量检查的核心。要实例化一个检查器您需要调用NewInspector()。
type Inspector struct {
	name                       string
	records                    *ScrollArray.ScrollArray
	autoValues                 map[string]func() *types.Value
	middlewares                []types.Middleware
	prefixOrders, suffixOrders []string
	visible                    bool
	showLabel                  bool
	record                     bool
	decolorate                 bool
	sep                        string
	writer                     io.Writer
	enable                     bool
}

func (insp *Inspector) SetAutoValue(label string, generator func() any, vt types.ValueType) error {
	if generator == nil {
		return errors.New(label + "(Auto): generator is nil")
	}
	insp.autoValues[label] = func() *types.Value {
		return vt(generator())
	}
	return nil
}

/*
// Inherit 继承其他检查器的类型
func (insp *Inspector) Inherit(inspectors ...*Inspector) error {
	for _,insp1:=range inspectors{
		for label,decos:=range insp1.vTypes{
			if gen,ok:=insp1.autoTypeGen[label];ok{
				insp.newType(true,label,gen,insp1.,)
			}

		}
	}
	return err
}
*/

// SetOrders 设定值顺序。仅此方法可以设定顺序，值顺序不会受调用Print()、Record()等方法时传参顺序影响
func (insp *Inspector) SetPrefixOrders(orders ...string) {
	insp.prefixOrders = orders
}

func (insp *Inspector) SetSuffixOrders(orders ...string) {
	insp.suffixOrders = orders
}

func (insp *Inspector) SetEnable(enable bool) {
	insp.enable = enable
}

func (insp *Inspector) initRecord(auto bool, values []*types.Value) types.Record {
	var allValues []*types.Value
	if auto {
		for _, label := range insp.prefixOrders {
			allValues = append(allValues, insp.autoValues[label]())
		}
	}

	allValues = append(allValues, values...)

	if auto {
		for _, label := range insp.suffixOrders {
			allValues = append(allValues, insp.autoValues[label]())
		}
	}

	return types.Record(allValues)
}

// Print 打印值与自动生成的值
func (insp *Inspector) Print(values ...*types.Value) int {
	return insp.printAndRecord(true, values...)
}

// JustPrint 此方法仅打印传入的值
func (insp *Inspector) JustPrint(values ...*types.Value) int {
	return insp.printAndRecord(false, values...)
}

func (insp *Inspector) printAndRecord(auto bool, values ...*types.Value) int {
	if !insp.enable {
		return -1
	}
	record := insp.initRecord(auto, values)
	for _, middleware := range insp.middlewares {
		record = middleware.Run(record)
	}
	if len(record) > 0 {
		if insp.visible && insp.writer != nil {
			insp.writer.Write([]byte(record.ToString(insp.sep, insp.showLabel, insp.decolorate)))
		}
		return int(insp.records.Append(record))
	} else {
		return -1
	}
}

// FetchRecord 取回获得id对应的记录。如果id对应的结果不存在或已被清除则返回nil
func (insp *Inspector) FetchRecord(id uint) types.Record {
	if r, ok := insp.records.LoadWithEid(id); ok {
		return r.(types.Record)
	} else {
		return nil
	}
}

// SetSeparator 设置打印时各参数间的分隔符。
func (insp *Inspector) SetSeparator(sep string) {
	insp.sep = sep
}

func (insp *Inspector) SetRecordMiddleware(rms ...types.Middleware) {
	insp.middlewares = rms
}

// SetVisible 设置输出可见性，如果false则Print()方法不做任何事
func (insp *Inspector) SetVisible(visible bool) {
	insp.visible = visible
}

// SetShowLabel 设置Label可见性
func (insp *Inspector) ShowLabel(show bool) {
	insp.showLabel = show
}

func (insp *Inspector) SetWriter(w io.Writer) {
	insp.writer = w
}

func (insp *Inspector) Range(f func(types.Record) bool) {
	insp.records.Range(func(r interface{}) bool {
		if r == nil {
			return true
		}
		record := r.(types.Record)
		return f(record)
	})
}

// NewInspector 实例化一个新的检查器，您需要为它命名并指定最大滚动储存的日志条数。
func NewInspector(name string, size uint) *Inspector {
	insp := &Inspector{
		name:       name,
		records:    ScrollArray.NewScrollArray(size),
		autoValues: map[string]func() *types.Value{},
		visible:    true,
		record:     true,
		decolorate: false,
		showLabel:  false,
		sep:        " ",
		writer:     os.Stdout,
		enable:     true,
	}

	insp.SetAutoValue("_start", func() any {
		return ">"
	}, useful.Text)

	insp.SetAutoValue("_time", func() any {
		return time.Now()
	}, useful.Time)

	insp.SetAutoValue("_end", func() any {
		return "\n"
	}, useful.Text)

	insp.SetPrefixOrders("_start", "_time")
	insp.SetSuffixOrders("_end")

	// _, _ = insp.newType(true, "_func", func() interface{} {
	// 	pc := make([]uintptr, 1)
	// 	runtime.Callers(4, pc)
	// 	f := runtime.FuncForPC(pc[0])
	// 	return f.Name()
	// }, func(v interface{}) string {
	// 	return v.(string)
	// })

	/*todo 当前文件
	_, _ = insp.newType(true, "_file", func() interface{} {
		return time.Now()
	}, func(v interface{}) string {
		return v.(time.Time).Format("2006/01/02 15:04:05")
	})
	*/

	return insp
}
