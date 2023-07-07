package Inspect

import (
	"errors"
	"fmt"
	"github.com/B9O2/Inspector/decorators"
	"github.com/B9O2/NStruct/ScrollArray"
	"time"
)

type VType func(interface{}, ...*decorators.Decorator) *Value

type Inspector struct {
	name        string
	records     *ScrollArray.ScrollArray
	vTypes      map[string][]*decorators.Decorator
	autoTypes   map[string]VType
	autoTypeGen map[string]func() interface{}
	rTypeOrders []string
	sep         string
}

func (insp *Inspector) NewType(label string, formatter func(interface{}) string, decos ...*decorators.Decorator) (VType, error) {
	if len(label) > 0 && label[0] == '_' {
		return nil, errors.New(fmt.Sprintf("[INSP::%s]: Label cannot start with '_'.", insp.name))
	}
	return insp.newType(false, label, nil, formatter, decos...)
}

func (insp *Inspector) NewAutoType(label string, generator func() interface{}, formatter func(interface{}) string, decos ...*decorators.Decorator) error {
	if len(label) > 0 && label[0] == '_' {
		return errors.New(fmt.Sprintf("[INSP::%s]: Label cannot start with '_'.", insp.name))
	}
	_, err := insp.newType(true, label, generator, formatter, decos...)
	return err
}

func (insp *Inspector) newTypeFunc(label string, formatter func(interface{}) string) VType {
	return func(value interface{}, decos ...*decorators.Decorator) *Value {
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

func (insp *Inspector) newType(auto bool, label string, generator func() interface{}, formatter func(interface{}) string, decos ...*decorators.Decorator) (VType, error) {
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
	if vType == nil {
		return "", false
	}
	label := ""
	switch vType.(type) {
	case VType:
		label = vType.(VType)(nil).typeLabel
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
			retRecord = append(retRecord, values[label][0])
		} else {
			retRecord = append(retRecord, values[label]...)
		}
	}
	return retRecord
}

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

func (insp *Inspector) SetAutoTypeFormatter(label string, formatter func(interface{}) string) error {
	if _, ok := insp.autoTypes[label]; ok {
		insp.autoTypes[label] = insp.newTypeFunc(label, formatter)
		return nil
	} else {
		return errors.New(fmt.Sprintf("[INSP::%s]: AutoType '%s' not exists.", insp.name, label))
	}
}

func (insp *Inspector) GetAutoType(label string) (VType, bool) {
	vType, ok := insp.autoTypes[label]
	return vType, ok
}

func (insp *Inspector) initRecord(values []*Value) Record {
	for _, value := range values {
		if decos, ok := insp.vTypes[value.typeLabel]; ok {
			//类型装饰器与额外装饰器
			for _, deco := range append(decos, value.extraDecorators...) {
				value.tags = append(value.tags, deco.Decorate(value.data))
			}
		} else {
			value.formatter = func(i interface{}) string {
				return "{" + value.typeLabel + " error: type not be registered in this inspector '" + insp.name + "'}"
			}
		}
	}
	for label, generator := range insp.autoTypeGen {
		values = append(values, insp.autoTypes[label](generator()))
	}
	return insp.order(values)
}

func (insp *Inspector) Record(values ...*Value) uint {
	return insp.records.Append(insp.initRecord(values))
}

func (insp *Inspector) Print(values ...*Value) {
	record := insp.initRecord(values)
	fmt.Print(record.ToString(insp.sep))
}

func (insp *Inspector) PrintAndRecord(values ...*Value) uint {
	record := insp.initRecord(values)
	fmt.Print(record.ToString(insp.sep))
	return insp.records.Append(record)
}

func (insp *Inspector) FetchRecord(id uint) Record {
	if r, ok := insp.records.LoadWithEid(id); ok {
		return r.(Record)
	} else {
		return nil
	}
}

func (insp *Inspector) SetSeparator(sep string) {
	insp.sep = sep
}

func (insp *Inspector) SetTypeDecorations(vType interface{}, decos ...*decorators.Decorator) error {
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

func NewInspector(name string, size uint) *Inspector {
	insp := &Inspector{
		name:        name,
		records:     ScrollArray.NewScrollArray(size),
		autoTypes:   map[string]VType{},
		autoTypeGen: map[string]func() interface{}{},
		vTypes:      map[string][]*decorators.Decorator{},
		sep:         " ",
	}

	_, _ = insp.newType(true, "_time", func() interface{} {
		return time.Now()
	}, func(v interface{}) string {
		return v.(time.Time).Format("2006/01/02 15:04:05")
	})

	/*todo 当前文件
	_, _ = insp.newType(true, "_file", func() interface{} {
		return time.Now()
	}, func(v interface{}) string {
		return v.(time.Time).Format("2006/01/02 15:04:05")
	})
	*/

	_, _ = insp.newType(true, "_start", func() interface{} {
		return " "
	}, func(v interface{}) string {
		return v.(string)
	})

	_, _ = insp.newType(true, "_end", func() interface{} {
		return "\n"
	}, func(v interface{}) string {
		return v.(string)
	})

	insp.SetOrders("_time") //初始化排序

	return insp
}
