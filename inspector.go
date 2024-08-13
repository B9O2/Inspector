package inspect

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/B9O2/Inspector/core"
	valuetypes "github.com/B9O2/Inspector/value_types"
	"github.com/B9O2/NStruct/ScrollArray"
)

type LogConifg struct {
	Level int
}

// Inspector 检查器是一切输出、记录、变量检查的核心。要实例化一个检查器您需要调用NewInspector()。
type Inspector struct {
	name                       string
	records                    *ScrollArray.ScrollArray
	autoValues                 map[string]func() *core.Value
	recordMiddleware           Middleware
	prefixOrders, suffixOrders []string
	visible                    bool
	showLabel                  bool
	record                     bool
	decolorate                 bool
	sep                        string
	writer                     io.Writer
}

func (insp *Inspector) SetAutoValue(label string, generator func() any, vt core.ValueType) error {
	if generator == nil {
		return errors.New(label + "(Auto): generator is nil")
	}
	insp.autoValues[label] = func() *core.Value {
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

func (insp *Inspector) initRecord(auto bool, values []*core.Value) core.Record {
	var allValues []*core.Value
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

	return core.Record(allValues)
}

// Print 打印值与自动生成的值
func (insp *Inspector) Print(values ...*core.Value) uint {
	return insp.printAndRecord(true, values...)
}

// JustPrint 此方法仅打印传入的值
func (insp *Inspector) JustPrint(values ...*core.Value) uint {
	return insp.printAndRecord(false, values...)
}

func (insp *Inspector) printAndRecord(auto bool, values ...*core.Value) uint {
	record := insp.initRecord(auto, values)
	if insp.recordMiddleware != nil {
		record = insp.recordMiddleware.Run(record)
	}
	if insp.visible && insp.writer != nil {
		insp.writer.Write([]byte(record.ToString(insp.sep, insp.showLabel, insp.decolorate)))
	}
	return insp.records.Append(record)
}

// FetchRecord 取回获得id对应的记录。如果id对应的结果不存在或已被清除则返回nil
func (insp *Inspector) FetchRecord(id uint) core.Record {
	if r, ok := insp.records.LoadWithEid(id); ok {
		return r.(core.Record)
	} else {
		return nil
	}
}

// SetSeparator 设置打印时各参数间的分隔符。
func (insp *Inspector) SetSeparator(sep string) {
	insp.sep = sep
}

func (insp *Inspector) SetRecordMiddleware(rm Middleware) {
	insp.recordMiddleware = rm
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

func (insp *Inspector) Range(f func(core.Record) bool) {
	insp.records.Range(func(r interface{}) bool {
		if r == nil {
			return true
		}
		record := r.(core.Record)
		return f(record)
	})
}

// NewInspector 实例化一个新的检查器，您需要为它命名并指定最大滚动储存的日志条数。
func NewInspector(name string, size uint) *Inspector {
	insp := &Inspector{
		name:       name,
		records:    ScrollArray.NewScrollArray(size),
		autoValues: map[string]func() *core.Value{},
		visible:    true,
		record:     true,
		decolorate: false,
		showLabel:  false,
		sep:        " ",
		writer:     os.Stdout,
	}

	insp.SetAutoValue("_start", func() any {
		return ">"
	}, valuetypes.Text)

	insp.SetAutoValue("_time", func() any {
		return time.Now()
	}, valuetypes.Time)

	insp.SetAutoValue("_end", func() any {
		return "\n"
	}, valuetypes.Text)

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
