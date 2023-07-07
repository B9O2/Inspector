package decorators

type Tag struct {
	label string
	data  interface{}
}

func (t Tag) Label() string {
	return t.label
}

func (t Tag) Data() interface{} {
	return t.data
}

func NewTag(label string, data interface{}) Tag {
	return Tag{
		label: label,
		data:  data,
	}
}

// Decorator 装饰器可以对值（Value）进行修饰
type Decorator struct {
	label string
	//todo 装饰器参数
	decoration func(i interface{}) interface{}
}

// Decorate 对值进行修饰
func (d Decorator) Decorate(i interface{}) (ret Tag) {
	defer func() {
		if r := recover(); r != nil {
			ret = Tag{
				label: d.label + ".error",
				data:  r.(error),
			}
		}
	}()
	data := d.decoration(i)
	return Tag{
		label: d.label,
		data:  data,
	}
}

func NewDecoration(label string, decoration func(i interface{}) interface{}) *Decorator {
	return &Decorator{
		label:      label,
		decoration: decoration,
	}
}
