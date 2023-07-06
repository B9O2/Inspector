package Inspect

import (
	"errors"
	"fmt"
	"github.com/B9O2/NStruct/ScrollArray"
	"time"
)

type VType func(value interface{}) *Value

type Inspector struct {
	name        string
	records     *ScrollArray.ScrollArray
	vTypes      map[string]*TestingGroup
	autoTypes   map[string]VType
	autoTypeGen map[string]func() interface{}
	rTypeOrders []string
	sep         string
}

func (insp *Inspector) NewType(label string, formatter func(interface{}) string) (VType, error) {
	if len(label) > 0 && label[0] == '_' {
		return nil, errors.New(fmt.Sprintf("[INSP::%s]: Label cannot start with '_'.", insp.name))
	}
	return insp.newType(false, label, nil, formatter)
}

func (insp *Inspector) NewAutoType(label string, generator func() interface{}, formatter func(interface{}) string) error {
	if len(label) > 0 && label[0] == '_' {
		return errors.New(fmt.Sprintf("[INSP::%s]: Label cannot start with '_'.", insp.name))
	}
	_, err := insp.newType(true, label, generator, formatter)
	return err
}

func (insp *Inspector) newType(auto bool, label string, generator func() interface{}, formatter func(interface{}) string) (VType, error) {

	typeFunc := func(value interface{}) *Value {
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
			data: value,
		}
	}

	if _, ok := insp.vTypes[label]; ok {
		return nil, errors.New(fmt.Sprintf("[INSP::%s]: Type '%s' already exists.", insp.name, label))
	} else {
		insp.vTypes[label] = nil
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
		retRecord = append(retRecord, values[label]...)
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
	insp.rTypeOrders = orders
}

func (insp *Inspector) SetAutoTypeFormatter(label string, formatter func(interface{}) string) error {
	if _, ok := insp.autoTypes[label]; ok {
		insp.autoTypes[label] = func(value interface{}) *Value {
			return &Value{
				typeLabel: label,
				formatter: formatter,
				data:      value,
			}
		}
		return nil
	} else {
		return errors.New(fmt.Sprintf("[INSP::%s]: AutoType '%s' not exists.", insp.name, label))
	}
}

func (insp *Inspector) initRecord(values []*Value) Record {
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

func (insp *Inspector) InspectAndRecord(values ...*Value) uint {
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

func NewInspector(name string, size uint) *Inspector {
	insp := &Inspector{
		name:        name,
		records:     ScrollArray.NewScrollArray(size),
		autoTypes:   map[string]VType{},
		autoTypeGen: map[string]func() interface{}{},
		vTypes:      map[string]*TestingGroup{},
		sep:         " ",
	}

	_, _ = insp.newType(true, "_time", func() interface{} {
		return time.Now()
	}, func(v interface{}) string {
		return v.(time.Time).Format("2006/01/02 15:04:05")
	})

	_, _ = insp.newType(true, "_start", func() interface{} {
		return nil
	}, func(v interface{}) string {
		return ">"
	})

	_, _ = insp.newType(true, "_end", func() interface{} {
		return nil
	}, func(v interface{}) string {
		return "\n"
	})

	insp.SetOrders("_time") //初始化排序

	return insp
}
