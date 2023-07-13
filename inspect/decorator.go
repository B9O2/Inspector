package inspect

import (
	"fmt"
	"strings"
)

type Tag struct {
	name    string
	tagType string
	data    interface{}
}

func (t Tag) Name() string {
	return t.name
}

func (t Tag) Label() string {
	return t.name + "." + t.tagType
}

func (t Tag) Type() string {
	return t.tagType
}
func (t Tag) Data() interface{} {
	return t.data
}

func NewTag(name string, tagType string, data interface{}) Tag {
	return Tag{
		name:    name,
		tagType: tagType,
		data:    data,
	}
}

// Decorator 装饰器可以对值（Value）进行装饰
type Decorator struct {
	label string
	//todo 装饰器参数
	decoration func(i interface{}) interface{}
}

func (d Decorator) parseLabel() (string, string) {
	suffix := ""
	parts := strings.Split(d.label, ".")
	if len(parts) > 1 {
		suffix = parts[len(parts)-1]
	}
	return parts[0], suffix
}
func (d Decorator) Decorate(i interface{}) (ret Tag) {
	name, tagType := d.parseLabel()
	defer func() {
		if r := recover(); r != nil {
			ret = Tag{
				name:    d.label,
				tagType: "error",
				data:    r.(error),
			}
		}
	}()
	data := d.decoration(i)
	return Tag{
		name:    name,
		tagType: tagType,
		data:    data,
	}
}

// NewDecoration 初始化一个新的装饰器。标签需要使用后缀标记其类型，例如："MyRed.color"中包含".color"后缀，因此decoration应当返回colors.Color类型。
func NewDecoration(label string, decoration func(i interface{}) interface{}) *Decorator {
	return &Decorator{
		label:      label,
		decoration: decoration,
	}
}

type TestingReport struct {
	Success bool
	Title   string
	Detail  string
}

func (tr TestingReport) String() string {
	checkSign := "○"
	tip := ""
	if tr.Success {
		checkSign = "◉"
		tip = "SUCCESS"
	} else {
		checkSign = "○"
		tip = "FAILED"
	}
	return fmt.Sprintf("%s [%s] %s %s",
		checkSign,
		tr.Title,
		tr.Detail,
		tip,
	)
}
