package inspect

import (
	"errors"
	"fmt"
	"github.com/B9O2/NStruct/ScrollArray"
	"runtime"
	"time"
)

type VType func(interface{}, ...*Decorator) *Value

// Inspector 检查器是一切输出、记录、变量检查的核心。要实例化一个检查器您需要调用NewInspector()。
type Inspector struct {
	name              string
	records           *ScrollArray.ScrollArray
	vTypes            map[string][]*Decorator
	autoTypes         map[string]VType
	autoTypeGen       map[string]func() interface{}
	frontDecorators   []*Decorator
	postDecorators    []*Decorator
	recordMiddleware  Middleware
	rTypeOrders       []string
	visibleConditions []*Decorator
	visible           bool
	sep               string
}

// NewType 创建类型方法可以声明一种新的类型描述函数(VType)，类型描述函数用于在输出、记录时标记值的类型。 您总是应该在有新类型的值时调用此方法。
func (insp *Inspector) NewType(label string, formatter func(interface{}) string, decos ...*Decorator) (VType, error) {
	if len(label) > 0 && label[0] == '_' {
		return nil, errors.New(fmt.Sprintf("[INSP::%s]: Label cannot start with '_'.", insp.name))
	}
	return insp.newType(false, label, nil, formatter, decos...)
}

// NewAutoType 创建自动类型方法是特殊的类型描述函数(VType)声明方法。此类型声明时需要实现generator()，此函数会在输出、记录时自动调用并将返回值提供给formatter()
func (insp *Inspector) NewAutoType(label string, generator func() interface{}, formatter func(interface{}) string, decos ...*Decorator) error {
	if len(label) > 0 && label[0] == '_' {
		return errors.New(fmt.Sprintf("[INSP::%s]: Label cannot start with '_'.", insp.name))
	}
	_, err := insp.newType(true, label, generator, formatter, decos...)
	return err
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

func (insp *Inspector) newTypeFunc(label string, formatter func(interface{}) string) VType {
	return func(value interface{}, decos ...*Decorator) *Value {
		return &Value{
			typeLabel: label,
			formatter: func(i interface{}) (str string) {
				defer func() {
					if r := recover(); r != nil {
						str = "{" + label + " error: " + r.(error).Error() + "}"
					}
				}()
				str = formatter(i)
				return
			},
			extraDecorators: decos,
			data:            value,
		}
	}
}

func (insp *Inspector) newType(auto bool, label string, generator func() interface{}, formatter func(interface{}) string, decos ...*Decorator) (VType, error) {
	typeFunc := insp.newTypeFunc(label, formatter)

	if _, ok := insp.vTypes[label]; ok {
		return nil, errors.New(fmt.Sprintf("[INSP::%s]: Type '%s' already exists.", insp.name, label))
	} else {
		insp.vTypes[label] = decos
		insp.rTypeOrders = append(insp.rTypeOrders, label)
	}

	if auto {
		insp.autoTypes[label] = typeFunc
		insp.autoTypeGen[label] = generator
		return nil, nil
	} else {
		return typeFunc, nil
	}
}

func (insp *Inspector) getLabel(vType interface{}) (string, bool) {
	label := ""
	switch vType.(type) {
	case VType:
		if vType.(VType) != nil {
			label = vType.(VType)(nil).typeLabel
		} else {
			return "", false
		}
	case string:
		label = vType.(string)
	default:
		return "", false
	}
	return label, true
}

func (insp *Inspector) order(record Record) Record {
	orders := []string{"_start"}
	orders = append(orders, insp.rTypeOrders...)
	orders = append(orders, "_end")
	values := map[string][]*Value{}
	for _, v := range record {
		values[v.typeLabel] = append(values[v.typeLabel], v)
	}

	retRecord := Record{}
	for _, label := range orders {
		if _, ok := insp.autoTypes[label]; ok {
			if len(values[label]) > 0 {
				retRecord = append(retRecord, values[label][0])
			}
		} else {
			retRecord = append(retRecord, values[label]...)
		}
	}
	return retRecord
}

// SetOrders 设定值顺序。仅此方法可以设定顺序，值顺序不会受调用Print()、Record()等方法时传参顺序影响
func (insp *Inspector) SetOrders(vTypes ...interface{}) {
	var orders []string
	m := map[string]byte{}
	for _, vType := range vTypes {
		if label, ok := insp.getLabel(vType); ok {
			if _, ok := m[label]; !ok {
				switch label {
				case "_start", "_end":
					continue
				default:
					m[label] = '_'
					orders = append(orders, label)
				}
			}
		}
	}
	for label := range insp.vTypes {
		if _, ok := m[label]; !ok {
			switch label {
			case "_start", "_end":
				continue
			default:
				m[label] = '_'
				orders = append(orders, label)
			}
		}
	}
	insp.rTypeOrders = orders
}

// SetAutoTypeFormatter 设定自动类型的格式化方法。此方法多用于设定内置自动类型。
func (insp *Inspector) SetAutoTypeFormatter(label string, formatter func(interface{}) string) error {
	if _, ok := insp.autoTypes[label]; ok {
		insp.autoTypes[label] = insp.newTypeFunc(label, formatter)
		return nil
	} else {
		return errors.New(fmt.Sprintf("[INSP::%s]: AutoType '%s' not exists.", insp.name, label))
	}
}

// GetAutoType 用于获取自动类型描述函数的方法。如希望主动提供自动类型的值需要调用此方法获得类型描述函数(VType)
func (insp *Inspector) GetAutoType(label string) (VType, bool) {
	vType, ok := insp.autoTypes[label]
	return vType, ok
}

func (insp *Inspector) initRecord(values []*Value, auto bool) Record {
	if auto {
		for label, generator := range insp.autoTypeGen {
			values = append(values, insp.autoTypes[label](generator()))
		}
	}
	for _, value := range values {
		if decos, ok := insp.vTypes[value.typeLabel]; ok {
			//在前置装饰器后加入类型装饰器
			decos = append(insp.frontDecorators, decos...)
			//在类型装饰器后加入额外装饰器
			decos = append(decos, value.extraDecorators...)
			//在额外装饰器后加入后置装饰器
			for _, deco := range append(decos, insp.postDecorators...) {
				value.tags = append(value.tags, deco.Decorate(value))
			}
		} else {
			value.formatter = func(i interface{}) string {
				return "{" + value.typeLabel + " error: type not be registered in this inspector '" + insp.name + "'}"
			}
		}
	}
	return insp.order(values)
}

// Record 记录参数并返回此记录的id。id永不重复但会在超出长度时清理掉过旧的记录
func (insp *Inspector) Record(values ...*Value) uint {
	return uint(insp.printAndRecord(false, true, true, values...))
}

// Print 打印值。此方法仅打印参数，不储存参数
func (insp *Inspector) Print(values ...*Value) {
	insp.printAndRecord(true, false, true, values...)
}

// JustPrint 此方法仅打印传入参数
func (insp *Inspector) JustPrint(values ...*Value) {
	insp.printAndRecord(true, false, false, values...)
}

// PrintAndRecord 此方法既打印参数又储存参数
func (insp *Inspector) PrintAndRecord(values ...*Value) uint {
	return uint(insp.printAndRecord(true, true, true, values...))
}

func (insp *Inspector) printAndRecord(isPrint, isRecord, auto bool, values ...*Value) int {
	var record Record
	if isRecord || insp.visible {
		record = insp.initRecord(values, auto)
		if isPrint && record.CalCondition(insp.visibleConditions...) {
			fmt.Print(record.ToString(insp.sep, false))
		}
		if isRecord {
			if insp.recordMiddleware != nil {
				record = insp.recordMiddleware.Run(record)
			}
			return int(insp.records.Append(record))
		}
	}
	return -1
}

// FetchRecord 取回获得id对应的记录。如果id对应的结果不存在或已被清除则返回nil
func (insp *Inspector) FetchRecord(id uint) Record {
	if r, ok := insp.records.LoadWithEid(id); ok {
		return r.(Record)
	} else {
		return nil
	}
}

// SetSeparator 设置打印时各参数间的分隔符。特殊的"_start"与"_end"与其他值之间没有分隔符。
func (insp *Inspector) SetSeparator(sep string) {
	insp.sep = sep
}

// SetTypeDecorations 设置类型装饰器。设定后所有被此类型描述的值都将具有这些装饰器。
func (insp *Inspector) SetTypeDecorations(vType interface{}, decos ...*Decorator) error {
	if label, ok := insp.getLabel(vType); ok {
		if _, ok := insp.vTypes[label]; ok {
			insp.vTypes[label] = decos
			return nil
		} else {
			return errors.New(fmt.Sprintf("[INSP::%s]: Type '%s' not exists.", insp.name, label))
		}
	}
	return errors.New(fmt.Sprintf("[INSP::%s]: '%v' is not a VType.", insp.name, vType))
}

// SetFrontDecorations 设置前置装饰器
func (insp *Inspector) SetFrontDecorations(decos ...*Decorator) {
	insp.frontDecorators = decos
}

// SetPostDecorations 设置后置装饰器
func (insp *Inspector) SetPostDecorations(decos ...*Decorator) {
	insp.postDecorators = decos
}

func (insp *Inspector) SetRecordMiddleware(rm Middleware) {
	insp.recordMiddleware = rm
}

// SetVisible 设置输出可见性，如果false则Print()方法不做任何事
func (insp *Inspector) SetVisible(visible bool) {
	insp.visible = visible
}

// SetVisibleConditions 设置可见条件，仅满足条件的记录(Record)会被打印
func (insp *Inspector) SetVisibleConditions(conditions ...*Decorator) {
	insp.visibleConditions = conditions
}

func (insp *Inspector) Range(f func(Record) bool, conditions ...*Decorator) {
	insp.records.Range(func(r interface{}) bool {
		if r == nil {
			return true
		}
		record := r.(Record)
		if record.CalCondition(conditions...) {
			return f(record)
		} else {
			return true
		}
	})
}

// NewInspector 实例化一个新的检查器，您需要为它命名并指定最大滚动储存的数目。
func NewInspector(name string, size uint) *Inspector {
	insp := &Inspector{
		name:        name,
		records:     ScrollArray.NewScrollArray(size),
		autoTypes:   map[string]VType{},
		autoTypeGen: map[string]func() interface{}{},
		vTypes:      map[string][]*Decorator{},
		visible:     true,
		sep:         " ",
	}

	_, _ = insp.newType(true, "_start", func() interface{} {
		return ""
	}, func(v interface{}) string {
		return v.(string)
	})

	_, _ = insp.newType(true, "_end", func() interface{} {
		return "\n"
	}, func(v interface{}) string {
		return v.(string)
	})

	_, _ = insp.newType(true, "_time", func() interface{} {
		return time.Now()
	}, func(v interface{}) string {
		return v.(time.Time).Format("2006/01/02 15:04:05")
	})

	_, _ = insp.newType(true, "_func", func() interface{} {
		pc := make([]uintptr, 1)
		runtime.Callers(4, pc)
		f := runtime.FuncForPC(pc[0])
		return f.Name()
	}, func(v interface{}) string {
		return v.(string)
	})

	/*todo 当前文件
	_, _ = insp.newType(true, "_file", func() interface{} {
		return time.Now()
	}, func(v interface{}) string {
		return v.(time.Time).Format("2006/01/02 15:04:05")
	})
	*/

	insp.SetOrders("_time", "_func") //初始化排序

	return insp
}
